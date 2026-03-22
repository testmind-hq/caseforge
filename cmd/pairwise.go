// cmd/pairwise.go
package cmd

import (
	"github.com/spf13/cobra"
)

var pairwiseCmd = &cobra.Command{
	Use:   "pairwise",
	Short: "Compute pairwise combinations for given parameters",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Week 4
		cmd.Println("pairwise: not yet implemented")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pairwiseCmd)
	pairwiseCmd.Flags().String("params", "", `Parameters in "name:v1,v2 name2:v3,v4" format`)
}
