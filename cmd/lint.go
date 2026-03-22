// cmd/lint.go
package cmd

import (
	"fmt"
	"os"

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

var lintSpec string

func init() {
	rootCmd.AddCommand(lintCmd)
	lintCmd.Flags().StringVar(&lintSpec, "spec", "", "OpenAPI spec file or URL (required)")
	_ = lintCmd.MarkFlagRequired("spec")
}

func runLint(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	loader := spec.NewLoader()
	parsedSpec, err := loader.Load(lintSpec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Failed to parse spec: %v\n", err)
		os.Exit(2)
	}

	issues := lint.RunAll(parsedSpec)
	if len(issues) == 0 {
		color.Green("✓ No lint issues found")
		return nil
	}

	hasError := false
	for _, iss := range issues {
		switch iss.Severity {
		case "error":
			color.Red("  [%s] %s: %s", iss.RuleID, iss.Path, iss.Message)
			hasError = true
		case "warning":
			color.Yellow("  [%s] %s: %s", iss.RuleID, iss.Path, iss.Message)
		}
	}

	shouldFail := hasError
	if cfg.Lint.FailOn == "warning" {
		shouldFail = len(issues) > 0
	}
	if shouldFail {
		os.Exit(3)
	}
	return nil
}
