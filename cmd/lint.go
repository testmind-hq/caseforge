// cmd/lint.go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/config"
	"github.com/testmind-hq/caseforge/internal/lint"
	"github.com/testmind-hq/caseforge/internal/spec"
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint an OpenAPI spec for quality issues",
	RunE:  runLint,
}

var (
	lintSpec      string
	lintMinScore  int
	lintFormat    string
	lintOutput    string
	lintSkipRules []string
)

func init() {
	rootCmd.AddCommand(lintCmd)
	lintCmd.Flags().StringVar(&lintSpec, "spec", "", "OpenAPI spec file or URL (required)")
	lintCmd.Flags().IntVar(&lintMinScore, "min-score", 0, "Fail if spec score is below this threshold (0 = disabled)")
	lintCmd.Flags().StringVar(&lintFormat, "format", "terminal", "Output format: terminal|json")
	lintCmd.Flags().StringVar(&lintOutput, "output", "", "Write lint-report.json to this directory")
	lintCmd.Flags().StringSliceVar(&lintSkipRules, "skip-rules", nil, "Comma-separated rule IDs to skip (e.g. L001,L003)")
	_ = lintCmd.MarkFlagRequired("spec")
}

func runLint(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Load .caseforgelint.yaml from working directory
	fileCfg, err := lint.LoadLintFileConfig(".")
	if err != nil {
		return fmt.Errorf("loading .caseforgelint.yaml: %w", err)
	}

	// Validate --format early, before any expensive work
	if lintFormat != "terminal" && lintFormat != "json" {
		return fmt.Errorf("unknown format %q: use terminal or json", lintFormat)
	}

	// Resolve fail_on: .caseforgelint.yaml > caseforge.yaml
	failOn := cfg.Lint.FailOn
	if fileCfg.FailOn != "" {
		failOn = fileCfg.FailOn
	}
	if failOn != "" && failOn != "error" && failOn != "warning" {
		return fmt.Errorf("invalid fail_on value %q: must be \"error\" or \"warning\"", failOn)
	}

	// Build skip set: union of caseforge.yaml + .caseforgelint.yaml + --skip-rules flag
	skip := make(map[string]bool)
	for _, id := range cfg.Lint.SkipRules {
		skip[strings.TrimSpace(id)] = true
	}
	for _, id := range fileCfg.SkipRules {
		skip[strings.TrimSpace(id)] = true
	}
	for _, id := range lintSkipRules {
		skip[strings.TrimSpace(id)] = true
	}

	loader := spec.NewLoader()
	parsedSpec, err := loader.Load(lintSpec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Failed to parse spec: %v\n", err)
		os.Exit(2)
	}

	issues := lint.RunAll(parsedSpec, skip)
	report := lint.NewReport(issues)

	// Serialize JSON once; reused for both --format json stdout and --output file.
	var reportJSON []byte
	if lintFormat == "json" || lintOutput != "" {
		reportJSON, err = report.ToJSON()
		if err != nil {
			return fmt.Errorf("serialising report: %w", err)
		}
	}

	// Render output
	if lintFormat == "json" {
		fmt.Println(string(reportJSON))
	} else {
		// terminal (coloured)
		for _, iss := range issues {
			switch iss.Severity {
			case "error":
				color.Red("  ✗ [%s] %s: %s", iss.RuleID, iss.Path, iss.Message)
			case "warning":
				color.Yellow("  ⚠ [%s] %s: %s", iss.RuleID, iss.Path, iss.Message)
			}
		}
		if len(issues) == 0 {
			color.Green("✓ No lint issues found")
		}
		fmt.Fprintf(os.Stderr, "\nSpec Score: %d/100  (%d errors, %d warnings)\n",
			report.Score, report.ErrorCount, report.WarningCount)
	}

	// Write file report if --output given
	if lintOutput != "" {
		if err := os.MkdirAll(lintOutput, 0755); err != nil {
			return fmt.Errorf("creating output dir: %w", err)
		}
		outPath := filepath.Join(lintOutput, "lint-report.json")
		if err := os.WriteFile(outPath, reportJSON, 0644); err != nil {
			return fmt.Errorf("writing lint-report.json: %w", err)
		}
		fmt.Fprintf(os.Stderr, "✓ Report written to %s\n", outPath)
	}

	// Determine exit
	shouldFail := report.ErrorCount > 0
	if failOn == "warning" {
		shouldFail = len(issues) > 0
	}
	if !shouldFail && lintMinScore > 0 && report.Score < lintMinScore {
		fmt.Fprintf(os.Stderr, "score %d < min-score %d\n", report.Score, lintMinScore)
		shouldFail = true
	}
	if shouldFail {
		os.Exit(1)
	}
	return nil
}
