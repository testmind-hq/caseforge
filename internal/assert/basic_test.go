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

func TestSchemaAssertions_EmailFormatUsesMatches(t *testing.T) {
	s := &specpkg.Schema{
		Type: "object",
		Properties: map[string]*specpkg.Schema{
			"email": {Type: "string", Format: "email"},
		},
	}
	assertions := SchemaAssertions("body", s)
	require.Len(t, assertions, 1)
	assert.Equal(t, "matches", assertions[0].Operator)
	assert.NotEmpty(t, assertions[0].Expected, "email matches assertion must include a regex pattern")
}

func TestSchemaAssertions_URIFormatUsesMatches(t *testing.T) {
	for _, format := range []string{"uri", "url", "uri-reference"} {
		s := &specpkg.Schema{
			Type: "object",
			Properties: map[string]*specpkg.Schema{
				"link": {Type: "string", Format: format},
			},
		}
		assertions := SchemaAssertions("body", s)
		require.Len(t, assertions, 1, "format=%s", format)
		assert.Equal(t, "matches", assertions[0].Operator, "format=%s", format)
	}
}

func TestSchemaAssertions_IPv4FormatUsesMatches(t *testing.T) {
	s := &specpkg.Schema{
		Type: "object",
		Properties: map[string]*specpkg.Schema{
			"ip": {Type: "string", Format: "ipv4"},
		},
	}
	assertions := SchemaAssertions("body", s)
	require.Len(t, assertions, 1)
	assert.Equal(t, "matches", assertions[0].Operator)
}

func TestSchemaAssertions_IPv6FormatUsesMatches(t *testing.T) {
	s := &specpkg.Schema{
		Type: "object",
		Properties: map[string]*specpkg.Schema{
			"ip": {Type: "string", Format: "ipv6"},
		},
	}
	assertions := SchemaAssertions("body", s)
	require.Len(t, assertions, 1)
	assert.Equal(t, "matches", assertions[0].Operator)
}

func TestRangeAssertions_MinimumGeneratesGte(t *testing.T) {
	min := 0.0
	s := &specpkg.Schema{
		Type: "object",
		Properties: map[string]*specpkg.Schema{
			"age": {Type: "integer", Minimum: &min},
		},
	}
	assertions := RangeAssertions("body", s)
	require.Len(t, assertions, 1)
	assert.Equal(t, "body.age", assertions[0].Target)
	assert.Equal(t, "gte", assertions[0].Operator)
	assert.Equal(t, 0.0, assertions[0].Expected)
}

func TestRangeAssertions_MaximumGeneratesLte(t *testing.T) {
	max := 150.0
	s := &specpkg.Schema{
		Type: "object",
		Properties: map[string]*specpkg.Schema{
			"age": {Type: "integer", Maximum: &max},
		},
	}
	assertions := RangeAssertions("body", s)
	require.Len(t, assertions, 1)
	assert.Equal(t, "lte", assertions[0].Operator)
	assert.Equal(t, 150.0, assertions[0].Expected)
}

func TestRangeAssertions_BothBounds_GeneratesTwo(t *testing.T) {
	min, max := 1.0, 100.0
	s := &specpkg.Schema{
		Type: "object",
		Properties: map[string]*specpkg.Schema{
			"score": {Type: "number", Minimum: &min, Maximum: &max},
		},
	}
	assertions := RangeAssertions("body", s)
	require.Len(t, assertions, 2)
	ops := map[string]bool{}
	for _, a := range assertions {
		ops[a.Operator] = true
	}
	assert.True(t, ops["gte"], "expected gte assertion")
	assert.True(t, ops["lte"], "expected lte assertion")
}

func TestRangeAssertions_NoConstraints_ReturnsEmpty(t *testing.T) {
	s := &specpkg.Schema{
		Type: "object",
		Properties: map[string]*specpkg.Schema{
			"name": {Type: "string"},
		},
	}
	assert.Empty(t, RangeAssertions("body", s))
}

func TestRangeAssertions_NilSchema(t *testing.T) {
	assert.Empty(t, RangeAssertions("body", nil))
}

func TestRangeAssertions_NonObjectSchema(t *testing.T) {
	s := &specpkg.Schema{Type: "array"}
	assert.Empty(t, RangeAssertions("body", s))
}
