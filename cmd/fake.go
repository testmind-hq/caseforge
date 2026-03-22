// cmd/fake.go
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/spec"
)

var fakeCmd = &cobra.Command{
	Use:   "fake",
	Short: "Generate fake data for a JSON schema",
	RunE:  runFake,
}

var fakeSchema string

func init() {
	rootCmd.AddCommand(fakeCmd)
	fakeCmd.Flags().StringVar(&fakeSchema, "schema", "", `Inline JSON schema (required)`)
	_ = fakeCmd.MarkFlagRequired("schema")
}

func runFake(cmd *cobra.Command, args []string) error {
	var s spec.Schema
	if err := json.Unmarshal([]byte(fakeSchema), &s); err != nil {
		return fmt.Errorf("parsing schema: %w", err)
	}
	g := datagen.NewGenerator(nil)
	val := g.Generate(&s, "")
	out, _ := json.MarshalIndent(val, "", "  ")
	fmt.Println(string(out))
	return nil
}
