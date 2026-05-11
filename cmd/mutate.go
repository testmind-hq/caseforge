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
	mutateCases             string
	mutateTarget            string
	mutateOutput            string
	mutateOperators         string
	mutateReportFormat      string
	mutateConcurrency       int
	mutateOperatorConcurrency int
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
	mutateCmd.Flags().StringVar(&mutateCases, "cases", "", "Directory containing index.json and .hurl files (required when not using --history)")
	mutateCmd.Flags().StringVar(&mutateTarget, "target", "", "API base URL, e.g. http://localhost:8080 (required when not using --history)")
	mutateCmd.Flags().StringVar(&mutateOutput, "output", "", "Directory to write mutation-report.json (optional)")
	mutateCmd.Flags().StringVar(&mutateOperators, "operator", "", "Comma-separated operator names to run (default: all 12)")
	mutateCmd.Flags().String("spec", "", "OpenAPI spec file (optional; passed to LLM in Phase 2)")
	mutateCmd.Flags().IntVar(&mutateConcurrency, "concurrency", 4, "Number of cases processed concurrently per operator")
	mutateCmd.Flags().IntVar(&mutateOperatorConcurrency, "operator-concurrency", 2, "Number of operators to run in parallel")
	mutateCmd.Flags().StringVar(&mutateReportFormat, "report-format", "json", `Comma-separated report formats: json,markdown,html,all`)
	mutateCmd.Flags().Bool("history", false, "Print mutation score history (does not run mutations; --target not required)")
	mutateCmd.Flags().Int("history-limit", 10, "Maximum number of historical runs to display")
	mutateCmd.Flags().Float64("min-score", 0,
		"Per-operation minimum mutation score (0.0–1.0); exit 6 if any operation scores below this")
	mutateCmd.Flags().Bool("feedback", false, "Run LLM feedback analysis on survivors (requires LLM provider in .caseforge.yaml)")
	mutateCmd.Flags().Bool("auto-fix", false, "Patch index.json with suggested assertions (requires --feedback)")
	mutateCmd.Flags().Bool("yes", false, "Skip confirmation prompt for --auto-fix")
}

func runMutate(cmd *cobra.Command, _ []string) error {
	if historyFlag, _ := cmd.Flags().GetBool("history"); historyFlag {
		limit, _ := cmd.Flags().GetInt("history-limit")
		runs, err := mutation.LoadHistory("", limit)
		if err != nil {
			return fmt.Errorf("loading history: %w", err)
		}
		fmt.Fprint(cmd.OutOrStdout(), mutation.RenderHistory(runs))
		return nil
	}

	// --cases and --target are required when not using --history
	if mutateCases == "" {
		return fmt.Errorf("required flag \"cases\" not set")
	}
	if mutateTarget == "" {
		return fmt.Errorf("required flag \"target\" not set")
	}

	feedbackFlag, _ := cmd.Flags().GetBool("feedback")
	autoFixFlag, _ := cmd.Flags().GetBool("auto-fix")
	if autoFixFlag && !feedbackFlag {
		return fmt.Errorf("--auto-fix requires --feedback")
	}

	// Validate --report-format unconditionally so a typo is always caught early.
	reportFormats, err := parseReportFormats(mutateReportFormat)
	if err != nil {
		return err
	}

	minScore, _ := cmd.Flags().GetFloat64("min-score")
	if minScore < 0 || minScore > 1.0 {
		return fmt.Errorf("--min-score must be between 0.0 and 1.0, got %.2f", minScore)
	}

	ops, err := resolveOperators(mutateOperators)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Running %d operator(s) × cases in %s...\n", len(ops), mutateCases)

	opts := mutation.RunOptions{
		Target:              mutateTarget,
		CasesDir:            mutateCases,
		Operators:           ops,
		Concurrency:         mutateConcurrency,
		OperatorConcurrency: mutateOperatorConcurrency,
	}

	run, err := mutation.Run(opts)
	if err != nil {
		return fmt.Errorf("mutation run: %w", err)
	}

	run.Clusters = mutation.ClusterSurvivors(run)

	if feedbackFlag && run.Survivors > 0 {
		cfg, cfgErr := config.Load()
		if cfgErr != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "warn: --feedback: failed to load config: %v\n", cfgErr)
		} else {
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
		if err := mutation.WriteReport(mutateOutput, run, reportFormats); err != nil {
			return fmt.Errorf("writing report: %w", err)
		}
		extMap := map[string]string{"json": ".json", "markdown": ".md", "html": ".html"}
		for _, f := range reportFormats {
			if ext, ok := extMap[f]; ok {
				fmt.Fprintf(cmd.ErrOrStderr(), "Report written to: %s\n",
					filepath.Join(mutateOutput, "mutation-report"+ext))
			}
		}
	}

	if minScore > 0 {
		var failing []mutation.OperationScore
		for _, op := range run.OperationScores {
			if op.MutationScore < minScore {
				failing = append(failing, op)
			}
		}
		if len(failing) > 0 {
			color.New(color.FgRed).Fprintf(out,
				"\n✗ %d operation(s) below --min-score %.0f%%:\n", len(failing), minScore*100)
			for _, op := range failing {
				fmt.Fprintf(out, "  %-52s %3.0f%%\n", op.Operation, op.MutationScore*100)
			}
			os.Exit(ExitPartialSuccess)
		}
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

func runAutoFix(run mutation.MutationRun, casesDir string, skipConfirm bool, out io.Writer) error {
	if len(run.Feedback) == 0 {
		return nil
	}

	fmt.Fprintf(out, "\nAuto-fix will append %d assertion(s) across %d case(s).\n",
		countSuggestions(run.Feedback), len(run.Feedback))

	if !skipConfirm {
		fmt.Fprint(out, "Apply? [y/N] ")
		var answer string
		fmt.Scanln(&answer)
		if answer != "y" && answer != "Y" {
			fmt.Fprintln(out, "Skipped.")
			return nil
		}
	}

	if err := mutation.PatchIndex(casesDir, run.Feedback); err != nil {
		return fmt.Errorf("patching index.json: %w", err)
	}
	fmt.Fprintln(out, "index.json updated.")
	return nil
}

func countSuggestions(items []mutation.FeedbackItem) int {
	n := 0
	for _, item := range items {
		n += len(item.SuggestedAssertions)
	}
	return n
}

// parseReportFormats splits the comma-separated format string, expands "all"
// to ["json", "markdown", "html"], and deduplicates. Returns ["json"] if empty.
// Returns an error for unrecognized format names.
func parseReportFormats(s string) ([]string, error) {
	if s == "" {
		return []string{"json"}, nil
	}
	seen := map[string]bool{}
	var out []string
	for fmtName := range strings.SplitSeq(s, ",") {
		fmtName = strings.ToLower(strings.TrimSpace(fmtName))
		if fmtName == "all" {
			for _, fn := range []string{"json", "markdown", "html"} {
				if !seen[fn] {
					seen[fn] = true
					out = append(out, fn)
				}
			}
		} else if fmtName != "" {
			switch fmtName {
			case "json", "markdown", "html":
				if !seen[fmtName] {
					seen[fmtName] = true
					out = append(out, fmtName)
				}
			default:
				return nil, fmt.Errorf("unknown report format %q; valid: json, markdown, html, all", fmtName)
			}
		}
	}
	if len(out) == 0 {
		return []string{"json"}, nil
	}
	return out, nil
}
