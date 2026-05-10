// internal/mutation/report.go
package mutation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// TextSummary returns a multi-line human-readable summary of a MutationRun.
func TextSummary(run MutationRun) string {
	var sb strings.Builder
	for _, r := range run.Results {
		if r.Survived {
			fmt.Fprintf(&sb, "  ✗ [%s] %s — mutation not caught\n", r.Operator, r.Title)
		} else {
			fmt.Fprintf(&sb, "  ✓ [%s] %s — caught\n", r.Operator, r.Title)
		}
	}
	pct := 0
	if run.TotalRuns > 0 {
		pct = int(run.MutationScore * 100)
	}
	fmt.Fprintf(&sb, "\nMutation Score: %d/%d killed (%d%%)\n", run.Killed, run.TotalRuns, pct)
	fmt.Fprintf(&sb, "Survivors: %d", run.Survivors)
	return sb.String()
}

// WriteReport writes mutation report files in the requested formats to outputDir.
// formats: slice containing any of "json", "markdown", "html".
func WriteReport(outputDir string, run MutationRun, formats []string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}
	for _, f := range formats {
		switch f {
		case "json":
			data, err := json.MarshalIndent(run, "", "  ")
			if err != nil {
				return err
			}
			if err := os.WriteFile(filepath.Join(outputDir, "mutation-report.json"), data, 0644); err != nil {
				return err
			}
		case "markdown":
			if err := os.WriteFile(filepath.Join(outputDir, "mutation-report.md"),
				[]byte(RenderMarkdown(run)), 0644); err != nil {
				return err
			}
		case "html":
			if err := os.WriteFile(filepath.Join(outputDir, "mutation-report.html"),
				[]byte(RenderHTML(run)), 0644); err != nil {
				return err
			}
		}
	}
	return nil
}

// Persist writes a timestamped run file to baseDir.
// If baseDir is empty, uses the default .caseforge/mutation/runs directory.
func Persist(baseDir string, run MutationRun) error {
	dir := baseDir
	if dir == "" {
		dir = persistDir()
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(run, "", "  ")
	if err != nil {
		return err
	}
	name := timestamp() + ".json"
	return os.WriteFile(filepath.Join(dir, name), data, 0644)
}

// ClusterSurvivors groups results by case ID, computes RiskScore, sorts descending.
func ClusterSurvivors(run MutationRun) []SurvivorCluster {
	byCase := map[string]*SurvivorCluster{}
	for _, r := range run.Results {
		if !r.Survived {
			continue
		}
		c, ok := byCase[r.CaseID]
		if !ok {
			byCase[r.CaseID] = &SurvivorCluster{CaseID: r.CaseID, Title: r.Title}
			c = byCase[r.CaseID]
		}
		c.Operators = append(c.Operators, r.Operator)
	}

	totalOps := len(run.Operators)
	if totalOps == 0 {
		totalOps = 1
	}

	clusters := make([]SurvivorCluster, 0, len(byCase))
	for _, c := range byCase {
		c.RiskScore = float64(len(c.Operators)) / float64(totalOps)
		clusters = append(clusters, *c)
	}
	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].RiskScore > clusters[j].RiskScore
	})
	return clusters
}
