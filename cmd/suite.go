// cmd/suite.go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	suitepkg "github.com/testmind-hq/caseforge/internal/suite"
	"github.com/testmind-hq/caseforge/internal/output/writer"
)

var suiteCmd = &cobra.Command{
	Use:   "suite",
	Short: "Manage TestSuite orchestration files (create, validate)",
}

// ── suite create ─────────────────────────────────────────────────────────────

var suiteCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new suite.json file",
	RunE:  runSuiteCreate,
}

var (
	suiteCreateID     string
	suiteCreateTitle  string
	suiteCreateKind   string
	suiteCreateCases  string
	suiteCreateOutput string
)

func init() {
	rootCmd.AddCommand(suiteCmd)
	suiteCmd.AddCommand(suiteCreateCmd)
	suiteCmd.AddCommand(suiteValidateCmd)

	suiteCreateCmd.Flags().StringVar(&suiteCreateID, "id", "", "Suite ID (required)")
	suiteCreateCmd.Flags().StringVar(&suiteCreateTitle, "title", "", "Suite title (required)")
	suiteCreateCmd.Flags().StringVar(&suiteCreateKind, "kind", "sequential", "Execution kind: sequential|chain")
	suiteCreateCmd.Flags().StringVar(&suiteCreateCases, "cases", "", "Comma-separated case IDs to include")
	suiteCreateCmd.Flags().StringVar(&suiteCreateOutput, "output", "suite.json", "Output file path")
	_ = suiteCreateCmd.MarkFlagRequired("id")
	_ = suiteCreateCmd.MarkFlagRequired("title")
}

func runSuiteCreate(cmd *cobra.Command, _ []string) error {
	if suiteCreateKind != "sequential" && suiteCreateKind != "chain" {
		return fmt.Errorf("invalid --kind %q: must be sequential or chain", suiteCreateKind)
	}

	var cases []schema.SuiteCase
	if suiteCreateCases != "" {
		for _, id := range splitTrimmed(suiteCreateCases) {
			cases = append(cases, schema.SuiteCase{CaseID: id})
		}
	}

	s := &schema.TestSuite{
		ID:    suiteCreateID,
		Title: suiteCreateTitle,
		Kind:  suiteCreateKind,
		Cases: cases,
	}

	// Write to output file (create parent directories if needed).
	outPath := suiteCreateOutput
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}
	if err := suitepkg.WriteSuiteFile(s, outPath); err != nil {
		return err
	}

	n := len(cases)
	fmt.Fprintf(cmd.OutOrStdout(), "Created suite %q (%s, %d case(s)) → %s\n",
		s.ID, s.Kind, n, outPath)
	return nil
}

// ── suite validate ───────────────────────────────────────────────────────────

var suiteValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a suite.json against its index.json",
	RunE:  runSuiteValidate,
}

var (
	suiteValidateSuite string
	suiteValidateCases string
)

func init() {
	suiteValidateCmd.Flags().StringVar(&suiteValidateSuite, "suite", "", "Path to suite.json (required)")
	suiteValidateCmd.Flags().StringVar(&suiteValidateCases, "cases", "", "Cases dir containing index.json (optional)")
	_ = suiteValidateCmd.MarkFlagRequired("suite")
}

func runSuiteValidate(cmd *cobra.Command, _ []string) error {
	s, err := suitepkg.LoadSuiteFile(suiteValidateSuite)
	if err != nil {
		return err
	}

	// Load known cases from index.json if --cases provided.
	var knownCases []schema.TestCase
	if suiteValidateCases != "" {
		indexPath := filepath.Join(suiteValidateCases, "index.json")
		knownCases, err = writer.NewJSONSchemaWriter().Read(indexPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not read %s: %v\n", indexPath, err)
		}
	}

	errs := suitepkg.Validate(s, knownCases)
	if len(errs) == 0 {
		// Also compute topological order to surface cycle errors even if Validate passes.
		order, topoErr := suitepkg.TopologicalOrder(s)
		if topoErr != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "✗ %s\n", topoErr)
			return topoErr
		}
		printSuiteValidationOK(cmd, s, order)
		return nil
	}

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Suite %q validation failed (%d error(s)):\n", s.ID, len(errs))
	for _, e := range errs {
		fmt.Fprintf(out, "  ✗ %s\n", e)
	}

	return fmt.Errorf("suite validation failed: %d error(s)", len(errs))
}

func printSuiteValidationOK(cmd *cobra.Command, s *schema.TestSuite, order []schema.SuiteCase) {
	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Suite %q is valid ✓\n", s.ID)
	fmt.Fprintf(out, "  kind:  %s\n", s.Kind)
	fmt.Fprintf(out, "  cases: %d\n", len(s.Cases))
	if len(order) > 0 {
		ids := make([]string, len(order))
		for i, sc := range order {
			ids[i] = sc.CaseID
		}
		fmt.Fprintf(out, "  order: %s\n", strings.Join(ids, " → "))
	}
}
