// internal/methodology/pairwise.go
package methodology

import (
	"fmt"
	"net/url"
	"sort"

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

type PairwiseTechnique struct{}

func NewPairwiseTechnique() *PairwiseTechnique { return &PairwiseTechnique{} }
func (t *PairwiseTechnique) Name() string      { return "pairwise" }

// Applies when there are 4 or more independent (enum or boolean) parameters for combinatorial coverage.
func (t *PairwiseTechnique) Applies(op *spec.Operation) bool {
	return countIndependentParams(op) >= 4
}

func (t *PairwiseTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	params := extractPairwiseParams(op)
	if len(params) < 4 {
		return nil, nil
	}
	rows := IPOG(params)
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
