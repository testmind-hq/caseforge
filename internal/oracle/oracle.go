package oracle

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// Constraint is a single verifiable response body property.
type Constraint struct {
	Type   string `json:"type"`   // "format"|"range"|"exists"|"pattern"|"enum"
	Field  string `json:"field"`
	Detail string `json:"detail"`
}

// ToAssertions converts the constraint to CaseForge schema assertions.
func (c Constraint) ToAssertions() []schema.Assertion {
	target := fmt.Sprintf("jsonpath $.%s", c.Field)
	switch c.Type {
	case "exists":
		return []schema.Assertion{{Target: target, Operator: schema.OperatorExists}}
	case "format":
		lower := strings.ToLower(c.Detail)
		switch {
		case strings.Contains(lower, "uuid"):
			return []schema.Assertion{{Target: target, Operator: schema.OperatorExists}}
		case strings.Contains(lower, "iso 8601") || strings.Contains(lower, "datetime") || strings.Contains(lower, "timestamp"):
			return []schema.Assertion{{Target: target, Operator: schema.OperatorExists}}
		default:
			return []schema.Assertion{{Target: target, Operator: schema.OperatorExists}}
		}
	case "range", "pattern", "enum":
		return []schema.Assertion{{Target: target, Operator: schema.OperatorExists}}
	}
	return nil
}

// Mine uses two-step Observation-Confirmation (OC) prompting to extract
// response body constraints from the operation's spec description.
// Returns nil if the provider is unavailable.
func Mine(ctx context.Context, op *spec.Operation, provider llm.LLMProvider) ([]Constraint, error) {
	if !provider.IsAvailable() {
		return nil, nil
	}

	respSchema := responseSchemaText(op)

	observePrompt := fmt.Sprintf(
		"Operation: %s %s\nSummary: %s\nDescription: %s\nResponse schema: %s\n\n"+
			"Identify response body constraints. Return JSON array with fields: type, field, detail.\n"+
			"Types: format, range, exists, pattern, enum.\n"+
			"Example: [{\"type\":\"format\",\"field\":\"email\",\"detail\":\"must be valid email\"}]\n"+
			"Respond with ONLY valid JSON array.",
		op.Method, op.Path, op.Summary, op.Description, respSchema,
	)

	observeResp, err := llm.Retry(ctx, 3, func() (*llm.CompletionResponse, error) {
		return provider.Complete(ctx, &llm.CompletionRequest{
			System:    "You are an API testing expert. Return only valid JSON.",
			Messages:  []llm.Message{{Role: "user", Content: observePrompt}},
			MaxTokens: 1024,
		})
	})
	if err != nil {
		return nil, fmt.Errorf("oracle observe: %w", err)
	}

	candidates := parseConstraints(observeResp.Text)
	if len(candidates) == 0 {
		return nil, nil
	}

	candidatesJSON, _ := json.Marshal(candidates)
	confirmPrompt := fmt.Sprintf(
		"Operation: %s %s\nDescription: %s\nResponse schema: %s\n\n"+
			"Candidate constraints:\n%s\n\n"+
			"Add \"confirmed\":true|false to each. confirmed=true only if explicitly stated or clearly implied.\n"+
			"Return the same JSON array with confirmed fields. Only valid JSON.",
		op.Method, op.Path, op.Description, respSchema, string(candidatesJSON),
	)

	confirmResp, err := llm.Retry(ctx, 3, func() (*llm.CompletionResponse, error) {
		return provider.Complete(ctx, &llm.CompletionRequest{
			System:    "You are an API testing expert. Return only valid JSON.",
			Messages:  []llm.Message{{Role: "user", Content: confirmPrompt}},
			MaxTokens: 1024,
		})
	})
	if err != nil {
		return candidates, nil // degrade gracefully
	}

	return parseConfirmedConstraints(confirmResp.Text), nil
}

// InjectIntoCase adds oracle assertions to the first step of tc for 2xx cases only.
func InjectIntoCase(tc schema.TestCase, constraints []Constraint) schema.TestCase {
	if len(constraints) == 0 || len(tc.Steps) == 0 {
		return tc
	}
	has2xx := false
	for _, a := range tc.Steps[0].Assertions {
		if a.Target == "status_code" {
			n, ok := oracleToInt(a.Expected)
			if !ok {
				continue
			}
			if (a.Operator == "lt" && n == 300) || (a.Operator == "gte" && n == 200) || (a.Operator == "eq" && n >= 200 && n < 300) {
				has2xx = true
			}
		}
	}
	if !has2xx {
		return tc
	}
	for _, c := range constraints {
		tc.Steps[0].Assertions = append(tc.Steps[0].Assertions, c.ToAssertions()...)
	}
	return tc
}

func oracleToInt(v any) (int, bool) {
	switch x := v.(type) {
	case int:
		return x, true
	case int64:
		return int(x), true
	case float64:
		return int(x), true
	}
	return 0, false
}

func parseConstraints(text string) []Constraint {
	extracted := llm.ExtractJSON(text)
	var constraints []Constraint
	_ = json.Unmarshal([]byte(extracted), &constraints)
	return constraints
}

type confirmedConstraint struct {
	Constraint
	Confirmed bool `json:"confirmed"`
}

func parseConfirmedConstraints(text string) []Constraint {
	extracted := llm.ExtractJSON(text)
	var confirmed []confirmedConstraint
	if err := json.Unmarshal([]byte(extracted), &confirmed); err != nil {
		return nil
	}
	var result []Constraint
	for _, c := range confirmed {
		if c.Confirmed {
			result = append(result, c.Constraint)
		}
	}
	return result
}

func responseSchemaText(op *spec.Operation) string {
	for code, resp := range op.Responses {
		n := 0
		fmt.Sscanf(code, "%d", &n)
		if n < 200 || n >= 300 {
			continue
		}
		if mt, ok := resp.Content["application/json"]; ok && mt != nil && mt.Schema != nil {
			b, _ := json.Marshal(mt.Schema)
			return string(b)
		}
	}
	return "{}"
}
