// cmd/dedupe.go
package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/dedupe"
)

var dedupeCmd = &cobra.Command{
	Use:   "dedupe",
	Short: "Detect and optionally remove duplicate test cases",
	Long: `Dedupe scans a directory of CaseForge JSON test case files and reports
exact and structural duplicates.

Exact duplicates share the same method, path, expected status, and request body.
Structural duplicates share the same method, path, and status, and their assertion
field names have a Jaccard similarity >= threshold.

Use --merge to auto-delete lower-scoring duplicates (keeps the case with the most
assertions; ties broken by lexicographic filename order).
Use --dry-run to report intended deletions without touching any files.

Exit codes:
  0 — no duplicates, or --dry-run (even when duplicates are found)
  1 — duplicates found and --dry-run is not set

Examples:
  caseforge dedupe --cases ./cases
  caseforge dedupe --cases ./cases --threshold 0.85 --merge
  caseforge dedupe --cases ./cases --dry-run --format json`,
	RunE:         runDedupe,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(dedupeCmd)
	dedupeCmd.Flags().String("cases", "./cases", "Directory of test case JSON files")
	dedupeCmd.Flags().Float64("threshold", 0.9, "Jaccard similarity threshold for structural duplicates (0.0–1.0)")
	dedupeCmd.Flags().Bool("merge", false, "Auto-delete lower-scoring duplicates")
	dedupeCmd.Flags().Bool("dry-run", false, "Report what would be deleted without deleting")
	dedupeCmd.Flags().String("format", "terminal", "Output format: terminal or json")
}

func runDedupe(cmd *cobra.Command, _ []string) error {
	casesDir, _ := cmd.Flags().GetString("cases")
	threshold, _ := cmd.Flags().GetFloat64("threshold")
	merge, _ := cmd.Flags().GetBool("merge")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	format, _ := cmd.Flags().GetString("format")

	stdout := cmd.OutOrStdout()
	stderr := cmd.ErrOrStderr()

	cases, err := dedupe.ScanCases(casesDir)
	if err != nil {
		return fmt.Errorf("scan cases: %w", err)
	}

	groups, err := dedupe.FindDuplicates(cases, threshold)
	if err != nil {
		return fmt.Errorf("find duplicates: %w", err)
	}

	exactCount, structCount := 0, 0
	for _, g := range groups {
		if g.Kind == dedupe.MatchExact {
			exactCount++
		} else {
			structCount++
		}
	}
	report := dedupe.DedupeReport{
		CasesDir:         casesDir,
		TotalScanned:     len(cases),
		ExactGroups:      exactCount,
		StructuralGroups: structCount,
		Groups:           groups,
		GeneratedAt:      time.Now(),
	}

	if format == "json" {
		data, marshalErr := dedupe.MarshalReportJSON(report)
		if marshalErr != nil {
			return fmt.Errorf("marshal report: %w", marshalErr)
		}
		fmt.Fprintln(stdout, string(data))
	} else {
		dedupe.PrintTerminal(stdout, report, dryRun)
	}

	if len(groups) > 0 && (merge || dryRun) {
		opts := dedupe.MergeOptions{DryRun: dryRun, Out: stderr}
		if _, mergeErr := dedupe.Merge(groups, opts); mergeErr != nil {
			return fmt.Errorf("merge: %w", mergeErr)
		}
	}

	if len(groups) > 0 && !dryRun {
		return fmt.Errorf("duplicates found: %d group(s)", len(groups))
	}
	return nil
}
