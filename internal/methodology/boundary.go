// internal/methodology/boundary.go
package methodology

import (
	"fmt"

	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

type BoundaryTechnique struct {
	gen *datagen.Generator
}

func NewBoundaryTechnique() *BoundaryTechnique {
	return &BoundaryTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *BoundaryTechnique) Name() string { return "boundary_value" }

func (t *BoundaryTechnique) Applies(op *spec.Operation) bool {
	return hasRangeConstraints(op)
}

func (t *BoundaryTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	s := getJSONSchema(op.RequestBody)
	if s == nil {
		return nil, nil
	}
	var cases []schema.TestCase
	base := t.buildValidBody(op)

	for fieldName, fieldSchema := range s.Properties {
		if !hasBoundary(fieldSchema) {
			continue
		}
		bounds := []struct {
			kind  datagen.BoundaryKind
			label string
			valid bool
		}{
			{datagen.BoundaryMin, "min_valid", true},
			{datagen.BoundaryMinMinusOne, "min_minus_one_invalid", false},
			{datagen.BoundaryMax, "max_valid", true},
			{datagen.BoundaryMaxPlusOne, "max_plus_one_invalid", false},
		}
		for _, b := range bounds {
			body := copyMap(base)
			body[fieldName] = t.gen.GenerateBoundary(fieldSchema, b.kind)
			tc := buildTestCase(op, body,
				fmt.Sprintf("%s_%s", fieldName, b.label),
				fmt.Sprintf("%s at %s boundary", fieldName, b.label),
				fmt.Sprintf("%s %s requestBody.properties.%s", op.Method, op.Path, fieldName))
			if b.valid {
				tc.Priority = "P1"
			} else {
				tc.Priority = "P2"
				tc.Steps[0].Assertions = []schema.Assertion{
					{Target: "status_code", Operator: "eq", Expected: 422},
				}
			}
			tc.Source = schema.CaseSource{
				Technique: "boundary_value",
				SpecPath:  fmt.Sprintf("%s %s requestBody.properties.%s", op.Method, op.Path, fieldName),
				Rationale: fmt.Sprintf("boundary value analysis: %s at %s", fieldName, b.label),
			}
			cases = append(cases, tc)
		}
	}
	return cases, nil
}

func (t *BoundaryTechnique) buildValidBody(op *spec.Operation) map[string]any {
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

func hasRangeConstraints(op *spec.Operation) bool {
	s := getJSONSchema(op.RequestBody)
	if s == nil {
		return false
	}
	for _, fieldSchema := range s.Properties {
		if hasBoundary(fieldSchema) {
			return true
		}
	}
	return false
}

func hasBoundary(s *spec.Schema) bool {
	return s.Minimum != nil || s.Maximum != nil || s.MinLength != nil || s.MaxLength != nil
}

func copyMap(m map[string]any) map[string]any {
	result := make(map[string]any, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}
