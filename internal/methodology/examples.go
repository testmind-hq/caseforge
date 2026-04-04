// internal/methodology/examples.go
package methodology

import (
	"fmt"
	"sort"

	assertpkg "github.com/testmind-hq/caseforge/internal/assert"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// ExampleTechnique generates test cases from examples declared in the OpenAPI spec.
// It covers both mediaType-level examples (mediaType.example and mediaType.examples)
// and the schema-level example (schema.example on the request body schema).
//
// For each extracted example:
//   - If the example passes schema validation → a test case expecting a 2xx response
//   - If the example fails schema validation (e.g., missing required fields) → a test
//     case expecting a 4xx response (the spec is documenting an invalid input)
type ExampleTechnique struct{}

func NewExampleTechnique() *ExampleTechnique { return &ExampleTechnique{} }

func (t *ExampleTechnique) Name() string { return "example_extraction" }

// Applies returns true when the operation has a JSON request body with at least
// one example (mediaType.example, mediaType.examples, or schema.example).
func (t *ExampleTechnique) Applies(op *spec.Operation) bool {
	mt := jsonMediaType(op)
	if mt == nil {
		return false
	}
	if mt.Example != nil {
		return true
	}
	if len(mt.Examples) > 0 {
		return true
	}
	if mt.Schema != nil && mt.Schema.Example != nil {
		return true
	}
	return false
}

func (t *ExampleTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	mt := jsonMediaType(op)
	if mt == nil {
		return nil, nil
	}

	var cases []schema.TestCase

	// Collect all examples as (label, value) pairs in deterministic order.
	type namedExample struct {
		label string
		value any
	}
	var examples []namedExample

	// 1. mediaType.examples (named map) — iterate in sorted key order.
	if len(mt.Examples) > 0 {
		names := make([]string, 0, len(mt.Examples))
		for n := range mt.Examples {
			names = append(names, n)
		}
		sort.Strings(names)
		for _, name := range names {
			ex := mt.Examples[name]
			if ex == nil || ex.Value == nil {
				continue
			}
			examples = append(examples, namedExample{name, ex.Value})
		}
	}

	// 2. mediaType.example (single value).
	if mt.Example != nil {
		examples = append(examples, namedExample{"inline", mt.Example})
	}

	// 3. schema.example as fallback when no mediaType-level examples exist.
	if len(examples) == 0 && mt.Schema != nil && mt.Schema.Example != nil {
		examples = append(examples, namedExample{"schema", mt.Schema.Example})
	}

	for _, ex := range examples {
		body, ok := ex.value.(map[string]any)
		if !ok {
			// Non-object examples: wrap in a simple structure for the test step.
			body = nil
		}

		errs := spec.ValidateExample(ex.value, mt.Schema)
		isValid := len(errs) == 0

		var title string
		if isValid {
			title = fmt.Sprintf("example %q — valid input", ex.label)
		} else {
			title = fmt.Sprintf("example %q — invalid input (%s)", ex.label, errs[0])
		}

		tc := buildTestCase(op, body, title,
			fmt.Sprintf("%s %s examples.%s", op.Method, op.Path, ex.label))

		if isValid {
			tc.Priority = "P1"
			tc.Steps[0].Assertions = append(
				assertpkg.BasicAssertions(op),
				responseSchemaAssertions(op)...,
			)
		} else {
			tc.Priority = "P2"
			tc.Steps[0].Assertions = []schema.Assertion{
				{Target: "status_code", Operator: "eq", Expected: 422},
			}
		}

		tc.Source = schema.CaseSource{
			Technique: "example_extraction",
			SpecPath:  fmt.Sprintf("%s %s requestBody examples.%s", op.Method, op.Path, ex.label),
			Rationale: title,
		}

		cases = append(cases, tc)
	}

	return cases, nil
}

// jsonMediaType returns the application/json MediaType for op's request body, or nil.
func jsonMediaType(op *spec.Operation) *spec.MediaType {
	if op.RequestBody == nil {
		return nil
	}
	mt, ok := op.RequestBody.Content["application/json"]
	if !ok {
		return nil
	}
	return mt
}
