package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestIdempotentAppliesForPOST(t *testing.T) {
	tech := NewIdempotentTechnique()
	assert.True(t, tech.Applies(&spec.Operation{Method: "POST"}))
	assert.True(t, tech.Applies(&spec.Operation{Method: "PUT"}))
	assert.True(t, tech.Applies(&spec.Operation{Method: "DELETE"}))
	assert.False(t, tech.Applies(&spec.Operation{Method: "GET"}))
}

func TestIdempotentGeneratesCase(t *testing.T) {
	tech := NewIdempotentTechnique()
	op := &spec.Operation{
		Method:    "POST",
		Path:      "/users",
		Responses: map[string]*spec.Response{"201": {Description: "Created"}},
	}
	cases, err := tech.Generate(op)
	require.NoError(t, err)
	require.Len(t, cases, 1)

	tc := cases[0]
	assert.Equal(t, "idempotency", tc.Source.Technique)
	assert.Equal(t, "idempotency", tc.Labels["type"])
	assert.Equal(t, "P2", tc.Priority)
	assert.Equal(t, "chain", tc.Kind)
	require.Len(t, tc.Steps, 2)
}

func TestIdempotentChainStructure(t *testing.T) {
	tech := NewIdempotentTechnique()
	op := &spec.Operation{
		Method:    "PUT",
		Path:      "/users/{id}",
		Responses: map[string]*spec.Response{"200": {Description: "OK"}},
	}
	cases, err := tech.Generate(op)
	require.NoError(t, err)
	require.Len(t, cases, 1)
	tc := cases[0]

	require.Len(t, tc.Steps, 2)
	setup := tc.Steps[0]
	test := tc.Steps[1]

	// setup step
	assert.Equal(t, "step-setup", setup.ID)
	assert.Equal(t, "setup", setup.Type)
	assert.Equal(t, op.Method, setup.Method)
	assert.Equal(t, op.Path, setup.Path)
	assert.Empty(t, setup.DependsOn)

	// test step — identical request, depends on setup
	assert.Equal(t, "step-test", test.ID)
	assert.Equal(t, "test", test.Type)
	assert.Equal(t, op.Method, test.Method)
	assert.Equal(t, op.Path, test.Path)
	assert.Equal(t, []string{"step-setup"}, test.DependsOn)

	// Both steps carry the same body and headers (idempotency: identical input)
	assert.Equal(t, setup.Body, test.Body)
	assert.Equal(t, setup.Headers, test.Headers)
}
