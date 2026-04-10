// internal/methodology/schema_violation.go
package methodology

import (
	"fmt"
	"strings"

	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// SchemaViolationTechnique generates one test case per schema constraint violation.
// Unlike boundary_value (numeric/string ranges only) this covers every OpenAPI
// schema constraint type: required field absence, type mismatches, range violations,
// string length violations, enum violations, format violations, array size violations.
// Each case isolates a single constraint violation; all other fields carry valid values.
type SchemaViolationTechnique struct {
	gen *datagen.Generator
}

func NewSchemaViolationTechnique() *SchemaViolationTechnique {
	return &SchemaViolationTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *SchemaViolationTechnique) Name() string { return "schema_violation" }

func (t *SchemaViolationTechnique) Applies(op *spec.Operation) bool {
	s := getJSONSchema(op.RequestBody)
	return s != nil && len(s.Properties) > 0
}

type svViolation struct {
	field     string
	value     any    // nil means "omit the field"
	label     string
	rationale string
}

func (t *SchemaViolationTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	s := getJSONSchema(op.RequestBody)
	if s == nil {
		return nil, nil
	}

	validBase := buildValidBody(t.gen, op)
	if validBase == nil {
		validBase = map[string]any{}
	}

	var violations []svViolation

	// Required field absence — one violation per required field
	for _, req := range s.Required {
		violations = append(violations, svViolation{
			field:     req,
			value:     nil,
			label:     fmt.Sprintf("%s_missing_required", req),
			rationale: fmt.Sprintf("required field %q is absent", req),
		})
	}

	// Per-field constraint violations
	for fieldName, fieldSchema := range s.Properties {
		violations = append(violations, extractSchemaViolations(fieldName, fieldSchema)...)
	}

	var cases []schema.TestCase
	for _, v := range violations {
		body := copyMap(validBase)
		if v.value == nil {
			delete(body, v.field)
		} else {
			body[v.field] = v.value
		}
		specPath := fmt.Sprintf("%s %s requestBody.properties.%s", op.Method, op.Path, v.field)
		tc := buildTestCase(op, body,
			fmt.Sprintf("[schema_violation] %s", v.label), specPath)
		tc.Priority = "P2"
		tc.Steps[0].Assertions = []schema.Assertion{
			{Target: "status_code", Operator: "eq", Expected: 422},
		}
		tc.Source = schema.CaseSource{
			Technique: "schema_violation",
			SpecPath:  specPath,
			Rationale: v.rationale,
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

// extractSchemaViolations returns one svViolation per constraint on the given field.
func extractSchemaViolations(fieldName string, s *spec.Schema) []svViolation {
	var vs []svViolation
	if s == nil {
		return vs
	}
	switch s.Type {
	case "integer", "number":
		if s.Minimum != nil {
			vs = append(vs, svViolation{
				field:     fieldName,
				value:     *s.Minimum - 1,
				label:     fmt.Sprintf("%s_below_minimum", fieldName),
				rationale: fmt.Sprintf("%s=%.0f violates minimum %.0f", fieldName, *s.Minimum-1, *s.Minimum),
			})
		}
		if s.Maximum != nil {
			vs = append(vs, svViolation{
				field:     fieldName,
				value:     *s.Maximum + 1,
				label:     fmt.Sprintf("%s_above_maximum", fieldName),
				rationale: fmt.Sprintf("%s=%.0f violates maximum %.0f", fieldName, *s.Maximum+1, *s.Maximum),
			})
		}
		vs = append(vs, svViolation{
			field:     fieldName,
			value:     "not_a_number",
			label:     fmt.Sprintf("%s_wrong_type", fieldName),
			rationale: fmt.Sprintf("%s is %s but received a string", fieldName, s.Type),
		})
	case "string":
		if s.MinLength != nil && *s.MinLength > 0 {
			vs = append(vs, svViolation{
				field:     fieldName,
				value:     "",
				label:     fmt.Sprintf("%s_too_short", fieldName),
				rationale: fmt.Sprintf("%s is empty, violates minLength %d", fieldName, *s.MinLength),
			})
		}
		if s.MaxLength != nil {
			vs = append(vs, svViolation{
				field:     fieldName,
				value:     strings.Repeat("x", int(*s.MaxLength)+1),
				label:     fmt.Sprintf("%s_too_long", fieldName),
				rationale: fmt.Sprintf("%s length %d violates maxLength %d", fieldName, *s.MaxLength+1, *s.MaxLength),
			})
		}
		if len(s.Enum) > 0 {
			vs = append(vs, svViolation{
				field:     fieldName,
				value:     "__invalid__",
				label:     fmt.Sprintf("%s_invalid_enum", fieldName),
				rationale: fmt.Sprintf("%s=\"__invalid__\" is not in enum %v", fieldName, s.Enum),
			})
		}
		if s.Format == "email" {
			vs = append(vs, svViolation{
				field:     fieldName,
				value:     "not-an-email",
				label:     fmt.Sprintf("%s_invalid_format_email", fieldName),
				rationale: fmt.Sprintf("%s=\"not-an-email\" violates format \"email\"", fieldName),
			})
		}
		if s.Format == "date" || s.Format == "date-time" {
			vs = append(vs, svViolation{
				field:     fieldName,
				value:     "not-a-date",
				label:     fmt.Sprintf("%s_invalid_format_%s", fieldName, s.Format),
				rationale: fmt.Sprintf("%s=\"not-a-date\" violates format %q", fieldName, s.Format),
			})
		}
	case "boolean":
		vs = append(vs, svViolation{
			field:     fieldName,
			value:     "not_a_boolean",
			label:     fmt.Sprintf("%s_wrong_type", fieldName),
			rationale: fmt.Sprintf("%s is boolean but received a string", fieldName),
		})
	case "array":
		if s.MinItems != nil && *s.MinItems > 0 {
			vs = append(vs, svViolation{
				field:     fieldName,
				value:     []any{},
				label:     fmt.Sprintf("%s_too_few_items", fieldName),
				rationale: fmt.Sprintf("%s=[] violates minItems %d", fieldName, *s.MinItems),
			})
		}
		if s.MaxItems != nil {
			over := make([]any, int(*s.MaxItems)+1)
			for i := range over {
				over[i] = i
			}
			vs = append(vs, svViolation{
				field:     fieldName,
				value:     over,
				label:     fmt.Sprintf("%s_too_many_items", fieldName),
				rationale: fmt.Sprintf("%s has %d items, violates maxItems %d", fieldName, *s.MaxItems+1, *s.MaxItems),
			})
		}
	}
	return vs
}
