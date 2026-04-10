// internal/methodology/pairwise_test.go
package methodology

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestPairwiseAppliesForFourOrMoreParams(t *testing.T) {
	tech := &PairwiseTechnique{}
	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "a", Schema: &spec.Schema{Enum: []any{"x", "y"}}},
			{Name: "b", Schema: &spec.Schema{Enum: []any{"1", "2"}}},
			{Name: "c", Schema: &spec.Schema{Type: "boolean"}},
			{Name: "d", Schema: &spec.Schema{Enum: []any{"p", "q"}}},
		},
	}
	assert.True(t, tech.Applies(op))
}

func TestPairwiseDoesNotApplyForThreeParams(t *testing.T) {
	tech := &PairwiseTechnique{}
	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "a", Schema: &spec.Schema{Enum: []any{"x", "y"}}},
			{Name: "b", Schema: &spec.Schema{Enum: []any{"1", "2"}}},
			{Name: "c", Schema: &spec.Schema{Type: "boolean"}},
		},
	}
	assert.False(t, tech.Applies(op))
}

func TestIPOGCoverageProperty(t *testing.T) {
	// 4 params, 2 values each → all pairs covered
	params := []PairwiseParam{
		{Name: "status", Values: []any{"active", "inactive"}},
		{Name: "role",   Values: []any{"admin", "user"}},
		{Name: "plan",   Values: []any{"free", "paid"}},
		{Name: "region", Values: []any{"us", "eu"}},
	}
	rows := IPOG(params)
	// Verify every pair is covered
	covered := make(map[string]bool)
	for _, row := range rows {
		for i := 0; i < len(params); i++ {
			for j := i + 1; j < len(params); j++ {
				key := fmt.Sprintf("%s=%v|%s=%v",
					params[i].Name, row[i],
					params[j].Name, row[j])
				covered[key] = true
			}
		}
	}
	// Total pairs = C(4,2) * 2*2 = 6 * 4 = 24
	assert.Equal(t, 24, len(covered), "all 24 pairs should be covered")
	// Rows should be fewer than full factorial (16)
	assert.Less(t, len(rows), 16)
}

func TestPairwiseGeneratesTestCases(t *testing.T) {
	tech := NewPairwiseTechnique()
	op := &spec.Operation{
		OperationID: "searchItems",
		Method:      "GET",
		Path:        "/items",
		Parameters: []*spec.Parameter{
			{Name: "status", Schema: &spec.Schema{Type: "string", Enum: []any{"active", "inactive"}}},
			{Name: "role",   Schema: &spec.Schema{Type: "string", Enum: []any{"admin", "user"}}},
			{Name: "plan",   Schema: &spec.Schema{Type: "string", Enum: []any{"free", "paid"}}},
			{Name: "region", Schema: &spec.Schema{Type: "string", Enum: []any{"us", "eu"}}},
		},
		Responses: map[string]*spec.Response{"200": {Description: "OK"}},
	}
	cases, err := tech.Generate(op)
	require.NoError(t, err)
	assert.NotEmpty(t, cases)
	for _, tc := range cases {
		assert.Equal(t, "pairwise", tc.Source.Technique)
	}
}

func TestIPOGt_Level3_CoversAllTriples(t *testing.T) {
	params := []PairwiseParam{
		{Name: "a", Values: []any{"a1", "a2"}},
		{Name: "b", Values: []any{"b1", "b2"}},
		{Name: "c", Values: []any{"c1", "c2"}},
		{Name: "d", Values: []any{"d1", "d2"}},
	}

	rows := IPOGt(params, 3)
	require.NotEmpty(t, rows)

	// Verify all triples are covered
	covered := make(map[[6]any]bool) // (colA, colB, colC, valA, valB, valC)

	for _, row := range rows {
		for i := 0; i < len(params)-2; i++ {
			for j := i + 1; j < len(params)-1; j++ {
				for k := j + 1; k < len(params); k++ {
					key := [6]any{i, j, k, row[i], row[j], row[k]}
					covered[key] = true
				}
			}
		}
	}

	// Check all expected triples are covered
	for i := 0; i < len(params)-2; i++ {
		for j := i + 1; j < len(params)-1; j++ {
			for k := j + 1; k < len(params); k++ {
				for _, va := range params[i].Values {
					for _, vb := range params[j].Values {
						for _, vc := range params[k].Values {
							key := [6]any{i, j, k, va, vb, vc}
							assert.True(t, covered[key],
								"triple (%d=%v, %d=%v, %d=%v) not covered",
								i, va, j, vb, k, vc)
						}
					}
				}
			}
		}
	}
}

func TestIPOGt_Level2_MatchesIPOG(t *testing.T) {
	params := []PairwiseParam{
		{Name: "x", Values: []any{1, 2}},
		{Name: "y", Values: []any{"a", "b"}},
		{Name: "z", Values: []any{true, false}},
	}
	rows2 := IPOG(params)
	rowsT := IPOGt(params, 2)
	// Both should cover the same pairs; row counts may differ slightly
	assert.Equal(t, len(rows2), len(rowsT),
		"IPOGt(t=2) should produce same number of rows as IPOG")
}

func TestPairwiseTechnique_TupleLevel3_GeneratesMoreCases(t *testing.T) {
	op := &spec.Operation{
		Method: "GET", Path: "/search",
		Parameters: []*spec.Parameter{
			{Name: "sort", In: "query", Schema: &spec.Schema{Type: "boolean"}},
			{Name: "filter", In: "query", Schema: &spec.Schema{Enum: []any{"all", "active"}}},
			{Name: "format", In: "query", Schema: &spec.Schema{Enum: []any{"json", "csv"}}},
			{Name: "lang", In: "query", Schema: &spec.Schema{Enum: []any{"en", "zh"}}},
		},
		Responses: map[string]*spec.Response{"200": {}},
	}

	tech2 := NewPairwiseTechniqueWithLevel(2)
	tech3 := NewPairwiseTechniqueWithLevel(3)

	cases2, err := tech2.Generate(op)
	require.NoError(t, err)
	cases3, err := tech3.Generate(op)
	require.NoError(t, err)

	// 3-way must produce at least as many cases as 2-way
	assert.GreaterOrEqual(t, len(cases3), len(cases2))

	// All cases must cite correct technique
	for _, tc := range cases3 {
		assert.Equal(t, "pairwise", tc.Source.Technique)
	}
}

func TestBuildPathWithQueryURLEncodesValues(t *testing.T) {
	path := buildPathWithQuery("/search", map[string]any{
		"q":      "hello world",
		"filter": "a+b",
	})
	assert.Contains(t, path, "q=hello+world", "spaces should be encoded as +")
	assert.Contains(t, path, "filter=a%2Bb", "literal + should be percent-encoded")
	assert.True(t, strings.HasPrefix(path, "/search?"))
}
