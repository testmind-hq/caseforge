// internal/methodology/pairwise.go
package methodology

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// PairwiseParam represents a single parameter with its possible values.
type PairwiseParam struct {
	Name   string
	Values []any
}

// IPOG implements the In-Parameter-Order-General algorithm for 2-way covering arrays.
// Returns a slice of rows, where each row is a slice of values in parameter order.
func IPOG(params []PairwiseParam) [][]any {
	if len(params) == 0 {
		return nil
	}
	if len(params) < 2 {
		// Need at least 2 parameters for pairwise coverage
		if len(params) == 1 {
			// Return one row per value of the single parameter
			var rows [][]any
			for _, v := range params[0].Values {
				rows = append(rows, []any{v})
			}
			return rows
		}
		return nil
	}
	// Start with the first two parameters: full Cartesian product
	var rows [][]any
	for _, v1 := range params[0].Values {
		for _, v2 := range params[1].Values {
			rows = append(rows, []any{v1, v2})
		}
	}

	// Iteratively extend each row with a new parameter
	for col := 2; col < len(params); col++ {
		param := params[col]
		// Track which pairs (existing col, new col) are still uncovered.
		// Key is [3]any{prevColIdx, prevValue, newValue} to avoid false matches
		// when different columns share the same enum value string.
		uncovered := make(map[[3]any]bool)
		for prevCol := 0; prevCol < col; prevCol++ {
			for _, v1 := range params[prevCol].Values {
				for _, v2 := range param.Values {
					uncovered[[3]any{prevCol, v1, v2}] = true
				}
			}
		}

		// Extend existing rows: pick value that covers the most new pairs
		for i := range rows {
			bestVal := param.Values[0]
			bestCoverage := -1
			for _, candidate := range param.Values {
				coverage := 0
				for prevCol := 0; prevCol < col; prevCol++ {
					key := [3]any{prevCol, rows[i][prevCol], candidate}
					if uncovered[key] {
						coverage++
					}
				}
				if coverage > bestCoverage {
					bestCoverage = coverage
					bestVal = candidate
				}
			}
			rows[i] = append(rows[i], bestVal)
			// Mark covered pairs
			for prevCol := 0; prevCol < col; prevCol++ {
				delete(uncovered, [3]any{prevCol, rows[i][prevCol], bestVal})
			}
		}

		// Add new rows for any still-uncovered pairs
		for len(uncovered) > 0 {
			row := make([]any, col+1)
			// Fill with wildcard (first value of each param) initially
			for k := 0; k < col; k++ {
				row[k] = params[k].Values[0]
			}
			// Pick a remaining uncovered pair and assign it
			for pair := range uncovered {
				// pair[0] is the prev col index, pair[1] is the prev col value,
				// pair[2] is the new col value
				prevCol := pair[0].(int)
				row[prevCol] = pair[1]
				row[col] = pair[2]
				break
			}
			rows = append(rows, row)
			// Mark covered pairs for this row
			for prevCol := 0; prevCol < col; prevCol++ {
				delete(uncovered, [3]any{prevCol, row[prevCol], row[col]})
			}
		}
	}

	return rows
}

// pendingTuple tracks a t-tuple that still needs to be covered.
type pendingTuple struct {
	prevCols []int
	prevVals []any
	newVal   any
}

// IPOGt implements In-Parameter-Order-General for t-way (N-way) covering arrays.
// t=2 produces the same result as IPOG (standard pairwise).
// t=3 ensures every triple of (col_i, col_j, col_k) value combinations is covered.
func IPOGt(params []PairwiseParam, t int) [][]any {
	if t < 2 {
		t = 2
	}
	if len(params) < t {
		return cartesianProductParams(params)
	}

	// Bootstrap: full Cartesian product of the first t parameters
	rows := cartesianProductParams(params[:t])

	// Extend one parameter at a time
	for col := t; col < len(params); col++ {
		param := params[col]

		// All (t-1)-subsets of columns 0..col-1
		prevColSets := columnSubsets(col, t-1)

		// Build list of uncovered t-tuples
		var uncovered []pendingTuple
		for _, prevCols := range prevColSets {
			for _, pvc := range enumValueCombinations(params, prevCols) {
				for _, newVal := range param.Values {
					uncovered = append(uncovered, pendingTuple{
						prevCols: prevCols,
						prevVals: pvc,
						newVal:   newVal,
					})
				}
			}
		}

		// Extend existing rows: pick value that covers the most uncovered tuples
		for i := range rows {
			bestVal := param.Values[0]
			bestCoverage := countPendingCoverage(rows[i], param.Values[0], col, uncovered)
			for _, candidate := range param.Values[1:] {
				if c := countPendingCoverage(rows[i], candidate, col, uncovered); c > bestCoverage {
					bestCoverage = c
					bestVal = candidate
				}
			}
			rows[i] = append(rows[i], bestVal)
			// Remove covered tuples
			uncovered = removeCoveredTuples(uncovered, rows[i], col, bestVal)
		}

		// Add new rows for still-uncovered t-tuples
		for len(uncovered) > 0 {
			// Pick the first uncovered tuple and seed the row from it
			seed := uncovered[0]
			row := make([]any, col+1)
			for k := 0; k < col; k++ {
				row[k] = params[k].Values[0]
			}
			for k, pc := range seed.prevCols {
				row[pc] = seed.prevVals[k]
			}
			row[col] = seed.newVal
			rows = append(rows, row)
			// Remove all tuples covered by this new row
			uncovered = removeCoveredTuples(uncovered, row, col, row[col])
		}
	}
	return rows
}

// countPendingCoverage counts how many pending tuples would be covered if column col
// is assigned candidate in the given row.
func countPendingCoverage(row []any, candidate any, col int, pending []pendingTuple) int {
	count := 0
	for _, pt := range pending {
		if pt.newVal != candidate {
			continue
		}
		match := true
		for k, pc := range pt.prevCols {
			if row[pc] != pt.prevVals[k] {
				match = false
				break
			}
		}
		if match {
			count++
		}
	}
	return count
}

// removeCoveredTuples returns the pending list with tuples covered by row[col]=newVal removed.
func removeCoveredTuples(pending []pendingTuple, row []any, col int, newVal any) []pendingTuple {
	result := pending[:0]
	for _, pt := range pending {
		if pt.newVal == newVal {
			match := true
			for k, pc := range pt.prevCols {
				if row[pc] != pt.prevVals[k] {
					match = false
					break
				}
			}
			if match {
				continue // covered, skip
			}
		}
		result = append(result, pt)
	}
	return result
}

// columnSubsets returns all k-element subsets of column indices 0..n-1.
func columnSubsets(n, k int) [][]int {
	if k == 0 {
		return [][]int{{}}
	}
	var result [][]int
	var recurse func(start int, current []int)
	recurse = func(start int, current []int) {
		if len(current) == k {
			cp := make([]int, k)
			copy(cp, current)
			result = append(result, cp)
			return
		}
		for i := start; i < n; i++ {
			recurse(i+1, append(current, i))
		}
	}
	recurse(0, nil)
	return result
}

// enumValueCombinations returns all value combinations for the given column indices.
func enumValueCombinations(params []PairwiseParam, cols []int) [][]any {
	if len(cols) == 0 {
		return [][]any{{}}
	}
	var result [][]any
	first := params[cols[0]].Values
	rest := enumValueCombinations(params, cols[1:])
	for _, v := range first {
		for _, r := range rest {
			row := make([]any, len(cols))
			row[0] = v
			copy(row[1:], r)
			result = append(result, row)
		}
	}
	return result
}


// cartesianProductParams returns the full Cartesian product of the given parameters.
func cartesianProductParams(params []PairwiseParam) [][]any {
	if len(params) == 0 {
		return nil
	}
	rows := [][]any{{}}
	for _, p := range params {
		var next [][]any
		for _, row := range rows {
			for _, v := range p.Values {
				newRow := make([]any, len(row)+1)
				copy(newRow, row)
				newRow[len(row)] = v
				next = append(next, newRow)
			}
		}
		rows = next
	}
	return rows
}

type PairwiseTechnique struct {
	tupleLevel int // 2 = pairwise (default), 3 = 3-way, etc.
}

func NewPairwiseTechnique() *PairwiseTechnique { return &PairwiseTechnique{tupleLevel: 2} }

// NewPairwiseTechniqueWithLevel creates a pairwise technique with explicit N-way level.
func NewPairwiseTechniqueWithLevel(level int) *PairwiseTechnique {
	if level < 2 {
		level = 2
	}
	return &PairwiseTechnique{tupleLevel: level}
}

func (t *PairwiseTechnique) Name() string { return "pairwise" }

// Applies when there are 4 or more independent (enum or boolean) parameters for combinatorial coverage.
func (t *PairwiseTechnique) Applies(op *spec.Operation) bool {
	return countIndependentParams(op) >= 4
}

func (t *PairwiseTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	params := extractPairwiseParams(op)
	if len(params) < 4 {
		return nil, nil
	}
	rows := IPOGt(params, t.tupleLevel)
	rows = filterConstrainedCombinations(rows, params, op)
	var cases []schema.TestCase
	for i, row := range rows {
		queryParams := map[string]any{}
		for j, p := range params {
			queryParams[p.Name] = row[j]
		}
		tc := buildTestCase(op, nil,
			fmt.Sprintf("pairwise combination %d", i+1),
			fmt.Sprintf("%s %s", op.Method, op.Path))
		tc.Priority = "P2"
		tc.Source = schema.CaseSource{
			Technique: "pairwise",
			SpecPath:  fmt.Sprintf("%s %s parameters", op.Method, op.Path),
			Rationale: fmt.Sprintf("IPOG pairwise row %d: covers all parameter pairs", i+1),
		}
		// Add query params to path in the step
		tc.Steps[0].Path = buildPathWithQuery(op.Path, queryParams)
		cases = append(cases, tc)
	}
	return cases, nil
}

func countIndependentParams(op *spec.Operation) int {
	count := 0
	for _, p := range op.Parameters {
		if p.Schema != nil && (len(p.Schema.Enum) > 0 || p.Schema.Type == "boolean") {
			count++
		}
	}
	return count
}

func extractPairwiseParams(op *spec.Operation) []PairwiseParam {
	var params []PairwiseParam
	for _, p := range op.Parameters {
		if p.Schema == nil {
			continue
		}
		if len(p.Schema.Enum) > 0 {
			params = append(params, PairwiseParam{Name: p.Name, Values: p.Schema.Enum})
		} else if p.Schema.Type == "boolean" {
			params = append(params, PairwiseParam{Name: p.Name, Values: []any{true, false}})
		}
	}
	return params
}

// filterConstrainedCombinations removes rows where a dependent parameter has an
// active (non-zero/non-false) value while its controlling parameter is disabled.
// If filtering would remove ALL rows, the original rows are returned unchanged.
func filterConstrainedCombinations(rows [][]any, params []PairwiseParam, op *spec.Operation) [][]any {
	groups := detectParamGroups(op)
	if len(groups) == 0 {
		return rows
	}

	// Build column index: param name → column index
	colIdx := make(map[string]int, len(params))
	for i, p := range params {
		colIdx[p.Name] = i
	}

	var filtered [][]any
	for _, row := range rows {
		feasible := true
		for _, g := range groups {
			ctrlCol, ok := colIdx[g.controller]
			if !ok {
				continue
			}
			if isDisabled(row[ctrlCol]) {
				for _, dep := range g.controlled {
					depCol, ok := colIdx[dep]
					if !ok {
						continue
					}
					if isActiveValue(row[depCol]) {
						feasible = false
						break
					}
				}
			}
			if !feasible {
				break
			}
		}
		if feasible {
			filtered = append(filtered, row)
		}
	}

	if len(filtered) == 0 {
		return rows // never filter everything out
	}
	return filtered
}

// isDisabled returns true for values that represent "off": false, "false", "", 0, "none", "no".
func isDisabled(v any) bool {
	switch val := v.(type) {
	case bool:
		return !val
	case string:
		lower := strings.ToLower(val)
		return lower == "false" || lower == "" || lower == "none" || lower == "no"
	case int, int64, float64:
		return fmt.Sprintf("%v", val) == "0"
	}
	return false
}

// isActiveValue returns true when a value is non-nil, non-false, non-empty.
func isActiveValue(v any) bool {
	if v == nil {
		return false
	}
	return !isDisabled(v)
}

func buildPathWithQuery(path string, params map[string]any) string {
	if len(params) == 0 {
		return path
	}
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	q := ""
	for _, k := range keys {
		if q != "" {
			q += "&"
		}
		q += fmt.Sprintf("%s=%s", k, url.QueryEscape(fmt.Sprintf("%v", params[k])))
	}
	return path + "?" + q
}
