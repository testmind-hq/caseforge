// cmd/run.go
package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/runner"
)

var runFormat string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute generated test cases (hurl or k6)",
	RunE: func(cmd *cobra.Command, args []string) error {
		varFlags, _ := cmd.Flags().GetStringArray("var")
		vars := runner.ParseVars(varFlags)

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

		for _, c := range result.Cases {
			if c.Passed {
				color.Green("  ✓ [%s] %s", c.ID, c.Title)
			} else {
				color.Red("  ✗ [%s] %s", c.ID, c.Title)
			}
		}
		total := result.Passed + result.Failed
		pct := 0
		if total > 0 {
			pct = result.Passed * 100 / total
		}
		if result.Failed == 0 {
			color.Green("✓ %d/%d tests passed (%d%%)", result.Passed, total, pct)
		} else {
			color.Yellow("✗ %d/%d tests passed (%d%%)", result.Passed, total, pct)
			os.Exit(6)
		}
		return nil
	},
}

var runCases string

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVar(&runCases, "cases", "", "Directory containing generated test files (required)")
	_ = runCmd.MarkFlagRequired("cases")
	runCmd.Flags().StringArray("var", nil, "Variables as key=value (repeatable)")
	runCmd.Flags().StringVar(&runFormat, "format", "hurl", "Test runner format: hurl|k6")
}
