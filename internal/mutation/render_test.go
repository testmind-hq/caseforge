// internal/mutation/render_test.go
package mutation_test

import (
	"strings"
	"testing"

	"github.com/testmind-hq/caseforge/internal/mutation"
)

func sampleRunWithClusters() mutation.MutationRun {
	run := sampleRun() // defined in report_test.go — same package mutation_test
	run.Clusters = mutation.ClusterSurvivors(run)
	return run
}

func sampleRunWithFeedback() mutation.MutationRun {
	run := sampleRunWithClusters()
	run.Feedback = []mutation.FeedbackItem{
		{
			CaseID:    "TC-0001",
			Title:     "GET /pets",
			RiskScore: 0.5,
			Diagnosis: "1 operator survived",
			SuggestedAssertions: []mutation.SuggestedAssertion{
				{Target: "jsonpath $.id", Operator: "exists"},
				{Target: "status_code", Operator: "eq", Expected: float64(200)},
			},
		},
	}
	return run
}

func TestRenderMarkdown_NoSurvivors(t *testing.T) {
	run := mutation.MutationRun{
		TotalRuns:     4,
		Killed:        4,
		Survivors:     0,
		MutationScore: 1.0,
		GeneratedAt:   "2026-05-10T12:00:00Z",
	}
	md := mutation.RenderMarkdown(run)
	if !strings.Contains(md, "100%") {
		t.Error("must show 100% when no survivors")
	}
	if strings.Contains(md, "## Survivor Summary") {
		t.Error("must not show Survivor Summary when no survivors")
	}
}

func TestRenderMarkdown_WithSurvivors(t *testing.T) {
	run := sampleRunWithClusters()
	md := mutation.RenderMarkdown(run)
	if !strings.Contains(md, "## Survivor Summary") {
		t.Error("must contain Survivor Summary section")
	}
	if !strings.Contains(md, "TC-0001") {
		t.Error("TC-0001 must appear in Survivor Summary")
	}
	if !strings.Contains(md, "TC-0003") {
		t.Error("TC-0003 must appear in Survivor Summary")
	}
}

func TestRenderMarkdown_WithFeedback(t *testing.T) {
	run := sampleRunWithFeedback()
	md := mutation.RenderMarkdown(run)
	if !strings.Contains(md, "## Suggested Assertions") {
		t.Error("must contain Suggested Assertions section when Feedback is non-empty")
	}
	if !strings.Contains(md, "jsonpath $.id") {
		t.Error("must list the suggested assertion target")
	}
}

func TestRenderHTML_ContainsHeatmap(t *testing.T) {
	run := sampleRunWithClusters()
	out := mutation.RenderHTML(run)
	if !strings.Contains(out, "<table") {
		t.Error("HTML must contain a table element")
	}
	if !strings.Contains(out, "field_drop") {
		t.Error("HTML must contain operator name field_drop")
	}
	if !strings.Contains(out, "TC-0001") {
		t.Error("HTML must contain case ID TC-0001")
	}
}

func TestRenderHTML_SelfContained(t *testing.T) {
	run := sampleRunWithClusters()
	out := mutation.RenderHTML(run)
	if strings.Contains(out, "http://") || strings.Contains(out, "https://") {
		t.Error("HTML report must not reference external URLs")
	}
}

func TestRenderHTML_NoSurvivors(t *testing.T) {
	run := mutation.MutationRun{
		Operators:     []string{"field_drop"},
		TotalRuns:     1,
		Killed:        1,
		Survivors:     0,
		MutationScore: 1.0,
		Results: []mutation.CaseMutationResult{
			{CaseID: "TC-0001", Title: "GET /pets", Operator: "field_drop", Survived: false},
		},
		GeneratedAt: "2026-05-10T12:00:00Z",
	}
	out := mutation.RenderHTML(run)
	if strings.Contains(out, "Survivors by Risk") {
		t.Error("must not show Survivors by Risk table when no clusters")
	}
}

func TestRenderMarkdown_OperationTable(t *testing.T) {
	run := sampleRun()
	run.OperationScores = []mutation.OperationScore{
		{Operation: "GET /pets", TotalRuns: 2, Killed: 1, Survivors: 1, MutationScore: 0.5},
		{Operation: "POST /pets", TotalRuns: 2, Killed: 2, Survivors: 0, MutationScore: 1.0},
	}
	out := mutation.RenderMarkdown(run)
	if !strings.Contains(out, "Per-Operation Mutation Score") {
		t.Error("expected per-operation section header")
	}
	if !strings.Contains(out, "GET /pets") {
		t.Error("expected GET /pets row")
	}
	if !strings.Contains(out, "50%") {
		t.Error("expected 50% score for GET /pets")
	}
	if !strings.Contains(out, "⚠") {
		t.Error("expected ⚠ badge for score < 70%")
	}
	if !strings.Contains(out, "✓") {
		t.Error("expected ✓ badge for 100% score")
	}
}

func TestRenderMarkdown_NoOperationScores(t *testing.T) {
	run := sampleRun()
	run.OperationScores = nil
	out := mutation.RenderMarkdown(run)
	if strings.Contains(out, "Per-Operation") {
		t.Error("per-operation section must be absent when OperationScores is nil")
	}
}

func TestRenderHTML_OperationTable(t *testing.T) {
	run := sampleRun()
	run.OperationScores = []mutation.OperationScore{
		{Operation: "GET /pets", TotalRuns: 2, Killed: 1, Survivors: 1, MutationScore: 0.5},
		{Operation: "DELETE /pets/{id}", TotalRuns: 2, Killed: 2, Survivors: 0, MutationScore: 1.0},
	}
	out := mutation.RenderHTML(run)
	if !strings.Contains(out, "Per-Operation") {
		t.Error("expected Per-Operation section in HTML")
	}
	if !strings.Contains(out, "GET /pets") {
		t.Error("expected GET /pets row in HTML")
	}
	if !strings.Contains(out, "fef2f2") {
		t.Error("expected light-red background for score < 70%")
	}
	if !strings.Contains(out, "f0fdf4") {
		t.Error("expected light-green background for 100% score")
	}
}
