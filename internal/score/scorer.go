// internal/score/scorer.go
// Package score computes a multi-dimensional quality score for a set of
// CaseForge test cases and produces improvement suggestions.
package score

import (
	"fmt"
	"sort"

	"github.com/testmind-hq/caseforge/internal/output/schema"
)

// opKey is the canonical key for a test operation (first step method + path).
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
			},
			Suggestions: []Suggestion{},
		}
	}

	// Group cases by the primary operation (method + path of the first step).
	opCases := make(map[opKey][]schema.TestCase)
	for _, c := range cases {
		if len(c.Steps) > 0 {
			k := opKey{c.Steps[0].Method, c.Steps[0].Path}
			opCases[k] = append(opCases[k], c)
		}
	}
	totalOps := len(opCases)

	breadthScore, breadthDetail := computeBreadth(opCases, totalOps)
	boundaryScore, boundaryDetail, opsMissingBoundary := computeBoundary(opCases, totalOps)
	secScore, secDetail := computeSecurity(opCases, totalOps)
	execScore, execDetail := computeExecutability(cases)

	// Weighted overall: breadth 30 %, boundary 25 %, security 25 %, exec 20 %
	overall := (breadthScore*30 + boundaryScore*25 + secScore*25 + execScore*20) / 100

	suggestions := buildSuggestions(secScore, opsMissingBoundary)

	return Report{
		Overall:    overall,
		TotalCases: len(cases),
		TotalOps:   totalOps,
		Dimensions: []Dimension{
			{Name: "Coverage Breadth", Score: breadthScore, Detail: breadthDetail},
			{Name: "Boundary Coverage", Score: boundaryScore, Detail: boundaryDetail},
			{Name: "Security Coverage", Score: secScore, Detail: secDetail},
			{Name: "Executability", Score: execScore, Detail: execDetail},
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
	detail := fmt.Sprintf("%d/%d operations have boundary/equivalence cases", covered, totalOps)
	return score, detail, missing
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
func buildSuggestions(secScore int, opsMissingBoundary []opKey) []Suggestion {
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
