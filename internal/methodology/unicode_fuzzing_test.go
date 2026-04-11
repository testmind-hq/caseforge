// internal/methodology/unicode_fuzzing_test.go
package methodology

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// makeStringOp builds a minimal POST operation with a JSON request body containing
// the given properties.
func makeStringOp(props map[string]*spec.Schema) *spec.Operation {
	return &spec.Operation{
		Method: "POST", Path: "/items",
		RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
			"application/json": {Schema: &spec.Schema{
				Type:       "object",
				Properties: props,
			}},
		}},
		Responses: map[string]*spec.Response{"400": {Description: "bad request"}},
	}
}

func TestUnicodeFuzzingTechnique_Applies_True(t *testing.T) {
	op := makeStringOp(map[string]*spec.Schema{
		"name": {Type: "string"},
	})
	assert.True(t, NewUnicodeFuzzingTechnique().Applies(op))
}

func TestUnicodeFuzzingTechnique_Applies_False(t *testing.T) {
	op := &spec.Operation{Method: "GET", Path: "/items"}
	assert.False(t, NewUnicodeFuzzingTechnique().Applies(op))
}

func TestUnicodeFuzzingTechnique_Applies_NoStringFields(t *testing.T) {
	op := makeStringOp(map[string]*spec.Schema{
		"count": {Type: "integer"},
		"price": {Type: "number"},
	})
	assert.False(t, NewUnicodeFuzzingTechnique().Applies(op))
}

func TestUnicodeFuzzingTechnique_Generate_ProducesExactly5PerStringField(t *testing.T) {
	op := makeStringOp(map[string]*spec.Schema{
		"name":  {Type: "string"},
		"email": {Type: "string"},
	})
	cases, err := NewUnicodeFuzzingTechnique().Generate(op)
	require.NoError(t, err)
	// 2 string fields × 5 mutations = 10 cases
	assert.Len(t, cases, 10)
}

func TestUnicodeFuzzingTechnique_Generate_SkipsNonStringFields(t *testing.T) {
	op := makeStringOp(map[string]*spec.Schema{
		"name":  {Type: "string"},
		"count": {Type: "integer"},
	})
	cases, err := NewUnicodeFuzzingTechnique().Generate(op)
	require.NoError(t, err)
	// Only string field mutated: 5 cases
	assert.Len(t, cases, 5)
	for _, c := range cases {
		// Ensure no case title references the integer field
		assert.Contains(t, c.Title, "name")
	}
}

func TestUnicodeFuzzingTechnique_Generate_AllHaveP3Priority(t *testing.T) {
	op := makeStringOp(map[string]*spec.Schema{
		"name": {Type: "string"},
	})
	cases, err := NewUnicodeFuzzingTechnique().Generate(op)
	require.NoError(t, err)
	assert.NotEmpty(t, cases)
	for _, c := range cases {
		assert.Equal(t, "P3", c.Priority)
	}
}

func TestUnicodeFuzzingTechnique_Generate_SourceScenario(t *testing.T) {
	op := makeStringOp(map[string]*spec.Schema{
		"name": {Type: "string"},
	})
	cases, err := NewUnicodeFuzzingTechnique().Generate(op)
	require.NoError(t, err)
	assert.NotEmpty(t, cases)
	for _, c := range cases {
		assert.Equal(t, "UNICODE_INJECTION", c.Source.Scenario)
		assert.Equal(t, "unicode_fuzzing", c.Source.Technique)
	}
}

func TestUnicodeFuzzingTechnique_Generate_ContainsControlChar(t *testing.T) {
	op := makeStringOp(map[string]*spec.Schema{
		"name": {Type: "string"},
	})
	cases, err := NewUnicodeFuzzingTechnique().Generate(op)
	require.NoError(t, err)

	found := false
	for _, c := range cases {
		require.Len(t, c.Steps, 1)
		body, ok := c.Steps[0].Body.(map[string]any)
		require.True(t, ok)
		if val, exists := body["name"]; exists {
			if s, ok := val.(string); ok {
				if strings.Contains(s, "\u0000") {
					found = true
					break
				}
			}
		}
	}
	assert.True(t, found, "expected at least one case body field to contain \\u0000")
}
