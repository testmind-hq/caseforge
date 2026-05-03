// internal/score/scorer.go
// Package score computes a multi-dimensional quality score for a set of
// CaseForge test cases and produces improvement suggestions.
package score

import (
	"fmt"
	"sort"
	"strings"

	"github.com/testmind-hq/caseforge/internal/output/schema"
)

// opKey is the canonical key for a spec operation (method + path from CaseSource.SpecPath, or from the first step as fallback).
type opKey struct{ method, path string }

// securityTechniques are the technique names that count towards security coverage.
var securityTechniques = map[string]bool{
	"owasp_api_top10":      true,
	"owasp_api_top10_spec": true,
}

// boundaryTechniques are the technique names that count towards boundary coverage.
var boundaryTechniques = map[string]bool{
	"boundary_value":              true,
	"equivalence_partitioning":    true,
	"decision_table":              true,
	"pairwise":                    true,
	"classification_tree":         true,
	"orthogonal_array":            true,
}

// Dimension is a single scored aspect of test-case quality.
type Dimension struct {
	Name   string `json:"name"`
	Score  int    `json:"score"`
	Detail string `json:"detail"`
}

// Suggestion is an improvement recommendation ordered by priority (1 = highest).
type Suggestion struct {
	Priority int    `json:"priority"`
	Message  string `json:"message"`
	Command  string `json:"command,omitempty"`
}

// Report is the complete quality report produced by Compute.
type Report struct {
	Overall     int          `json:"overall"`
	TotalCases  int          `json:"total_cases"`
	TotalOps    int          `json:"total_ops"`
	Dimensions  []Dimension  `json:"dimensions"`
	Suggestions []Suggestion `json:"suggestions"`
}

// Compute scores cases across four dimensions and returns a Report.
// When len(cases) == 0 all scores are 0.
func Compute(cases []schema.TestCase) Report {
	if len(cases) == 0 {
		zero := func(name string) Dimension {
			return Dimension{Name: name, Score: 0, Detail: "no test cases found"}
		}
		return Report{
			Dimensions: []Dimension{
				zero("Coverage Breadth"),
				zero("Boundary Coverage"),
				zero("Security Coverage"),
				zero("Executability"),
				zero("Status Coverage"),
			},
			Suggestions: []Suggestion{},
		}
	}

	// Group cases by the canonical spec operation key derived from CaseSource.SpecPath.
	// Using SpecPath (e.g. "DELETE /api/tokens/{id}") rather than the step's actual
	// path prevents OWASP-injected attack payloads (SQLi, XSS, path traversal,
	// BOLA placeholders) from inflating the distinct-operation count.
	opCases := make(map[opKey][]schema.TestCase)
	for _, c := range cases {
		k := canonicalOpKey(c)
		if k.method != "" {
			opCases[k] = append(opCases[k], c)
		}
	}
	totalOps := len(opCases)

	breadthScore, breadthDetail := computeBreadth(opCases, totalOps)
	boundaryScore, boundaryDetail, opsMissingBoundary := computeBoundary(opCases, totalOps)
	secScore, secDetail := computeSecurity(opCases, totalOps)
	execScore, execDetail := computeExecutability(cases)
	statusScore, statusDetail := computeStatusCoverage(opCases, totalOps)

	// Weighted overall: breadth 25 %, boundary 20 %, security 20 %, exec 15 %, status_coverage 20 %
	overall := (breadthScore*25 + boundaryScore*20 + secScore*20 + execScore*15 + statusScore*20) / 100

	suggestions := buildSuggestions(secScore, statusScore, opsMissingBoundary)

	return Report{
		Overall:    overall,
		TotalCases: len(cases),
		TotalOps:   totalOps,
		Dimensions: []Dimension{
			{Name: "Coverage Breadth", Score: breadthScore, Detail: breadthDetail},
			{Name: "Boundary Coverage", Score: boundaryScore, Detail: boundaryDetail},
			{Name: "Security Coverage", Score: secScore, Detail: secDetail},
			{Name: "Executability", Score: execScore, Detail: execDetail},
			{Name: "Status Coverage", Score: statusScore, Detail: statusDetail},
		},
		Suggestions: suggestions,
	}
}

// computeBreadth scores technique diversity: avg distinct techniques per operation.
// 4+ distinct techniques per operation = 100; each missing technique step costs 25 pts.
func computeBreadth(opCases map[opKey][]schema.TestCase, totalOps int) (int, string) {
	if totalOps == 0 {
		return 0, "no operations found"
	}
	totalTech := 0
	for _, cases := range opCases {
		seen := make(map[string]bool)
		for _, c := range cases {
			if c.Source.Technique != "" {
				seen[c.Source.Technique] = true
			}
		}
		totalTech += len(seen)
	}
	avg := float64(totalTech) / float64(totalOps)
	score := int(avg * 25)
	if score > 100 {
		score = 100
	}
	detail := fmt.Sprintf("%d distinct operation(s) covered; avg %.1f technique(s)/operation", totalOps, avg)
	return score, detail
}

// computeBoundary returns the % of operations that have ≥ 1 boundary/equivalence case.
func computeBoundary(opCases map[opKey][]schema.TestCase, totalOps int) (int, string, []opKey) {
	if totalOps == 0 {
		return 0, "no operations found", nil
	}
	covered := 0
	var missing []opKey
	for op, cases := range opCases {
		hasBoundary := false
		for _, c := range cases {
			if boundaryTechniques[c.Source.Technique] {
				hasBoundary = true
				break
			}
		}
		if hasBoundary {
			covered++
		} else {
			missing = append(missing, op)
		}
	}
	// Sort for deterministic output.
	sort.Slice(missing, func(i, j int) bool {
		if missing[i].path != missing[j].path {
			return missing[i].path < missing[j].path
		}
		return missing[i].method < missing[j].method
	})
	score := covered * 100 / totalOps

	// Flatten all cases for scenario analysis
	var allCases []schema.TestCase
	for _, cs := range opCases {
		allCases = append(allCases, cs...)
	}
	coveredScen, missingScen := computeScenarioCoverage(allCases)

	detail := fmt.Sprintf("%d/%d operations have boundary/equivalence cases", covered, totalOps)
	if len(coveredScen) > 0 {
		detail += fmt.Sprintf("; covered scenarios: %s", strings.Join(coveredScen, ", "))
	}
	if len(missingScen) > 0 {
		detail += fmt.Sprintf("; missing scenarios: %s", strings.Join(missingScen, ", "))
	}
	return score, detail, missing
}

// computeScenarioCoverage collects all Scenario values from cases and returns
// (covered, missing) lists relative to the tracked scenario set.
func computeScenarioCoverage(cases []schema.TestCase) (covered, missing []string) {
	seen := make(map[string]bool)
	for _, c := range cases {
		if c.Source.Scenario != "" {
			seen[c.Source.Scenario] = true
		}
	}
	for _, s := range trackedScenarios {
		if seen[s] {
			covered = append(covered, s)
		} else {
			missing = append(missing, s)
		}
	}
	return
}

// computeSecurity returns the % of operations that have ≥ 1 security test case.
// Spec-level security techniques (owasp_api_top10_spec) count for all operations.
func computeSecurity(opCases map[opKey][]schema.TestCase, totalOps int) (int, string) {
	if totalOps == 0 {
		return 0, "no operations found"
	}
	// Check if any spec-level security case exists — if so, all ops are covered.
	specLevelSecurity := false
	for _, cases := range opCases {
		for _, c := range cases {
			if c.Source.Technique == "owasp_api_top10_spec" {
				specLevelSecurity = true
				break
			}
		}
		if specLevelSecurity {
			break
		}
	}
	if specLevelSecurity {
		return 100, "spec-level OWASP security cases cover all operations"
	}

	covered := 0
	for _, cases := range opCases {
		for _, c := range cases {
			if securityTechniques[c.Source.Technique] {
				covered++
				break
			}
		}
	}
	score := covered * 100 / totalOps
	detail := fmt.Sprintf("%d/%d operations have security (OWASP) cases (%d%%)", covered, totalOps, score)
	return score, detail
}

// computeStatusCoverage returns the % of operations that have both a 2xx-asserting
// case (happy path) and a 4xx/5xx-asserting case (error path).
//
// The metric concept — that an operation is "covered" only when both happy and
// error paths are tested — is informed by EvoMaster's research on coverage
// targets. This implementation is independent: it operates on caseforge's
// TestCase / Step / Assertion schema, with no source code derived from
// EvoMaster's Kotlin codebase. See NOTICE for full attribution.
func computeStatusCoverage(opCases map[opKey][]schema.TestCase, totalOps int) (int, string) {
	if totalOps == 0 {
		return 0, "no operations found"
	}
	covered := 0
	for _, cases := range opCases {
		has2xx := false
		has4xx := false
		for _, c := range cases {
			for _, step := range c.Steps {
				for _, a := range step.Assertions {
					if a.Target != "status_code" {
						continue
					}
					if is2xxAssertion(a) {
						has2xx = true
					}
					if is4xxAssertion(a) {
						has4xx = true
					}
				}
			}
		}
		if has2xx && has4xx {
			covered++
		}
	}
	score := covered * 100 / totalOps
	detail := fmt.Sprintf("%d/%d operations have both 2xx and 4xx/5xx cases", covered, totalOps)
	return score, detail
}

func is2xxAssertion(a schema.Assertion) bool {
	n, ok := toAssertInt(a.Expected)
	if !ok {
		return false
	}
	switch a.Operator {
	case schema.OperatorLt:
		// "status_code lt 300" is the canonical happy-path assertion; n must be exactly 300.
		return n == 300
	case schema.OperatorEq:
		return n >= 200 && n < 300
	}
	return false
}

func is4xxAssertion(a schema.Assertion) bool {
	n, ok := toAssertInt(a.Expected)
	if !ok {
		return false
	}
	switch a.Operator {
	case schema.OperatorGte:
		return n >= 400
	case schema.OperatorEq:
		return n >= 400
	}
	return false
}

func toAssertInt(v any) (int, bool) {
	switch x := v.(type) {
	case int:
		return x, true
	case int64:
		return int(x), true
	case float64:
		return int(x), true
	}
	return 0, false
}

// computeExecutability returns the % of cases that have ≥ 1 assertion.
func computeExecutability(cases []schema.TestCase) (int, string) {
	if len(cases) == 0 {
		return 0, "no test cases found"
	}
	executable := 0
	for _, c := range cases {
		for _, s := range c.Steps {
			if len(s.Assertions) > 0 {
				executable++
				break
			}
		}
	}
	score := executable * 100 / len(cases)
	detail := fmt.Sprintf("%d/%d cases have at least one assertion", executable, len(cases))
	return score, detail
}

// buildSuggestions assembles improvement suggestions ordered by priority.
// Returns an empty (non-nil) slice so JSON output is always "[]", never "null".
func buildSuggestions(secScore, statusScore int, opsMissingBoundary []opKey) []Suggestion {
	out := []Suggestion{}
	p := 1
	if secScore < 80 {
		out = append(out, Suggestion{
			Priority: p,
			Message:  fmt.Sprintf("Add OWASP security cases (security score: %d/100)", secScore),
			Command:  "caseforge gen --technique owasp_api_top10",
		})
		p++
	}
	if statusScore < 80 {
		out = append(out, Suggestion{
			Priority: p,
			Message:  fmt.Sprintf("Add error-path cases (status coverage: %d/100)", statusScore),
			Command:  "caseforge gen --technique mutation,isolated_negative,constraint_mutation",
		})
		p++
	}
	for i, op := range opsMissingBoundary {
		if i >= 3 {
			break
		}
		// Note: --operations accepts operationId values from the spec, which are
		// not stored in index.json. The suggestion therefore omits --operations and
		// lets the user target the operation manually via --spec filtering.
		out = append(out, Suggestion{
			Priority: p,
			Message:  fmt.Sprintf("%s %s is missing boundary/equivalence test cases", op.method, op.path),
			Command:  "caseforge gen --technique boundary_value,equivalence_partitioning",
		})
		p++
	}
	return out
}

// canonicalOpKey returns the operation key for a test case using CaseSource.SpecPath
// when available, falling back to the first step's method+path.
//
// CaseSource.SpecPath has the format "METHOD /path [field.subfield...]" — we take
// only the first two whitespace-separated tokens (method + path). This means OWASP
// cases that inject attack payloads into step paths (SQLi, XSS, path traversal,
// {{other_resource_id}}) are still grouped under their original spec operation,
// preventing artificial inflation of the distinct-operation count.
func canonicalOpKey(c schema.TestCase) opKey {
	if c.Source.SpecPath != "" {
		parts := strings.Fields(c.Source.SpecPath)
		if len(parts) >= 2 {
			return opKey{method: parts[0], path: parts[1]}
		}
	}
	if len(c.Steps) > 0 {
		return opKey{method: c.Steps[0].Method, path: c.Steps[0].Path}
	}
	return opKey{}
}
