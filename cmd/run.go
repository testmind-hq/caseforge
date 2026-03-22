// cmd/run.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/runner"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute generated test cases using hurl",
	RunE: func(cmd *cobra.Command, args []string) error {
		varFlags, _ := cmd.Flags().GetStringArray("var")
		vars := runner.ParseVars(varFlags)

		r := runner.NewHurlRunner()
		passed, failed, err := r.Run(runCases, vars)
		if err != nil {
			fmt.Fprintf(os.Stderr, "✗ %v\n", err)
			os.Exit(6)
		}
		total := passed + failed
		fmt.Fprintf(os.Stderr, "✓ %d/%d tests passed\n", passed, total)
		if failed > 0 {
			os.Exit(6)
		}
		return nil
	},
}

var runCases string

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVar(&runCases, "cases", "", "Directory containing .hurl files (required)")
	_ = runCmd.MarkFlagRequired("cases")
	runCmd.Flags().StringArray("var", nil, "Variables as key=value (repeatable)")
}
