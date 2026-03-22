// internal/methodology/pairwise.go
package methodology

import (
	"fmt"

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
		// Track which pairs (existing col, new col) are still uncovered
		uncovered := make(map[[2]any]bool)
		for prevCol := 0; prevCol < col; prevCol++ {
			for _, v1 := range params[prevCol].Values {
				for _, v2 := range param.Values {
					uncovered[[2]any{v1, v2}] = true
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
					key := [2]any{rows[i][prevCol], candidate}
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
				delete(uncovered, [2]any{rows[i][prevCol], bestVal})
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
				// pair[0] is the value from some prev col, pair[1] is the new col value
				// Find which prev col pair[0] belongs to
				for prevCol := 0; prevCol < col; prevCol++ {
					for _, v := range params[prevCol].Values {
						if v == pair[0] {
							row[prevCol] = pair[0]
							row[col] = pair[1]
							goto assigned
						}
					}
				}
			assigned:
				break
			}
			rows = append(rows, row)
			// Mark covered pairs for this row
			for prevCol := 0; prevCol < col; prevCol++ {
				delete(uncovered, [2]any{row[prevCol], row[col]})
			}
		}
	}

	return rows
}

type PairwiseTechnique struct{}

func NewPairwiseTechnique() *PairwiseTechnique { return &PairwiseTechnique{} }
func (t *PairwiseTechnique) Name() string      { return "pairwise" }

// Applies when there are 4 or more parameters (independent parameters for combinatorial coverage).
func (t *PairwiseTechnique) Applies(op *spec.Operation) bool {
	return len(op.Parameters) >= 4
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
		tc := buildTestCase(op, nil, fmt.Sprintf("combo_%d", i),
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
	q := ""
	for k, v := range params {
		if q != "" {
			q += "&"
		}
		q += fmt.Sprintf("%s=%v", k, v)
	}
	return path + "?" + q
}
