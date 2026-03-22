package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestStateAppliesWhenHasStateMachine(t *testing.T) {
	tech := &StateTechnique{}
	op := &spec.Operation{
		SemanticInfo: &spec.SemanticAnnotation{
			HasStateMachine: true,
			StateField:      "status",
		},
	}
	assert.True(t, tech.Applies(op))
}

func TestStateDoesNotApplyWithoutStateMachine(t *testing.T) {
	tech := &StateTechnique{}
	op := &spec.Operation{}
	assert.False(t, tech.Applies(op))
}

func TestStateGeneratesTransitionCases(t *testing.T) {
	tech := &StateTechnique{}
	op := &spec.Operation{
		Method: "PUT",
		Path:   "/orders/{id}",
		SemanticInfo: &spec.SemanticAnnotation{
			HasStateMachine: true,
			StateField:      "status",
		},
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {Schema: &spec.Schema{
					Type: "object",
					Properties: map[string]*spec.Schema{
						"status": {Enum: []any{"pending", "confirmed", "shipped", "cancelled"}},
					},
				}},
			},
		},
		Responses: map[string]*spec.Response{"200": {Description: "OK"}},
	}
	cases, err := tech.Generate(op)
	require.NoError(t, err)
	assert.Equal(t, 4, len(cases), "one case per state value")

	for _, tc := range cases {
		assert.Equal(t, "state_transition", tc.Source.Technique)
		assert.NotEmpty(t, tc.Source.Rationale)
	}
}
