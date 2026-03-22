// internal/methodology/equivalence.go
package methodology

import (
	"fmt"

	"github.com/testmind-hq/caseforge/internal/assert"
	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

type EquivalenceTechnique struct {
	gen *datagen.Generator
}

func NewEquivalenceTechnique() *EquivalenceTechnique {
	return &EquivalenceTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *EquivalenceTechnique) Name() string { return "equivalence_partitioning" }

// Applies returns true for all operations — equivalence partitioning is always applicable.
func (t *EquivalenceTechnique) Applies(_ *spec.Operation) bool { return true }

func (t *EquivalenceTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	var cases []schema.TestCase

	// 1. Happy path: all required fields with valid data
	happyBody := buildValidBody(t.gen, op)
	happyCase := buildTestCase(op, happyBody,
		"valid request with all required fields",
		fmt.Sprintf("%s %s", op.Method, op.Path))
	happyCase.Priority = "P0"
	happyCase.Source = schema.CaseSource{
		Technique: "equivalence_partitioning",
		SpecPath:  fmt.Sprintf("%s %s", op.Method, op.Path),
		Rationale: "valid equivalence class: all required fields present with correct types",
	}
	happyCase.Steps[0].Assertions = append(
		assert.BasicAssertions(op),
		responseSchemaAssertions(op)...,
	)
	cases = append(cases, happyCase)

	// 2. Missing required field cases
	if op.RequestBody != nil {
		reqSchema := getJSONSchema(op.RequestBody)
		if reqSchema != nil {
			for _, requiredField := range reqSchema.Required {
				body := t.buildBodyWithout(op, requiredField)
				tc := buildTestCase(op, body,
					fmt.Sprintf("missing required field %q", requiredField),
					fmt.Sprintf("%s %s requestBody.%s", op.Method, op.Path, requiredField))
				tc.Priority = "P1"
				tc.Source = schema.CaseSource{
					Technique: "equivalence_partitioning",
					SpecPath:  fmt.Sprintf("%s %s requestBody.properties.%s", op.Method, op.Path, requiredField),
					Rationale: fmt.Sprintf("invalid equivalence class: required field %q is absent", requiredField),
				}
				tc.Steps[0].Assertions = []schema.Assertion{
					{Target: "status_code", Operator: "eq", Expected: 422},
				}
				cases = append(cases, tc)
			}
		}
	}

	return cases, nil
}

func (t *EquivalenceTechnique) buildBodyWithout(op *spec.Operation, excludeField string) map[string]any {
	body := buildValidBody(t.gen, op)
	if body != nil {
		delete(body, excludeField)
	}
	return body
}
