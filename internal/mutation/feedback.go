// internal/mutation/feedback.go
package mutation

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/output/render"
	"github.com/testmind-hq/caseforge/internal/output/writer"
)

// Analyze runs LLM OC-prompting on survivor clusters and returns FeedbackItems.
// Returns nil, nil if provider is unavailable or there are no clusters.
func Analyze(ctx context.Context, run MutationRun, provider llm.LLMProvider) ([]FeedbackItem, error) {
	if provider == nil || !provider.IsAvailable() {
		return nil, nil
	}
	if len(run.Clusters) == 0 {
		return nil, nil
	}

	var items []FeedbackItem
	for _, cluster := range run.Clusters {
		item, err := analyzeCluster(ctx, cluster, provider)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: feedback for %s: %v\n", cluster.CaseID, err)
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func analyzeCluster(ctx context.Context, cluster SurvivorCluster, provider llm.LLMProvider) (FeedbackItem, error) {
	observePrompt := BuildObservePrompt(cluster, nil)

	observeResp, err := llm.Retry(ctx, 3, func() (*llm.CompletionResponse, error) {
		return provider.Complete(ctx, &llm.CompletionRequest{
			System:    "You are an API testing expert. Return only valid JSON arrays.",
			Messages:  []llm.Message{{Role: "user", Content: observePrompt}},
			MaxTokens: 512,
		})
	})
	if err != nil {
		return FeedbackItem{}, fmt.Errorf("observe LLM call: %w", err)
	}

	suggestions := parseSuggestedAssertions(observeResp.Text)

	return FeedbackItem{
		CaseID:              cluster.CaseID,
		Title:               cluster.Title,
		RiskScore:           cluster.RiskScore,
		Diagnosis:           diagnosisFromCluster(cluster),
		SuggestedAssertions: suggestions,
	}, nil
}

// BuildObservePrompt constructs the OC observation prompt for a cluster.
// Exported for testing.
func BuildObservePrompt(cluster SurvivorCluster, existingAssertions []string) string {
	existing := strings.Join(existingAssertions, "\n- ")
	if existing == "" {
		existing = "(none)"
	}
	return fmt.Sprintf(
		"Test case [%s] \"%s\" survived the following HTTP response mutations:\n%s\n\n"+
			"Existing assertions that did NOT catch these mutations:\n- %s\n\n"+
			"Return a JSON array of additional assertions that would detect these mutations.\n"+
			"Each assertion must have: \"target\" (e.g. \"jsonpath $.field\", \"status_code\", \"header Name\"),\n"+
			"\"operator\" (exists|eq|ne|lt|gt|gte|lte|contains|matches|is_iso8601|is_uuid), and optionally \"expected\".\n"+
			"Example: [{\"target\":\"jsonpath $.id\",\"operator\":\"exists\"}]\n"+
			"Return ONLY valid JSON. No explanation.",
		cluster.CaseID,
		cluster.Title,
		"- "+strings.Join(cluster.Operators, "\n- "),
		existing,
	)
}

func parseSuggestedAssertions(text string) []SuggestedAssertion {
	extracted := llm.ExtractJSON(text)
	var raw []struct {
		Target   string `json:"target"`
		Operator string `json:"operator"`
		Expected any    `json:"expected,omitempty"`
	}
	if err := json.Unmarshal([]byte(extracted), &raw); err != nil {
		return nil
	}
	out := make([]SuggestedAssertion, len(raw))
	for i, r := range raw {
		out[i] = SuggestedAssertion{Target: r.Target, Operator: r.Operator, Expected: r.Expected}
	}
	return out
}

func diagnosisFromCluster(cluster SurvivorCluster) string {
	return fmt.Sprintf("%d operator(s) survived (%s) — assertions do not cover response mutations",
		len(cluster.Operators), strings.Join(cluster.Operators, ", "))
}

// PatchIndex appends suggested assertions from FeedbackItems to matching test cases
// in index.json and writes the updated file back.
func PatchIndex(casesDir string, items []FeedbackItem) error {
	indexPath := filepath.Join(casesDir, "index.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return fmt.Errorf("reading index.json: %w", err)
	}

	var idx struct {
		TestCases []json.RawMessage `json:"test_cases"`
	}
	if err := json.Unmarshal(data, &idx); err != nil {
		return fmt.Errorf("parsing index.json: %w", err)
	}

	patches := map[string][]SuggestedAssertion{}
	for _, item := range items {
		patches[item.CaseID] = item.SuggestedAssertions
	}

	updated := make([]json.RawMessage, len(idx.TestCases))
	for i, raw := range idx.TestCases {
		var tc map[string]json.RawMessage
		if err := json.Unmarshal(raw, &tc); err != nil {
			updated[i] = raw
			continue
		}
		var id string
		if idRaw, ok := tc["id"]; !ok {
			updated[i] = raw
			continue
		} else {
			_ = json.Unmarshal(idRaw, &id)
		}
		newAssertions, hasPatch := patches[id]
		if !hasPatch {
			updated[i] = raw
			continue
		}

		var steps []json.RawMessage
		if err := json.Unmarshal(tc["steps"], &steps); err != nil || len(steps) == 0 {
			updated[i] = raw
			continue
		}
		var step map[string]json.RawMessage
		if err := json.Unmarshal(steps[0], &step); err != nil {
			updated[i] = raw
			continue
		}
		var assertions []json.RawMessage
		_ = json.Unmarshal(step["assertions"], &assertions)
		for _, a := range newAssertions {
			ab, _ := json.Marshal(a)
			assertions = append(assertions, json.RawMessage(ab))
		}
		ab, _ := json.Marshal(assertions)
		step["assertions"] = json.RawMessage(ab)
		sb, _ := json.Marshal(step)
		steps[0] = json.RawMessage(sb)
		stepsBytes, _ := json.Marshal(steps)
		tc["steps"] = json.RawMessage(stepsBytes)
		tcBytes, _ := json.Marshal(tc)
		updated[i] = json.RawMessage(tcBytes)
	}

	idx.TestCases = updated
	out, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(indexPath, out, 0644); err != nil {
		return err
	}

	// Re-render .hurl files so patched assertions take effect
	cases, err := writer.NewJSONSchemaWriter().Read(indexPath)
	if err != nil {
		return fmt.Errorf("reading patched index.json: %w", err)
	}
	if err := render.NewHurlRenderer("").Render(cases, casesDir); err != nil {
		return fmt.Errorf("index.json patched but .hurl re-render failed: %w", err)
	}
	return nil
}
