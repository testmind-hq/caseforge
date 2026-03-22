// internal/methodology/equivalence_test.go
package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestEquivalenceAppliesAlways(t *testing.T) {
	tech := &EquivalenceTechnique{}
	op := &spec.Operation{Method: "GET", Path: "/users"}
	assert.True(t, tech.Applies(op))
}

func TestEquivalenceGeneratesPositiveAndNegativeCases(t *testing.T) {
	tech := &EquivalenceTechnique{}
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
							"email": {Type: "string", Format: "email"},
							"age":   {Type: "integer", Minimum: floatPtr(0), Maximum: floatPtr(120)},
						},
						Required: []string{"email"},
					},
				},
			},
		},
		Responses: map[string]*spec.Response{"201": {Description: "Created"}},
	}
	cases, err := tech.Generate(op)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(cases), 2, "should generate at least happy path and one negative case")

	// Verify source tracing
	for _, tc := range cases {
		assert.Equal(t, "equivalence_partitioning", tc.Source.Technique)
		assert.NotEmpty(t, tc.Source.Rationale)
		assert.NotEmpty(t, tc.ID)
	}
}

func floatPtr(f float64) *float64 { return &f }
