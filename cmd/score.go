// cmd/score.go
package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/output/writer"
	"github.com/testmind-hq/caseforge/internal/score"
)

var scoreCmd = &cobra.Command{
	Use:   "score",
	Short: "Score the quality of generated test cases",
	Long: `Score reads index.json from a cases directory and evaluates quality
across four dimensions: coverage breadth, boundary coverage, security
coverage, and executability. It also generates improvement suggestions.

Examples:
  caseforge score
  caseforge score --cases ./cases --format json`,
	RunE:         runScore,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(scoreCmd)
	scoreCmd.Flags().String("cases", "./cases", "Directory containing index.json")
	scoreCmd.Flags().String("format", "terminal", "Output format: terminal or json")
	scoreCmd.Flags().Int("min-score", 0, "Exit with code 1 if Overall score is below this threshold (0 = disabled)")
	scoreCmd.Flags().Bool("save-history", false, "Append current score to .caseforge-conformance.json")
	scoreCmd.Flags().String("spec", "", "OpenAPI spec path (required for --fill-gaps)")
	scoreCmd.Flags().Bool("fill-gaps", false, "Auto-generate cases for operations missing 2xx or 4xx coverage")
}

const conformanceHistoryFile = ".caseforge-conformance.json"

// scoreOutput wraps Report with a conformance block for JSON output.
type scoreOutput struct {
	score.Report
	Conformance score.ConformanceResult `json:"conformance"`
}

func runScore(cmd *cobra.Command, _ []string) error {
	casesDir, _ := cmd.Flags().GetString("cases")
	format, _ := cmd.Flags().GetString("format")
	minScore, _ := cmd.Flags().GetInt("min-score")
	saveHistory, _ := cmd.Flags().GetBool("save-history")

	if format != "terminal" && format != "json" {
		return fmt.Errorf("unknown --format %q: must be terminal or json", format)
	}

	indexPath := filepath.Join(casesDir, "index.json")
	w := writer.NewJSONSchemaWriter()
	cases, err := w.Read(indexPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", indexPath, err)
	}

	report := score.Compute(cases)

	// Load history and compute conformance.
	history, err := score.LoadHistory(conformanceHistoryFile)
	if err != nil {
		return fmt.Errorf("loading conformance history: %w", err)
	}
	conformance := score.ComputeConformance(report, history)

	// Optionally save history.
	if saveHistory {
		if saveErr := score.SaveHistory(conformanceHistoryFile, history, report.Overall); saveErr != nil {
			return fmt.Errorf("saving conformance history: %w", saveErr)
		}
	}

	fillGaps, _ := cmd.Flags().GetBool("fill-gaps")
	specPath, _ := cmd.Flags().GetString("spec")
	if fillGaps {
		if specPath == "" {
			return fmt.Errorf("--fill-gaps requires --spec <spec-file>")
		}
		if err := runFillGaps(cmd, cases, specPath, casesDir); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "fill-gaps: %v\n", err)
		}
	}

	if format == "json" {
		out := scoreOutput{Report: report, Conformance: conformance}
		data, marshalErr := json.MarshalIndent(out, "", "  ")
		if marshalErr != nil {
			return fmt.Errorf("marshaling report: %w", marshalErr)
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
	} else {
		printScoreReport(cmd, report, casesDir)
		fmt.Fprintf(cmd.OutOrStdout(), "Conformance score: %d   trend: %s\n", conformance.Score, conformance.Trend)
	}

	// CI gate: exit 1 if below threshold.
	if minScore > 0 && report.Overall < minScore {
		return fmt.Errorf("score %d is below minimum threshold %d", report.Overall, minScore)
	}

	return nil
}

func printScoreReport(cmd *cobra.Command, r score.Report, casesDir string) {
	out := cmd.OutOrStdout()
	sep := strings.Repeat("─", 44)

	fmt.Fprintf(out, "\nTest case quality score (%s)\n", casesDir)
	fmt.Fprintf(out, "Overall: %d / 100\n", r.Overall)
	fmt.Fprintln(out, sep)

	for _, d := range r.Dimensions {
		fmt.Fprintf(out, "%-20s  %3d   %s\n", d.Name, d.Score, d.Detail)
	}

	if len(r.Suggestions) > 0 {
		fmt.Fprintln(out, "\nSuggestions (by priority):")
		for _, s := range r.Suggestions {
			fmt.Fprintf(out, "  %d. %s\n", s.Priority, s.Message)
			if s.Command != "" {
				fmt.Fprintf(out, "     → %s\n", s.Command)
			}
		}
	}

	fmt.Fprintf(out, "\n%d case(s) across %d operation(s)\n", r.TotalCases, r.TotalOps)
}

func runFillGaps(cmd *cobra.Command, cases []schema.TestCase, specPath, casesDir string) error {
	gaps := score.ComputeGaps(cases)
	if len(gaps) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No coverage gaps found.")
		return nil
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Found %d operation(s) with coverage gaps. Generating...\n", len(gaps))

	seen := map[string]bool{}
	var techniques []string
	for _, gap := range gaps {
		if !gap.Has2xx {
			for _, t := range []string{"positive_examples", "equivalence_partitioning"} {
				if !seen[t] {
					seen[t] = true
					techniques = append(techniques, t)
				}
			}
		}
		if !gap.Has4xx {
			for _, t := range []string{"mutation", "isolated_negative", "required_omission"} {
				if !seen[t] {
					seen[t] = true
					techniques = append(techniques, t)
				}
			}
		}
	}
	techniqueFlag := strings.Join(techniques, ",")
	genArgs := []string{
		"--spec", specPath,
		"--no-ai",
		"--technique", techniqueFlag,
		"--output", casesDir,
	}
	fmt.Fprintf(cmd.OutOrStdout(), "  → caseforge gen %s\n", strings.Join(genArgs, " "))
	genCmd.SetArgs(genArgs)
	if err := genCmd.Execute(); err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "  warn: gen failed: %v\n", err)
	}
	return nil
}
