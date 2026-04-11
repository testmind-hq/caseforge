// internal/methodology/constraint_mutation_test.go
package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestConstraintMutationTechnique_Name(t *testing.T) {
	assert.Equal(t, "constraint_mutation", NewConstraintMutationTechnique().Name())
}

func TestConstraintMutationTechnique_Applies_WithProperties(t *testing.T) {
	op := &spec.Operation{
		Method: "POST", Path: "/pets",
		RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
			"application/json": {Schema: &spec.Schema{
				Type:       "object",
				Properties: map[string]*spec.Schema{"name": {Type: "string"}},
			}},
		}},
	}
	assert.True(t, NewConstraintMutationTechnique().Applies(op))
}

func TestConstraintMutationTechnique_Applies_NoBody(t *testing.T) {
	assert.False(t, NewConstraintMutationTechnique().Applies(&spec.Operation{Method: "GET", Path: "/pets"}))
}

func TestConstraintMutationTechnique_Applies_EmptyProperties(t *testing.T) {
	op := &spec.Operation{
		Method: "POST", Path: "/empty",
		RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
			"application/json": {Schema: &spec.Schema{Type: "object", Properties: map[string]*spec.Schema{}}},
		}},
	}
	assert.False(t, NewConstraintMutationTechnique().Applies(op))
}

func TestConstraintMutationTechnique_Generate_NullInjection(t *testing.T) {
	op := &spec.Operation{
		Method: "POST", Path: "/pets",
		RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
			"application/json": {Schema: &spec.Schema{
				Type: "object",
				Properties: map[string]*spec.Schema{
					"name": {Type: "string"},                 // not nullable → should produce null injection case
					"tag":  {Type: "string", Nullable: true}, // nullable → skip
				},
			}},
		}},
		Responses: map[string]*spec.Response{"422": {Description: "invalid"}},
	}
	cases, err := NewConstraintMutationTechnique().Generate(op)
	require.NoError(t, err)
	// 1 null-injection (name only, tag is nullable) + 1 wrong-content-type
	assert.Len(t, cases, 2)
	for _, c := range cases {
		assert.Equal(t, "constraint_mutation", c.Source.Technique)
		assert.Equal(t, "P2", c.Priority)
	}
}

func TestConstraintMutationTechnique_Generate_NullInjection_Body(t *testing.T) {
	op := &spec.Operation{
		Method: "POST", Path: "/pets",
		RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
			"application/json": {Schema: &spec.Schema{
				Type:       "object",
				Properties: map[string]*spec.Schema{"count": {Type: "integer"}},
			}},
		}},
		Responses: map[string]*spec.Response{"422": {}},
	}
	cases, err := NewConstraintMutationTechnique().Generate(op)
	require.NoError(t, err)
	// Find null-injection case
	var nullCase *schema.TestCase
	for i, c := range cases {
		if body, ok := c.Steps[0].Body.(map[string]any); ok {
			if v, exists := body["count"]; exists && v == nil {
				nullCase = &cases[i]
				break
			}
		}
	}
	require.NotNil(t, nullCase, "expected null injection case for field 'count'")
	require.Len(t, nullCase.Steps[0].Assertions, 1)
	assert.Equal(t, 422, nullCase.Steps[0].Assertions[0].Expected)
}

func TestConstraintMutationTechnique_Generate_WrongContentType(t *testing.T) {
	op := &spec.Operation{
		Method: "POST", Path: "/pets",
		RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
			"application/json": {Schema: &spec.Schema{
				Type:       "object",
				Properties: map[string]*spec.Schema{"name": {Type: "string"}},
			}},
		}},
		Responses: map[string]*spec.Response{"415": {Description: "unsupported media type"}},
	}
	cases, err := NewConstraintMutationTechnique().Generate(op)
	require.NoError(t, err)
	var ctCase *schema.TestCase
	for i, c := range cases {
		if c.Steps[0].Headers["Content-Type"] == "text/plain" {
			ctCase = &cases[i]
			break
		}
	}
	require.NotNil(t, ctCase, "expected wrong-content-type case")
	require.Len(t, ctCase.Steps[0].Assertions, 1)
	assert.Equal(t, 415, ctCase.Steps[0].Assertions[0].Expected)
}

func TestConstraintMutationTechnique_Generate_AllNullable_OnlyWrongCT(t *testing.T) {
	op := &spec.Operation{
		Method: "POST", Path: "/pets",
		RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
			"application/json": {Schema: &spec.Schema{
				Type: "object",
				Properties: map[string]*spec.Schema{
					"name": {Type: "string", Nullable: true}, // all nullable
				},
			}},
		}},
		Responses: map[string]*spec.Response{"200": {}},
	}
	cases, err := NewConstraintMutationTechnique().Generate(op)
	require.NoError(t, err)
	// Only the wrong-content-type case; no null injections since all fields nullable
	assert.Len(t, cases, 1)
	assert.Equal(t, "text/plain", cases[0].Steps[0].Headers["Content-Type"])
}
