// cmd/lint.go
package cmd

import (
	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint an OpenAPI spec for quality issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Week 4
		cmd.Println("lint: not yet implemented")
		return nil
	},
}

var lintSpec string

func init() {
	rootCmd.AddCommand(lintCmd)
	lintCmd.Flags().StringVar(&lintSpec, "spec", "", "OpenAPI spec file or URL (required)")
	_ = lintCmd.MarkFlagRequired("spec")
}
