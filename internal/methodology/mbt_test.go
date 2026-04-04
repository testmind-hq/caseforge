// internal/methodology/mbt_test.go
package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestClassificationTreeApplies_TrueWhenEnumParam(t *testing.T) {
	tech := NewClassificationTreeTechnique()
	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "status", Schema: &spec.Schema{Enum: []any{"active", "inactive"}}},
		},
	}
	assert.True(t, tech.Applies(op))
}

func TestClassificationTreeApplies_TrueWhenBooleanParam(t *testing.T) {
	tech := NewClassificationTreeTechnique()
	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "enabled", Schema: &spec.Schema{Type: "boolean"}},
		},
	}
	assert.True(t, tech.Applies(op))
}

func TestClassificationTreeApplies_FalseWhenNoDiscreteParams(t *testing.T) {
	tech := NewClassificationTreeTechnique()
	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "name", Schema: &spec.Schema{Type: "string"}},
			{Name: "age", Schema: &spec.Schema{Type: "integer"}},
		},
	}
	assert.False(t, tech.Applies(op))
}

func TestClassificationTreeApplies_FalseWhenEnumHasSingleValue(t *testing.T) {
	tech := NewClassificationTreeTechnique()
	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "type", Schema: &spec.Schema{Enum: []any{"fixed"}}},
		},
	}
	assert.False(t, tech.Applies(op))
}

func TestClassificationTree_ECTCoverage(t *testing.T) {
	// 3 classifications: status=[a,b,c], role=[x,y], active=[true,false]
	// max leaves = 3 → 3 test cases
	// Every leaf must appear at least once.
	tech := NewClassificationTreeTechnique()
	op := &spec.Operation{
		Method: "GET",
		Path:   "/items",
		Parameters: []*spec.Parameter{
			{Name: "status", Schema: &spec.Schema{Type: "string", Enum: []any{"active", "inactive", "pending"}}},
			{Name: "role", Schema: &spec.Schema{Type: "string", Enum: []any{"admin", "user"}}},
			{Name: "verified", Schema: &spec.Schema{Type: "boolean"}},
		},
		Responses: map[string]*spec.Response{"200": {Description: "OK"}},
	}
	cases, err := tech.Generate(op)
	require.NoError(t, err)
	require.Len(t, cases, 3, "ECT row count = max leaves across classifications = 3")

	// Verify each leaf appears in at least one test case path (checked per-path to
	// avoid false positives from path segment substring matches).
	for _, leaf := range []string{"active", "inactive", "pending", "admin", "user", "true", "false"} {
		found := false
		for _, tc := range cases {
			if containsParamValue(tc.Steps[0].Path, leaf) {
				found = true
				break
			}
		}
		assert.True(t, found, "leaf %q must appear in at least one test case path", leaf)
	}
}

func TestClassificationTree_SingleClassification(t *testing.T) {
	// Only one boolean param → 2 test cases (one per leaf).
	tech := NewClassificationTreeTechnique()
	op := &spec.Operation{
		Method: "GET",
		Path:   "/users",
		Parameters: []*spec.Parameter{
			{Name: "active", Schema: &spec.Schema{Type: "boolean"}},
		},
		Responses: map[string]*spec.Response{"200": {Description: "OK"}},
	}
	cases, err := tech.Generate(op)
	require.NoError(t, err)
	require.Len(t, cases, 2)
}

func TestClassificationTree_SetsSourceTechnique(t *testing.T) {
	tech := NewClassificationTreeTechnique()
	op := &spec.Operation{
		Method: "GET",
		Path:   "/items",
		Parameters: []*spec.Parameter{
			{Name: "status", Schema: &spec.Schema{Enum: []any{"a", "b"}}},
			{Name: "flag", Schema: &spec.Schema{Type: "boolean"}},
		},
		Responses: map[string]*spec.Response{"200": {Description: "OK"}},
	}
	cases, err := tech.Generate(op)
	require.NoError(t, err)
	for _, tc := range cases {
		assert.Equal(t, "classification_tree", tc.Source.Technique)
		assert.Equal(t, "P2", tc.Priority)
	}
}

func TestClassificationTree_NoDiscreteParams_ReturnsEmpty(t *testing.T) {
	tech := NewClassificationTreeTechnique()
	op := &spec.Operation{
		Method:     "POST",
		Path:       "/users",
		Parameters: []*spec.Parameter{},
		Responses:  map[string]*spec.Response{"201": {Description: "Created"}},
	}
	cases, err := tech.Generate(op)
	require.NoError(t, err)
	assert.Empty(t, cases)
}

func TestBuildClassifications_EnumAndBoolean(t *testing.T) {
	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "status", Schema: &spec.Schema{Enum: []any{"a", "b", "c"}}},
			{Name: "flag", Schema: &spec.Schema{Type: "boolean"}},
			{Name: "name", Schema: &spec.Schema{Type: "string"}},         // excluded
			{Name: "single", Schema: &spec.Schema{Enum: []any{"only"}}}, // excluded (1 value)
		},
	}
	classes := buildClassifications(op)
	require.Len(t, classes, 2)
	assert.Equal(t, "status", classes[0].name)
	assert.Equal(t, []any{"a", "b", "c"}, classes[0].leaves)
	assert.Equal(t, "flag", classes[1].name)
	assert.Equal(t, []any{true, false}, classes[1].leaves)
}
