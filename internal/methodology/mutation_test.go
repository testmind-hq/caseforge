// internal/methodology/mutation_test.go
package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func mutationOp() *spec.Operation {
	return &spec.Operation{
		OperationID: "createUser",
		Method:      "POST",
		Path:        "/users",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {Schema: &spec.Schema{
					Type:     "object",
					Required: []string{"email", "age"},
					Properties: map[string]*spec.Schema{
						"email": {Type: "string", Format: "email"},
						"age":   {Type: "integer", Minimum: floatPtr(18), Maximum: floatPtr(120)},
						"name":  {Type: "string"},
					},
				}},
			},
		},
		Responses: map[string]*spec.Response{"201": {Description: "Created"}},
	}
}

func TestMutationTechnique_Applies(t *testing.T) {
	tech := NewMutationTechnique()
	assert.True(t, tech.Applies(mutationOp()), "applies to POST with JSON body")

	noBody := &spec.Operation{
		Method: "GET", Path: "/users",
		Responses: map[string]*spec.Response{"200": {}},
	}
	assert.False(t, tech.Applies(noBody), "does not apply to GET without body")
}

func TestMutationTechnique_GeneratesPerFieldCases(t *testing.T) {
	tech := NewMutationTechnique()
	cases, err := tech.Generate(mutationOp())
	require.NoError(t, err)
	require.NotEmpty(t, cases)

	for _, tc := range cases {
		assert.Equal(t, "mutation", tc.Source.Technique)
		assert.NotEmpty(t, tc.Source.Rationale)
		assert.NotEmpty(t, tc.Steps)
		// Each case must assert 4xx
		found4xx := false
		for _, a := range tc.Steps[0].Assertions {
			if a.Target == "status_code" {
				if code, ok := a.Expected.(int); ok && code >= 400 && code < 500 {
					found4xx = true
				}
			}
		}
		assert.True(t, found4xx, "mutation case must assert 4xx status: %s", tc.Title)
	}
}

func TestMutationTechnique_MaxCasesLimit(t *testing.T) {
	tech := NewMutationTechniqueWithMax(3)
	cases, err := tech.Generate(mutationOp())
	require.NoError(t, err)
	assert.LessOrEqual(t, len(cases), 3, "must respect max cases limit")
}

func TestMutationTechnique_RationaleDescribesMutation(t *testing.T) {
	tech := NewMutationTechnique()
	cases, err := tech.Generate(mutationOp())
	require.NoError(t, err)
	require.NotEmpty(t, cases)

	for _, tc := range cases {
		assert.NotEmpty(t, tc.Source.Rationale, "rationale must describe the mutation")
	}
}
