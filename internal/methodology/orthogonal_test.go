// internal/methodology/orthogonal_test.go
package methodology

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestOrthogonalArrayApplies_TrueForThreeParams(t *testing.T) {
	tech := NewOrthogonalArrayTechnique()
	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "a", Schema: &spec.Schema{Enum: []any{"x", "y"}}},
			{Name: "b", Schema: &spec.Schema{Enum: []any{"1", "2"}}},
			{Name: "c", Schema: &spec.Schema{Type: "boolean"}},
		},
	}
	assert.True(t, tech.Applies(op))
}

func TestOrthogonalArrayApplies_FalseForTwoParams(t *testing.T) {
	tech := NewOrthogonalArrayTechnique()
	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "a", Schema: &spec.Schema{Enum: []any{"x", "y"}}},
			{Name: "b", Schema: &spec.Schema{Type: "boolean"}},
		},
	}
	assert.False(t, tech.Applies(op))
}

func TestOrthogonalArrayApplies_FalseForNoDiscreteParams(t *testing.T) {
	tech := NewOrthogonalArrayTechnique()
	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "name", Schema: &spec.Schema{Type: "string"}},
		},
	}
	assert.False(t, tech.Applies(op))
}

func TestOrthogonalArray_L4_ThreeParams(t *testing.T) {
	// 3 two-level params → L4 → 4 test cases
	tech := NewOrthogonalArrayTechnique()
	op := &spec.Operation{
		Method: "GET",
		Path:   "/search",
		Parameters: []*spec.Parameter{
			{Name: "status", Schema: &spec.Schema{Enum: []any{"active", "inactive"}}},
			{Name: "role", Schema: &spec.Schema{Enum: []any{"admin", "user"}}},
			{Name: "verified", Schema: &spec.Schema{Type: "boolean"}},
		},
		Responses: map[string]*spec.Response{"200": {Description: "OK"}},
	}
	cases, err := tech.Generate(op)
	require.NoError(t, err)
	require.Len(t, cases, 4, "L4 produces 4 rows for 3 factors")
}

func TestOrthogonalArray_L8_SevenParams(t *testing.T) {
	// 7 two-level params → L8 → 8 test cases
	tech := NewOrthogonalArrayTechnique()
	params := make([]*spec.Parameter, 7)
	for i := range params {
		params[i] = &spec.Parameter{
			Name:   fmt.Sprintf("p%d", i),
			Schema: &spec.Schema{Type: "boolean"},
		}
	}
	op := &spec.Operation{
		Method:     "GET",
		Path:       "/items",
		Parameters: params,
		Responses:  map[string]*spec.Response{"200": {Description: "OK"}},
	}
	cases, err := tech.Generate(op)
	require.NoError(t, err)
	require.Len(t, cases, 8, "L8 produces 8 rows for 7 factors")
}

func TestOrthogonalArray_SetsSourceTechnique(t *testing.T) {
	tech := NewOrthogonalArrayTechnique()
	op := &spec.Operation{
		Method: "GET",
		Path:   "/items",
		Parameters: []*spec.Parameter{
			{Name: "a", Schema: &spec.Schema{Enum: []any{"x", "y"}}},
			{Name: "b", Schema: &spec.Schema{Enum: []any{"1", "2"}}},
			{Name: "c", Schema: &spec.Schema{Type: "boolean"}},
		},
		Responses: map[string]*spec.Response{"200": {Description: "OK"}},
	}
	cases, err := tech.Generate(op)
	require.NoError(t, err)
	for _, tc := range cases {
		assert.Equal(t, "orthogonal_array", tc.Source.Technique)
		assert.Equal(t, "P2", tc.Priority)
	}
}

func TestOrthogonalArray_BalancedCoverage_L4(t *testing.T) {
	// L4 property: each factor value appears exactly twice in 4 rows (2 levels, 4 rows each).
	tech := NewOrthogonalArrayTechnique()
	op := &spec.Operation{
		Method: "GET",
		Path:   "/items",
		Parameters: []*spec.Parameter{
			{Name: "status", Schema: &spec.Schema{Enum: []any{"on", "off"}}},
			{Name: "mode", Schema: &spec.Schema{Enum: []any{"fast", "slow"}}},
			{Name: "debug", Schema: &spec.Schema{Type: "boolean"}},
		},
		Responses: map[string]*spec.Response{"200": {Description: "OK"}},
	}
	cases, err := tech.Generate(op)
	require.NoError(t, err)
	require.Len(t, cases, 4)

	// Count occurrences of each value in the paths.
	counts := map[string]int{}
	for _, tc := range cases {
		path := tc.Steps[0].Path
		for _, v := range []string{"on", "off", "fast", "slow", "true", "false"} {
			if containsParamValue(path, v) {
				counts[v]++
			}
		}
	}
	// Each value should appear exactly 2 times (balanced L4 property).
	for _, v := range []string{"on", "off", "fast", "slow", "true", "false"} {
		assert.Equal(t, 2, counts[v], "L4: value %q must appear exactly 2 times", v)
	}
}

func TestSelectOAArray(t *testing.T) {
	cases := []struct {
		n    int
		rows int
		cols int
		name string
	}{
		{3, 4, 3, "L4(2^3)"},
		{4, 8, 7, "L8(2^7)"},
		{7, 8, 7, "L8(2^7)"},
		{8, 27, 13, "L27(3^13)"},
		{13, 27, 13, "L27(3^13)"},
	}
	for _, tc := range cases {
		arr := selectOAArray(tc.n)
		require.NotNil(t, arr, "n=%d", tc.n)
		assert.Len(t, arr, tc.rows, "n=%d should use %d-row array", tc.n, tc.rows)
		assert.Len(t, arr[0], tc.cols, "n=%d should have %d columns", tc.n, tc.cols)
		assert.Equal(t, tc.name, arr.name(), "n=%d", tc.n)
	}
}

func TestSelectOAArray_TooManyParams_ReturnsNil(t *testing.T) {
	assert.Nil(t, selectOAArray(14))
	assert.Nil(t, selectOAArray(100))
}

func TestOrthogonalArrayApplies_FalseForFourteenParams(t *testing.T) {
	tech := NewOrthogonalArrayTechnique()
	params := make([]*spec.Parameter, 14)
	for i := range params {
		params[i] = &spec.Parameter{
			Name:   fmt.Sprintf("p%d", i),
			Schema: &spec.Schema{Type: "boolean"},
		}
	}
	op := &spec.Operation{Parameters: params}
	assert.False(t, tech.Applies(op), "14 params exceeds L27 column capacity")
}

func TestOrthogonalArray_L27_EightParams(t *testing.T) {
	// 8 boolean params → L27 → 27 test cases; paths must contain valid boolean values.
	tech := NewOrthogonalArrayTechnique()
	params := make([]*spec.Parameter, 8)
	for i := range params {
		params[i] = &spec.Parameter{
			Name:   fmt.Sprintf("flag%d", i),
			Schema: &spec.Schema{Type: "boolean"},
		}
	}
	op := &spec.Operation{
		Method:     "GET",
		Path:       "/items",
		Parameters: params,
		Responses:  map[string]*spec.Response{"200": {Description: "OK"}},
	}
	cases, err := tech.Generate(op)
	require.NoError(t, err)
	require.Len(t, cases, 27, "L27 produces 27 rows for 8–13 factors")

	// Each case path must contain only valid boolean values ("true" or "false").
	for i, tc := range cases {
		path := tc.Steps[0].Path
		assert.True(t,
			strings.Contains(path, "=true") || strings.Contains(path, "=false"),
			"case %d path %q must contain a boolean param value", i, path)
	}

	// Rationale should mention the approximately-balanced limitation.
	for _, tc := range cases {
		assert.Contains(t, tc.Source.Rationale, "approximately balanced",
			"L27 rationale must mention approximate balance for 2-level params")
	}
}

func TestExtractOAParams_EnumAndBoolean(t *testing.T) {
	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "status", Schema: &spec.Schema{Enum: []any{"a", "b", "c"}}},
			{Name: "flag", Schema: &spec.Schema{Type: "boolean"}},
			{Name: "name", Schema: &spec.Schema{Type: "string"}},           // excluded
			{Name: "single", Schema: &spec.Schema{Enum: []any{"only"}}},   // excluded (1 value)
			{Name: "big", Schema: &spec.Schema{Enum: []any{1, 2, 3, 4}}},  // coerced to 3 levels
		},
	}
	params := extractOAParams(op)
	require.Len(t, params, 3) // status, flag, big
	assert.Equal(t, "status", params[0].name)
	assert.Len(t, params[0].values, 3)
	assert.Equal(t, "flag", params[1].name)
	assert.Equal(t, []any{true, false}, params[1].values)
	assert.Equal(t, "big", params[2].name)
	assert.Len(t, params[2].values, 3, "4-value enum coerced to 3 levels")
}

func TestLevelToValue_Wraps(t *testing.T) {
	vals := []any{"a", "b"}
	assert.Equal(t, "a", levelToValue(vals, 0))
	assert.Equal(t, "b", levelToValue(vals, 1))
	assert.Equal(t, "a", levelToValue(vals, 2), "level 2 wraps to index 0 for 2-level param")
}

// containsParamValue checks if the URL path query string contains param=value.
func containsParamValue(path, value string) bool {
	return strings.Contains(path, "="+value+"&") || strings.HasSuffix(path, "="+value)
}
