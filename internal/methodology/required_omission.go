// internal/methodology/required_omission.go
package methodology

import (
	"fmt"
	"maps"
	"slices"

	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// RequiredOmissionTechnique generates one test case per required field in the
// JSON request body, where that field is completely absent (not null) from the
// payload. All other fields are present with valid values.
type RequiredOmissionTechnique struct {
	gen *datagen.Generator
}

func NewRequiredOmissionTechnique() *RequiredOmissionTechnique {
	return &RequiredOmissionTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *RequiredOmissionTechnique) Name() string { return "required_omission" }

func (t *RequiredOmissionTechnique) Applies(op *spec.Operation) bool {
	s := getJSONSchema(op.RequestBody)
	return s != nil && len(s.Required) > 0
}

func (t *RequiredOmissionTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	s := getJSONSchema(op.RequestBody)
	if s == nil || len(s.Required) == 0 {
		return nil, nil
	}

	// Iterate in sorted order for determinism
	requiredFields := slices.Sorted(slices.Values(s.Required))

	var cases []schema.TestCase
	for _, fieldName := range requiredFields {
		body := bestExampleForField(op, fieldName, t.gen)
		// Delete the field entirely — not null, actually absent
		delete(body, fieldName)

		specPath := fmt.Sprintf("%s %s requestBody.%s", op.Method, op.Path, fieldName)
		title := fmt.Sprintf("[required_omission] %s absent", fieldName)
		tc := buildTestCase(op, body, title, specPath)
		tc.Priority = "P2"
		tc.Steps[0].Assertions = []schema.Assertion{
			{Target: "status_code", Operator: "gte", Expected: 400},
			{Target: "status_code", Operator: "lt", Expected: 500},
		}
		tc.Source = schema.CaseSource{
			Technique: "required_omission",
			SpecPath:  specPath,
			Rationale: fmt.Sprintf("required field %q omitted entirely (not null) — server must reject with 4xx", fieldName),
			Scenario:  "REQUIRED_OMISSION",
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

// bestExampleForField picks the best seed body for omitting fieldName.
// It prefers a named example that contains the field, then inline example,
// then falls back to buildValidBody.
func bestExampleForField(op *spec.Operation, fieldName string, gen *datagen.Generator) map[string]any {
	mt := jsonMediaType(op)
	if mt != nil {
		// Try named examples in sorted order for determinism
		names := slices.Sorted(maps.Keys(mt.Examples))
		for _, name := range names {
			ex := mt.Examples[name]
			if ex == nil || ex.Value == nil {
				continue
			}
			if m, ok := ex.Value.(map[string]any); ok {
				if _, has := m[fieldName]; has {
					return copyMap(m)
				}
			}
		}
		// Try inline example
		if mt.Example != nil {
			if m, ok := mt.Example.(map[string]any); ok {
				if _, has := m[fieldName]; has {
					return copyMap(m)
				}
			}
		}
	}
	// Fallback: generate valid body
	body := buildValidBody(gen, op)
	if body == nil {
		body = map[string]any{}
	}
	return body
}

