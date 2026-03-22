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
	// Generate one test case per enum value for each enum field
	for fieldName, fieldSchema := range s.Properties {
		if len(fieldSchema.Enum) == 0 {
			continue
		}
		for i, enumVal := range fieldSchema.Enum {
			body := t.buildValidBody(op)
			body[fieldName] = enumVal
			tc := buildTestCase(op, body,
				fmt.Sprintf("%s_enum_%d", fieldName, i),
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
	}
	return cases, nil
}

func (t *DecisionTechnique) buildValidBody(op *spec.Operation) map[string]any {
	s := getJSONSchema(op.RequestBody)
	if s == nil {
		return map[string]any{}
	}
	body := map[string]any{}
	for name, fieldSchema := range s.Properties {
		body[name] = t.gen.Generate(fieldSchema, name)
	}
	return body
}
