// internal/methodology/isolated_negative.go
package methodology

import (
	"fmt"
	"strings"

	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// IsolatedNegativeTechnique generates one test case per invalid field or missing
// required parameter. Each case has exactly one source of invalidity — all other
// fields/params carry valid values. This mirrors Tcases' isolated failure principle:
// test failures can be attributed to a single root cause.
type IsolatedNegativeTechnique struct {
	gen *datagen.Generator
}

func NewIsolatedNegativeTechnique() *IsolatedNegativeTechnique {
	return &IsolatedNegativeTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *IsolatedNegativeTechnique) Name() string { return "isolated_negative" }

func (t *IsolatedNegativeTechnique) Applies(op *spec.Operation) bool {
	s := getJSONSchema(op.RequestBody)
	if s != nil && (len(s.Required) > 0 || hasConstrainedField(s)) {
		return true
	}
	for _, p := range op.Parameters {
		if p.Required {
			return true
		}
	}
	return false
}

func (t *IsolatedNegativeTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	var cases []schema.TestCase

	s := getJSONSchema(op.RequestBody)
	if s != nil {
		validBase := buildValidBody(t.gen, op)
		if validBase == nil {
			validBase = map[string]any{}
		}

		// One case per required field: omit only that field
		for _, req := range s.Required {
			body := copyMap(validBase)
			delete(body, req)
			specPath := fmt.Sprintf("%s %s requestBody.properties.%s", op.Method, op.Path, req)
			tc := buildTestCase(op, body,
				fmt.Sprintf("missing required field %q", req), specPath)
			tc.Priority = "P1"
			tc.Steps[0].Assertions = []schema.Assertion{
				{Target: "status_code", Operator: "eq", Expected: 422},
			}
			tc.Source = schema.CaseSource{
				Technique: "isolated_negative",
				SpecPath:  specPath,
				Rationale: fmt.Sprintf("isolated failure: only %q is absent; all other fields valid", req),
				Scenario:  "MISSING_REQUIRED",
			}
			cases = append(cases, tc)
		}

		// One case per field with an easily-violated constraint
		for fieldName, fieldSchema := range s.Properties {
			invalid, reason := firstIsolatedInvalid(fieldSchema)
			if invalid == nil {
				continue
			}
			body := copyMap(validBase)
			body[fieldName] = invalid
			specPath := fmt.Sprintf("%s %s requestBody.properties.%s", op.Method, op.Path, fieldName)
			tc := buildTestCase(op, body,
				fmt.Sprintf("invalid %s: %s", fieldName, reason), specPath)
			tc.Priority = "P2"
			tc.Steps[0].Assertions = []schema.Assertion{
				{Target: "status_code", Operator: "eq", Expected: 422},
			}
			tc.Source = schema.CaseSource{
				Technique: "isolated_negative",
				SpecPath:  specPath,
				Rationale: fmt.Sprintf("isolated failure: only %q is invalid (%s); all other fields valid", fieldName, reason),
			}
			cases = append(cases, tc)
		}
	}

	// One case per required parameter: omit only that parameter
	for _, p := range op.Parameters {
		if !p.Required {
			continue
		}
		specPath := fmt.Sprintf("%s %s parameters.%s", op.Method, op.Path, p.Name)
		tc := buildTestCase(op, nil,
			fmt.Sprintf("missing required param %q", p.Name), specPath)
		tc.Priority = "P1"
		tc.Steps[0].Path = buildPathWithoutParam(op, p.Name)
		tc.Steps[0].Assertions = []schema.Assertion{
			{Target: "status_code", Operator: "eq", Expected: 422},
		}
		tc.Source = schema.CaseSource{
			Technique: "isolated_negative",
			SpecPath:  specPath,
			Rationale: fmt.Sprintf("isolated failure: required param %q is absent", p.Name),
			Scenario:  "MISSING_REQUIRED",
		}
		cases = append(cases, tc)
	}

	return cases, nil
}

// firstIsolatedInvalid returns the simplest invalid value for a schema and a short reason.
// Returns (nil, "") when no violation is trivially constructable.
func firstIsolatedInvalid(s *spec.Schema) (any, string) {
	if s == nil {
		return nil, ""
	}
	switch s.Type {
	case "integer", "number":
		if s.Minimum != nil {
			v := *s.Minimum - 1
			return v, fmt.Sprintf("below minimum (%.0f)", *s.Minimum)
		}
		if s.Maximum != nil {
			v := *s.Maximum + 1
			return v, fmt.Sprintf("above maximum (%.0f)", *s.Maximum)
		}
		return "not_a_number", "wrong type (string for numeric)"
	case "string":
		if s.MinLength != nil && *s.MinLength > 0 {
			return "", fmt.Sprintf("empty string violates minLength %d", *s.MinLength)
		}
		if len(s.Enum) > 0 {
			return "__invalid_enum__", "value not in enum"
		}
		if s.Format == "email" {
			return "not-an-email", "invalid email format"
		}
	case "boolean":
		return "not_a_boolean", "wrong type (string for boolean)"
	case "array":
		if s.MinItems != nil && *s.MinItems > 0 {
			return []any{}, fmt.Sprintf("empty array violates minItems %d", *s.MinItems)
		}
	}
	return nil, ""
}

// hasConstrainedField returns true if any property has a constraint that can be violated.
func hasConstrainedField(s *spec.Schema) bool {
	for _, prop := range s.Properties {
		if v, _ := firstIsolatedInvalid(prop); v != nil {
			return true
		}
	}
	return false
}

// buildPathWithoutParam builds the operation path+query omitting the named parameter.
func buildPathWithoutParam(op *spec.Operation, omit string) string {
	params := map[string]any{}
	for _, p := range op.Parameters {
		if p.Name == omit || p.In == "path" {
			continue
		}
		if p.Schema != nil {
			params[p.Name] = "valid"
		}
	}
	path := op.Path
	// Replace path params with placeholders
	for _, p := range op.Parameters {
		if p.In == "path" {
			path = strings.ReplaceAll(path, "{"+p.Name+"}", "1")
		}
	}
	if len(params) > 0 {
		return buildPathWithQuery(path, params)
	}
	return path
}
