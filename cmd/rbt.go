// cmd/rbt.go
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/config"
	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/methodology"
	"github.com/testmind-hq/caseforge/internal/output/render"
	"github.com/testmind-hq/caseforge/internal/output/writer"
	"github.com/testmind-hq/caseforge/internal/rbt"
	"github.com/testmind-hq/caseforge/internal/spec"
)

var rbtCmd = &cobra.Command{
	Use:   "rbt",
	Short: "Risk-Based Testing: assess which API operations are at risk from recent changes",
	Long: `RBT analyses the git diff between two commits and cross-references changed
source files with your OpenAPI spec and existing test cases to produce a risk
report. Operations with no tests covering changed code are flagged as high risk.

Examples:
  caseforge rbt --spec openapi.yaml
  caseforge rbt --spec openapi.yaml --base HEAD~3 --head HEAD --format json
  caseforge rbt --spec openapi.yaml --dry-run --output ./reports
  caseforge rbt --spec openapi.yaml --generate --gen-format hurl`,
	RunE: runRBT,
}

func init() {
	rootCmd.AddCommand(rbtCmd)
	rbtCmd.Flags().String("spec", "", "OpenAPI spec file (required)")
	rbtCmd.Flags().String("cases", "./cases", "Directory containing test case JSON files")
	rbtCmd.Flags().String("src", "./", "Source code root directory")
	rbtCmd.Flags().String("base", "HEAD~1", "Base git ref for diff")
	rbtCmd.Flags().String("head", "HEAD", "Head git ref for diff")
	rbtCmd.Flags().Bool("generate", false, "Generate test cases for high-risk uncovered operations")
	rbtCmd.Flags().Bool("no-ai", false, "Disable LLM for both route inference and test case generation; use algorithm-only mode")
	rbtCmd.Flags().String("gen-format", "hurl", "Format for generated test cases: hurl|postman|k6|markdown|csv")
	rbtCmd.Flags().String("output", "./reports", "Output directory for rbt-report.json")
	rbtCmd.Flags().String("format", "terminal", "Output format: terminal or json")
	rbtCmd.Flags().String("fail-on", "high", "Exit non-zero if any operation has risk >= level (none|low|medium|high)")
	rbtCmd.Flags().String("map", "", "Path to caseforge-map.yaml (default: <src>/caseforge-map.yaml)")
	rbtCmd.Flags().Bool("dry-run", false, "Skip git diff and tree-sitter; report all operations as risk=none")
}

func runRBT(cmd *cobra.Command, _ []string) error {
	specPath, _ := cmd.Flags().GetString("spec")
	casesDir, _ := cmd.Flags().GetString("cases")
	srcDir, _ := cmd.Flags().GetString("src")
	base, _ := cmd.Flags().GetString("base")
	head, _ := cmd.Flags().GetString("head")
	outputDir, _ := cmd.Flags().GetString("output")
	format, _ := cmd.Flags().GetString("format")
	failOn, _ := cmd.Flags().GetString("fail-on")
	mapFile, _ := cmd.Flags().GetString("map")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	generate, _ := cmd.Flags().GetBool("generate")
	noAI, _ := cmd.Flags().GetBool("no-ai")

	if specPath == "" {
		return fmt.Errorf("--spec is required")
	}

	// Load spec
	loader := spec.NewLoader()
	parsedSpec, err := loader.Load(specPath)
	if err != nil {
		return fmt.Errorf("load spec: %w", err)
	}

	var provider llm.LLMProvider
	var report rbt.RiskReport

	if dryRun {
		// Build a report with all operations at RiskNone
		ops := make([]rbt.OperationCoverage, 0, len(parsedSpec.Operations))
		for _, op := range parsedSpec.Operations {
			ops = append(ops, rbt.OperationCoverage{
				OperationID: op.OperationID,
				Method:      op.Method,
				Path:        op.Path,
				Affected:    false,
				Risk:        rbt.RiskNone,
			})
		}
		report = rbt.RiskReport{
			DiffBase:    base,
			DiffHead:    head,
			Operations:  ops,
			GeneratedAt: time.Now(),
		}
	} else {
		// Load LLM provider; used by both the parser chain and --generate.
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		providerName := cfg.AI.Provider
		if noAI {
			providerName = "noop"
		}
		provider = llm.NewProviderWithConfig(cfg.AI.APIKey, providerName, cfg.AI.Model, cfg.AI.BaseURL)
		if !provider.IsAvailable() {
			provider = llm.NewProviderWithConfig("", "noop", "", "")
		}

		// Git diff
		changedFiles, err := rbt.Diff(srcDir, base, head)
		if err != nil {
			return fmt.Errorf("git diff: %w", err)
		}

		// Build parser chain
		if mapFile == "" {
			mapFile = filepath.Join(srcDir, "caseforge-map.yaml")
		}
		parsers := []rbt.SourceParser{
			rbt.NewMapFileParser(mapFile),
			rbt.NewTreeSitterParser(),
			rbt.NewRegexParser(),
			rbt.NewLLMParser(provider, parsedSpec.Operations),
		}

		// Map changed files to route mappings
		mappings, err := rbt.MapChain(parsers, srcDir, changedFiles)
		if err != nil {
			return fmt.Errorf("map chain: %w", err)
		}

		// Scan test cases
		caseIndex, err := rbt.ScanCases(casesDir)
		if err != nil {
			return fmt.Errorf("scan cases: %w", err)
		}

		// Assess risk
		report = rbt.Assess(parsedSpec, mappings, caseIndex, base, head, changedFiles)
	}

	// Output
	out := cmd.OutOrStdout()
	errOut := cmd.ErrOrStderr()
	if format == "json" {
		data, err := rbt.MarshalReportJSON(report)
		if err != nil {
			return fmt.Errorf("marshal report: %w", err)
		}
		fmt.Fprintln(out, string(data))
	} else {
		rbt.PrintTerminal(out, report)
	}

	// Always write report JSON
	reportPath, err := rbt.WriteReportJSON(outputDir, report)
	if err != nil {
		return fmt.Errorf("write report: %w", err)
	}
	// Print to stderr when --format json so stdout is valid JSON for piping.
	if format == "json" {
		fmt.Fprintln(errOut, "Report written to:", reportPath)
	} else {
		fmt.Fprintln(out, "Report written to:", reportPath)
	}

	// --generate: auto-generate tests for HIGH-risk operations.
	// Skipped in dry-run mode because all operations are RiskNone there.
	if generate {
		if dryRun {
			fmt.Fprintln(out, "--generate is ignored with --dry-run (no HIGH-risk operations in dry-run mode)")
		} else {
			if err := runGenerate(cmd, out, errOut, specPath, parsedSpec, report, casesDir, provider); err != nil {
				return err
			}
		}
	}

	// Exit non-zero for high-risk (only in non-dry-run path)
	if !dryRun && rbt.ShouldFail(report, failOn) {
		os.Exit(1)
	}

	return nil
}

// runGenerate generates test cases for all HIGH-risk operations in the report.
// It runs the same methodology pipeline as `caseforge gen`, restricted to the
// high-risk operations identified by the assessor.
// specPath is the spec file path used to compute the spec hash for index.json.
// errOut is used for the success message so stdout stays clean when --format=json.
func runGenerate(cmd *cobra.Command, w, errOut io.Writer, specPath string, parsedSpec *spec.ParsedSpec, report rbt.RiskReport, casesDir string, provider llm.LLMProvider) error {
	highRiskOps := rbt.HighRiskOperations(report, parsedSpec)
	if len(highRiskOps) == 0 {
		fmt.Fprintln(w, "No HIGH-risk operations to generate tests for.")
		return nil
	}

	genFormat, _ := cmd.Flags().GetString("gen-format")

	// Build a filtered spec containing only the high-risk operations.
	filteredSpec := &spec.ParsedSpec{
		Title:           parsedSpec.Title,
		Version:         parsedSpec.Version,
		Operations:      highRiskOps,
		Schemas:         parsedSpec.Schemas,
		SecuritySchemes: parsedSpec.SecuritySchemes,
		GlobalSecurity:  parsedSpec.GlobalSecurity,
	}

	// Run the standard methodology engine on the filtered spec.
	engine := methodology.NewEngine(provider,
		methodology.NewEquivalenceTechnique(),
		methodology.NewBoundaryTechnique(),
		methodology.NewDecisionTechnique(),
		methodology.NewStateTechnique(),
		methodology.NewIdempotentTechnique(),
		methodology.NewPairwiseTechnique(),
		methodology.NewSecurityTechnique(),
	)
	engine.AddSpecTechnique(methodology.NewChainTechnique())
	engine.AddSpecTechnique(methodology.NewSecuritySpecTechnique())

	cases, err := engine.Generate(filteredSpec)
	if err != nil {
		return fmt.Errorf("generating test cases: %w", err)
	}

	// Write index.json to casesDir; include spec hash for provenance.
	specHash, _ := writer.HashFile(specPath)
	caseWriter := writer.NewJSONSchemaWriter()
	if err := caseWriter.Write(cases, casesDir, writer.WriteOptions{
		SpecHash:         specHash,
		CaseforgeVersion: Version,
	}); err != nil {
		return fmt.Errorf("writing generated cases: %w", err)
	}

	// Render to the requested format.
	var renderer render.Renderer
	switch genFormat {
	case "markdown":
		renderer = render.NewMarkdownRenderer()
	case "csv":
		renderer = render.NewCSVRenderer()
	case "postman":
		renderer = render.NewPostmanRenderer()
	case "k6":
		renderer = render.NewK6Renderer()
	default: // "hurl" and unrecognised values
		renderer = render.NewHurlRenderer("")
	}
	if err := renderer.Render(cases, casesDir); err != nil {
		return fmt.Errorf("rendering generated cases: %w", err)
	}

	// Write to errOut (stderr) so stdout stays clean when --format=json.
	fmt.Fprintf(errOut, "✓ Generated %d test cases for %d HIGH-risk operation(s) → %s\n",
		len(cases), len(highRiskOps), casesDir)
	return nil
}
