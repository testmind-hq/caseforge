// internal/methodology/boundary_test.go
package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestBoundaryAppliesWhenRangeConstraints(t *testing.T) {
	tech := &BoundaryTechnique{}
	op := &spec.Operation{
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type: "object",
						Properties: map[string]*spec.Schema{
							"age": {Type: "integer", Minimum: floatPtr(0), Maximum: floatPtr(120)},
						},
					},
				},
			},
		},
	}
	assert.True(t, tech.Applies(op))
}

func TestBoundaryDoesNotApplyWithoutConstraints(t *testing.T) {
	tech := &BoundaryTechnique{}
	op := &spec.Operation{Method: "GET", Path: "/users"}
	assert.False(t, tech.Applies(op))
}

func TestBoundaryGeneratesMinMaxCases(t *testing.T) {
	tech := NewBoundaryTechnique()
	op := &spec.Operation{
		OperationID: "createUser",
		Method:      "POST",
		Path:        "/users",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type: "object",
						Properties: map[string]*spec.Schema{
							"age": {Type: "integer", Minimum: floatPtr(18), Maximum: floatPtr(100)},
						},
					},
				},
			},
		},
		Responses: map[string]*spec.Response{"201": {Description: "Created"}},
	}
	cases, err := tech.Generate(op)
	require.NoError(t, err)
	// Expect: min valid, min-1 invalid, max valid, max+1 invalid = 4 cases for "age"
	assert.GreaterOrEqual(t, len(cases), 4)
	for _, tc := range cases {
		assert.Equal(t, "boundary_value", tc.Source.Technique)
		assert.NotEmpty(t, tc.Source.Rationale)
		assert.NotEmpty(t, tc.ID)
	}
}
