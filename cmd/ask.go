// cmd/ask.go
package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/ask"
	"github.com/testmind-hq/caseforge/internal/config"
	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/output/render"
	"github.com/testmind-hq/caseforge/internal/output/writer"
)

var askCmd = &cobra.Command{
	Use:   "ask <description>",
	Short: "Generate test cases from a natural language description",
	Long:  `Generate test cases from a natural language description using an LLM provider.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  runAsk,
}

var (
	askOutput string
	askFormat string
)

func init() {
	rootCmd.AddCommand(askCmd)
	askCmd.Flags().StringVar(&askOutput, "output", "./cases", "Output directory")
	askCmd.Flags().StringVar(&askFormat, "format", "hurl", "Output format: hurl|markdown|csv|postman|k6")
}

func runAsk(cmd *cobra.Command, args []string) error {
	description := strings.Join(args, " ")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if cmd.Flags().Changed("format") {
		cfg.Output.DefaultFormat = askFormat
	}

	provider := llm.NewProviderWithConfig(llm.ProviderConfig{
		APIKey:   cfg.AI.APIKey,
		Provider: cfg.AI.Provider,
		Model:    cfg.AI.Model,
		BaseURL:  cfg.AI.BaseURL,
		Region:   cfg.AI.Region,
	})

	gen := ask.NewGenerator(provider)

	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	cases, err := gen.Generate(ctx, description)
	if err != nil {
		return err
	}

	w := writer.NewJSONSchemaWriter()
	if err := w.Write(cases, askOutput, writer.WriteOptions{CaseforgeVersion: Version}); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	var renderer render.Renderer
	switch cfg.Output.DefaultFormat {
	case "markdown":
		renderer = render.NewMarkdownRenderer()
	case "csv":
		renderer = render.NewCSVRenderer()
	case "postman":
		renderer = render.NewPostmanRenderer()
	case "k6":
		renderer = render.NewK6Renderer()
	default: // "hurl" and anything unrecognised
		renderer = render.NewHurlRenderer("")
	}
	if err := renderer.Render(cases, askOutput); err != nil {
		return fmt.Errorf("rendering output: %w", err)
	}

	fmt.Fprintf(cmd.ErrOrStderr(), "✓ Generated %d test cases → %s\n", len(cases), askOutput)
	return nil
}
