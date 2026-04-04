// cmd/run.go
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/runner"
)

var runFormat string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute generated test cases (hurl or k6)",
	Long: `Run executes generated test case files against a live API.

Supports hurl (.hurl files) and k6 (k6_tests.js) formats.
Use --target to set the API base URL (injected as BASE_URL variable).
Use --output to write a run-report.json with structured results.

Exit codes:
  0 — all tests passed
  6 — one or more tests failed

Examples:
  caseforge run --cases ./cases --target http://localhost:8080
  caseforge run --cases ./cases --target http://localhost:8080 --format k6
  caseforge run --cases ./cases --target http://localhost:8080 --output ./reports`,
	RunE: runRun,
}

var runCases string

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVar(&runCases, "cases", "", "Directory containing generated test files (required)")
	_ = runCmd.MarkFlagRequired("cases")
	runCmd.Flags().StringArray("var", nil, "Variables as key=value (repeatable)")
	runCmd.Flags().StringVar(&runFormat, "format", "hurl", "Test runner format: hurl|k6")
	runCmd.Flags().String("target", "", "API base URL, e.g. http://localhost:8080 (injected as BASE_URL)")
	runCmd.Flags().String("output", "", "Directory to write run-report.json (optional)")
}

func runRun(cmd *cobra.Command, _ []string) error {
	varFlags, _ := cmd.Flags().GetStringArray("var")
	target, _ := cmd.Flags().GetString("target")
	outputDir, _ := cmd.Flags().GetString("output")

	vars := runner.ParseVars(varFlags)
	if target != "" {
		if _, alreadySet := vars["BASE_URL"]; alreadySet {
			fmt.Fprintf(cmd.ErrOrStderr(), "warning: --target overrides BASE_URL set via --var\n")
		}
		vars["BASE_URL"] = target
	}

	var r runner.Runner
	switch runFormat {
	case "k6":
		r = runner.NewK6Runner()
	default:
		r = runner.NewHurlRunner()
	}

	result, err := r.Run(runCases, vars)
	if err != nil {
		return fmt.Errorf("running tests: %w", err)
	}

	out := cmd.OutOrStdout()
	for _, c := range result.Cases {
		if c.Passed {
			color.New(color.FgGreen).Fprintf(out, "  ✓ [%s] %s\n", c.ID, c.Title)
		} else {
			color.New(color.FgRed).Fprintf(out, "  ✗ [%s] %s\n", c.ID, c.Title)
		}
	}
	total := result.Passed + result.Failed
	pct := 0
	if total > 0 {
		pct = result.Passed * 100 / total
	}
	if result.Failed == 0 {
		color.New(color.FgGreen).Fprintf(out, "✓ %d/%d tests passed (%d%%)\n", result.Passed, total, pct)
	} else {
		color.New(color.FgYellow).Fprintf(out, "✗ %d/%d tests passed (%d%%)\n", result.Passed, total, pct)
	}

	if outputDir != "" {
		if err := writeRunReport(outputDir, target, runCases, runFormat, result); err != nil {
			return fmt.Errorf("write report: %w", err)
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "Report written to: %s\n", filepath.Join(outputDir, "run-report.json"))
	}

	if result.Failed > 0 {
		os.Exit(6)
	}
	return nil
}

type runReport struct {
	Target      string            `json:"target"`
	CasesDir    string            `json:"cases_dir"`
	Format      string            `json:"format"`
	Passed      int               `json:"passed"`
	Failed      int               `json:"failed"`
	Total       int               `json:"total"`
	PassPct     int               `json:"pass_pct"`
	Cases       []runner.CaseResult `json:"cases"`
	GeneratedAt string            `json:"generated_at"`
}

func writeRunReport(outputDir, target, casesDir, format string, result runner.RunResult) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}
	total := result.Passed + result.Failed
	pct := 0
	if total > 0 {
		pct = result.Passed * 100 / total
	}
	report := runReport{
		Target:      target,
		CasesDir:    casesDir,
		Format:      format,
		Passed:      result.Passed,
		Failed:      result.Failed,
		Total:       total,
		PassPct:     pct,
		Cases:       result.Cases,
		GeneratedAt: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outputDir, "run-report.json"), data, 0644)
}
