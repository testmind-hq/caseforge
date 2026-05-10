// internal/mutation/feedback.go
package mutation

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/testmind-hq/caseforge/internal/llm"
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
