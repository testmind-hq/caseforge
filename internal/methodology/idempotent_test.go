package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestIdempotentAppliesForPOST(t *testing.T) {
	tech := &IdempotentTechnique{}
	assert.True(t, tech.Applies(&spec.Operation{Method: "POST"}))
	assert.True(t, tech.Applies(&spec.Operation{Method: "PUT"}))
	assert.True(t, tech.Applies(&spec.Operation{Method: "DELETE"}))
	assert.False(t, tech.Applies(&spec.Operation{Method: "GET"}))
}

func TestIdempotentGeneratesCase(t *testing.T) {
	tech := &IdempotentTechnique{}
	op := &spec.Operation{
		Method: "POST",
		Path:   "/users",
		Responses: map[string]*spec.Response{"201": {Description: "Created"}},
	}
	cases, err := tech.Generate(op)
	require.NoError(t, err)
	require.Len(t, cases, 1)

	tc := cases[0]
	assert.Equal(t, "idempotency", tc.Source.Technique)
	assert.Equal(t, "idempotency", tc.Labels["type"])
	assert.Equal(t, "P2", tc.Priority)
}
