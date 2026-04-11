// internal/datagen/generator_test.go
package datagen

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestGenerateForEmailFormat(t *testing.T) {
	g := NewGenerator(nil)
	val := g.Generate(&spec.Schema{Type: "string", Format: "email"}, "email")
	s, ok := val.(string)
	assert.True(t, ok)
	assert.Contains(t, s, "@")
}

func TestGenerateForFieldNameSemantic(t *testing.T) {
	g := NewGenerator(nil)
	// "email" field name should produce an email even without format hint
	val := g.Generate(&spec.Schema{Type: "string"}, "email")
	s, ok := val.(string)
	assert.True(t, ok)
	assert.Contains(t, s, "@")
}

func TestGenerateForIntegerType(t *testing.T) {
	g := NewGenerator(nil)
	val := g.Generate(&spec.Schema{Type: "integer"}, "count")
	_, ok := val.(int64)
	assert.True(t, ok)
}

func TestGenerateForBooleanType(t *testing.T) {
	g := NewGenerator(nil)
	val := g.Generate(&spec.Schema{Type: "boolean"}, "active")
	_, ok := val.(bool)
	assert.True(t, ok)
}

func TestGenerateWithEnumPicksFromEnum(t *testing.T) {
	g := NewGenerator(nil)
	enum := []any{"active", "inactive", "pending"}
	val := g.Generate(&spec.Schema{Type: "string", Enum: enum}, "status")
	assert.Contains(t, enum, val)
}

func TestGenerateBoundaryMin(t *testing.T) {
	g := NewGenerator(nil)
	min := 18.0
	val := g.GenerateBoundary(&spec.Schema{Type: "integer", Minimum: &min}, BoundaryMin)
	assert.Equal(t, int64(18), val)
}

func TestGenerateBoundaryMaxPlus1(t *testing.T) {
	g := NewGenerator(nil)
	max := 100.0
	val := g.GenerateBoundary(&spec.Schema{Type: "integer", Maximum: &max}, BoundaryMaxPlusOne)
	assert.Equal(t, int64(101), val)
}

func TestGenerateByPattern_Digits(t *testing.T) {
	val, ok := generateByPattern(`^[0-9]+$`)
	assert.True(t, ok, "expected ok=true for digit pattern")
	assert.Regexp(t, regexp.MustCompile(`[0-9]+`), val)
}

func TestGenerateByPattern_Alphanumeric(t *testing.T) {
	val, ok := generateByPattern(`^[a-zA-Z0-9]+$`)
	assert.True(t, ok, "expected ok=true for alphanumeric pattern")
	assert.Regexp(t, regexp.MustCompile(`^[a-zA-Z0-9]+$`), val)
}

func TestGenerateByPattern_InvalidPattern(t *testing.T) {
	// Should not panic; returns fallback word
	val, ok := generateByPattern(`[invalid`)
	assert.False(t, ok, "expected ok=false for invalid pattern")
	assert.NotEmpty(t, val, "fallback should return non-empty string")
}

func TestGenerator_Generate_UsesPattern(t *testing.T) {
	g := NewGenerator(nil)
	schema := &spec.Schema{Type: "string", Pattern: `^[0-9]{3}$`}
	val := g.Generate(schema, "code")
	s, ok := val.(string)
	assert.True(t, ok)
	assert.Regexp(t, regexp.MustCompile(`^[0-9]{3}$`), s, "generated value should match pattern")
}

func TestGenerator_Generate_PatternFallback(t *testing.T) {
	g := NewGenerator(nil)
	// Use a pattern that our simple generator cannot satisfy but is a valid regex
	// The fallback should still return a non-empty string without panicking
	schema := &spec.Schema{Type: "string", Pattern: `^(?:foo|bar){2,3}baz$`}
	val := g.Generate(schema, "complex")
	s, ok := val.(string)
	assert.True(t, ok)
	assert.NotEmpty(t, s, "fallback should return non-empty string")
}
