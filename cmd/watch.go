// cmd/watch.go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/config"
	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/methodology"
	"github.com/testmind-hq/caseforge/internal/output/render"
	"github.com/testmind-hq/caseforge/internal/output/writer"
	"github.com/testmind-hq/caseforge/internal/spec"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch an OpenAPI spec file and regenerate cases on change",
	Long: `Watch monitors an OpenAPI spec file for changes and automatically
regenerates test cases whenever the spec is updated.

Press Ctrl-C to stop watching.

Examples:
  caseforge watch --spec ./openapi.yaml
  caseforge watch --spec ./openapi.yaml --output ./cases --format hurl --no-ai`,
	RunE:         runWatch,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(watchCmd)
	watchCmd.Flags().String("spec", "", "OpenAPI spec file path to watch (required, local file only)")
	watchCmd.Flags().String("output", "./cases", "Output directory for generated cases")
	watchCmd.Flags().String("format", "", "Output format: hurl|k6|postman|markdown|csv (default from config)")
	watchCmd.Flags().Bool("no-ai", false, "Disable LLM, use pure algorithm mode")
	_ = watchCmd.MarkFlagRequired("spec")
}

func runWatch(cmd *cobra.Command, _ []string) error {
	specPath, _ := cmd.Flags().GetString("spec")
	outputDir, _ := cmd.Flags().GetString("output")
	format, _ := cmd.Flags().GetString("format")
	noAI, _ := cmd.Flags().GetBool("no-ai")

	if strings.HasPrefix(specPath, "http://") || strings.HasPrefix(specPath, "https://") {
		return fmt.Errorf("watch only supports local files, not URLs")
	}

	absSpec, err := filepath.Abs(specPath)
	if err != nil {
		return fmt.Errorf("resolving spec path: %w", err)
	}
	if _, err := os.Stat(absSpec); err != nil {
		return fmt.Errorf("spec file not found: %s", absSpec)
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	if noAI {
		cfg.AI.Provider = "noop"
	}
	if format != "" {
		cfg.Output.DefaultFormat = format
	}

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "✓ 监听 %s...\n\n", absSpec)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("creating watcher: %w", err)
	}
	defer watcher.Close()

	if err := watcher.Add(absSpec); err != nil {
		return fmt.Errorf("watching %s: %w", absSpec, err)
	}

	// Debounce: skip duplicate events within 200ms
	var lastEvent time.Time

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
				continue
			}
			if time.Since(lastEvent) < 200*time.Millisecond {
				continue
			}
			lastEvent = time.Now()

			ts := time.Now().Format("15:04:05")
			fmt.Fprintf(out, "[%s] 检测到变更：%s\n", ts, filepath.Base(absSpec))

			count, genErr := watchRegenerate(cfg, absSpec, outputDir)
			if genErr != nil {
				fmt.Fprintf(out, "[%s] ✗ 生成失败：%v\n\n", ts, genErr)
				continue
			}
			fmt.Fprintf(out, "[%s] ✓ 更新了 %d 个用例 → %s\n\n",
				time.Now().Format("15:04:05"), count, outputDir)

		case watchErr, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "watch error: %v\n", watchErr)
		}
	}
}

// watchRegenerate runs gen for the given spec and returns the number of cases produced.
func watchRegenerate(cfg *config.Config, specPath, outputDir string) (int, error) {
	loader := spec.NewLoader()
	parsedSpec, err := loader.Load(specPath)
	if err != nil {
		return 0, fmt.Errorf("loading spec: %w", err)
	}

	provider := llm.NewProviderWithConfig(cfg.AI.APIKey, cfg.AI.Provider, cfg.AI.Model, cfg.AI.BaseURL)

	allTechniques := []methodology.Technique{
		methodology.NewEquivalenceTechnique(),
		methodology.NewBoundaryTechnique(),
		methodology.NewDecisionTechnique(),
		methodology.NewStateTechnique(),
		methodology.NewIdempotentTechnique(),
		methodology.NewPairwiseTechnique(),
		methodology.NewSecurityTechnique(),
		methodology.NewExampleTechnique(),
	}
	engine := methodology.NewEngine(provider, allTechniques...)
	engine.AddSpecTechnique(methodology.NewChainTechnique())
	engine.AddSpecTechnique(methodology.NewSecuritySpecTechnique())

	cases, err := engine.Generate(parsedSpec)
	if err != nil {
		return 0, fmt.Errorf("generating cases: %w", err)
	}

	// Write index.json
	specHash, _ := writer.HashFile(specPath)
	w := writer.NewJSONSchemaWriter()
	if err := w.Write(cases, outputDir, writer.WriteOptions{
		SpecHash:         specHash,
		CaseforgeVersion: Version,
	}); err != nil {
		return 0, fmt.Errorf("writing index.json: %w", err)
	}

	// Render to format
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
	default:
		renderer = render.NewHurlRenderer("")
	}
	if err := renderer.Render(cases, outputDir); err != nil {
		return 0, fmt.Errorf("rendering: %w", err)
	}

	return len(cases), nil
}
