// cmd/diff.go
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/diff"
	"github.com/testmind-hq/caseforge/internal/output/schema"
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
	diffOld    string
	diffNew    string
	diffCases  string
	diffFormat string
)

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().StringVar(&diffOld, "old", "", "Old spec file (required)")
	diffCmd.Flags().StringVar(&diffNew, "new", "", "New spec file (required)")
	diffCmd.Flags().StringVar(&diffCases, "cases", "", "Cases output dir (optional; reads index.json to infer affected cases)")
	diffCmd.Flags().StringVar(&diffFormat, "format", "text", "Output format: text|json")
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
		cases, readErr := loadCasesFromIndex(indexPath)
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

	if result.HasBreaking() {
		return errBreakingChanges
	}
	return nil
}

func loadCasesFromIndex(indexPath string) ([]schema.TestCase, error) {
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, err
	}
	var cases []schema.TestCase
	if err := json.Unmarshal(data, &cases); err != nil {
		return nil, err
	}
	return cases, nil
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
