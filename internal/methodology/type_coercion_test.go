// internal/methodology/type_coercion_test.go
package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// makeTypedOp builds a minimal POST operation with a JSON request body containing
// the given properties.
func makeTypedOp(props map[string]*spec.Schema) *spec.Operation {
	return &spec.Operation{
		Method: "POST", Path: "/items",
		RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
			"application/json": {Schema: &spec.Schema{
				Type:       "object",
				Properties: props,
			}},
		}},
		Responses: map[string]*spec.Response{"422": {Description: "invalid"}},
	}
}

func TestTypeCoercionTechnique_Applies_True(t *testing.T) {
	op := makeTypedOp(map[string]*spec.Schema{
		"name": {Type: "string"},
	})
	assert.True(t, NewTypeCoercionTechnique().Applies(op))
}

func TestTypeCoercionTechnique_Applies_False(t *testing.T) {
	op := &spec.Operation{Method: "GET", Path: "/items"}
	assert.False(t, NewTypeCoercionTechnique().Applies(op))
}

func TestTypeCoercionTechnique_Applies_False_NoTypedFields(t *testing.T) {
	// All fields have no Type set — should not apply.
	op := makeTypedOp(map[string]*spec.Schema{
		"flexible": {Type: ""},
	})
	assert.False(t, NewTypeCoercionTechnique().Applies(op))
}

func TestTypeCoercionTechnique_Generate_StringField(t *testing.T) {
	op := makeTypedOp(map[string]*spec.Schema{
		"name": {Type: "string"},
	})
	cases, err := NewTypeCoercionTechnique().Generate(op)
	require.NoError(t, err)
	// string → 2 mutations (integer + boolean)
	assert.Len(t, cases, 2)
	for _, c := range cases {
		assert.Equal(t, "type_coercion", c.Source.Technique)
		assert.Equal(t, "WRONG_TYPE", c.Source.Scenario)
		assert.Equal(t, "P2", c.Priority)
	}
}

func TestTypeCoercionTechnique_Generate_NumberField(t *testing.T) {
	op := makeTypedOp(map[string]*spec.Schema{
		"count": {Type: "number"},
	})
	cases, err := NewTypeCoercionTechnique().Generate(op)
	require.NoError(t, err)
	// number → 2 mutations (string + boolean)
	assert.Len(t, cases, 2)
	for _, c := range cases {
		assert.Equal(t, "WRONG_TYPE", c.Source.Scenario)
	}
}

func TestTypeCoercionTechnique_Generate_IntegerField(t *testing.T) {
	op := makeTypedOp(map[string]*spec.Schema{
		"age": {Type: "integer"},
	})
	cases, err := NewTypeCoercionTechnique().Generate(op)
	require.NoError(t, err)
	// integer → 2 mutations (string + boolean)
	assert.Len(t, cases, 2)
}

func TestTypeCoercionTechnique_Generate_BooleanField(t *testing.T) {
	op := makeTypedOp(map[string]*spec.Schema{
		"active": {Type: "boolean"},
	})
	cases, err := NewTypeCoercionTechnique().Generate(op)
	require.NoError(t, err)
	// boolean → 2 mutations (string + integer)
	assert.Len(t, cases, 2)
	for _, c := range cases {
		assert.Equal(t, "WRONG_TYPE", c.Source.Scenario)
	}
}

func TestTypeCoercionTechnique_Generate_ArrayField(t *testing.T) {
	op := makeTypedOp(map[string]*spec.Schema{
		"tags": {Type: "array"},
	})
	cases, err := NewTypeCoercionTechnique().Generate(op)
	require.NoError(t, err)
	// array → 1 mutation (string)
	assert.Len(t, cases, 1)
}

func TestTypeCoercionTechnique_Generate_ObjectField(t *testing.T) {
	op := makeTypedOp(map[string]*spec.Schema{
		"meta": {Type: "object"},
	})
	cases, err := NewTypeCoercionTechnique().Generate(op)
	require.NoError(t, err)
	// object → 1 mutation (string)
	assert.Len(t, cases, 1)
}

func TestTypeCoercionTechnique_Generate_AllExpect422(t *testing.T) {
	op := makeTypedOp(map[string]*spec.Schema{
		"name":   {Type: "string"},
		"count":  {Type: "integer"},
		"active": {Type: "boolean"},
		"tags":   {Type: "array"},
		"meta":   {Type: "object"},
	})
	cases, err := NewTypeCoercionTechnique().Generate(op)
	require.NoError(t, err)
	assert.NotEmpty(t, cases)
	for _, c := range cases {
		require.Len(t, c.Steps, 1)
		require.Len(t, c.Steps[0].Assertions, 1)
		a := c.Steps[0].Assertions[0]
		assert.Equal(t, "status_code", a.Target)
		assert.Equal(t, "eq", a.Operator)
		assert.Equal(t, 422, a.Expected)
	}
}

func TestTypeCoercionTechnique_Generate_SourceScenario(t *testing.T) {
	op := makeTypedOp(map[string]*spec.Schema{
		"name":   {Type: "string"},
		"active": {Type: "boolean"},
	})
	cases, err := NewTypeCoercionTechnique().Generate(op)
	require.NoError(t, err)
	assert.NotEmpty(t, cases)
	for _, c := range cases {
		assert.Equal(t, "WRONG_TYPE", c.Source.Scenario)
	}
}

func TestTypeCoercionTechnique_Generate_SkipsUntypedFields(t *testing.T) {
	op := makeTypedOp(map[string]*spec.Schema{
		"name":      {Type: "string"},
		"flexible":  {Type: ""},   // no type → skip
		"nilSchema": nil,          // nil schema → skip
	})
	cases, err := NewTypeCoercionTechnique().Generate(op)
	require.NoError(t, err)
	// Only string field mutated: 2 cases
	assert.Len(t, cases, 2)
}

func TestTypeCoercionTechnique_Generate_SkipsEnumMatch(t *testing.T) {
	// A string field whose enum includes integer 123 — the integer mutation should be skipped.
	op := makeTypedOp(map[string]*spec.Schema{
		"status": {
			Type: "string",
			Enum: []any{123, "active"},
		},
	})
	cases, err := NewTypeCoercionTechnique().Generate(op)
	require.NoError(t, err)
	// integer 123 is in enum → skip that mutation; boolean true is not → 1 case
	assert.Len(t, cases, 1)
}

func TestTypeCoercionTechnique_Name(t *testing.T) {
	assert.Equal(t, "type_coercion", NewTypeCoercionTechnique().Name())
}
