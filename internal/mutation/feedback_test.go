// internal/mutation/feedback_test.go
package mutation_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
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

func TestPatchIndex(t *testing.T) {
	dir := t.TempDir()
	indexData := []byte(`{"test_cases":[{"id":"TC-0001","title":"GET /pets","steps":[{"id":"step-1","assertions":[{"target":"status_code","operator":"eq","expected":200}]}]}]}`)
	if err := os.WriteFile(filepath.Join(dir, "index.json"), indexData, 0644); err != nil {
		t.Fatal(err)
	}

	items := []mutation.FeedbackItem{{
		CaseID: "TC-0001",
		SuggestedAssertions: []mutation.SuggestedAssertion{
			{Target: "jsonpath $.id", Operator: "exists"},
		},
	}}

	if err := mutation.PatchIndex(dir, items); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "index.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte("jsonpath $.id")) {
		t.Fatalf("patched index.json must contain new assertion:\n%s", data)
	}
}
