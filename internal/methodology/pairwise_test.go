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

func TestBuildPathWithQueryURLEncodesValues(t *testing.T) {
	path := buildPathWithQuery("/search", map[string]any{
		"q":      "hello world",
		"filter": "a+b",
	})
	assert.Contains(t, path, "q=hello+world", "spaces should be encoded as +")
	assert.Contains(t, path, "filter=a%2Bb", "literal + should be percent-encoded")
	assert.True(t, strings.HasPrefix(path, "/search?"))
}
