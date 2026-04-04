// internal/assert/basic_test.go
package assert

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	specpkg "github.com/testmind-hq/caseforge/internal/spec"
)

func TestStatusCodeAssertion(t *testing.T) {
	op := &specpkg.Operation{
		Responses: map[string]*specpkg.Response{
			"200": {Description: "OK"},
		},
	}
	assertions := BasicAssertions(op)
	var found bool
	for _, a := range assertions {
		if a.Target == "status_code" && a.Operator == "eq" && a.Expected == 200 {
			found = true
		}
	}
	assert.True(t, found, "expected status_code eq 200 assertion")
}

func TestDurationAssertion(t *testing.T) {
	op := &specpkg.Operation{}
	assertions := BasicAssertions(op)
	var found bool
	for _, a := range assertions {
		if a.Target == "duration_ms" && a.Operator == "lt" {
			found = true
		}
	}
	assert.True(t, found, "expected duration_ms lt assertion")
}

func TestSchemaAssertionsForObjectResponse(t *testing.T) {
	props := map[string]*specpkg.Schema{
		"id":   {Type: "string", Format: "uuid"},
		"name": {Type: "string"},
	}
	s := &specpkg.Schema{Type: "object", Properties: props}
	assertions := SchemaAssertions("body", s)
	targets := make(map[string]bool)
	for _, a := range assertions {
		targets[a.Target] = true
	}
	assert.True(t, targets["body.id"])
	assert.True(t, targets["body.name"])
}

func TestSchemaAssertions_UUIDFieldUsesIsUUID(t *testing.T) {
	s := &specpkg.Schema{
		Type: "object",
		Properties: map[string]*specpkg.Schema{
			"id": {Type: "string", Format: "uuid"},
		},
	}
	assertions := SchemaAssertions("body", s)
	require.Len(t, assertions, 1)
	assert.Equal(t, "body.id", assertions[0].Target)
	assert.Equal(t, "is_uuid", assertions[0].Operator)
}

func TestSchemaAssertions_DateTimeFieldUsesIsISO8601(t *testing.T) {
	for _, format := range []string{"date-time", "date", "time"} {
		s := &specpkg.Schema{
			Type: "object",
			Properties: map[string]*specpkg.Schema{
				"ts": {Type: "string", Format: format},
			},
		}
		assertions := SchemaAssertions("body", s)
		require.Len(t, assertions, 1, "format=%s", format)
		assert.Equal(t, "is_iso8601", assertions[0].Operator, "format=%s", format)
	}
}

func TestSchemaAssertions_PlainStringFieldUsesExists(t *testing.T) {
	s := &specpkg.Schema{
		Type: "object",
		Properties: map[string]*specpkg.Schema{
			"name": {Type: "string"},
		},
	}
	assertions := SchemaAssertions("body", s)
	require.Len(t, assertions, 1)
	assert.Equal(t, "exists", assertions[0].Operator)
}

func TestSchemaAssertions_NilSchema(t *testing.T) {
	assert.Empty(t, SchemaAssertions("body", nil))
}

func TestSchemaAssertions_NonObjectSchema(t *testing.T) {
	s := &specpkg.Schema{Type: "array"}
	assert.Empty(t, SchemaAssertions("body", s))
}
