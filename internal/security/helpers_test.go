package security_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testmind-hq/caseforge/internal/security"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestHasIDPathParam(t *testing.T) {
	assert.True(t, security.HasIDPathParam(&spec.Operation{Path: "/users/{userId}"}))
	assert.True(t, security.HasIDPathParam(&spec.Operation{Path: "/items/{id}"}))
	assert.False(t, security.HasIDPathParam(&spec.Operation{Path: "/users"}))
	assert.False(t, security.HasIDPathParam(&spec.Operation{Path: "/health"}))
}

func TestFindSensitiveFields(t *testing.T) {
	schema := &spec.Schema{Properties: map[string]*spec.Schema{
		"email":        {Type: "string"},
		"passwordHash": {Type: "string"},
		"accessToken":  {Type: "string"},
		"name":         {Type: "string"},
	}}
	fields := security.FindSensitiveFields(schema)
	assert.Contains(t, fields, "passwordHash")
	assert.Contains(t, fields, "accessToken")
	assert.NotContains(t, fields, "email")
	assert.NotContains(t, fields, "name")
}

func TestIsAuthRequired(t *testing.T) {
	assert.True(t, security.IsAuthRequired(&spec.Operation{Security: []string{"bearerAuth"}}))
	assert.False(t, security.IsAuthRequired(&spec.Operation{}))
}

func TestFindVersionedPaths(t *testing.T) {
	ops := []*spec.Operation{
		{Path: "/v1/users"},
		{Path: "/v2/users"},
		{Path: "/health"},
	}
	v1, v2 := security.FindVersionedPaths(ops)
	assert.NotEmpty(t, v1)
	assert.NotEmpty(t, v2)
}

func TestFindVersionedPaths_NoPairs(t *testing.T) {
	ops := []*spec.Operation{{Path: "/users"}, {Path: "/orders"}}
	v1, v2 := security.FindVersionedPaths(ops)
	assert.Empty(t, v1)
	assert.Empty(t, v2)
}
