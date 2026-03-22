// cmd/fake.go
package cmd

import (
	"github.com/spf13/cobra"
)

var fakeCmd = &cobra.Command{
	Use:   "fake",
	Short: "Generate fake data for a given JSON schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Week 4
		cmd.Println("fake: not yet implemented")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fakeCmd)
	fakeCmd.Flags().String("schema", "", "Inline JSON schema")
}
