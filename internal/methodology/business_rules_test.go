package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func makeBusinessRuleOp() *spec.Operation {
	return &spec.Operation{
		Method: "POST", Path: "/users",
		RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
			"application/json": {Schema: &spec.Schema{
				Type: "object",
				Properties: map[string]*spec.Schema{
					"email":    {Type: "string"},
					"username": {Type: "string"},
				},
				Required: []string{"email"},
			}},
		}},
		Responses:    map[string]*spec.Response{"201": {Description: "created"}},
		SemanticInfo: &spec.SemanticAnnotation{
			ImplicitRules: []string{
				"email must be unique per account",
				"username must not contain spaces",
			},
		},
	}
}

func TestBusinessRuleTechnique_Applies_WithRules(t *testing.T) {
	op := makeBusinessRuleOp()
	assert.True(t, NewBusinessRuleTechnique().Applies(op))
}

func TestBusinessRuleTechnique_Applies_NoSemanticInfo(t *testing.T) {
	op := &spec.Operation{Method: "GET", Path: "/users"}
	assert.False(t, NewBusinessRuleTechnique().Applies(op))
}

func TestBusinessRuleTechnique_Applies_EmptyRules(t *testing.T) {
	op := &spec.Operation{
		Method: "GET", Path: "/users",
		SemanticInfo: &spec.SemanticAnnotation{ImplicitRules: []string{}},
	}
	assert.False(t, NewBusinessRuleTechnique().Applies(op))
}

func TestBusinessRuleTechnique_Generate_OnePerRule(t *testing.T) {
	op := makeBusinessRuleOp()
	cases, err := NewBusinessRuleTechnique().Generate(op)
	require.NoError(t, err)
	assert.Len(t, cases, 2)
}

func TestBusinessRuleTechnique_Generate_Expects4xx(t *testing.T) {
	op := makeBusinessRuleOp()
	cases, err := NewBusinessRuleTechnique().Generate(op)
	require.NoError(t, err)
	for _, tc := range cases {
		require.Len(t, tc.Steps, 1)
		found4xx := false
		for _, a := range tc.Steps[0].Assertions {
			if a.Target == "status_code" && a.Operator == "gte" {
				v, ok := a.Expected.(int)
				if ok && v >= 400 {
					found4xx = true
				}
			}
		}
		assert.True(t, found4xx, "expected 4xx assertion in case %s", tc.ID)
	}
}

func TestBusinessRuleTechnique_Generate_PriorityP2(t *testing.T) {
	op := makeBusinessRuleOp()
	cases, err := NewBusinessRuleTechnique().Generate(op)
	require.NoError(t, err)
	for _, tc := range cases {
		assert.Equal(t, "P2", tc.Priority)
	}
}

func TestBusinessRuleTechnique_Generate_ScenarioPopulated(t *testing.T) {
	op := makeBusinessRuleOp()
	cases, err := NewBusinessRuleTechnique().Generate(op)
	require.NoError(t, err)
	for _, tc := range cases {
		assert.Equal(t, "BUSINESS_RULE_VIOLATION", tc.Source.Scenario)
	}
}
