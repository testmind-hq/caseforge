// cmd/gen.go
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/output/schema"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate test cases from an OpenAPI spec",
	RunE:  runGen,
}

var (
	genSpec   string
	genOutput string
	genNoAI   bool
	genFormat string
)

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.Flags().StringVar(&genSpec, "spec", "", "OpenAPI spec file or URL (required)")
	genCmd.Flags().StringVar(&genOutput, "output", "./cases", "Output directory")
	genCmd.Flags().BoolVar(&genNoAI, "no-ai", false, "Disable LLM, use pure algorithm mode")
	genCmd.Flags().StringVar(&genFormat, "format", "hurl", "Output format: hurl|markdown|csv")
	_ = genCmd.MarkFlagRequired("spec")
}

func runGen(cmd *cobra.Command, args []string) error {
	// STUB: return one hardcoded test case to prove the pipeline compiles
	tc := schema.TestCase{
		Schema:   schema.SchemaBaseURL,
		Version:  "1",
		ID:       "TC-0001",
		Title:    "STUB: spec not yet parsed",
		Kind:     "single",
		Priority: "P1",
		Source: schema.CaseSource{
			Technique: "stub",
			SpecPath:  genSpec,
			Rationale: "stub output — real implementation in Week 4",
		},
		Steps: []schema.Step{
			{
				ID:     "step-main",
				Title:  "stub step",
				Type:   "test",
				Method: "GET",
				Path:   "/stub",
				Assertions: []schema.Assertion{
					{Target: "status_code", Operator: "eq", Expected: 200},
				},
			},
		},
		GeneratedAt: time.Now(),
	}

	out, _ := json.MarshalIndent(tc, "", "  ")
	fmt.Fprintln(os.Stdout, string(out))
	return nil
}
