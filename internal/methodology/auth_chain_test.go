// internal/methodology/auth_chain_test.go
package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func authSpec() *spec.ParsedSpec {
	return &spec.ParsedSpec{
		GlobalSecurity:  []string{"BearerAuth"},
		SecuritySchemes: []string{"BearerAuth"},
		Operations: []*spec.Operation{
			{
				OperationID: "login",
				Method:      "POST",
				Path:        "/auth/login",
				RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
					"application/json": {Schema: &spec.Schema{
						Type: "object",
						Properties: map[string]*spec.Schema{
							"username": {Type: "string"},
							"password": {Type: "string"},
						},
					}},
				}},
				Responses: map[string]*spec.Response{"200": {
					Content: map[string]*spec.MediaType{
						"application/json": {Schema: &spec.Schema{
							Type: "object",
							Properties: map[string]*spec.Schema{
								"access_token": {Type: "string"},
								"expires_in":   {Type: "integer"},
							},
						}},
					},
				}},
			},
			{
				OperationID: "getProfile",
				Method:      "GET",
				Path:        "/profile",
				Security:    []string{"BearerAuth"},
				Responses:   map[string]*spec.Response{"200": {}},
			},
			{
				OperationID: "listItems",
				Method:      "GET",
				Path:        "/items",
				Security:    []string{"BearerAuth"},
				Responses:   map[string]*spec.Response{"200": {}},
			},
		},
	}
}

func TestAuthChainTechnique_DetectsAuthOperation(t *testing.T) {
	tech := NewAuthChainTechnique()
	cases, err := tech.Generate(authSpec())
	require.NoError(t, err)
	require.NotEmpty(t, cases, "should generate auth chain cases for secured operations")
}

func TestAuthChainTechnique_GeneratesTwoStepChain(t *testing.T) {
	tech := NewAuthChainTechnique()
	cases, err := tech.Generate(authSpec())
	require.NoError(t, err)
	require.NotEmpty(t, cases)

	for _, tc := range cases {
		assert.Equal(t, "chain", tc.Kind)
		assert.Equal(t, 2, len(tc.Steps), "auth chain must have exactly 2 steps: auth + test")
		assert.Equal(t, "setup", tc.Steps[0].Type)
		assert.Equal(t, "test", tc.Steps[1].Type)
	}
}

func TestAuthChainTechnique_CapturesToken(t *testing.T) {
	tech := NewAuthChainTechnique()
	cases, err := tech.Generate(authSpec())
	require.NoError(t, err)
	require.NotEmpty(t, cases)

	authStep := cases[0].Steps[0]
	require.NotEmpty(t, authStep.Captures)
	assert.Equal(t, "authToken", authStep.Captures[0].Name)
	assert.Contains(t, authStep.Captures[0].From, "access_token")
}

func TestAuthChainTechnique_InjectsTokenInTestStep(t *testing.T) {
	tech := NewAuthChainTechnique()
	cases, err := tech.Generate(authSpec())
	require.NoError(t, err)
	require.NotEmpty(t, cases)

	testStep := cases[0].Steps[1]
	authHeader, ok := testStep.Headers["Authorization"]
	assert.True(t, ok, "test step must have Authorization header")
	assert.Contains(t, authHeader, "{{authToken}}")
}

func TestAuthChainTechnique_NoAuthOp_NoOutput(t *testing.T) {
	ps := &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{Method: "GET", Path: "/users", Responses: map[string]*spec.Response{"200": {}}},
		},
	}
	tech := NewAuthChainTechnique()
	cases, err := tech.Generate(ps)
	require.NoError(t, err)
	assert.Empty(t, cases)
}

func TestAuthChainTechnique_SourceAnnotation(t *testing.T) {
	tech := NewAuthChainTechnique()
	cases, err := tech.Generate(authSpec())
	require.NoError(t, err)
	require.NotEmpty(t, cases)
	for _, tc := range cases {
		assert.Equal(t, "auth_chain", tc.Source.Technique)
	}
}
