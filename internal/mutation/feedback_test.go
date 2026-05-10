// internal/mutation/feedback_test.go
package mutation_test

import (
	"context"
	"strings"
	"testing"

	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/mutation"
)

func TestAnalyze_UnavailableProvider(t *testing.T) {
	run := sampleRun()
	run.Clusters = mutation.ClusterSurvivors(run)

	// "noop" provider name → NoopProvider → IsAvailable() = false
	provider := llm.NewProvider("", "noop", "")
	items, err := mutation.Analyze(context.Background(), run, provider)
	if err != nil {
		t.Fatal(err)
	}
	if items != nil {
		t.Fatal("unavailable provider must return nil feedback")
	}
}

func TestBuildObservePrompt(t *testing.T) {
	cluster := mutation.SurvivorCluster{
		CaseID:    "TC-0001",
		Title:     "GET /pets",
		Operators: []string{"field_drop", "array_to_null"},
		RiskScore: 0.5,
	}
	existingAssertions := []string{"status_code gte 200", "status_code lt 300"}
	prompt := mutation.BuildObservePrompt(cluster, existingAssertions)
	if len(prompt) < 50 {
		t.Fatalf("prompt too short: %q", prompt)
	}
	if !strings.Contains(prompt, "field_drop") {
		t.Error("prompt must mention surviving operators")
	}
	if !strings.Contains(prompt, "TC-0001") {
		t.Error("prompt must mention case ID")
	}
}
