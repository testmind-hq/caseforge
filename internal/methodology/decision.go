// internal/methodology/decision.go
package methodology

import (
	"fmt"

	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

type DecisionTechnique struct {
	gen *datagen.Generator
}

func NewDecisionTechnique() *DecisionTechnique {
	return &DecisionTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *DecisionTechnique) Name() string { return "decision_table" }

// Applies when there are enum fields or boolean conditions (multi-condition logic).
func (t *DecisionTechnique) Applies(op *spec.Operation) bool {
	s := getJSONSchema(op.RequestBody)
	if s == nil {
		return false
	}
	enumCount := 0
	for _, fieldSchema := range s.Properties {
		if len(fieldSchema.Enum) > 0 || fieldSchema.Type == "boolean" {
			enumCount++
		}
	}
	return enumCount >= 2
}

func (t *DecisionTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	s := getJSONSchema(op.RequestBody)
	if s == nil {
		return nil, nil
	}
	var cases []schema.TestCase
	// Generate one test case per enum value for each enum field, and two cases (true/false) for boolean fields
	for fieldName, fieldSchema := range s.Properties {
		if len(fieldSchema.Enum) > 0 {
			for _, enumVal := range fieldSchema.Enum {
				body := buildValidBody(t.gen, op)
				if body == nil {
					body = map[string]any{}
				}
				body[fieldName] = enumVal
				tc := buildTestCase(op, body,
					fmt.Sprintf("%s = %v", fieldName, enumVal),
					fmt.Sprintf("%s %s requestBody.properties.%s", op.Method, op.Path, fieldName))
				tc.Priority = "P1"
				tc.Source = schema.CaseSource{
					Technique: "decision_table",
					SpecPath:  fmt.Sprintf("%s %s requestBody.properties.%s", op.Method, op.Path, fieldName),
					Rationale: fmt.Sprintf("decision table: %s takes enum value %v", fieldName, enumVal),
				}
				cases = append(cases, tc)
			}
		} else if fieldSchema.Type == "boolean" {
			for _, val := range []any{true, false} {
				body := buildValidBody(t.gen, op)
				if body == nil {
					body = map[string]any{}
				}
				body[fieldName] = val
				tc := buildTestCase(op, body,
					fmt.Sprintf("%s = %v", fieldName, val),
					fmt.Sprintf("%s %s requestBody.properties.%s", op.Method, op.Path, fieldName))
				tc.Priority = "P1"
				tc.Source = schema.CaseSource{
					Technique: "decision_table",
					SpecPath:  fmt.Sprintf("%s %s requestBody.properties.%s", op.Method, op.Path, fieldName),
					Rationale: fmt.Sprintf("decision table: %s takes boolean value %v", fieldName, val),
				}
				cases = append(cases, tc)
			}
		}
	}
	return cases, nil
}

