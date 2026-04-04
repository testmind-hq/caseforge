// internal/dedupe/report.go
package dedupe

import (
	"encoding/json"
	"fmt"
	"io"
)

// PrintTerminal writes the DedupeReport in human-readable format to w.
func PrintTerminal(w io.Writer, report DedupeReport, dryRun bool) {
	total := report.ExactGroups + report.StructuralGroups
	fmt.Fprintf(w, "Duplicate Groups Found: %d\n\n", total)

	for i, g := range report.Groups {
		fmt.Fprintf(w, "Group %d (%s, similarity: %.2f)\n", i+1, g.Kind, g.Similarity)
		for _, cs := range g.Cases {
			action := "DELETE"
			if cs.Keep {
				action = "KEEP  "
			}
			fmt.Fprintf(w, "  %s %s  [score: %d assertions]\n",
				action, cs.FilePath, cs.AssertionCount)
		}
		fmt.Fprintln(w)
	}

	dupCount := 0
	for _, g := range report.Groups {
		dupCount += len(g.Cases) - 1
	}

	mergeHint := ""
	if !dryRun && dupCount > 0 {
		mergeHint = " Run with --merge to delete duplicates."
	}

	fmt.Fprintf(w, "Summary: %d duplicate(s) found (%d exact, %d structural).%s\n",
		dupCount, report.ExactGroups, report.StructuralGroups, mergeHint)
}

// MarshalReportJSON serializes a DedupeReport to indented JSON with snake_case keys.
func MarshalReportJSON(report DedupeReport) ([]byte, error) {
	type out struct {
		CasesDir         string           `json:"cases_dir"`
		TotalScanned     int              `json:"total_scanned"`
		ExactGroups      int              `json:"exact_groups"`
		StructuralGroups int              `json:"structural_groups"`
		Groups           []DuplicateGroup `json:"groups"`
		GeneratedAt      string           `json:"generated_at"`
	}
	return json.MarshalIndent(out{
		CasesDir:         report.CasesDir,
		TotalScanned:     report.TotalScanned,
		ExactGroups:      report.ExactGroups,
		StructuralGroups: report.StructuralGroups,
		Groups:           report.Groups,
		GeneratedAt:      report.GeneratedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}, "", "  ")
}
