package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestDecisionAppliesWithEnumFields(t *testing.T) {
	tech := &DecisionTechnique{}
	op := &spec.Operation{
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {Schema: &spec.Schema{
					Type: "object",
					Properties: map[string]*spec.Schema{
						"status": {Enum: []any{"active", "inactive"}},
						"role":   {Enum: []any{"admin", "user"}},
					},
				}},
			},
		},
	}
	assert.True(t, tech.Applies(op))
}

func TestDecisionDoesNotApplyWithOneEnumField(t *testing.T) {
	tech := &DecisionTechnique{}
	op := &spec.Operation{
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {Schema: &spec.Schema{
					Type: "object",
					Properties: map[string]*spec.Schema{
						"status": {Enum: []any{"active", "inactive"}},
					},
				}},
			},
		},
	}
	assert.False(t, tech.Applies(op))
}

func TestDecisionGeneratesBooleanCases(t *testing.T) {
	tech := &DecisionTechnique{}
	op := &spec.Operation{
		Method: "POST",
		Path:   "/items",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {Schema: &spec.Schema{
					Type: "object",
					Properties: map[string]*spec.Schema{
						"active":  {Type: "boolean"},
						"visible": {Type: "boolean"},
					},
				}},
			},
		},
		Responses: map[string]*spec.Response{"200": {Description: "OK"}},
	}
	cases, err := tech.Generate(op)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(cases), 4, "2 boolean fields * 2 values each = at least 4 cases")

	for _, tc := range cases {
		assert.Equal(t, "decision_table", tc.Source.Technique)
		assert.NotEmpty(t, tc.Source.Rationale)
	}
}
