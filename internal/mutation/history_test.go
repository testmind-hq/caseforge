package mutation_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/testmind-hq/caseforge/internal/mutation"
)

func writeRunFile(t *testing.T, dir, name string, run mutation.MutationRun) {
	t.Helper()
	data, err := json.Marshal(run)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), data, 0644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadHistory_Empty(t *testing.T) {
	runs, err := mutation.LoadHistory(t.TempDir()+"/nonexistent", 10)
	if err != nil {
		t.Fatalf("expected nil error for missing dir, got %v", err)
	}
	if runs != nil {
		t.Fatalf("expected nil slice for missing dir, got %v", runs)
	}
}

func TestLoadHistory_ReadsFiles(t *testing.T) {
	dir := t.TempDir()
	for i := 1; i <= 3; i++ {
		writeRunFile(t, dir, fmt.Sprintf("2026-05-0%d.json", i), mutation.MutationRun{
			GeneratedAt: fmt.Sprintf("2026-05-0%dT12:00:00Z", i),
			TotalRuns:   10,
			Killed:      i * 2,
			MutationScore: float64(i*2) / 10.0,
		})
	}
	runs, err := mutation.LoadHistory(dir, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(runs) != 3 {
		t.Fatalf("expected 3 runs, got %d", len(runs))
	}
	// Newest first: 2026-05-03 should come before 2026-05-01
	if runs[0].GeneratedAt < runs[2].GeneratedAt {
		t.Error("runs must be ordered newest-first")
	}
}

func TestLoadHistory_Limit(t *testing.T) {
	dir := t.TempDir()
	for i := 1; i <= 5; i++ {
		writeRunFile(t, dir, fmt.Sprintf("2026-05-%02d.json", i), mutation.MutationRun{
			GeneratedAt: fmt.Sprintf("2026-05-%02dT12:00:00Z", i),
		})
	}
	runs, err := mutation.LoadHistory(dir, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(runs) != 3 {
		t.Fatalf("expected 3 runs with limit=3, got %d", len(runs))
	}
}

func TestLoadHistory_SkipsCorrupt(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "corrupt.json"), []byte("{bad json"), 0644); err != nil {
		t.Fatal(err)
	}
	writeRunFile(t, dir, "valid.json", mutation.MutationRun{GeneratedAt: "2026-05-10T12:00:00Z"})
	runs, err := mutation.LoadHistory(dir, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(runs) != 1 {
		t.Fatalf("expected 1 valid run, got %d", len(runs))
	}
}

func TestRenderHistory_Empty(t *testing.T) {
	out := mutation.RenderHistory(nil)
	if !strings.Contains(out, "No mutation history found") {
		t.Error("empty history must contain 'No mutation history found'")
	}
}

func TestRenderHistory_Trend(t *testing.T) {
	runs := []mutation.MutationRun{
		{GeneratedAt: "2026-05-03T12:00:00Z", TotalRuns: 10, Killed: 8, MutationScore: 0.8},
		{GeneratedAt: "2026-05-02T12:00:00Z", TotalRuns: 10, Killed: 6, MutationScore: 0.6},
		{GeneratedAt: "2026-05-01T12:00:00Z", TotalRuns: 10, Killed: 6, MutationScore: 0.6},
	}
	out := mutation.RenderHistory(runs)
	// First run (newest) should show ↑+20% compared to second
	if !strings.Contains(out, "↑+20%") {
		t.Errorf("expected ↑+20%% delta, got:\n%s", out)
	}
	// Second and third are equal → show —
	if !strings.Contains(out, "—") {
		t.Errorf("expected — delta for equal scores, got:\n%s", out)
	}
	// Last run shows — (no prior)
}

func TestRenderHistory_7DayAvg(t *testing.T) {
	now := time.Now().UTC()
	recent := now.AddDate(0, 0, -1).Format("2006-01-02T15:04:05Z")
	old := now.AddDate(0, 0, -10).Format("2006-01-02T15:04:05Z")
	runs := []mutation.MutationRun{
		{GeneratedAt: recent, TotalRuns: 10, Killed: 8, MutationScore: 0.8},
		{GeneratedAt: old, TotalRuns: 10, Killed: 4, MutationScore: 0.4},
	}
	out := mutation.RenderHistory(runs)
	// Only the recent run (80%) should be in the 7-day average
	if !strings.Contains(out, "Avg score (last 7 days): 80%") {
		t.Errorf("expected 7-day avg of 80%%, got:\n%s", out)
	}
}
