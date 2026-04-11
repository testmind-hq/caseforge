// internal/methodology/constraint_mutation.go
package methodology

import (
	"fmt"

	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// ConstraintMutationTechnique generates negative test cases by applying
// two schemathesis-style mutations not covered by schema_violation:
//
//  1. Null injection: for each non-nullable field, sends null and expects 422.
//  2. Wrong content-type: sends a valid body with Content-Type: text/plain
//     and expects 415 Unsupported Media Type.
//
// Each case isolates a single mutation; all other fields carry valid values.
type ConstraintMutationTechnique struct {
	gen *datagen.Generator
}

func NewConstraintMutationTechnique() *ConstraintMutationTechnique {
	return &ConstraintMutationTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *ConstraintMutationTechnique) Name() string { return "constraint_mutation" }

func (t *ConstraintMutationTechnique) Applies(op *spec.Operation) bool {
	s := getJSONSchema(op.RequestBody)
	return s != nil && len(s.Properties) > 0
}

func (t *ConstraintMutationTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	s := getJSONSchema(op.RequestBody)
	if s == nil {
		return nil, nil
	}

	validBase := buildValidBody(t.gen, op)
	if validBase == nil {
		validBase = map[string]any{}
	}

	var cases []schema.TestCase

	// Mutation 1: null injection for non-nullable fields
	for fieldName, fieldSchema := range s.Properties {
		if fieldSchema != nil && fieldSchema.Nullable {
			continue // field explicitly accepts null — skip
		}
		body := copyMap(validBase)
		body[fieldName] = nil // inject null
		specPath := fmt.Sprintf("%s %s requestBody.properties.%s", op.Method, op.Path, fieldName)
		tc := buildTestCase(op, body,
			fmt.Sprintf("null injection: %s", fieldName), specPath)
		tc.Priority = "P2"
		tc.Steps[0].Assertions = []schema.Assertion{
			{Target: "status_code", Operator: "eq", Expected: 422},
		}
		tc.Source = schema.CaseSource{
			Technique: "constraint_mutation",
			SpecPath:  specPath,
			Rationale: fmt.Sprintf("field %q is non-nullable but receives null — server must reject with 422", fieldName),
			Scenario:  "NULL_INJECTION",
		}
		cases = append(cases, tc)
	}

	// Mutation 2: wrong content-type — valid JSON body sent as text/plain
	specPath := fmt.Sprintf("%s %s requestBody", op.Method, op.Path)
	body := copyMap(validBase)
	tc := buildTestCase(op, body, "wrong content-type (text/plain)", specPath)
	tc.Priority = "P2"
	// Override the content-type set by buildTestCase
	if tc.Steps[0].Headers == nil {
		tc.Steps[0].Headers = make(map[string]string)
	}
	tc.Steps[0].Headers["Content-Type"] = "text/plain"
	tc.Steps[0].Assertions = []schema.Assertion{
		{Target: "status_code", Operator: "eq", Expected: 415},
	}
	tc.Source = schema.CaseSource{
		Technique: "constraint_mutation",
		SpecPath:  specPath,
		Rationale: "valid JSON body sent with Content-Type: text/plain — server must return 415 Unsupported Media Type",
		Scenario:  "WRONG_CONTENT_TYPE",
	}
	cases = append(cases, tc)

	return cases, nil
}
