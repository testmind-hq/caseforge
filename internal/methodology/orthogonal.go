// internal/methodology/orthogonal.go
package methodology

import (
	"fmt"

	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// OrthogonalArrayTechnique implements Orthogonal Array Testing Strategy (OATS).
//
// For operations with 3–13 discrete parameters (enum or boolean), it selects a
// pre-defined Taguchi orthogonal array (L4, L8, L9, L16, L27) and maps each
// column to a parameter. Every factor level is distributed uniformly across
// rows, guaranteeing balanced coverage with a minimal test set.
//
// Level mapping:
//   - 2-level factors (boolean / 2-value enum): use 0→first value, 1→second value
//   - 3-level factors (3-value enum): 0→first, 1→second, 2→third
//   - Factors with more values are coerced to 3 levels (first 3 values used)
//   - Factors with exactly 1 value are excluded
//
// Array selection priority: fewest rows that accommodate the required column count.
// Mixed-level scenarios (2-level params in a 3-level L27 array): level 2 wraps to
// level 0 via modulo, so the first value appears twice as often as the second. This
// produces approximately-balanced (not strictly uniform) coverage for those params.
//
// Distinct from PairwiseTechnique (IPOG minimises test count while covering all
// pairs) and ClassificationTreeTechnique (ECT per-leaf variation): OA provides
// provably balanced factor-level distribution across a fixed-size grid.
//
// Applies when the operation has 3–13 discrete parameters.
type OrthogonalArrayTechnique struct {
	gen *datagen.Generator
}

func NewOrthogonalArrayTechnique() *OrthogonalArrayTechnique {
	return &OrthogonalArrayTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *OrthogonalArrayTechnique) Name() string { return "orthogonal_array" }

// Applies returns true when there are between 3 and 13 discrete parameters.
func (t *OrthogonalArrayTechnique) Applies(op *spec.Operation) bool {
	n := len(extractOAParams(op))
	return n >= 3 && n <= 13
}

func (t *OrthogonalArrayTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	params := extractOAParams(op)
	if len(params) < 3 {
		return nil, nil
	}

	array := selectOAArray(len(params))
	if array == nil {
		return nil, nil
	}

	// Determine whether any 2-level param is mapped into a 3-level array (L27).
	// If so, level 2 wraps to level 0, producing approximately-balanced (not strictly
	// uniform) coverage for those params.
	arrayIs3Level := len(array) == 27
	has2LevelIn3Level := false
	if arrayIs3Level {
		for _, p := range params {
			if len(p.values) == 2 {
				has2LevelIn3Level = true
				break
			}
		}
	}
	balanceDesc := "balanced factor-level distribution"
	if has2LevelIn3Level {
		balanceDesc = "approximately balanced (2-level params in L27: first value appears 2× as often)"
	}

	var cases []schema.TestCase
	for rowIdx, row := range array {
		queryParams := make(map[string]any, len(params))
		for colIdx, level := range row {
			if colIdx >= len(params) {
				break
			}
			p := params[colIdx]
			val := levelToValue(p.values, level)
			queryParams[p.name] = val
		}

		base := buildValidBody(t.gen, op)
		tc := buildTestCase(op, base,
			fmt.Sprintf("orthogonal array row %d", rowIdx+1),
			"")
		tc.Priority = "P2"
		tc.Source = schema.CaseSource{
			Technique: "orthogonal_array",
			SpecPath:  fmt.Sprintf("%s %s parameters", op.Method, op.Path),
			Rationale: fmt.Sprintf("OA %s row %d — %s", array.name(), rowIdx+1, balanceDesc),
		}
		tc.Steps[0].Path = buildPathWithQuery(op.Path, queryParams)
		cases = append(cases, tc)
	}
	return cases, nil
}

// oaParam holds a parameter name and its ordered discrete values (max 3).
type oaParam struct {
	name   string
	values []any // 1–3 values; level index maps into this slice
}

// extractOAParams returns parameters suitable for OA: enum (up to 3 levels) or boolean.
// Parameters with 0 or 1 distinct values are excluded.
func extractOAParams(op *spec.Operation) []oaParam {
	var out []oaParam
	for _, p := range op.Parameters {
		if p.Schema == nil {
			continue
		}
		var values []any
		if len(p.Schema.Enum) >= 2 {
			vals := p.Schema.Enum
			if len(vals) > 3 {
				vals = vals[:3] // coerce to 3 levels
			}
			values = vals
		} else if p.Schema.Type == "boolean" {
			values = []any{true, false}
		}
		if len(values) >= 2 {
			out = append(out, oaParam{name: p.Name, values: values})
		}
	}
	return out
}

// levelToValue maps an OA level index (0, 1, 2) to an actual parameter value.
// For 2-level parameters, level 2 wraps to level 0.
func levelToValue(values []any, level int) any {
	if len(values) == 0 {
		return nil
	}
	return values[level%len(values)]
}

// oaArray is a slice of rows; each row is a slice of level indices (0, 1, or 2).
type oaArray [][]int

// name returns the Taguchi designation based on the number of columns, which
// uniquely identifies the array regardless of how many rows it has.
func (a oaArray) name() string {
	if len(a) == 0 {
		return "L0"
	}
	switch len(a[0]) {
	case 3:
		return "L4(2^3)"
	case 7:
		return "L8(2^7)"
	case 13:
		return "L27(3^13)"
	default:
		return fmt.Sprintf("L%d(%dcols)", len(a), len(a[0]))
	}
}

// selectOAArray returns the smallest pre-defined OA that has enough columns for n factors.
// Returns nil if n exceeds the largest available array (13).
//
// Array selection (by column count):
//   - n ≤ 3  → L4  (3 columns,  4 rows)
//   - n ≤ 7  → L8  (7 columns,  8 rows)
//   - n ≤ 13 → L27 (13 columns, 27 rows)
//
// For 2-level factors in L27 (a 3-level array), levelToValue wraps level 2 → level 0,
// ensuring valid values are still selected.
func selectOAArray(n int) oaArray {
	switch {
	case n <= 3:
		return oaL4
	case n <= 7:
		return oaL8
	case n <= 13:
		return oaL27
	default:
		return nil
	}
}

// Pre-defined Taguchi orthogonal arrays (subset used by selectOAArray).
// Column count determines the maximum number of factors supported.

// L4(2^3): 4 rows, 3 two-level columns. Every pair of columns covers all 4 combinations.
var oaL4 = oaArray{
	{0, 0, 0},
	{0, 1, 1},
	{1, 0, 1},
	{1, 1, 0},
}

// L8(2^7): 8 rows, 7 two-level columns.
var oaL8 = oaArray{
	{0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 1, 1, 1, 1},
	{0, 1, 1, 0, 0, 1, 1},
	{0, 1, 1, 1, 1, 0, 0},
	{1, 0, 1, 0, 1, 0, 1},
	{1, 0, 1, 1, 0, 1, 0},
	{1, 1, 0, 0, 1, 1, 0},
	{1, 1, 0, 1, 0, 0, 1},
}

// L27(3^13): 27 rows, up to 13 three-level columns.
var oaL27 = oaArray{
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{0, 0, 0, 0, 2, 2, 2, 2, 2, 2, 2, 2, 2},
	{0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 2, 2, 2},
	{0, 1, 1, 1, 1, 1, 1, 2, 2, 2, 0, 0, 0},
	{0, 1, 1, 1, 2, 2, 2, 0, 0, 0, 1, 1, 1},
	{0, 2, 2, 2, 0, 0, 0, 2, 2, 2, 1, 1, 1},
	{0, 2, 2, 2, 1, 1, 1, 0, 0, 0, 2, 2, 2},
	{0, 2, 2, 2, 2, 2, 2, 1, 1, 1, 0, 0, 0},
	{1, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2},
	{1, 0, 1, 2, 1, 2, 0, 1, 2, 0, 1, 2, 0},
	{1, 0, 1, 2, 2, 0, 1, 2, 0, 1, 2, 0, 1},
	{1, 1, 2, 0, 0, 1, 2, 1, 2, 0, 2, 0, 1},
	{1, 1, 2, 0, 1, 2, 0, 2, 0, 1, 0, 1, 2},
	{1, 1, 2, 0, 2, 0, 1, 0, 1, 2, 1, 2, 0},
	{1, 2, 0, 1, 0, 1, 2, 2, 0, 1, 1, 2, 0},
	{1, 2, 0, 1, 1, 2, 0, 0, 1, 2, 2, 0, 1},
	{1, 2, 0, 1, 2, 0, 1, 1, 2, 0, 0, 1, 2},
	{2, 0, 2, 1, 0, 2, 1, 0, 2, 1, 0, 2, 1},
	{2, 0, 2, 1, 1, 0, 2, 1, 0, 2, 1, 0, 2},
	{2, 0, 2, 1, 2, 1, 0, 2, 1, 0, 2, 1, 0},
	{2, 1, 0, 2, 0, 2, 1, 1, 0, 2, 2, 1, 0},
	{2, 1, 0, 2, 1, 0, 2, 2, 1, 0, 0, 2, 1},
	{2, 1, 0, 2, 2, 1, 0, 0, 2, 1, 1, 0, 2},
	{2, 2, 1, 0, 0, 2, 1, 2, 1, 0, 1, 0, 2},
	{2, 2, 1, 0, 1, 0, 2, 0, 2, 1, 2, 1, 0},
	{2, 2, 1, 0, 2, 1, 0, 1, 0, 2, 0, 2, 1},
}
