// cmd/explore.go
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/spf13/cobra"
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
	exploreCmd.Flags().String("export-pool", "", "Write observed field values to a JSON data pool file (from 2xx responses)")
	exploreCmd.Flags().Bool("prioritize-uncovered", false, "Two-pass probe scheduling: cover all ops in pass 1, focus budget on non-2xx ops in pass 2")
	exploreCmd.Flags().String("include-path", "", "Regex to include operations by path")
	exploreCmd.Flags().String("exclude-path", "", "Regex to exclude operations by path")
	exploreCmd.Flags().String("include-tag", "", "Comma-separated tags to include")
	exploreCmd.Flags().String("exclude-tag", "", "Comma-separated tags to exclude")
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

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Loading spec: %s\n", specPath)
	loader := spec.NewLoader()
	parsedSpec, err := loader.Load(specPath)
	if err != nil {
		return fmt.Errorf("load spec: %w", err)
	}
	fmt.Fprintf(out, "✓ Loaded %d operations\n", len(parsedSpec.Operations))

	includePath, _ := cmd.Flags().GetString("include-path")
	excludePath, _ := cmd.Flags().GetString("exclude-path")
	includeTag, _ := cmd.Flags().GetString("include-tag")
	excludeTag, _ := cmd.Flags().GetString("exclude-tag")
	opFilter := buildFilterSet(includePath, excludePath, includeTag, excludeTag)
	if !opFilter.IsEmpty() {
		if err := opFilter.Validate(); err != nil {
			return err
		}
		parsedSpec.Operations = opFilter.Apply(parsedSpec.Operations)
		if len(parsedSpec.Operations) == 0 {
			fmt.Fprintln(out, "warning: path/tag filters matched no operations")
		}
	}

	explorer := dea.NewExplorer(targetURL, maxProbes)
	explorer.DryRun = dryRun
	if prio, _ := cmd.Flags().GetBool("prioritize-uncovered"); prio {
		explorer.PrioritizeUncovered = true
	}

	if dryRun {
		fmt.Fprintln(out, "Dry-run mode: seeding hypotheses without HTTP probes...")
	} else {
		fmt.Fprintf(out, "Exploring %s (max %d probes)...\n", targetURL, maxProbes)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	report, err := explorer.Explore(ctx, parsedSpec)
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

	if poolPath, _ := cmd.Flags().GetString("export-pool"); poolPath != "" {
		if err := explorer.DataPool().Save(poolPath); err != nil {
			return fmt.Errorf("write data pool: %w", err)
		}
		fmt.Fprintf(out, "✓ Data pool written to %s (%d field(s))\n", poolPath, explorer.DataPool().Len())
	}
	return nil
}
