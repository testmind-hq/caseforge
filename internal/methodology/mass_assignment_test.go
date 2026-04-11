// internal/methodology/mass_assignment_test.go
package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// makeMassAssignmentOp builds a minimal POST operation with a JSON request body.
func makeMassAssignmentOp() *spec.Operation {
	return &spec.Operation{
		Method: "POST", Path: "/users",
		RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
			"application/json": {Schema: &spec.Schema{
				Type: "object",
				Properties: map[string]*spec.Schema{
					"username": {Type: "string"},
					"email":    {Type: "string"},
				},
				Required: []string{"username"},
			}},
		}},
		Responses: map[string]*spec.Response{"200": {Description: "ok"}},
	}
}

func TestMassAssignmentTechnique_Applies_True(t *testing.T) {
	op := makeMassAssignmentOp()
	assert.True(t, NewMassAssignmentTechnique().Applies(op))
}

func TestMassAssignmentTechnique_Applies_False(t *testing.T) {
	op := &spec.Operation{Method: "GET", Path: "/users"}
	assert.False(t, NewMassAssignmentTechnique().Applies(op))
}

func TestMassAssignmentTechnique_Generate_ProducesExactly4Cases(t *testing.T) {
	op := makeMassAssignmentOp()
	cases, err := NewMassAssignmentTechnique().Generate(op)
	require.NoError(t, err)
	assert.Len(t, cases, 4)
}

func TestMassAssignmentTechnique_Generate_AllHaveP2Priority(t *testing.T) {
	op := makeMassAssignmentOp()
	cases, err := NewMassAssignmentTechnique().Generate(op)
	require.NoError(t, err)
	assert.NotEmpty(t, cases)
	for _, c := range cases {
		assert.Equal(t, "P2", c.Priority)
	}
}

func TestMassAssignmentTechnique_Generate_ScenariosPresent(t *testing.T) {
	op := makeMassAssignmentOp()
	cases, err := NewMassAssignmentTechnique().Generate(op)
	require.NoError(t, err)

	scenarios := make(map[string]bool)
	for _, c := range cases {
		scenarios[c.Source.Scenario] = true
	}

	assert.True(t, scenarios["MASS_ASSIGNMENT_PRIVILEGE"], "MASS_ASSIGNMENT_PRIVILEGE scenario missing")
	assert.True(t, scenarios["MASS_ASSIGNMENT_STATUS"], "MASS_ASSIGNMENT_STATUS scenario missing")
	assert.True(t, scenarios["MASS_ASSIGNMENT_FINANCIAL"], "MASS_ASSIGNMENT_FINANCIAL scenario missing")
	assert.True(t, scenarios["MASS_ASSIGNMENT_IDENTITY"], "MASS_ASSIGNMENT_IDENTITY scenario missing")
}

func TestMassAssignmentTechnique_Generate_InjectsProbeFields(t *testing.T) {
	op := makeMassAssignmentOp()
	cases, err := NewMassAssignmentTechnique().Generate(op)
	require.NoError(t, err)

	found := false
	for _, c := range cases {
		if c.Source.Scenario == "MASS_ASSIGNMENT_PRIVILEGE" {
			require.Len(t, c.Steps, 1)
			body, ok := c.Steps[0].Body.(map[string]any)
			require.True(t, ok)
			_, hasAdmin := body["admin"]
			if hasAdmin {
				found = true
			}
			break
		}
	}
	assert.True(t, found, "expected privilege case body to contain 'admin' key")
}

func TestMassAssignmentTechnique_Generate_DoesNotOverrideSchemaFields(t *testing.T) {
	op := makeMassAssignmentOp()
	cases, err := NewMassAssignmentTechnique().Generate(op)
	require.NoError(t, err)

	for _, c := range cases {
		require.Len(t, c.Steps, 1)
		body, ok := c.Steps[0].Body.(map[string]any)
		require.True(t, ok)
		// "username" is a schema field and must still be present
		_, hasUsername := body["username"]
		assert.True(t, hasUsername, "schema field 'username' should still be present in case %s", c.Source.Scenario)
	}
}
