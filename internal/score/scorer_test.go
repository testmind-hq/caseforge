// internal/score/scorer_test.go
package score

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testmind-hq/caseforge/internal/output/schema"
)

// helpers

func makeCase(method, path, technique string, assertions int) schema.TestCase {
	step := schema.Step{Method: method, Path: path}
	for i := 0; i < assertions; i++ {
		step.Assertions = append(step.Assertions, schema.Assertion{
			Target:   "status_code",
			Operator: schema.OperatorEq,
			Expected: 200,
		})
	}
	return schema.TestCase{
		Source: schema.CaseSource{
			Technique: technique,
			SpecPath:  method + " " + path,
		},
		Steps: []schema.Step{step},
	}
}

// makeCaseWithSpecPath creates a case where the step path differs from the spec path
// (as OWASP techniques do when injecting attack payloads).
func makeCaseWithSpecPath(specMethod, specPath, stepMethod, stepPath, technique string, assertions int) schema.TestCase {
	step := schema.Step{Method: stepMethod, Path: stepPath}
	for i := 0; i < assertions; i++ {
		step.Assertions = append(step.Assertions, schema.Assertion{
			Target:   "status_code",
			Operator: schema.OperatorEq,
			Expected: 200,
		})
	}
	return schema.TestCase{
		Source: schema.CaseSource{
			Technique: technique,
			SpecPath:  specMethod + " " + specPath,
		},
		Steps: []schema.Step{step},
	}
}

func TestCompute_EmptyCases(t *testing.T) {
	r := Compute(nil)
	assert.Equal(t, 0, r.Overall)
	assert.Equal(t, 0, r.TotalCases)
	assert.Len(t, r.Dimensions, 4)
	for _, d := range r.Dimensions {
		assert.Equal(t, 0, d.Score)
	}
}

func TestCompute_FullCoverage(t *testing.T) {
	cases := []schema.TestCase{
		makeCase("GET", "/pets", "equivalence_partitioning", 1),
		makeCase("GET", "/pets", "boundary_value", 2),
		makeCase("GET", "/pets", "owasp_api_top10", 1),
		makeCase("GET", "/pets", "pairwise", 1),
		makeCase("POST", "/pets", "equivalence_partitioning", 1),
		makeCase("POST", "/pets", "boundary_value", 1),
		makeCase("POST", "/pets", "owasp_api_top10", 1),
		makeCase("POST", "/pets", "decision_table", 1),
	}
	r := Compute(cases)
	assert.Equal(t, 2, r.TotalOps)
	assert.Equal(t, len(cases), r.TotalCases)
	// All cases have assertions → executability = 100
	exec := r.Dimensions[3]
	assert.Equal(t, 100, exec.Score)
	// All ops have boundary → boundary = 100
	boundary := r.Dimensions[1]
	assert.Equal(t, 100, boundary.Score)
	// All ops have owasp → security = 100
	sec := r.Dimensions[2]
	assert.Equal(t, 100, sec.Score)
	// Overall should be fairly high
	assert.GreaterOrEqual(t, r.Overall, 70)
}

func TestCompute_NoSecurityCases(t *testing.T) {
	cases := []schema.TestCase{
		makeCase("GET", "/pets", "equivalence_partitioning", 1),
		makeCase("GET", "/pets", "boundary_value", 1),
		makeCase("POST", "/pets", "equivalence_partitioning", 1),
	}
	r := Compute(cases)
	sec := r.Dimensions[2]
	assert.Equal(t, 0, sec.Score, "no security cases → 0")
	// Should have a suggestion to add OWASP cases
	assert.NotEmpty(t, r.Suggestions)
	assert.Contains(t, r.Suggestions[0].Command, "owasp_api_top10")
}

func TestCompute_SpecLevelSecurityCoversAll(t *testing.T) {
	cases := []schema.TestCase{
		makeCase("GET", "/pets", "equivalence_partitioning", 1),
		makeCase("POST", "/pets", "owasp_api_top10_spec", 1),
	}
	r := Compute(cases)
	sec := r.Dimensions[2]
	assert.Equal(t, 100, sec.Score, "spec-level security counts for all ops")
}

func TestCompute_NoBoundaryCases_GeneratesSuggestion(t *testing.T) {
	cases := []schema.TestCase{
		makeCase("GET", "/pets", "owasp_api_top10", 1),
		makeCase("POST", "/orders", "owasp_api_top10", 1),
	}
	r := Compute(cases)
	boundary := r.Dimensions[1]
	assert.Equal(t, 0, boundary.Score)
	// Should have boundary improvement suggestions (one per missing-boundary op)
	boundaryMsgs := 0
	for _, s := range r.Suggestions {
		if s.Message != "" && s.Command == "caseforge gen --technique boundary_value,equivalence_partitioning" {
			boundaryMsgs++
		}
	}
	assert.GreaterOrEqual(t, boundaryMsgs, 1)
}

func TestCompute_NoAssertions_ExecZero(t *testing.T) {
	c := makeCase("GET", "/pets", "equivalence_partitioning", 0)
	r := Compute([]schema.TestCase{c})
	exec := r.Dimensions[3]
	assert.Equal(t, 0, exec.Score)
}

func TestCompute_SuggestionsAtMostThreeBoundaryOps(t *testing.T) {
	// 5 operations, none with boundary cases, no security
	ops := []string{"/a", "/b", "/c", "/d", "/e"}
	var cases []schema.TestCase
	for _, p := range ops {
		cases = append(cases, makeCase("GET", p, "state_transition", 1))
	}
	r := Compute(cases)
	// Boundary suggestions capped at 3 ops; plus 1 security suggestion = max 4
	assert.LessOrEqual(t, len(r.Suggestions), 4)
}

func TestCompute_BreadthScalesByTechniqueDiversity(t *testing.T) {
	// 1 operation with 4 distinct techniques → breadth = 100
	cases := []schema.TestCase{
		makeCase("GET", "/a", "equivalence_partitioning", 1),
		makeCase("GET", "/a", "boundary_value", 1),
		makeCase("GET", "/a", "owasp_api_top10", 1),
		makeCase("GET", "/a", "pairwise", 1),
	}
	r := Compute(cases)
	breadth := r.Dimensions[0]
	assert.Equal(t, 100, breadth.Score)
}

func TestCompute_BreadthLowWhenSingleTechnique(t *testing.T) {
	cases := []schema.TestCase{
		makeCase("GET", "/a", "equivalence_partitioning", 1),
		makeCase("POST", "/b", "equivalence_partitioning", 1),
	}
	r := Compute(cases)
	breadth := r.Dimensions[0]
	// avg 1 technique/op → score = 1*25 = 25
	assert.Equal(t, 25, breadth.Score)
}

func TestCompute_SuggestionsNeverNullInJSON(t *testing.T) {
	// Full-coverage case: no suggestions generated.
	// Verify JSON encodes "suggestions":[] not "suggestions":null.
	cases := []schema.TestCase{
		makeCase("GET", "/a", "equivalence_partitioning", 1),
		makeCase("GET", "/a", "boundary_value", 1),
		makeCase("GET", "/a", "owasp_api_top10", 1),
		makeCase("GET", "/a", "pairwise", 1),
	}
	r := Compute(cases)
	assert.NotNil(t, r.Suggestions, "Suggestions must be non-nil slice")
	// Zero suggestions when coverage is perfect.
	assert.Empty(t, r.Suggestions)

	// Empty-cases path must also produce non-nil slice.
	r2 := Compute(nil)
	assert.NotNil(t, r2.Suggestions)
}

func TestCompute_DimensionOrder(t *testing.T) {
	r := Compute([]schema.TestCase{makeCase("GET", "/x", "equivalence_partitioning", 1)})
	assert.Equal(t, "Coverage Breadth", r.Dimensions[0].Name)
	assert.Equal(t, "Boundary Coverage", r.Dimensions[1].Name)
	assert.Equal(t, "Security Coverage", r.Dimensions[2].Name)
	assert.Equal(t, "Executability", r.Dimensions[3].Name)
}

// TestCompute_OWASPInjectedPathsDoNotInflateOpCount is a regression test for the
// bug where OWASP techniques inject attack payloads into step paths (SQLi, XSS,
// path traversal, {{other_resource_id}}), causing the scorer to treat each injected
// URL as a new distinct "operation" instead of grouping under the original spec op.
func TestCompute_OWASPInjectedPathsDoNotInflateOpCount(t *testing.T) {
	cases := []schema.TestCase{
		// One real spec operation: DELETE /api/tokens/{id}
		makeCase("DELETE", "/api/tokens/{id}", "equivalence_partitioning", 1),
		// OWASP cases for the same operation — step paths are injected, spec_path is canonical.
		makeCaseWithSpecPath("DELETE", "/api/tokens/{id}", "DELETE", "/api/tokens/%27+OR+1%3D1--", "owasp_api_top10", 1),
		makeCaseWithSpecPath("DELETE", "/api/tokens/{id}", "DELETE", "/api/tokens/%22%3E%3Cscript%3E", "owasp_api_top10", 1),
		makeCaseWithSpecPath("DELETE", "/api/tokens/{id}", "DELETE", "/api/tokens/{{other_resource_id}}", "owasp_api_top10", 1),
		// OPTIONS CORS preflight for same path
		makeCaseWithSpecPath("OPTIONS", "/api/tokens/{id}", "OPTIONS", "/api/tokens/{id}", "owasp_api_top10", 1),
	}
	r := Compute(cases)
	// Without the fix: TotalOps would be 5 (one per unique step path).
	// With the fix: TotalOps should be 2 (DELETE + OPTIONS are distinct spec_path methods).
	assert.LessOrEqual(t, r.TotalOps, 2, "OWASP injected paths must not inflate op count beyond the distinct spec operations")
}

// TestCompute_CanonicalOpKey_FallbackToStepWhenSpecPathEmpty verifies that cases
// with an empty SpecPath still group correctly by step method+path.
func TestCompute_CanonicalOpKey_FallbackToStepWhenSpecPathEmpty(t *testing.T) {
	c := schema.TestCase{
		Source: schema.CaseSource{Technique: "equivalence_partitioning"}, // SpecPath intentionally empty
		Steps: []schema.Step{
			{Method: "GET", Path: "/fallback"},
		},
	}
	r := Compute([]schema.TestCase{c})
	assert.Equal(t, 1, r.TotalOps)
}
