// internal/methodology/idor_test.go
package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// makeIntPathParamOp builds an operation with an integer path parameter named "id".
func makeIntPathParamOp() *spec.Operation {
	return &spec.Operation{
		Method: "GET", Path: "/users/{id}",
		Parameters: []*spec.Parameter{
			{Name: "id", In: "path", Required: true, Schema: &spec.Schema{Type: "integer"}},
		},
		Responses: map[string]*spec.Response{"200": {Description: "ok"}},
	}
}

// makeUUIDQueryParamOp builds an operation with a UUID query parameter named "userId".
func makeUUIDQueryParamOp() *spec.Operation {
	return &spec.Operation{
		Method: "GET", Path: "/orders",
		Parameters: []*spec.Parameter{
			{Name: "userId", In: "query", Required: false, Schema: &spec.Schema{Type: "string", Format: "uuid"}},
		},
		Responses: map[string]*spec.Response{"200": {Description: "ok"}},
	}
}

// makeNoIDParamOp builds an operation with no ID-like parameters.
func makeNoIDParamOp() *spec.Operation {
	return &spec.Operation{
		Method: "GET", Path: "/health",
		Parameters: []*spec.Parameter{
			{Name: "format", In: "query", Required: false, Schema: &spec.Schema{Type: "string"}},
		},
		Responses: map[string]*spec.Response{"200": {Description: "ok"}},
	}
}

func TestIDORTechnique_Applies_IntPathParam(t *testing.T) {
	op := makeIntPathParamOp()
	assert.True(t, NewIDORTechnique().Applies(op))
}

func TestIDORTechnique_Applies_UUIDQueryParam(t *testing.T) {
	op := makeUUIDQueryParamOp()
	assert.True(t, NewIDORTechnique().Applies(op))
}

func TestIDORTechnique_Applies_False(t *testing.T) {
	op := makeNoIDParamOp()
	assert.False(t, NewIDORTechnique().Applies(op))
}

func TestIDORTechnique_Generate_IntegerID_Produces2Cases(t *testing.T) {
	op := makeIntPathParamOp()
	cases, err := NewIDORTechnique().Generate(op)
	require.NoError(t, err)
	assert.Len(t, cases, 2)
}

func TestIDORTechnique_Generate_UUIDID_Produces2Cases(t *testing.T) {
	op := makeUUIDQueryParamOp()
	cases, err := NewIDORTechnique().Generate(op)
	require.NoError(t, err)
	assert.Len(t, cases, 2)

	// Verify the UUID values are substituted correctly
	paths := make([]string, 0, 2)
	for _, c := range cases {
		require.Len(t, c.Steps, 1)
		paths = append(paths, c.Steps[0].Path)
	}
	assert.Contains(t, paths[0]+paths[1], idorAltUUID)
	assert.Contains(t, paths[0]+paths[1], idorNilUUID)
}

func TestIDORTechnique_Generate_AllHaveP1Priority(t *testing.T) {
	op := makeIntPathParamOp()
	cases, err := NewIDORTechnique().Generate(op)
	require.NoError(t, err)
	assert.NotEmpty(t, cases)
	for _, c := range cases {
		assert.Equal(t, "P1", c.Priority)
	}
}

func TestIDORTechnique_Generate_AllExpect403(t *testing.T) {
	op := makeIntPathParamOp()
	cases, err := NewIDORTechnique().Generate(op)
	require.NoError(t, err)
	assert.NotEmpty(t, cases)
	for _, c := range cases {
		require.Len(t, c.Steps, 1)
		require.Len(t, c.Steps[0].Assertions, 1)
		assert.Equal(t, "status_code", c.Steps[0].Assertions[0].Target)
		assert.Equal(t, "eq", c.Steps[0].Assertions[0].Operator)
		assert.Equal(t, 403, c.Steps[0].Assertions[0].Expected)
	}
}

func TestIDORTechnique_Generate_ScenarioIDOR(t *testing.T) {
	op := makeIntPathParamOp()
	cases, err := NewIDORTechnique().Generate(op)
	require.NoError(t, err)
	assert.NotEmpty(t, cases)
	for _, c := range cases {
		assert.Equal(t, "IDOR_PARAM", c.Source.Scenario)
	}
}

func TestIDORTechnique_Generate_IntegerPathSubstitution(t *testing.T) {
	op := makeIntPathParamOp()
	cases, err := NewIDORTechnique().Generate(op)
	require.NoError(t, err)
	require.Len(t, cases, 2)

	paths := map[string]bool{}
	for _, c := range cases {
		paths[c.Steps[0].Path] = true
	}
	assert.True(t, paths["/users/99999"], "expected /users/99999 path")
	assert.True(t, paths["/users/0"], "expected /users/0 path")
}

func TestIDORTechnique_IsIDParam_NameHeuristic(t *testing.T) {
	assert.True(t, isIDLike("user_id"))
	assert.True(t, isIDLike("userId"))
	assert.True(t, isIDLike("id"))
	assert.True(t, isIDLike("orderId"))
	assert.True(t, isIDLike("USERID"))
	assert.False(t, isIDLike("name"))
	assert.False(t, isIDLike("email"))
	assert.False(t, isIDLike("status"))
	assert.False(t, isIDLike("format"))
}
