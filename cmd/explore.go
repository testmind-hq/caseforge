// cmd/explore.go
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/config"
	"github.com/testmind-hq/caseforge/internal/dea"
	"github.com/testmind-hq/caseforge/internal/spec"
)

var exploreCmd = &cobra.Command{
	Use:   "explore",
	Short: "Dynamically explore a live API and infer implicit rules (DEA)",
	Long: `Explore runs a Dynamic Exploration Agent against a live API.

It builds a hypothesis tree from your OpenAPI spec, designs minimal HTTP probes
for each hypothesis, executes them against the target URL, and infers rules about
actual API behavior — including constraints not declared in the spec.

Examples:
  caseforge explore --spec openapi.yaml --target http://localhost:8080
  caseforge explore --spec openapi.yaml --target http://api.example.com --dry-run
  caseforge explore --spec openapi.yaml --target http://localhost:8080 --max-probes 100 --output ./reports`,
	RunE: runExplore,
}

func init() {
	rootCmd.AddCommand(exploreCmd)
	exploreCmd.Flags().String("spec", "", "OpenAPI spec file (required)")
	exploreCmd.Flags().String("target", "", "Target API base URL, e.g. http://localhost:8080 (required without --dry-run)")
	exploreCmd.Flags().Int("max-probes", 50, "Maximum number of HTTP probes per run")
	exploreCmd.Flags().String("output", "./reports", "Output directory for dea-report.json")
	exploreCmd.Flags().Bool("dry-run", false, "Seed hypotheses without executing HTTP probes")
}

func runExplore(cmd *cobra.Command, _ []string) error {
	specPath, _ := cmd.Flags().GetString("spec")
	targetURL, _ := cmd.Flags().GetString("target")
	maxProbes, _ := cmd.Flags().GetInt("max-probes")
	outputDir, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	if specPath == "" {
		return fmt.Errorf("--spec is required")
	}
	if targetURL == "" && !dryRun {
		return fmt.Errorf("--target is required (or use --dry-run to skip HTTP execution)")
	}

	cfg, _ := config.Load()
	_ = cfg

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Loading spec: %s\n", specPath)
	loader := spec.NewLoader()
	parsedSpec, err := loader.Load(specPath)
	if err != nil {
		return fmt.Errorf("load spec: %w", err)
	}
	fmt.Fprintf(out, "✓ Loaded %d operations\n", len(parsedSpec.Operations))

	explorer := dea.NewExplorer(targetURL, maxProbes)
	explorer.DryRun = dryRun

	if dryRun {
		fmt.Fprintln(out, "Dry-run mode: seeding hypotheses without HTTP probes...")
	} else {
		fmt.Fprintf(out, "Exploring %s (max %d probes)...\n", targetURL, maxProbes)
	}

	report, err := explorer.Explore(context.Background(), parsedSpec)
	if err != nil {
		return fmt.Errorf("exploration failed: %w", err)
	}

	fmt.Fprintf(out, "\n✓ Exploration complete: %d probes, %d rules discovered\n",
		report.TotalProbes, len(report.Rules))

	implicitCount := 0
	mismatchCount := 0
	for _, r := range report.Rules {
		if r.Implicit {
			implicitCount++
		}
		if r.Category == dea.CategorySpecMismatch {
			mismatchCount++
		}
	}
	if implicitCount > 0 {
		fmt.Fprintf(out, "  ⚠  %d implicit rules (constraints not in spec)\n", implicitCount)
	}
	if mismatchCount > 0 {
		fmt.Fprintf(out, "  ✗  %d spec mismatches (server behavior differs from spec)\n", mismatchCount)
	}

	report.SpecPath = specPath
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	reportPath := filepath.Join(outputDir, "dea-report.json")
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal report: %w", err)
	}
	if err := os.WriteFile(reportPath, data, 0644); err != nil {
		return fmt.Errorf("write report: %w", err)
	}
	fmt.Fprintf(out, "\n✓ Report written to %s\n", reportPath)
	return nil
}
