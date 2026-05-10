// cmd/mutate.go
package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/config"
	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/mutation"
)

var (
	mutateCases       string
	mutateTarget      string
	mutateOutput      string
	mutateOperators   string
	mutateConcurrency int
)

var mutateCmd = &cobra.Command{
	Use:   "mutate",
	Short: "Run HTTP boundary mutations to find weak test assertions",
	Long: `Mutate starts a local reverse proxy between hurl and your API.
For each mutation operator × test case, the proxy alters the response
before hurl evaluates assertions. Cases where hurl still passes are
"survivors" — mutations your assertions failed to catch.

Requires hurl on PATH. Test cases must be previously generated with 'caseforge gen'.

Exit codes:
  0 — run complete, no survivors
  6 — one or more mutations survived

Examples:
  caseforge mutate --cases ./cases --target http://localhost:8080
  caseforge mutate --cases ./cases --target http://localhost:8080 --output ./reports
  caseforge mutate --cases ./cases --target http://localhost:8080 --operator field_drop,status_swap_2xx`,
	RunE:         runMutate,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(mutateCmd)
	mutateCmd.Flags().StringVar(&mutateCases, "cases", "", "Directory containing index.json and .hurl files (required)")
	_ = mutateCmd.MarkFlagRequired("cases")
	mutateCmd.Flags().StringVar(&mutateTarget, "target", "", "API base URL, e.g. http://localhost:8080 (required)")
	_ = mutateCmd.MarkFlagRequired("target")
	mutateCmd.Flags().StringVar(&mutateOutput, "output", "", "Directory to write mutation-report.json (optional)")
	mutateCmd.Flags().StringVar(&mutateOperators, "operator", "", "Comma-separated operator names to run (default: all 12)")
	mutateCmd.Flags().String("spec", "", "OpenAPI spec file (optional; passed to LLM in Phase 2)")
	mutateCmd.Flags().IntVar(&mutateConcurrency, "concurrency", 4, "Number of cases processed concurrently per operator")
	mutateCmd.Flags().Bool("feedback", false, "Run LLM feedback analysis on survivors (requires LLM provider in .caseforge.yaml)")
	mutateCmd.Flags().Bool("auto-fix", false, "Patch index.json with suggested assertions (requires --feedback)")
	mutateCmd.Flags().Bool("yes", false, "Skip confirmation prompt for --auto-fix")
}

func runMutate(cmd *cobra.Command, _ []string) error {
	ops, err := resolveOperators(mutateOperators)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Running %d operator(s) × cases in %s...\n", len(ops), mutateCases)

	opts := mutation.RunOptions{
		Target:      mutateTarget,
		CasesDir:    mutateCases,
		Operators:   ops,
		Concurrency: mutateConcurrency,
	}

	run, err := mutation.Run(opts)
	if err != nil {
		return fmt.Errorf("mutation run: %w", err)
	}

	run.Clusters = mutation.ClusterSurvivors(run)

	feedbackFlag, _ := cmd.Flags().GetBool("feedback")
	if feedbackFlag && run.Survivors > 0 {
		cfg, cfgErr := config.Load()
		if cfgErr == nil {
			provider := llm.NewProviderWithConfig(llm.ProviderConfig{
				APIKey:   cfg.AI.APIKey,
				Provider: cfg.AI.Provider,
				Model:    cfg.AI.Model,
				BaseURL:  cfg.AI.BaseURL,
				Region:   cfg.AI.Region,
			})
			items, fbErr := mutation.Analyze(context.Background(), run, provider)
			if fbErr == nil && len(items) > 0 {
				run.Feedback = items
				fmt.Fprintf(out, "\nFeedback (%d cases with weak assertions):\n", len(items))
				for _, item := range items {
					fmt.Fprintf(out, "  ⚠ [%s] %s  risk=%.2f  → %d suggested assertion(s)\n",
						item.CaseID, item.Title, item.RiskScore, len(item.SuggestedAssertions))
				}
				autoFix, _ := cmd.Flags().GetBool("auto-fix")
				if autoFix {
					yes, _ := cmd.Flags().GetBool("yes")
					if err := runAutoFix(run, mutateCases, yes, out); err != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "auto-fix: %v\n", err)
					}
				} else {
					fmt.Fprintln(out, "Run with --feedback --auto-fix to patch index.json")
				}
			}
		}
	}

	for _, r := range run.Results {
		if r.Survived {
			color.New(color.FgRed).Fprintf(out, "  ✗ [%s] %s — not caught\n", r.Operator, r.Title)
		} else {
			color.New(color.FgGreen).Fprintf(out, "  ✓ [%s] %s — caught\n", r.Operator, r.Title)
		}
	}

	pct := 0
	if run.TotalRuns > 0 {
		pct = int(run.MutationScore * 100)
	}
	if run.Survivors == 0 {
		color.New(color.FgGreen).Fprintf(out, "\nMutation Score: %d/%d killed (%d%%) — no survivors\n",
			run.Killed, run.TotalRuns, pct)
	} else {
		color.New(color.FgYellow).Fprintf(out, "\nMutation Score: %d/%d killed (%d%%)\nSurvivors: %d\n",
			run.Killed, run.TotalRuns, pct, run.Survivors)
	}

	_ = mutation.Persist("", run)

	if mutateOutput != "" {
		if err := mutation.WriteReport(mutateOutput, run); err != nil {
			return fmt.Errorf("writing report: %w", err)
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "Report written to: %s\n",
			filepath.Join(mutateOutput, "mutation-report.json"))
	}

	if run.Survivors > 0 {
		os.Exit(ExitPartialSuccess)
	}
	return nil
}

func resolveOperators(filter string) ([]mutation.Operator, error) {
	all := mutation.Registry()
	if filter == "" {
		return all, nil
	}
	names := strings.Split(filter, ",")
	nameSet := map[string]bool{}
	for _, n := range names {
		nameSet[strings.TrimSpace(n)] = true
	}
	var ops []mutation.Operator
	for _, op := range all {
		if nameSet[op.Name()] {
			ops = append(ops, op)
		}
	}
	if len(ops) == 0 {
		return nil, fmt.Errorf("no operators matched --operator %q; valid names: %s",
			filter, operatorNames(all))
	}
	return ops, nil
}

func operatorNames(ops []mutation.Operator) string {
	names := make([]string, len(ops))
	for i, op := range ops {
		names[i] = op.Name()
	}
	return strings.Join(names, ", ")
}

func runAutoFix(_ mutation.MutationRun, _ string, _ bool, _ io.Writer) error {
	return fmt.Errorf("--auto-fix not yet implemented")
}
