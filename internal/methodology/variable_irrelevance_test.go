// internal/methodology/variable_irrelevance_test.go
package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestVariableIrrelevance_DetectsParamGroups(t *testing.T) {
	op := &spec.Operation{
		Method: "GET", Path: "/items",
		Parameters: []*spec.Parameter{
			{Name: "sort", In: "query", Schema: &spec.Schema{Type: "boolean"}},
			{Name: "sort_field", In: "query", Schema: &spec.Schema{Type: "string"}},
			{Name: "sort_order", In: "query", Schema: &spec.Schema{Enum: []any{"asc", "desc"}}},
			{Name: "q", In: "query", Schema: &spec.Schema{Type: "string"}},
		},
		Responses: map[string]*spec.Response{"200": {}},
	}

	groups := detectParamGroups(op)
	require.Len(t, groups, 1, "expected one dependency group (sort → sort_field, sort_order)")
	assert.Equal(t, "sort", groups[0].controller)
	assert.ElementsMatch(t, []string{"sort_field", "sort_order"}, groups[0].controlled)
}

func TestVariableIrrelevance_GeneratesCasesForGroups(t *testing.T) {
	op := &spec.Operation{
		Method: "GET", Path: "/items",
		Parameters: []*spec.Parameter{
			{Name: "sort", In: "query", Schema: &spec.Schema{Type: "boolean"}},
			{Name: "sort_field", In: "query", Schema: &spec.Schema{Type: "string", Enum: []any{"name", "date"}}},
			{Name: "sort_order", In: "query", Schema: &spec.Schema{Enum: []any{"asc", "desc"}}},
			{Name: "limit", In: "query", Schema: &spec.Schema{Type: "integer"}},
		},
		Responses: map[string]*spec.Response{"200": {}},
	}

	tech := NewVariableIrrelevanceTechnique()
	require.True(t, tech.Applies(op))

	cases, err := tech.Generate(op)
	require.NoError(t, err)
	require.NotEmpty(t, cases)

	for _, tc := range cases {
		assert.Equal(t, "variable_irrelevance", tc.Source.Technique)
		assert.NotEmpty(t, tc.Source.Rationale)
		// Case path must contain sort=false (controller disabled)
		assert.Contains(t, tc.Steps[0].Path, "sort=false",
			"expected controller param to be disabled in path")
		// Dependent params should NOT appear in the path
		assert.NotContains(t, tc.Steps[0].Path, "sort_field",
			"dependent param sort_field must be absent when sort=false")
		assert.NotContains(t, tc.Steps[0].Path, "sort_order",
			"dependent param sort_order must be absent when sort=false")
	}
}

func TestVariableIrrelevance_DoesNotApplyWithoutDependencies(t *testing.T) {
	op := &spec.Operation{
		Method: "GET", Path: "/items",
		Parameters: []*spec.Parameter{
			{Name: "q", In: "query", Schema: &spec.Schema{Type: "string"}},
			{Name: "limit", In: "query", Schema: &spec.Schema{Type: "integer"}},
		},
		Responses: map[string]*spec.Response{"200": {}},
	}
	tech := NewVariableIrrelevanceTechnique()
	assert.False(t, tech.Applies(op))
}

func TestVariableIrrelevance_MultipleGroups(t *testing.T) {
	op := &spec.Operation{
		Method: "GET", Path: "/items",
		Parameters: []*spec.Parameter{
			{Name: "sort", In: "query", Schema: &spec.Schema{Type: "boolean"}},
			{Name: "sort_field", In: "query", Schema: &spec.Schema{Type: "string"}},
			{Name: "filter", In: "query", Schema: &spec.Schema{Type: "boolean"}},
			{Name: "filter_type", In: "query", Schema: &spec.Schema{Type: "string"}},
		},
		Responses: map[string]*spec.Response{"200": {}},
	}

	groups := detectParamGroups(op)
	assert.Len(t, groups, 2, "expected two groups: sort→sort_field and filter→filter_type")
}
