// cmd/run.go
package cmd

import (
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute generated test cases using hurl",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Week 4
		cmd.Println("run: not yet implemented")
		return nil
	},
}

var runCases string

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVar(&runCases, "cases", "", "Directory containing generated cases (required)")
	_ = runCmd.MarkFlagRequired("cases")
	runCmd.Flags().StringArray("var", nil, "Variables in key=value format (repeatable)")
}
