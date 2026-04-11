// internal/methodology/positive_examples_test.go
package methodology

import (
	"strings"
	"testing"

	"github.com/testmind-hq/caseforge/internal/spec"
)

// helper: build an operation with path parameters that have named examples.
func opWithNamedParamExamples() *spec.Operation {
	return &spec.Operation{
		Method: "GET",
		Path:   "/users/{id}",
		Parameters: []*spec.Parameter{
			{
				Name:     "id",
				In:       "path",
				Required: true,
				Schema:   &spec.Schema{Type: "integer"},
				Examples: map[string]*spec.Example{
					"activeUser": {Summary: "Active user", Value: 42},
					"adminUser":  {Summary: "Admin user", Value: 1},
				},
			},
		},
	}
}

// helper: build an operation with a single-value (non-named) parameter example.
func opWithSingleParamExample() *spec.Operation {
	return &spec.Operation{
		Method: "GET",
		Path:   "/search",
		Parameters: []*spec.Parameter{
			{
				Name:     "q",
				In:       "query",
				Required: true,
				Schema:   &spec.Schema{Type: "string"},
				Example:  "hello world",
			},
		},
	}
}

// helper: build an operation with no parameter examples.
func opWithNoParamExamples() *spec.Operation {
	return &spec.Operation{
		Method: "GET",
		Path:   "/items",
		Parameters: []*spec.Parameter{
			{
				Name:     "limit",
				In:       "query",
				Required: false,
				Schema:   &spec.Schema{Type: "integer"},
			},
		},
	}
}

func TestPositiveExamplesTechnique_Applies_WithNamedParamExamples(t *testing.T) {
	tech := NewPositiveExamplesTechnique()
	op := opWithNamedParamExamples()
	if !tech.Applies(op) {
		t.Error("expected Applies to return true for operation with named parameter examples")
	}
}

func TestPositiveExamplesTechnique_Applies_WithSingleParamExample(t *testing.T) {
	tech := NewPositiveExamplesTechnique()
	op := opWithSingleParamExample()
	if !tech.Applies(op) {
		t.Error("expected Applies to return true for operation with single parameter example")
	}
}

func TestPositiveExamplesTechnique_Applies_False_NoExamples(t *testing.T) {
	tech := NewPositiveExamplesTechnique()
	op := opWithNoParamExamples()
	if tech.Applies(op) {
		t.Error("expected Applies to return false for operation with no parameter examples")
	}
}

func TestPositiveExamplesTechnique_Generate_OneCasePerNamedExample(t *testing.T) {
	tech := NewPositiveExamplesTechnique()
	op := opWithNamedParamExamples()

	cases, err := tech.Generate(op)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if len(cases) != 2 {
		t.Errorf("expected 2 cases (one per named example), got %d", len(cases))
	}
}

func TestPositiveExamplesTechnique_Generate_PathParamSubstituted(t *testing.T) {
	tech := NewPositiveExamplesTechnique()
	op := opWithNamedParamExamples()

	cases, err := tech.Generate(op)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	// Collect all step paths
	var paths []string
	for _, tc := range cases {
		for _, step := range tc.Steps {
			paths = append(paths, step.Path)
		}
	}

	// At least one case should have /users/42 substituted
	found42 := false
	for _, p := range paths {
		if p == "/users/42" {
			found42 = true
		}
	}
	if !found42 {
		t.Errorf("expected one case with path /users/42, got paths: %v", paths)
	}
}

func TestPositiveExamplesTechnique_Generate_QueryParamAppended(t *testing.T) {
	tech := NewPositiveExamplesTechnique()
	op := opWithSingleParamExample()

	cases, err := tech.Generate(op)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if len(cases) != 1 {
		t.Fatalf("expected 1 case for single-value example, got %d", len(cases))
	}

	stepPath := cases[0].Steps[0].Path
	if !strings.Contains(stepPath, "q=") {
		t.Errorf("expected query param q= in path, got %q", stepPath)
	}
}

func TestPositiveExamplesTechnique_Generate_Expects2xx(t *testing.T) {
	tech := NewPositiveExamplesTechnique()
	op := opWithNamedParamExamples()

	cases, err := tech.Generate(op)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	for _, tc := range cases {
		for _, step := range tc.Steps {
			hasGte200 := false
			hasLt300 := false
			for _, a := range step.Assertions {
				if a.Target == "status_code" && a.Operator == "gte" && a.Expected == 200 {
					hasGte200 = true
				}
				if a.Target == "status_code" && a.Operator == "lt" && a.Expected == 300 {
					hasLt300 = true
				}
			}
			if !hasGte200 || !hasLt300 {
				t.Errorf("case %q step %q: expected 2xx assertions (gte 200 AND lt 300), got %+v",
					tc.ID, step.ID, step.Assertions)
			}
		}
	}
}

func TestPositiveExamplesTechnique_Generate_ScenarioPopulated(t *testing.T) {
	tech := NewPositiveExamplesTechnique()
	op := opWithNamedParamExamples()

	cases, err := tech.Generate(op)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	for _, tc := range cases {
		if tc.Source.Scenario != "POSITIVE_EXAMPLE" {
			t.Errorf("expected Source.Scenario = POSITIVE_EXAMPLE, got %q", tc.Source.Scenario)
		}
		if tc.Source.Technique != "positive_examples" {
			t.Errorf("expected Source.Technique = positive_examples, got %q", tc.Source.Technique)
		}
	}
}
