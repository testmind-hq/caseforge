// internal/datagen/generator_test.go
package datagen

import (
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
