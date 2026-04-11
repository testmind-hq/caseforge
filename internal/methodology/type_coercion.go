// internal/methodology/type_coercion.go
package methodology

import (
	"fmt"

	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// TypeCoercionTechnique generates test cases that send a value of the wrong JSON
// type for each typed field in the request body.
//
// Mutations by field type:
//   - string  → integer 123 and boolean true (2 cases)
//   - integer/number → string "not_a_number" and boolean false (2 cases)
//   - boolean → string "not_a_boolean" and integer 1 (2 cases)
//   - array   → string "not_an_array" (1 case)
//   - object  → string "not_an_object" (1 case)
//
// Fields with no Type set (type-polymorphic) are skipped.
// Fields where the injected value would match an allowed enum value are skipped.
// All cases: Priority P2, expected status 422, Scenario "WRONG_TYPE".
type TypeCoercionTechnique struct {
	gen *datagen.Generator
}

// NewTypeCoercionTechnique returns a new TypeCoercionTechnique.
func NewTypeCoercionTechnique() *TypeCoercionTechnique {
	return &TypeCoercionTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *TypeCoercionTechnique) Name() string { return "type_coercion" }

// Applies returns true if the operation has a JSON request body with at least
// one property that has a declared Type.
func (t *TypeCoercionTechnique) Applies(op *spec.Operation) bool {
	s := getJSONSchema(op.RequestBody)
	if s == nil {
		return false
	}
	for _, fieldSchema := range s.Properties {
		if fieldSchema != nil && fieldSchema.Type != "" {
			return true
		}
	}
	return false
}

// wrongTypeMutation describes a single type-coercion mutation.
type wrongTypeMutation struct {
	value any
	label string
}

// Generate produces one test case per wrong-type mutation per field.
func (t *TypeCoercionTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	s := getJSONSchema(op.RequestBody)
	if s == nil {
		return nil, nil
	}

	validBase := buildValidBody(t.gen, op)
	if validBase == nil {
		validBase = map[string]any{}
	}

	var cases []schema.TestCase

	for fieldName, fieldSchema := range s.Properties {
		if fieldSchema == nil || fieldSchema.Type == "" {
			continue // skip type-polymorphic fields
		}

		mutations := wrongTypeMutationsFor(fieldSchema.Type)

		for _, m := range mutations {
			// Skip if the injected value happens to be in the allowed enum.
			if isInEnum(m.value, fieldSchema.Enum) {
				continue
			}

			body := copyMap(validBase)
			body[fieldName] = m.value

			specPath := fmt.Sprintf("%s %s requestBody.properties.%s", op.Method, op.Path, fieldName)
			tc := buildTestCase(op, body,
				fmt.Sprintf("[type_coercion] %s %s", fieldName, m.label), specPath)
			tc.Priority = "P2"
			tc.Steps[0].Assertions = []schema.Assertion{
				{Target: "status_code", Operator: "eq", Expected: 422},
			}
			tc.Source = schema.CaseSource{
				Technique: "type_coercion",
				SpecPath:  specPath,
				Rationale: fmt.Sprintf("field %q is %s but receives %s — server must reject with 422", fieldName, fieldSchema.Type, m.label),
				Scenario:  "WRONG_TYPE",
			}
			cases = append(cases, tc)
		}
	}

	return cases, nil
}

// wrongTypeMutationsFor returns the wrong-type mutations for the given field type.
func wrongTypeMutationsFor(fieldType string) []wrongTypeMutation {
	switch fieldType {
	case "string":
		return []wrongTypeMutation{
			{value: 123, label: "wrong_type_integer"},
			{value: true, label: "wrong_type_boolean"},
		}
	case "integer", "number":
		return []wrongTypeMutation{
			{value: "not_a_number", label: "wrong_type_string"},
			{value: false, label: "wrong_type_boolean"},
		}
	case "boolean":
		return []wrongTypeMutation{
			{value: "not_a_boolean", label: "wrong_type_string"},
			{value: 1, label: "wrong_type_integer"},
		}
	case "array":
		return []wrongTypeMutation{
			{value: "not_an_array", label: "wrong_type_string"},
		}
	case "object":
		return []wrongTypeMutation{
			{value: "not_an_object", label: "wrong_type_string"},
		}
	}
	return nil
}

// isInEnum returns true if value is present in the enum slice.
func isInEnum(value any, enum []any) bool {
	for _, e := range enum {
		if e == value {
			return true
		}
	}
	return false
}
