// internal/methodology/engine_test.go
package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestEngineGeneratesFromParsedSpec(t *testing.T) {
	noop := &llm.NoopProvider{}
	engine := NewEngine(noop,
		NewEquivalenceTechnique(),
		NewBoundaryTechnique(),
		NewDecisionTechnique(),
		NewStateTechnique(),
		NewIdempotentTechnique(),
		NewPairwiseTechnique(),
	)

	ps := &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{
				OperationID: "createUser",
				Method:      "POST",
				Path:        "/users",
				Summary:     "Create a user",
				RequestBody: &spec.RequestBody{
					Content: map[string]*spec.MediaType{
						"application/json": {
							Schema: &spec.Schema{
								Type:     "object",
								Required: []string{"email"},
								Properties: map[string]*spec.Schema{
									"email": {Type: "string", Format: "email"},
									"age":   {Type: "integer", Minimum: floatPtr(0), Maximum: floatPtr(120)},
								},
							},
						},
					},
				},
				Responses: map[string]*spec.Response{"201": {Description: "Created"}},
			},
		},
	}

	cases, err := engine.Generate(ps)
	require.NoError(t, err)
	assert.NotEmpty(t, cases)

	// All cases must have traceability
	for _, tc := range cases {
		assert.NotEmpty(t, tc.ID)
		assert.NotEmpty(t, tc.Source.Technique)
		assert.NotEmpty(t, tc.Source.SpecPath)
		assert.NotEmpty(t, tc.Source.Rationale)
		assert.Equal(t, "1", tc.Version)
	}
}
