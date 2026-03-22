// cmd/gen.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/methodology"
	"github.com/testmind-hq/caseforge/internal/output/render"
	"github.com/testmind-hq/caseforge/internal/output/writer"
	"github.com/testmind-hq/caseforge/internal/spec"
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
	// 1. Build LLM provider
	var provider llm.LLMProvider
	if genNoAI {
		provider = &llm.NoopProvider{}
	} else {
		// NewProvider falls back to NoopProvider if no API key is set
		provider = llm.NewProvider("", "anthropic", "")
	}

	// 2. Load the spec
	loader := spec.NewLoader()
	parsedSpec, err := loader.Load(genSpec)
	if err != nil {
		return fmt.Errorf("loading spec: %w", err)
	}

	// 3. Build the methodology engine with all 6 techniques
	engine := methodology.NewEngine(
		provider,
		methodology.NewEquivalenceTechnique(),
		methodology.NewBoundaryTechnique(),
		methodology.NewDecisionTechnique(),
		methodology.NewStateTechnique(),
		methodology.NewIdempotentTechnique(),
		methodology.NewPairwiseTechnique(),
	)

	// 4. Generate test cases
	cases, err := engine.Generate(parsedSpec)
	if err != nil {
		return fmt.Errorf("generating test cases: %w", err)
	}

	// 5. Write index.json
	w := writer.NewJSONSchemaWriter()
	if err := w.Write(cases, genOutput); err != nil {
		return fmt.Errorf("writing index.json: %w", err)
	}

	// 6. Render output files
	var renderer render.Renderer
	switch genFormat {
	case "hurl":
		renderer = render.NewHurlRenderer("")
	case "markdown":
		renderer = render.NewMarkdownRenderer()
	case "csv":
		renderer = render.NewCSVRenderer()
	default:
		fmt.Fprintf(os.Stderr, "Warning: unknown format %q; falling back to hurl\n", genFormat)
		renderer = render.NewHurlRenderer("")
	}

	if err := renderer.Render(cases, genOutput); err != nil {
		return fmt.Errorf("rendering %s files: %w", renderer.Format(), err)
	}

	fmt.Fprintf(os.Stderr, "✓ Generated %d test cases → %s\n", len(cases), genOutput)
	return nil
}
