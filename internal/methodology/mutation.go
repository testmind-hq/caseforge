// internal/methodology/mutation.go
package methodology

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// MutationTechnique generates test cases by injecting invalid values into individual
// fields of the request body. Each case mutates exactly one field; all others remain valid.
// The expected response is 4xx (API must reject invalid input — never 5xx).
type MutationTechnique struct {
	gen      *datagen.Generator
	maxCases int // 0 = use default of 10
}

func NewMutationTechnique() *MutationTechnique {
	return &MutationTechnique{gen: datagen.NewGenerator(nil), maxCases: 10}
}

// NewMutationTechniqueWithMax creates a mutation technique with an explicit case limit.
func NewMutationTechniqueWithMax(max int) *MutationTechnique {
	return &MutationTechnique{gen: datagen.NewGenerator(nil), maxCases: max}
}

func (t *MutationTechnique) Name() string { return "mutation" }

// Applies when the operation has a JSON request body with at least one property.
func (t *MutationTechnique) Applies(op *spec.Operation) bool {
	s := getJSONSchema(op.RequestBody)
	return s != nil && len(s.Properties) > 0
}

func (t *MutationTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	s := getJSONSchema(op.RequestBody)
	if s == nil {
		return nil, nil
	}

	base := buildValidBody(t.gen, op)

	var cases []schema.TestCase
	for fieldName, fieldSchema := range s.Properties {
		for _, m := range mutationsForField(fieldName, fieldSchema) {
			if t.maxCases > 0 && len(cases) >= t.maxCases {
				return cases, nil
			}
			body := cloneMap(base)
			body[fieldName] = m.value
			tc := schema.TestCase{
				Schema:   schema.SchemaBaseURL,
				Version:  "1",
				ID:       fmt.Sprintf("TC-%s", uuid.New().String()[:8]),
				Title:    fmt.Sprintf("%s %s - mutation: %s %s", op.Method, op.Path, fieldName, m.desc),
				Kind:     "single",
				Priority: "P2",
				Tags:     op.Tags,
				Source: schema.CaseSource{
					Technique: "mutation",
					SpecPath:  fmt.Sprintf("%s %s requestBody.%s", op.Method, op.Path, fieldName),
					Rationale: fmt.Sprintf("field %q mutated with %s; API must reject with 4xx", fieldName, m.desc),
				},
				Steps: []schema.Step{{
					ID:      "step-main",
					Title:   fmt.Sprintf("mutation: %s → %s", fieldName, m.desc),
					Type:    "test",
					Method:  op.Method,
					Path:    op.Path,
					Headers: map[string]string{"Content-Type": "application/json"},
					Body:    body,
					Assertions: []schema.Assertion{
						{Target: "status_code", Operator: schema.OperatorGte, Expected: 400},
						{Target: "status_code", Operator: schema.OperatorLt, Expected: 500},
					},
				}},
				GeneratedAt: time.Now(),
			}
			cases = append(cases, tc)
		}
	}
	return cases, nil
}

type mutation struct {
	value any
	desc  string
}

// mutationsForField returns the set of invalid values to try for a given field.
func mutationsForField(fieldName string, s *spec.Schema) []mutation {
	ms := []mutation{
		{value: nil, desc: "null value"},
	}
	switch s.Type {
	case "string":
		ms = append(ms,
			mutation{value: "", desc: "empty string"},
			mutation{value: 12345, desc: "integer instead of string"},
			mutation{value: repeatChar('a', 300), desc: "oversized string (300 chars)"},
		)
		if s.Format == "email" {
			ms = append(ms, mutation{value: "not-an-email", desc: "invalid email format"})
		}
		if s.Format == "date" || s.Format == "date-time" {
			ms = append(ms, mutation{value: "not-a-date", desc: "invalid date format"})
		}
	case "integer", "number":
		ms = append(ms,
			mutation{value: "not-a-number", desc: "string instead of number"},
			mutation{value: -999999, desc: "extreme negative integer"},
			mutation{value: 9999999999, desc: "extreme large integer"},
		)
	case "boolean":
		ms = append(ms,
			mutation{value: "yes", desc: "string instead of boolean"},
			mutation{value: 1, desc: "integer instead of boolean"},
		)
	case "array":
		ms = append(ms,
			mutation{value: "not-an-array", desc: "string instead of array"},
			mutation{value: map[string]any{}, desc: "object instead of array"},
		)
	case "object":
		ms = append(ms,
			mutation{value: "not-an-object", desc: "string instead of object"},
			mutation{value: []any{}, desc: "array instead of object"},
		)
	}
	return ms
}

func repeatChar(c byte, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = c
	}
	return string(b)
}

func cloneMap(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
