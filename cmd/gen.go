// cmd/gen.go
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/checkpoint"
	"github.com/testmind-hq/caseforge/internal/config"
	"github.com/testmind-hq/caseforge/internal/event"
	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/methodology"
	"github.com/testmind-hq/caseforge/internal/output/render"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/output/writer"
	"github.com/testmind-hq/caseforge/internal/webhook"
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
	genConcurrency    int
	genResume         bool
	genTupleLevel     int
	genSeed           int64
	genMaxCasesPerOp  int
)

// allTechniqueNames is the canonical list used for --technique completion.
var allTechniqueNames = []string{
	"equivalence_partitioning",
	"boundary_value",
	"decision_table",
	"state_transition",
	"idempotent",
	"pairwise",
	"classification_tree",
	"orthogonal_array",
	"owasp_api_top10",
	"examples",
	"chain",
	"owasp_api_top10_spec",
	"isolated_negative",
	"schema_violation",
	"variable_irrelevance",
}

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
	genCmd.Flags().BoolVar(&genResume, "resume", false, "Resume an interrupted run; skips completed operations. Cases for skipped ops are taken from the last complete run's output.")
	genCmd.Flags().IntVar(&genTupleLevel, "tuple-level", 2, "N-way coverage level for pairwise technique (2=pairwise, 3=3-way, max 4)")
	genCmd.Flags().Int64Var(&genSeed, "seed", 0, "Seed for deterministic generation (0 = random)")
	genCmd.Flags().IntVar(&genMaxCasesPerOp, "max-cases-per-op", 0, "Maximum test cases per operation (0 = unlimited, P0 cases take priority)")
	_ = genCmd.MarkFlagRequired("spec")

	// Dynamic completion: --operations reads the spec and suggests operationIds.
	_ = genCmd.RegisterFlagCompletionFunc("operations", func(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		specFile, _ := cmd.Flags().GetString("spec")
		if specFile == "" {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		ps, err := spec.NewLoader().Load(specFile)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var ids []string
		for _, op := range ps.Operations {
			if op.OperationID != "" {
				ids = append(ids, op.OperationID)
			}
		}
		return ids, cobra.ShellCompDirectiveNoFileComp
	})

	// Dynamic completion: --technique returns all known technique names.
	_ = genCmd.RegisterFlagCompletionFunc("technique", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return allTechniqueNames, cobra.ShellCompDirectiveNoFileComp
	})

	// Dynamic completion: --format returns all valid formats.
	_ = genCmd.RegisterFlagCompletionFunc("format", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"hurl", "markdown", "csv", "postman", "k6"}, cobra.ShellCompDirectiveNoFileComp
	})

	// Dynamic completion: --priority returns valid priority levels.
	_ = genCmd.RegisterFlagCompletionFunc("priority", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"P0", "P1", "P2", "P3"}, cobra.ShellCompDirectiveNoFileComp
	})
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

	// Validate --tuple-level
	if genTupleLevel < 2 || genTupleLevel > 4 {
		return fmt.Errorf("invalid --tuple-level %d: must be between 2 and 4", genTupleLevel)
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

	// Checkpoint / resume logic.
	// HashFile fails for URL specs — if the hash is empty, disable checkpointing
	// for this run to avoid false hash matches on resume.
	specHash, hashErr := writer.HashFile(genSpec)
	if hashErr != nil && genResume {
		fmt.Fprintln(os.Stderr, "warning: cannot hash spec (URL or unreadable file) — --resume disabled for this run")
		genResume = false
	}

	ckptMgr := checkpoint.NewManager(genOutput)
	var ckptState *checkpoint.State
	var resumedCases []schema.TestCase

	if genResume {
		ckptState, err = ckptMgr.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not load checkpoint: %v\n", err)
			ckptState = nil
		}
		if ckptState != nil {
			if ckptState.SpecHash != specHash {
				fmt.Fprintln(os.Stderr, "warning: spec changed since last run — starting fresh (checkpoint discarded)")
				ckptState = nil
				_ = ckptMgr.Delete()
			} else {
				// Load previously generated cases from index.json.
				// Note: these are cases from the last *complete* run. Cases generated
				// during the interrupted run itself are not preserved because index.json
				// is only written on successful completion.
				resumedCases = loadExistingCases(genOutput)
				// Filter out already-completed operations.
				var remaining []*spec.Operation
				for _, op := range parsedSpec.Operations {
					key := checkpoint.OperationKey(op.Method, op.Path)
					if !ckptState.Completed[key] {
						remaining = append(remaining, op)
					}
				}
				skipped := len(parsedSpec.Operations) - len(remaining)
				if skipped > 0 {
					fmt.Fprintf(os.Stderr, "↩ Resuming: skipping %d already-completed operation(s)\n", skipped)
					if len(resumedCases) == 0 {
						fmt.Fprintln(os.Stderr, "warning: no prior output found — skipped operations will not appear in output")
					}
				}
				parsedSpec.Operations = remaining
			}
		}
	}

	if ckptState == nil {
		ckptState = checkpoint.NewState(specHash)
	}

	// Set up event bus
	bus := event.NewBus()

	// Subscribe a checkpoint sink that persists state after each operation.
	// A mutex guards ckptState.Completed to prevent data races when
	// --concurrency > 1 causes concurrent EventOperationDone deliveries.
	// ckptMu serialises writes to ckptState.Completed and the Clone call.
	// Save receives a snapshot so json.MarshalIndent never races with a concurrent
	// map write from another goroutine (safe with --concurrency > 1).
	var ckptMu sync.Mutex
	bus.Subscribe(event.SinkFunc(func(e event.Event) {
		if e.Type != event.EventOperationDone {
			return
		}
		if p, ok := e.Payload.(event.OperationDonePayload); ok {
			ckptMu.Lock()
			ckptState.Completed[checkpoint.OperationKey(p.Method, p.Path)] = true
			snap := ckptState.Clone()
			ckptMu.Unlock()
			_ = ckptMgr.Save(snap) // best-effort; ignore write errors
		}
	}))

	// Wire webhook sink if any endpoints are configured.
	if len(cfg.Webhooks) > 0 {
		whSink := webhook.New(cfg.Webhooks)
		whSink.SetOutputDir(genOutput)
		bus.Subscribe(whSink)
	}

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
		methodology.NewPairwiseTechniqueWithLevel(genTupleLevel),
		methodology.NewClassificationTreeTechnique(),
		methodology.NewOrthogonalArrayTechnique(),
		methodology.NewSecurityTechnique(),
		methodology.NewExampleTechnique(),
		methodology.NewIsolatedNegativeTechnique(),
		methodology.NewSchemaViolationTechnique(),
		methodology.NewVariableIrrelevanceTechnique(),
	}
	allSpecTechniques := []methodology.SpecTechnique{
		methodology.NewChainTechnique(),
		methodology.NewSecuritySpecTechnique(),
	}
	selectedTechniques, selectedSpec := filterTechniques(allTechniques, allSpecTechniques, genTechnique)
	if genTechnique != "" && len(selectedTechniques) == 0 && len(selectedSpec) == 0 {
		fmt.Fprintf(os.Stderr, "warning: --technique %q matched no known technique names\n", genTechnique)
	}

	// Generate test cases for remaining operations
	engine := methodology.NewEngine(provider, selectedTechniques...)
	for _, st := range selectedSpec {
		engine.AddSpecTechnique(st)
	}
	engine.SetSink(bus)
	engine.SetConcurrency(genConcurrency)
	if genSeed != 0 {
		engine.SetSeed(genSeed)
	}
	if genMaxCasesPerOp > 0 {
		engine.SetMaxCasesPerOp(genMaxCasesPerOp)
	}
	newCases, err := engine.Generate(parsedSpec)
	if err != nil {
		return fmt.Errorf("generating test cases: %w", err)
	}

	// Merge resumed cases with newly generated cases
	cases := append(resumedCases, newCases...)

	// --priority: keep cases whose priority is at least as high as the requested threshold.
	if genPriority != "" {
		cases = filterByPriority(cases, genPriority)
	}

	// Write index.json
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

	// Remove checkpoint on successful completion
	_ = ckptMgr.Delete()

	fmt.Fprintf(os.Stderr, "✓ Generated %d test cases → %s\n", len(cases), genOutput)
	return nil
}

// loadExistingCases reads previously generated cases from index.json in outputDir.
// Returns nil if the file does not exist. Logs a warning if the file exists but
// cannot be parsed, so the user knows prior output is not being carried forward.
func loadExistingCases(outputDir string) []schema.TestCase {
	indexPath := filepath.Join(outputDir, "index.json")
	r := writer.NewJSONSchemaWriter()
	cases, err := r.Read(indexPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "warning: could not read existing cases from %s: %v\n", indexPath, err)
		}
		return nil
	}
	return cases
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
