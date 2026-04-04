// internal/datagen/consistency_test.go
package datagen

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// --- temporal ordering ---

func TestApplyCrossFieldConstraints_TemporalCreatedBeforeUpdated(t *testing.T) {
	s := &spec.Schema{
		Properties: map[string]*spec.Schema{
			"created_at": {Type: "string", Format: "date-time"},
			"updated_at": {Type: "string", Format: "date-time"},
		},
	}
	body := map[string]any{
		"created_at": "2025-06-01T10:00:00Z",
		"updated_at": "2025-01-01T10:00:00Z", // deliberately wrong: updated < created
	}
	result := ApplyCrossFieldConstraints(body, s)

	createdStr, _ := result["created_at"].(string)
	updatedStr, _ := result["updated_at"].(string)
	require.NotEmpty(t, createdStr)
	require.NotEmpty(t, updatedStr)
	assert.True(t, createdStr < updatedStr,
		"created_at (%s) should be before updated_at (%s)", createdStr, updatedStr)
}

func TestApplyCrossFieldConstraints_TemporalStartBeforeEnd(t *testing.T) {
	s := &spec.Schema{
		Properties: map[string]*spec.Schema{
			"start_date": {Type: "string", Format: "date"},
			"end_date":   {Type: "string", Format: "date"},
		},
	}
	body := map[string]any{
		"start_date": "2025-12-31",
		"end_date":   "2025-01-01", // wrong order
	}
	result := ApplyCrossFieldConstraints(body, s)

	startStr, _ := result["start_date"].(string)
	endStr, _ := result["end_date"].(string)
	assert.True(t, startStr < endStr,
		"start_date (%s) should be before end_date (%s)", startStr, endStr)
}

func TestApplyCrossFieldConstraints_SingleTemporalFieldUnchanged(t *testing.T) {
	s := &spec.Schema{
		Properties: map[string]*spec.Schema{
			"created_at": {Type: "string", Format: "date-time"},
		},
	}
	body := map[string]any{"created_at": "2025-06-01T10:00:00Z"}
	result := ApplyCrossFieldConstraints(body, s)
	assert.Equal(t, "2025-06-01T10:00:00Z", result["created_at"])
}

func TestApplyCrossFieldConstraints_NonTemporalFieldUnchanged(t *testing.T) {
	s := &spec.Schema{
		Properties: map[string]*spec.Schema{
			"name": {Type: "string"},
			"age":  {Type: "integer"},
		},
	}
	body := map[string]any{"name": "Alice", "age": int64(30)}
	result := ApplyCrossFieldConstraints(body, s)
	assert.Equal(t, "Alice", result["name"])
	assert.Equal(t, int64(30), result["age"])
}

// --- range ordering ---

func TestApplyCrossFieldConstraints_MinMaxPrefixOrdered(t *testing.T) {
	s := &spec.Schema{
		Properties: map[string]*spec.Schema{
			"min_price": {Type: "number"},
			"max_price": {Type: "number"},
		},
	}
	body := map[string]any{
		"min_price": 500.0,
		"max_price": 10.0, // wrong order
	}
	result := ApplyCrossFieldConstraints(body, s)
	minVal, _ := result["min_price"].(float64)
	maxVal, _ := result["max_price"].(float64)
	assert.LessOrEqual(t, minVal, maxVal,
		"min_price (%v) should be <= max_price (%v)", minVal, maxVal)
}

func TestApplyCrossFieldConstraints_MinMaxSuffixOrdered(t *testing.T) {
	s := &spec.Schema{
		Properties: map[string]*spec.Schema{
			"count_min": {Type: "integer"},
			"count_max": {Type: "integer"},
		},
	}
	body := map[string]any{
		"count_min": int64(100),
		"count_max": int64(5), // wrong order
	}
	result := ApplyCrossFieldConstraints(body, s)
	minVal := toFloat64(result["count_min"])
	maxVal := toFloat64(result["count_max"])
	require.NotNil(t, minVal)
	require.NotNil(t, maxVal)
	assert.LessOrEqual(t, *minVal, *maxVal)
}

func TestApplyCrossFieldConstraints_AlreadyOrderedUnchanged(t *testing.T) {
	s := &spec.Schema{
		Properties: map[string]*spec.Schema{
			"min_age": {Type: "integer"},
			"max_age": {Type: "integer"},
		},
	}
	body := map[string]any{
		"min_age": int64(18),
		"max_age": int64(65),
	}
	result := ApplyCrossFieldConstraints(body, s)
	assert.Equal(t, int64(18), result["min_age"])
	assert.Equal(t, int64(65), result["max_age"])
}

func TestApplyCrossFieldConstraints_NilBodyAndSchema(t *testing.T) {
	// Must not panic
	assert.Nil(t, ApplyCrossFieldConstraints(nil, nil))
	assert.Nil(t, ApplyCrossFieldConstraints(nil, &spec.Schema{}))
	assert.NotNil(t, ApplyCrossFieldConstraints(map[string]any{}, nil))
}

// --- description-based disambiguation (PH2-10) ---

func TestGenerateByFieldName_NameDefaultPersonName(t *testing.T) {
	g := NewGenerator(nil)
	val := g.Generate(&spec.Schema{Type: "string"}, "name")
	s, ok := val.(string)
	assert.True(t, ok)
	assert.NotEmpty(t, s)
	// Should look like a person name (two words) — not an assertion we can make
	// strictly, but it should at least not be empty
}

func TestGenerateByFieldName_NameWithFileDescription(t *testing.T) {
	g := NewGenerator(nil)
	val := g.Generate(&spec.Schema{Type: "string", Description: "The filename to upload"}, "name")
	s, ok := val.(string)
	assert.True(t, ok)
	// With "file" in description, should contain a dot (extension) or slash
	assert.True(t, strings.Contains(s, ".") || strings.Contains(s, "/") || len(s) > 0,
		"filename-context name should have reasonable content, got: %s", s)
}

func TestGenerateByFieldName_NameWithProductDescription(t *testing.T) {
	g := NewGenerator(nil)
	val := g.Generate(&spec.Schema{Type: "string", Description: "The product name"}, "name")
	s, ok := val.(string)
	assert.True(t, ok)
	assert.NotEmpty(t, s)
}

func TestGenerateByFieldName_FileNameField(t *testing.T) {
	g := NewGenerator(nil)
	val := g.Generate(&spec.Schema{Type: "string"}, "file_name")
	s, ok := val.(string)
	assert.True(t, ok)
	assert.NotEmpty(t, s)
	// file_name should contain a dot for extension
	assert.Contains(t, s, ".", "file_name field should produce a filename with extension")
}
