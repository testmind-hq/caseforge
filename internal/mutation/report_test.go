// internal/mutation/report_test.go
package mutation_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/testmind-hq/caseforge/internal/mutation"
)

func sampleRun() mutation.MutationRun {
	return mutation.MutationRun{
		Target:        "http://api:8080",
		CasesDir:      "./cases",
		Operators:     []string{"field_drop", "status_swap_2xx"},
		TotalCases:    3,
		TotalRuns:     6,
		Killed:        4,
		Survivors:     2,
		MutationScore: 4.0 / 6.0,
		Results: []mutation.CaseMutationResult{
			{CaseID: "TC-0001", Title: "GET /pets", Operator: "field_drop", Survived: true},
			{CaseID: "TC-0001", Title: "GET /pets", Operator: "status_swap_2xx", Survived: false},
			{CaseID: "TC-0002", Title: "POST /pets", Operator: "field_drop", Survived: false},
			{CaseID: "TC-0002", Title: "POST /pets", Operator: "status_swap_2xx", Survived: false},
			{CaseID: "TC-0003", Title: "DELETE /pets/{id}", Operator: "field_drop", Survived: true},
			{CaseID: "TC-0003", Title: "DELETE /pets/{id}", Operator: "status_swap_2xx", Survived: false},
		},
		GeneratedAt: "2026-05-10T12:00:00Z",
	}
}

func TestTextSummary(t *testing.T) {
	run := sampleRun()
	summary := mutation.TextSummary(run)
	if !strings.Contains(summary, "Mutation Score") {
		t.Error("summary must contain 'Mutation Score'")
	}
	if !strings.Contains(summary, "Survivors: 2") {
		t.Error("summary must contain 'Survivors: 2'")
	}
}

func TestWriteReport(t *testing.T) {
	dir := t.TempDir()
	run := sampleRun()
	if err := mutation.WriteReport(dir, run); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "mutation-report.json"))
	if err != nil {
		t.Fatal(err)
	}
	var decoded mutation.MutationRun
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if decoded.Survivors != 2 {
		t.Fatalf("expected 2 survivors in report, got %d", decoded.Survivors)
	}
}

func TestPersist(t *testing.T) {
	dir := t.TempDir()
	run := sampleRun()
	if err := mutation.Persist(dir, run); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 run file, got %d", len(entries))
	}
	if !strings.HasSuffix(entries[0].Name(), ".json") {
		t.Fatalf("run file must be .json, got %s", entries[0].Name())
	}
}

func TestClusterSurvivors(t *testing.T) {
	run := sampleRun()
	clusters := mutation.ClusterSurvivors(run)
	if len(clusters) != 2 {
		t.Fatalf("expected 2 clusters (TC-0001 and TC-0003), got %d", len(clusters))
	}
	// Clusters sorted by RiskScore descending
	if clusters[0].RiskScore < clusters[1].RiskScore {
		t.Error("clusters must be sorted by RiskScore descending")
	}
}
