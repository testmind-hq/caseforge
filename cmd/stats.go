// cmd/stats.go
package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/output/writer"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show test case statistics for a cases directory",
	Long: `Stats reads index.json from a cases directory and prints a summary
dashboard: total case count, coverage by technique and priority, and
generation metadata.

Examples:
  caseforge stats
  caseforge stats --cases ./cases --format json`,
	RunE:         runStats,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(statsCmd)
	statsCmd.Flags().String("cases", "./cases", "Directory containing index.json")
	statsCmd.Flags().String("format", "terminal", "Output format: terminal or json")
}

func runStats(cmd *cobra.Command, _ []string) error {
	casesDir, _ := cmd.Flags().GetString("cases")
	format, _ := cmd.Flags().GetString("format")

	if format != "terminal" && format != "json" {
		return fmt.Errorf("unknown --format %q: must be terminal or json", format)
	}

	indexPath := filepath.Join(casesDir, "index.json")
	w := writer.NewJSONSchemaWriter()
	index, err := w.ReadFull(indexPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", indexPath, err)
	}

	if format == "json" {
		type statsJSON struct {
			CasesDir    string         `json:"cases_dir"`
			Total       int            `json:"total"`
			GeneratedAt time.Time      `json:"generated_at"`
			ByTechnique map[string]int `json:"by_technique,omitempty"`
			ByPriority  map[string]int `json:"by_priority,omitempty"`
			ByKind      map[string]int `json:"by_kind,omitempty"`
		}
		payload := statsJSON{
			CasesDir:    casesDir,
			Total:       len(index.TestCases),
			GeneratedAt: index.GeneratedAt,
			ByTechnique: index.Meta.ByTechnique,
			ByPriority:  index.Meta.ByPriority,
			ByKind:      index.Meta.ByKind,
		}
		data, marshalErr := json.MarshalIndent(payload, "", "  ")
		if marshalErr != nil {
			return fmt.Errorf("marshaling stats: %w", marshalErr)
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return nil
	}

	// Terminal output
	out := cmd.OutOrStdout()
	total := len(index.TestCases)
	sep := strings.Repeat("─", 44)

	fmt.Fprintf(out, "\nTest case stats (%s)\n", casesDir)
	fmt.Fprintln(out, sep)
	fmt.Fprintf(out, "Total cases:    %d\n", total)
	fmt.Fprintf(out, "Generated at:   %s\n", index.GeneratedAt.Format("2006-01-02 15:04:05"))

	if index.Meta.CaseforgeVersion != "" {
		fmt.Fprintf(out, "Version:        %s\n", index.Meta.CaseforgeVersion)
	}

	if len(index.Meta.ByTechnique) > 0 {
		fmt.Fprintln(out, "\nTechnique distribution:")
		keys := sortedKeys(index.Meta.ByTechnique)
		// Use actual total; fall back to sum of technique counts when test_cases is empty
		denominator := total
		if denominator == 0 {
			for _, n := range index.Meta.ByTechnique {
				denominator += n
			}
		}
		for _, k := range keys {
			n := index.Meta.ByTechnique[k]
			var pct float64
			if denominator > 0 {
				pct = float64(n) * 100 / float64(denominator)
			}
			barLen := int(pct / 5)
			if barLen < 0 {
				barLen = 0
			}
			bar := strings.Repeat("█", barLen)
			fmt.Fprintf(out, "  %-28s %4d  %-16s %.0f%%\n", k, n, bar, pct)
		}
	}

	if len(index.Meta.ByPriority) > 0 {
		fmt.Fprintln(out, "\nPriority distribution:")
		for _, p := range []string{"P0", "P1", "P2", "P3"} {
			if n, ok := index.Meta.ByPriority[p]; ok {
				fmt.Fprintf(out, "  %s  %d\n", p, n)
			}
		}
	}

	if len(index.Meta.ByKind) > 0 {
		fmt.Fprintln(out, "\nCase kinds:")
		keys := sortedKeys(index.Meta.ByKind)
		for _, k := range keys {
			fmt.Fprintf(out, "  %-12s %d\n", k, index.Meta.ByKind[k])
		}
	}

	fmt.Fprintln(out)
	return nil
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Primary: descending count. Secondary: alphabetical for determinism on ties.
	sort.Slice(keys, func(i, j int) bool {
		if m[keys[i]] != m[keys[j]] {
			return m[keys[i]] > m[keys[j]]
		}
		return keys[i] < keys[j]
	})
	return keys
}
