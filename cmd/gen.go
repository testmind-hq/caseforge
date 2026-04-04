// cmd/gen.go
package cmd

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/config"
	"github.com/testmind-hq/caseforge/internal/event"
	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/methodology"
	"github.com/testmind-hq/caseforge/internal/output/render"
	"github.com/testmind-hq/caseforge/internal/output/schema"
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
	genSpec        string
	genOutput      string
	genNoAI        bool
	genFormat      string
	genTechnique   string
	genPriority    string
	genOperations  string
	genConcurrency int
)

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.Flags().StringVar(&genSpec, "spec", "", "OpenAPI spec file or URL (required)")
	genCmd.Flags().StringVar(&genOutput, "output", "./cases", "Output directory")
	genCmd.Flags().BoolVar(&genNoAI, "no-ai", false, "Disable LLM, use pure algorithm mode")
	genCmd.Flags().StringVar(&genFormat, "format", "hurl", "Output format: hurl|markdown|csv|postman|k6")
	genCmd.Flags().StringVar(&genTechnique, "technique", "", "Only run the named technique(s), comma-separated (e.g. equivalence_partitioning,boundary_value)")
	genCmd.Flags().StringVar(&genPriority, "priority", "", "Filter output by minimum priority: P0|P1|P2|P3 (P0 = highest)")
	genCmd.Flags().StringVar(&genOperations, "operations", "", "Comma-separated operationIds to process (default: all)")
	genCmd.Flags().IntVar(&genConcurrency, "concurrency", 1, "Number of operations processed concurrently (default 1)")
	_ = genCmd.MarkFlagRequired("spec")
}

// priorityRank maps P0..P3 to a numeric rank where lower = higher priority.
var priorityRank = map[string]int{"P0": 0, "P1": 1, "P2": 2, "P3": 3}

func runGen(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Override config from flags
	if genNoAI {
		cfg.AI.Provider = "noop"
	}
	if cmd.Flags().Changed("format") || genFormat != "" {
		cfg.Output.DefaultFormat = genFormat
	}

	// Validate --priority
	if genPriority != "" {
		if _, ok := priorityRank[strings.ToUpper(genPriority)]; !ok {
			return fmt.Errorf("invalid --priority %q: must be P0, P1, P2, or P3", genPriority)
		}
		genPriority = strings.ToUpper(genPriority)
	}

	// Validate --concurrency
	if genConcurrency < 1 {
		return fmt.Errorf("invalid --concurrency %d: must be ≥ 1", genConcurrency)
	}

	// Resolve LLM provider
	provider := llm.NewProviderWithConfig(cfg.AI.APIKey, cfg.AI.Provider, cfg.AI.Model, cfg.AI.BaseURL)
	if cfg.AI.Provider != "noop" && !provider.IsAvailable() {
		if !genNoAI {
			fmt.Fprintln(os.Stderr, "✗ LLM provider unavailable. Use --no-ai to run in algorithm-only mode.")
			os.Exit(ExitNoOutput)
		}
	}

	// Load spec
	loader := spec.NewLoader()
	parsedSpec, err := loader.Load(genSpec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Failed to parse spec: %v\n", err)
		os.Exit(ExitSpecParseError)
	}

	// --operations: filter parsedSpec to only requested operationIds
	if genOperations != "" {
		allowed := splitTrimmed(genOperations)
		allowedSet := make(map[string]bool, len(allowed))
		for _, id := range allowed {
			allowedSet[id] = true
		}
		var filtered []*spec.Operation
		for _, op := range parsedSpec.Operations {
			if allowedSet[op.OperationID] {
				filtered = append(filtered, op)
			}
		}
		parsedSpec.Operations = filtered
		if len(filtered) == 0 {
			fmt.Fprintf(os.Stderr, "warning: --operations %q matched no operationIds in the spec\n", genOperations)
		}
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

	// --technique: select which per-operation techniques to register
	allTechniques := []methodology.Technique{
		methodology.NewEquivalenceTechnique(),
		methodology.NewBoundaryTechnique(),
		methodology.NewDecisionTechnique(),
		methodology.NewStateTechnique(),
		methodology.NewIdempotentTechnique(),
		methodology.NewPairwiseTechnique(),
		methodology.NewClassificationTreeTechnique(),
		methodology.NewOrthogonalArrayTechnique(),
		methodology.NewSecurityTechnique(),
		methodology.NewExampleTechnique(),
	}
	allSpecTechniques := []methodology.SpecTechnique{
		methodology.NewChainTechnique(),
		methodology.NewSecuritySpecTechnique(),
	}
	selectedTechniques, selectedSpec := filterTechniques(allTechniques, allSpecTechniques, genTechnique)
	if genTechnique != "" && len(selectedTechniques) == 0 && len(selectedSpec) == 0 {
		fmt.Fprintf(os.Stderr, "warning: --technique %q matched no known technique names\n", genTechnique)
	}

	// Generate test cases
	engine := methodology.NewEngine(provider, selectedTechniques...)
	for _, st := range selectedSpec {
		engine.AddSpecTechnique(st)
	}
	engine.SetSink(bus)
	engine.SetConcurrency(genConcurrency)
	cases, err := engine.Generate(parsedSpec)
	if err != nil {
		return fmt.Errorf("generating test cases: %w", err)
	}

	// --priority: keep cases whose priority is at least as high as the requested threshold.
	if genPriority != "" {
		cases = filterByPriority(cases, genPriority)
	}

	// Write index.json
	specHash, _ := writer.HashFile(genSpec)
	w := writer.NewJSONSchemaWriter()
	if err := w.Write(cases, genOutput, writer.WriteOptions{
		SpecHash:         specHash,
		CaseforgeVersion: Version,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "✗ Failed to write output: %v\n", err)
		os.Exit(ExitWriteError)
	}

	// Render to target format
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
	if err := renderer.Render(cases, genOutput); err != nil {
		fmt.Fprintf(os.Stderr, "✗ Render failed: %v\n", err)
		os.Exit(ExitWriteError)
	}

	bus.Emit(event.Event{Type: event.EventRenderDone})

	if tuiDone != nil {
		<-tuiDone
	}

	fmt.Fprintf(os.Stderr, "✓ Generated %d test cases → %s\n", len(cases), genOutput)
	return nil
}

// filterTechniques returns the subsets of per-operation and spec techniques
// that match the comma-separated names in the filter string. An empty filter
// returns all techniques unchanged.
func filterTechniques(
	ops []methodology.Technique,
	specs []methodology.SpecTechnique,
	filter string,
) ([]methodology.Technique, []methodology.SpecTechnique) {
	if filter == "" {
		return ops, specs
	}
	names := make(map[string]bool)
	for _, n := range splitTrimmed(filter) {
		names[n] = true
	}
	var filteredOps []methodology.Technique
	for _, t := range ops {
		if names[t.Name()] {
			filteredOps = append(filteredOps, t)
		}
	}
	var filteredSpec []methodology.SpecTechnique
	for _, t := range specs {
		if names[t.Name()] {
			filteredSpec = append(filteredSpec, t)
		}
	}
	return filteredOps, filteredSpec
}

// filterByPriority keeps cases whose priority is at least as high as minPriority.
// Because P0 has rank 0 (highest) and P3 has rank 3 (lowest), a case passes
// when its numeric rank is <= the threshold rank. Cases with unrecognised or
// empty priority fields are excluded.
func filterByPriority(cases []schema.TestCase, minPriority string) []schema.TestCase {
	threshold := priorityRank[minPriority]
	var out []schema.TestCase
	for _, c := range cases {
		rank, ok := priorityRank[c.Priority]
		if ok && rank <= threshold {
			out = append(out, c)
		}
	}
	return out
}

// splitTrimmed splits s on commas and trims whitespace from each token.
func splitTrimmed(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
