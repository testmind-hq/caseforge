package mutation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// LoadHistory reads up to limit most-recent MutationRun files from baseDir.
// If baseDir is empty it defaults to ".caseforge/mutation/runs".
// Returns nil slice (not error) when the directory does not exist yet.
func LoadHistory(baseDir string, limit int) ([]MutationRun, error) {
	if baseDir == "" {
		baseDir = filepath.Join(".caseforge", "mutation", "runs")
	}
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	// sort descending by filename (timestamp prefix ensures chronological order)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() > entries[j].Name()
	})
	if limit > 0 && len(entries) > limit {
		entries = entries[:limit]
	}
	runs := make([]MutationRun, 0, len(entries))
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(baseDir, e.Name()))
		if err != nil {
			continue // skip unreadable files
		}
		var run MutationRun
		if err := json.Unmarshal(data, &run); err != nil {
			continue // skip corrupt files
		}
		runs = append(runs, run)
	}
	return runs, nil
}

// RenderHistory returns a terminal-friendly trend table.
// Runs must be ordered newest-first (as returned by LoadHistory).
func RenderHistory(runs []MutationRun) string {
	if len(runs) == 0 {
		return "No mutation history found. Run `caseforge mutate` first.\n"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Mutation History (last %d runs)\n", len(runs))
	fmt.Fprintf(&b, "%s\n", strings.Repeat("─", 70))
	fmt.Fprintf(&b, "  %-26s  %5s  %10s  %9s  %4s\n",
		"Run", "Score", "Killed", "Survivors", "Δ")
	for i, run := range runs {
		pct := 0
		if run.TotalRuns > 0 {
			pct = int(run.MutationScore * 100)
		}
		delta := ""
		if i+1 < len(runs) {
			prev := runs[i+1]
			prevPct := 0
			if prev.TotalRuns > 0 {
				prevPct = int(prev.MutationScore * 100)
			}
			diff := pct - prevPct
			if diff > 0 {
				delta = fmt.Sprintf("↑+%d%%", diff)
			} else if diff < 0 {
				delta = fmt.Sprintf("↓%d%%", diff)
			} else {
				delta = "—"
			}
		} else {
			delta = "—"
		}
		ts := run.GeneratedAt
		if ts == "" {
			ts = "(unknown)"
		}
		fmt.Fprintf(&b, "  %-26s  %4d%%  %4d/%-4d  %9d  %s\n",
			ts, pct, run.Killed, run.TotalRuns, run.Survivors, delta)
	}
	fmt.Fprintf(&b, "%s\n", strings.Repeat("─", 70))
	// 7-day average
	cutoff := time.Now().AddDate(0, 0, -7)
	var recent []int
	for _, run := range runs {
		t, err := time.Parse("2006-01-02T15:04:05Z", run.GeneratedAt)
		if err == nil && t.After(cutoff) && run.TotalRuns > 0 {
			recent = append(recent, int(run.MutationScore*100))
		}
	}
	if len(recent) > 0 {
		sum := 0
		for _, v := range recent {
			sum += v
		}
		fmt.Fprintf(&b, "Avg score (last 7 days): %d%%\n", sum/len(recent))
	}
	return b.String()
}
