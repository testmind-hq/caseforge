package methodology

import (
	"fmt"
	"strings"

	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/score"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// BusinessRuleTechnique generates negative test cases for each LLM-inferred
// implicit business rule in op.SemanticInfo.ImplicitRules.
type BusinessRuleTechnique struct {
	gen *datagen.Generator
}

func NewBusinessRuleTechnique() *BusinessRuleTechnique {
	return &BusinessRuleTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *BusinessRuleTechnique) Name() string { return "business_rule_violation" }

func (t *BusinessRuleTechnique) Applies(op *spec.Operation) bool {
	return op.SemanticInfo != nil && len(op.SemanticInfo.ImplicitRules) > 0
}

func (t *BusinessRuleTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	if op.SemanticInfo == nil {
		return nil, nil
	}
	var cases []schema.TestCase
	for _, rule := range op.SemanticInfo.ImplicitRules {
		body := buildViolationBody(t.gen, op, rule)
		if body == nil {
			body = map[string]any{}
		}
		specPath := fmt.Sprintf("%s %s", op.Method, op.Path)
		title := fmt.Sprintf("[business_rule_violation] violates: %s", rule)
		tc := buildTestCase(op, body, title, specPath)
		tc.Priority = "P2"
		tc.Steps[0].Assertions = []schema.Assertion{
			{Target: "status_code", Operator: "gte", Expected: 400},
			{Target: "status_code", Operator: "lt", Expected: 500},
		}
		tc.Source = schema.CaseSource{
			Technique: "business_rule_violation",
			SpecPath:  specPath,
			Rationale: fmt.Sprintf("LLM-inferred business rule violated: %q", rule),
			Scenario:  score.ScenarioBusinessRuleViolation,
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

func buildViolationBody(gen *datagen.Generator, op *spec.Operation, rule string) map[string]any {
	body := buildValidBody(gen, op)
	if body == nil {
		body = map[string]any{}
	}
	s := getJSONSchema(op.RequestBody)
	if s == nil {
		return body
	}
	ruleLower := strings.ToLower(rule)
	for fieldName := range s.Properties {
		if !strings.Contains(ruleLower, strings.ToLower(fieldName)) {
			continue
		}
		switch {
		case strings.Contains(ruleLower, "unique") || strings.Contains(ruleLower, "duplicate"):
			body[fieldName] = fmt.Sprintf("duplicate-%s@example.com", fieldName)
		case strings.Contains(ruleLower, "space") || strings.Contains(ruleLower, "whitespace"):
			body[fieldName] = "invalid value with spaces"
		case strings.Contains(ruleLower, "negative") || strings.Contains(ruleLower, "positive"):
			body[fieldName] = -1
		case strings.Contains(ruleLower, "future") || strings.Contains(ruleLower, "past"):
			body[fieldName] = "1970-01-01T00:00:00Z"
		default:
			body[fieldName] = ""
		}
		break
	}
	return body
}
