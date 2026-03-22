// cmd/gen.go
package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/config"
	"github.com/testmind-hq/caseforge/internal/event"
	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/methodology"
	"github.com/testmind-hq/caseforge/internal/output/render"
	"github.com/testmind-hq/caseforge/internal/output/writer"
	"github.com/testmind-hq/caseforge/internal/spec"
	"github.com/testmind-hq/caseforge/internal/tui"
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
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Override config from flags
	if genNoAI {
		cfg.AI.Provider = "noop"
	}
	if cmd.Flags().Changed("format") {
		cfg.Output.DefaultFormat = genFormat
	}

	// Resolve LLM provider
	provider := llm.NewProvider(cfg.AI.APIKey, cfg.AI.Provider, cfg.AI.Model)
	if cfg.AI.Provider != "noop" && !provider.IsAvailable() {
		fmt.Fprintln(os.Stderr, "warning: LLM provider unavailable, degrading to --no-ai mode")
		provider = &llm.NoopProvider{}
	}

	// Load spec
	loader := spec.NewLoader()
	parsedSpec, err := loader.Load(genSpec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Failed to parse spec: %v\n", err)
		os.Exit(2)
	}

	// Set up event bus
	bus := event.NewBus()

	// Wire TUI if stderr is a terminal
	var tuiDone <-chan struct{}
	if isatty.IsTerminal(os.Stderr.Fd()) {
		model := tui.NewProgressModel(len(parsedSpec.Operations))
		prog := tea.NewProgram(model, tea.WithOutput(os.Stderr))
		sink := tui.NewTUISink(prog)
		bus.Subscribe(sink)
		doneCh := make(chan struct{})
		go func() {
			defer close(doneCh)
			_, _ = prog.Run()
		}()
		tuiDone = doneCh
	}

	// Generate test cases
	engine := methodology.NewEngine(provider,
		methodology.NewEquivalenceTechnique(),
		methodology.NewBoundaryTechnique(),
		methodology.NewDecisionTechnique(),
		methodology.NewStateTechnique(),
		methodology.NewIdempotentTechnique(),
		methodology.NewPairwiseTechnique(),
	)
	engine.AddSpecTechnique(methodology.NewChainTechnique())
	engine.SetSink(bus)
	cases, err := engine.Generate(parsedSpec)
	if err != nil {
		return fmt.Errorf("generating test cases: %w", err)
	}

	// Write index.json
	w := writer.NewJSONSchemaWriter()
	if err := w.Write(cases, genOutput); err != nil {
		fmt.Fprintf(os.Stderr, "✗ Failed to write output: %v\n", err)
		os.Exit(5)
	}

	// Render to target format
	var renderer render.Renderer
	switch cfg.Output.DefaultFormat {
	case "markdown":
		renderer = render.NewMarkdownRenderer()
	case "csv":
		renderer = render.NewCSVRenderer()
	default: // "hurl" and anything unrecognised
		renderer = render.NewHurlRenderer("")
	}
	if err := renderer.Render(cases, genOutput); err != nil {
		fmt.Fprintf(os.Stderr, "✗ Render failed: %v\n", err)
		os.Exit(5)
	}

	bus.Emit(event.Event{Type: event.EventRenderDone})

	if tuiDone != nil {
		<-tuiDone
	}

	fmt.Fprintf(os.Stderr, "✓ Generated %d test cases → %s\n", len(cases), genOutput)
	return nil
}
