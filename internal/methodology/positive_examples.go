// internal/methodology/positive_examples.go
package methodology

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/score"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// PositiveExamplesTechnique generates happy-path test cases by assembling complete
// requests from parameter-level OpenAPI examples. Complements ExampleTechnique which
// only handles request body examples.
//
// For each named example found in any parameter, it builds a complete request URL
// with that example's value substituted, expecting a 2xx response.
type PositiveExamplesTechnique struct{}

func NewPositiveExamplesTechnique() *PositiveExamplesTechnique {
	return &PositiveExamplesTechnique{}
}

func (t *PositiveExamplesTechnique) Name() string { return "positive_examples" }

// Applies returns true if any parameter has Example != nil OR len(Examples) > 0.
func (t *PositiveExamplesTechnique) Applies(op *spec.Operation) bool {
	for _, p := range op.Parameters {
		if p.Example != nil || len(p.Examples) > 0 {
			return true
		}
	}
	return false
}

// Generate produces one TestCase per named example (or one "inline" case when only
// single-value examples exist on parameters).
func (t *PositiveExamplesTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	names := collectExampleNames(op)

	var cases []schema.TestCase

	if len(names) > 0 {
		for _, exName := range names {
			tc := buildPositiveExampleCase(op, exName)
			cases = append(cases, tc)
		}
	} else {
		// No named examples; check for single-value (inline) examples on any parameter.
		hasInline := false
		for _, p := range op.Parameters {
			if p.Example != nil {
				hasInline = true
				break
			}
		}
		if hasInline {
			tc := buildPositiveExampleCase(op, "inline")
			cases = append(cases, tc)
		}
	}

	return cases, nil
}

// buildPositiveExampleCase assembles a test case for the given example name.
func buildPositiveExampleCase(op *spec.Operation, exampleName string) schema.TestCase {
	requestPath := buildPositiveRequestPath(op, exampleName)
	title := fmt.Sprintf("positive example %q", exampleName)
	id := fmt.Sprintf("TC-%s", uuid.New().String()[:8])

	tc := schema.TestCase{
		Schema:   schema.SchemaBaseURL,
		Version:  "1",
		ID:       id,
		Title:    fmt.Sprintf("%s %s - %s", op.Method, op.Path, title),
		Kind:     "single",
		Priority: "P1",
		Tags:     op.Tags,
		Steps: []schema.Step{
			{
				ID:    "step-main",
				Title: title,
				Type:  "test",
				Method: op.Method,
				Path:  requestPath,
				Assertions: []schema.Assertion{
					{Target: "status_code", Operator: "gte", Expected: 200},
					{Target: "status_code", Operator: "lt", Expected: 300},
				},
			},
		},
		Source: schema.CaseSource{
			Technique: "positive_examples",
			SpecPath:  fmt.Sprintf("%s %s parameters.%s", op.Method, op.Path, exampleName),
			Scenario:  score.ScenarioPositiveExample,
			Rationale: title,
		},
		GeneratedAt: time.Now(),
	}
	return tc
}

// collectExampleNames collects all named example names across all parameters (sorted).
func collectExampleNames(op *spec.Operation) []string {
	seen := map[string]bool{}
	for _, p := range op.Parameters {
		for name := range p.Examples {
			seen[name] = true
		}
	}
	names := make([]string, 0, len(seen))
	for n := range seen {
		names = append(names, n)
	}
	slices.Sort(names)
	return names
}

// exampleValueFor returns the example value for a parameter given an example name.
// Falls back to the single p.Example if no named example matches.
func exampleValueFor(p *spec.Parameter, exampleName string) any {
	if ex, ok := p.Examples[exampleName]; ok && ex != nil {
		return ex.Value
	}
	if p.Example != nil {
		return p.Example
	}
	return nil
}

// buildPositiveRequestPath builds the request path with parameter substitutions and
// query string appended.
func buildPositiveRequestPath(op *spec.Operation, exampleName string) string {
	path := op.Path
	queryParams := map[string]any{}
	for _, p := range op.Parameters {
		val := exampleValueFor(p, exampleName)
		if val == nil {
			continue
		}
		if p.In == "path" {
			path = strings.ReplaceAll(path, "{"+p.Name+"}", fmt.Sprintf("%v", val))
		} else if p.In == "query" {
			queryParams[p.Name] = val
		}
	}
	return buildPathWithQuery(path, queryParams)
}
