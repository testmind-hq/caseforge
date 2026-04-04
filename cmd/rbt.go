// cmd/rbt.go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/rbt"
	"github.com/testmind-hq/caseforge/internal/spec"
)

var rbtCmd = &cobra.Command{
	Use:   "rbt",
	Short: "Risk-Based Testing: assess which API operations are at risk from recent changes",
	Long: `RBT analyses the git diff between two commits and cross-references changed
source files with your OpenAPI spec and existing test cases to produce a risk
report. Operations with no tests covering changed code are flagged as high risk.

Examples:
  caseforge rbt --spec openapi.yaml
  caseforge rbt --spec openapi.yaml --base HEAD~3 --head HEAD --format json
  caseforge rbt --spec openapi.yaml --dry-run --output ./reports`,
	RunE: runRBT,
}

func init() {
	rootCmd.AddCommand(rbtCmd)
	rbtCmd.Flags().String("spec", "", "OpenAPI spec file (required)")
	rbtCmd.Flags().String("cases", "./cases", "Directory containing test case JSON files")
	rbtCmd.Flags().String("src", "./", "Source code root directory")
	rbtCmd.Flags().String("base", "HEAD~1", "Base git ref for diff")
	rbtCmd.Flags().String("head", "HEAD", "Head git ref for diff")
	rbtCmd.Flags().Bool("generate", false, "Generate test cases for high-risk uncovered operations")
	rbtCmd.Flags().String("output", "./reports", "Output directory for rbt-report.json")
	rbtCmd.Flags().String("format", "terminal", "Output format: terminal or json")
	rbtCmd.Flags().String("fail-on", "high", "Exit non-zero if any operation has risk >= level (none|low|medium|high)")
	rbtCmd.Flags().String("map", "", "Path to caseforge-map.yaml (default: <src>/caseforge-map.yaml)")
	rbtCmd.Flags().Bool("dry-run", false, "Skip git diff and tree-sitter; report all operations as risk=none")
}

func runRBT(cmd *cobra.Command, _ []string) error {
	specPath, _ := cmd.Flags().GetString("spec")
	casesDir, _ := cmd.Flags().GetString("cases")
	srcDir, _ := cmd.Flags().GetString("src")
	base, _ := cmd.Flags().GetString("base")
	head, _ := cmd.Flags().GetString("head")
	outputDir, _ := cmd.Flags().GetString("output")
	format, _ := cmd.Flags().GetString("format")
	failOn, _ := cmd.Flags().GetString("fail-on")
	mapFile, _ := cmd.Flags().GetString("map")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	generate, _ := cmd.Flags().GetBool("generate")

	if specPath == "" {
		return fmt.Errorf("--spec is required")
	}

	// Load spec
	loader := spec.NewLoader()
	parsedSpec, err := loader.Load(specPath)
	if err != nil {
		return fmt.Errorf("load spec: %w", err)
	}

	var report rbt.RiskReport

	if dryRun {
		// Build a report with all operations at RiskNone
		ops := make([]rbt.OperationCoverage, 0, len(parsedSpec.Operations))
		for _, op := range parsedSpec.Operations {
			ops = append(ops, rbt.OperationCoverage{
				OperationID: op.OperationID,
				Method:      op.Method,
				Path:        op.Path,
				Affected:    false,
				Risk:        rbt.RiskNone,
			})
		}
		report = rbt.RiskReport{
			DiffBase:    base,
			DiffHead:    head,
			Operations:  ops,
			GeneratedAt: time.Now(),
		}
	} else {
		// Git diff
		changedFiles, err := rbt.Diff(srcDir, base, head)
		if err != nil {
			return fmt.Errorf("git diff: %w", err)
		}

		// Build parser chain
		if mapFile == "" {
			mapFile = filepath.Join(srcDir, "caseforge-map.yaml")
		}
		parsers := []rbt.SourceParser{
			rbt.NewMapFileParser(mapFile),
			rbt.NewTreeSitterParser(),
			rbt.NewRegexParser(),
			// LLM provider is nil in V1 — LLMParser gracefully returns empty when no
		// provider is configured. Wire a real provider here in V2 for LLM inference.
		rbt.NewLLMParser(nil, ""),
		}

		// Map changed files to route mappings
		mappings, err := rbt.MapChain(parsers, srcDir, changedFiles)
		if err != nil {
			return fmt.Errorf("map chain: %w", err)
		}

		// Scan test cases
		caseIndex, err := rbt.ScanCases(casesDir)
		if err != nil {
			return fmt.Errorf("scan cases: %w", err)
		}

		// Assess risk
		report = rbt.Assess(parsedSpec, mappings, caseIndex, base, head, changedFiles)
	}

	// Output
	out := cmd.OutOrStdout()
	errOut := cmd.ErrOrStderr()
	if format == "json" {
		data, err := rbt.MarshalReportJSON(report)
		if err != nil {
			return fmt.Errorf("marshal report: %w", err)
		}
		fmt.Fprintln(out, string(data))
	} else {
		rbt.PrintTerminal(out, report)
	}

	// Always write report JSON
	reportPath, err := rbt.WriteReportJSON(outputDir, report)
	if err != nil {
		return fmt.Errorf("write report: %w", err)
	}
	// Print to stderr when --format json so stdout is valid JSON for piping.
	if format == "json" {
		fmt.Fprintln(errOut, "Report written to:", reportPath)
	} else {
		fmt.Fprintln(out, "Report written to:", reportPath)
	}

	// --generate: auto-generate tests for HIGH-risk operations
	if generate && !dryRun {
		if err := rbt.GenerateForHighRisk(out, report, casesDir); err != nil {
			return err
		}
	}

	// Exit non-zero for high-risk (only in non-dry-run path)
	if !dryRun && rbt.ShouldFail(report, failOn) {
		os.Exit(1)
	}

	return nil
}
