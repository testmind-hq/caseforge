// internal/methodology/semantic_annotation.go
package methodology

import (
	"fmt"
	"maps"
	"strings"

	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/score"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// SemanticAnnotationTechnique generates test cases that validate OpenAPI semantic
// field annotations (nullable, readOnly, writeOnly).
//
// Category A — Nullable acceptance: for nullable fields, send null and expect 2xx.
// Category B — ReadOnly write rejection: for readOnly fields sent in POST/PUT/PATCH, expect 4xx.
// Category C — WriteOnly read suppression: for writeOnly fields in GET response schema,
//              assert the field is absent from the response body.
type SemanticAnnotationTechnique struct {
	gen *datagen.Generator
}

func NewSemanticAnnotationTechnique() *SemanticAnnotationTechnique {
	return &SemanticAnnotationTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *SemanticAnnotationTechnique) Name() string { return "semantic_annotation" }

func (t *SemanticAnnotationTechnique) Applies(op *spec.Operation) bool {
	// Category A: nullable fields in request body
	if s := getJSONSchema(op.RequestBody); s != nil {
		for _, prop := range s.Properties {
			if prop.Nullable {
				return true
			}
		}
	}

	// Category B: readOnly fields in POST/PUT/PATCH request body
	method := strings.ToUpper(op.Method)
	if method == "POST" || method == "PUT" || method == "PATCH" {
		if s := getJSONSchema(op.RequestBody); s != nil {
			for _, prop := range s.Properties {
				if prop.ReadOnly {
					return true
				}
			}
		}
	}

	// Category C: writeOnly fields in GET 2xx response schema
	if strings.ToUpper(op.Method) == "GET" {
		for code, resp := range op.Responses {
			n := 0
			fmt.Sscanf(code, "%d", &n)
			if n >= 200 && n < 300 {
				if mt, ok := resp.Content["application/json"]; ok && mt.Schema != nil {
					for _, prop := range mt.Schema.Properties {
						if prop.WriteOnly {
							return true
						}
					}
				}
			}
		}
	}

	return false
}

func (t *SemanticAnnotationTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	var cases []schema.TestCase

	// Category A: Nullable acceptance
	cases = append(cases, t.generateNullableCases(op)...)

	// Category B: ReadOnly write rejection
	cases = append(cases, t.generateReadOnlyCases(op)...)

	// Category C: WriteOnly read suppression
	cases = append(cases, t.generateWriteOnlyCases(op)...)

	return cases, nil
}

// generateNullableCases creates one test case per nullable field in the JSON request body.
// Each case sends null for that field in an otherwise valid body and expects 2xx.
func (t *SemanticAnnotationTechnique) generateNullableCases(op *spec.Operation) []schema.TestCase {
	s := getJSONSchema(op.RequestBody)
	if s == nil {
		return nil
	}

	validBase := buildValidBody(t.gen, op)
	if validBase == nil {
		validBase = map[string]any{}
	}

	var cases []schema.TestCase
	for fieldName, prop := range s.Properties {
		if !prop.Nullable {
			continue
		}
		body := maps.Clone(validBase)
		body[fieldName] = nil

		specPath := fmt.Sprintf("%s %s requestBody.%s", op.Method, op.Path, fieldName)
		tc := buildTestCase(op, body,
			fmt.Sprintf("[semantic_annotation] nullable field %q accepts null", fieldName), specPath)
		tc.Priority = "P1"
		tc.Steps[0].Assertions = []schema.Assertion{
			{Target: "status_code", Operator: schema.OperatorGte, Expected: 200},
			{Target: "status_code", Operator: schema.OperatorLt, Expected: 300},
		}
		tc.Source = schema.CaseSource{
			Technique: "semantic_annotation",
			SpecPath:  specPath,
			Rationale: fmt.Sprintf("field %q is nullable; server MUST accept null value", fieldName),
			Scenario:  string(score.ScenarioNullableAcceptance),
		}
		cases = append(cases, tc)
	}
	return cases
}

// generateReadOnlyCases creates one test case per readOnly field for POST/PUT/PATCH operations.
// Each case sends a valid body that includes the readOnly field and expects 4xx.
func (t *SemanticAnnotationTechnique) generateReadOnlyCases(op *spec.Operation) []schema.TestCase {
	method := strings.ToUpper(op.Method)
	if method != "POST" && method != "PUT" && method != "PATCH" {
		return nil
	}

	s := getJSONSchema(op.RequestBody)
	if s == nil {
		return nil
	}

	validBase := buildValidBody(t.gen, op)
	if validBase == nil {
		validBase = map[string]any{}
	}

	var cases []schema.TestCase
	for fieldName, prop := range s.Properties {
		if !prop.ReadOnly {
			continue
		}

		// Generate a plausible value for the readOnly field
		value := t.gen.Generate(prop, fieldName)

		body := maps.Clone(validBase)
		body[fieldName] = value

		specPath := fmt.Sprintf("%s %s requestBody.%s", op.Method, op.Path, fieldName)
		tc := buildTestCase(op, body,
			fmt.Sprintf("[semantic_annotation] readOnly field %q write rejected", fieldName), specPath)
		tc.Priority = "P2"
		tc.Steps[0].Assertions = []schema.Assertion{
			{Target: "status_code", Operator: schema.OperatorGte, Expected: 400},
			{Target: "status_code", Operator: schema.OperatorLt, Expected: 500},
		}
		tc.Source = schema.CaseSource{
			Technique: "semantic_annotation",
			SpecPath:  specPath,
			Rationale: fmt.Sprintf("field %q is readOnly; server SHOULD reject writes to it", fieldName),
			Scenario:  string(score.ScenarioReadOnlyWrite),
		}
		cases = append(cases, tc)
	}
	return cases
}

// generateWriteOnlyCases creates one test case per writeOnly field found in GET 2xx response schemas.
// Each case asserts that the writeOnly field does NOT appear in the response body.
func (t *SemanticAnnotationTechnique) generateWriteOnlyCases(op *spec.Operation) []schema.TestCase {
	if strings.ToUpper(op.Method) != "GET" {
		return nil
	}

	var cases []schema.TestCase
	for code, resp := range op.Responses {
		n := 0
		fmt.Sscanf(code, "%d", &n)
		if n < 200 || n >= 300 {
			continue
		}
		mt, ok := resp.Content["application/json"]
		if !ok || mt.Schema == nil {
			continue
		}
		for fieldName, prop := range mt.Schema.Properties {
			if !prop.WriteOnly {
				continue
			}

			specPath := fmt.Sprintf("%s %s responses.%s.%s", op.Method, op.Path, code, fieldName)
			tc := buildTestCase(op, nil,
				fmt.Sprintf("[semantic_annotation] writeOnly field %q absent from response", fieldName), specPath)
			tc.Priority = "P2"
			tc.Steps[0].Assertions = []schema.Assertion{
				{Target: "status_code", Operator: schema.OperatorEq, Expected: 200},
				{Target: fmt.Sprintf("jsonpath $.%s", fieldName), Operator: schema.OperatorNotExists},
			}
			tc.Source = schema.CaseSource{
				Technique: "semantic_annotation",
				SpecPath:  specPath,
				Rationale: fmt.Sprintf("field %q is writeOnly; server MUST NOT include it in GET responses", fieldName),
				Scenario:  string(score.ScenarioWriteOnlyRead),
			}
			cases = append(cases, tc)
		}
	}
	return cases
}
