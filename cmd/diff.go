// cmd/diff.go
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/diff"
	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/methodology"
	"github.com/testmind-hq/caseforge/internal/output/writer"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// errBreakingChanges is returned by runDiff when breaking changes are detected.
// root.Execute() converts any non-nil error to os.Exit(1).
var errBreakingChanges = errors.New("breaking changes detected")

var diffCmd = &cobra.Command{
	Use:           "diff",
	Short:         "Compare two OpenAPI specs and classify breaking changes",
	RunE:          runDiff,
	SilenceErrors: true, // suppress cobra's "Error: breaking changes detected" message
	SilenceUsage:  true, // suppress usage on error
}

var (
	diffOld      string
	diffNew      string
	diffCases    string
	diffFormat   string
	diffGenCases string
)

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().StringVar(&diffOld, "old", "", "Old spec file (required)")
	diffCmd.Flags().StringVar(&diffNew, "new", "", "New spec file (required)")
	diffCmd.Flags().StringVar(&diffCases, "cases", "", "Cases output dir (optional; reads index.json to infer affected cases)")
	diffCmd.Flags().StringVar(&diffFormat, "format", "text", "Output format: text|json")
	diffCmd.Flags().StringVar(&diffGenCases, "gen-cases", "", "Generate test cases for breaking operations into this directory")
	_ = diffCmd.MarkFlagRequired("old")
	_ = diffCmd.MarkFlagRequired("new")
}

func runDiff(cmd *cobra.Command, args []string) error {
	loader := spec.NewLoader()

	oldSpec, err := loader.Load(diffOld)
	if err != nil {
		return fmt.Errorf("loading old spec: %w", err)
	}
	newSpec, err := loader.Load(diffNew)
	if err != nil {
		return fmt.Errorf("loading new spec: %w", err)
	}

	result := diff.Diff(oldSpec, newSpec)

	// Load affected cases if --cases provided
	var affected []diff.AffectedCase
	if diffCases != "" {
		indexPath := filepath.Join(diffCases, "index.json")
		cases, readErr := writer.NewJSONSchemaWriter().Read(indexPath)
		if readErr != nil {
			fmt.Fprintf(os.Stderr, "warning: could not read %s: %v\n", indexPath, readErr)
		} else {
			affected = diff.Suggest(result, cases)
		}
	}

	switch diffFormat {
	case "json":
		if err := printDiffJSON(cmd, result, affected); err != nil {
			return err
		}
	default:
		printDiffText(cmd, result, affected)
	}

	// --gen-cases: generate test cases for breaking operations into the given dir
	if diffGenCases != "" && result.HasBreaking() {
		if err := generateCasesForBreakingChanges(result, newSpec, diffGenCases, cmd.OutOrStdout()); err != nil {
			fmt.Fprintf(os.Stderr, "warning: gen-cases failed: %v\n", err)
		}
	}

	if result.HasBreaking() {
		return errBreakingChanges
	}
	return nil
}

// generateCasesForBreakingChanges filters newSpec to operations that have
// breaking or potentially-breaking changes, generates test cases for them
// using the algorithm-only engine (no API key required), and writes index.json
// to outDir.
func generateCasesForBreakingChanges(result diff.DiffResult, newSpec *spec.ParsedSpec, outDir string, out io.Writer) error {
	// Collect operations with breaking changes.
	breakingOps := make(map[string]bool)
	for _, c := range result.Changes {
		if (c.Kind == diff.Breaking || c.Kind == diff.PotentiallyBreaking) && c.Method != "" {
			breakingOps[c.Method+" "+c.Path] = true
		}
	}
	if len(breakingOps) == 0 {
		return nil
	}

	// Filter newSpec to only the affected operations.
	var ops []*spec.Operation
	for _, op := range newSpec.Operations {
		if breakingOps[op.Method+" "+op.Path] {
			ops = append(ops, op)
		}
	}
	if len(ops) == 0 {
		return nil
	}

	filteredSpec := &spec.ParsedSpec{Operations: ops}

	// Use noop provider — no API key required, algorithm-only generation.
	provider := &llm.NoopProvider{}
	engine := methodology.NewEngine(provider,
		methodology.NewEquivalenceTechnique(),
		methodology.NewBoundaryTechnique(),
	)
	cases, err := engine.Generate(filteredSpec)
	if err != nil {
		return fmt.Errorf("generating cases: %w", err)
	}

	w := writer.NewJSONSchemaWriter()
	if err := w.Write(cases, outDir, writer.WriteOptions{}); err != nil {
		return fmt.Errorf("writing cases: %w", err)
	}

	fmt.Fprintf(out, "\nGenerated %d test case(s) for %d breaking operation(s) → %s\n",
		len(cases), len(ops), outDir)
	return nil
}

func printDiffText(cmd *cobra.Command, result diff.DiffResult, affected []diff.AffectedCase) {
	// Group by kind
	var breaking, potBreaking, nonBreaking []diff.Change
	for _, c := range result.Changes {
		switch c.Kind {
		case diff.Breaking:
			breaking = append(breaking, c)
		case diff.PotentiallyBreaking:
			potBreaking = append(potBreaking, c)
		case diff.NonBreaking:
			nonBreaking = append(nonBreaking, c)
		}
	}

	out := cmd.OutOrStdout()

	if len(breaking) > 0 {
		fmt.Fprintf(out, "\nBREAKING (%d):\n", len(breaking))
		for _, c := range breaking {
			loc := ""
			if c.Location != "" {
				loc = " " + c.Location
			}
			fmt.Fprintf(out, "  x %-8s %-30s %s\n", c.Method, c.Path+loc, c.Description)
		}
	}
	if len(potBreaking) > 0 {
		fmt.Fprintf(out, "\nPOTENTIALLY BREAKING (%d):\n", len(potBreaking))
		for _, c := range potBreaking {
			loc := ""
			if c.Location != "" {
				loc = " " + c.Location
			}
			fmt.Fprintf(out, "  ! %-8s %-30s %s\n", c.Method, c.Path+loc, c.Description)
		}
	}
	if len(nonBreaking) > 0 {
		fmt.Fprintf(out, "\nNON-BREAKING (%d):\n", len(nonBreaking))
		for _, c := range nonBreaking {
			fmt.Fprintf(out, "  + %-8s %-30s %s\n", c.Method, c.Path, c.Description)
		}
	}

	if len(affected) > 0 {
		fmt.Fprintf(out, "\nAffected test cases:\n")
		for _, a := range affected {
			fmt.Fprintf(out, "  %s  %s - %s\n", a.ID, a.Title, a.Reason)
		}
	}

	if len(result.Changes) == 0 {
		fmt.Fprintln(out, "No changes detected.")
	}
}

type jsonDiffOutput struct {
	Summary struct {
		Breaking            int `json:"breaking"`
		PotentiallyBreaking int `json:"potentially_breaking"`
		NonBreaking         int `json:"non_breaking"`
	} `json:"summary"`
	Changes       []diff.Change       `json:"changes"`
	AffectedCases []diff.AffectedCase `json:"affected_cases,omitempty"`
}

func printDiffJSON(cmd *cobra.Command, result diff.DiffResult, affected []diff.AffectedCase) error {
	out := jsonDiffOutput{}
	out.Changes = result.Changes
	out.AffectedCases = affected
	for _, c := range result.Changes {
		switch c.Kind {
		case diff.Breaking:
			out.Summary.Breaking++
		case diff.PotentiallyBreaking:
			out.Summary.PotentiallyBreaking++
		case diff.NonBreaking:
			out.Summary.NonBreaking++
		}
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return nil
}
