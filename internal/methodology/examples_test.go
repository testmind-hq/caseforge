// internal/methodology/examples_test.go
package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func makeExampleOp(mediaType *spec.MediaType) *spec.Operation {
	return &spec.Operation{
		Method: "POST",
		Path:   "/items",
		RequestBody: &spec.RequestBody{
			Required: true,
			Content:  map[string]*spec.MediaType{"application/json": mediaType},
		},
		Responses: map[string]*spec.Response{
			"201": {Description: "created"},
		},
	}
}

func TestExampleTechnique_AppliesWhenMediaTypeExample(t *testing.T) {
	et := NewExampleTechnique()
	op := makeExampleOp(&spec.MediaType{
		Schema:  &spec.Schema{Type: "object"},
		Example: map[string]any{"name": "Alice"},
	})
	assert.True(t, et.Applies(op))
}

func TestExampleTechnique_AppliesWhenNamedExamples(t *testing.T) {
	et := NewExampleTechnique()
	op := makeExampleOp(&spec.MediaType{
		Schema: &spec.Schema{Type: "object"},
		Examples: map[string]*spec.Example{
			"ex1": {Value: map[string]any{"name": "Bob"}},
		},
	})
	assert.True(t, et.Applies(op))
}

func TestExampleTechnique_AppliesWhenSchemaExample(t *testing.T) {
	et := NewExampleTechnique()
	op := makeExampleOp(&spec.MediaType{
		Schema: &spec.Schema{
			Type:    "object",
			Example: map[string]any{"name": "Charlie"},
		},
	})
	assert.True(t, et.Applies(op))
}

func TestExampleTechnique_DoesNotApplyWithoutExamples(t *testing.T) {
	et := NewExampleTechnique()
	op := makeExampleOp(&spec.MediaType{
		Schema: &spec.Schema{Type: "object"},
	})
	assert.False(t, et.Applies(op))
}

func TestExampleTechnique_DoesNotApplyWithoutRequestBody(t *testing.T) {
	et := NewExampleTechnique()
	op := &spec.Operation{Method: "GET", Path: "/items"}
	assert.False(t, et.Applies(op))
}

func TestExampleTechnique_ValidExampleExpects2xx(t *testing.T) {
	et := NewExampleTechnique()
	op := makeExampleOp(&spec.MediaType{
		Schema: &spec.Schema{
			Type:     "object",
			Required: []string{"name"},
			Properties: map[string]*spec.Schema{
				"name": {Type: "string"},
			},
		},
		Examples: map[string]*spec.Example{
			"valid": {
				Summary: "A valid item",
				Value:   map[string]any{"name": "Widget"},
			},
		},
	})
	cases, err := et.Generate(op)
	require.NoError(t, err)
	require.Len(t, cases, 1)
	tc := cases[0]
	assert.Equal(t, "P1", tc.Priority)
	assert.Equal(t, "example_extraction", tc.Source.Technique)
	// Should have at least the status_code assertion (basic assertions)
	require.NotEmpty(t, tc.Steps[0].Assertions)
}

func TestExampleTechnique_InvalidExampleExpects4xx(t *testing.T) {
	et := NewExampleTechnique()
	op := makeExampleOp(&spec.MediaType{
		Schema: &spec.Schema{
			Type:     "object",
			Required: []string{"name", "email"},
			Properties: map[string]*spec.Schema{
				"name":  {Type: "string"},
				"email": {Type: "string"},
			},
		},
		Examples: map[string]*spec.Example{
			"missing_email": {
				Summary: "Invalid: missing email",
				Value:   map[string]any{"name": "Widget"}, // email missing
			},
		},
	})
	cases, err := et.Generate(op)
	require.NoError(t, err)
	require.Len(t, cases, 1)
	tc := cases[0]
	assert.Equal(t, "P2", tc.Priority)
	require.Len(t, tc.Steps[0].Assertions, 1)
	assert.Equal(t, 422, tc.Steps[0].Assertions[0].Expected)
}

func TestExampleTechnique_NamedExamplesInDeterministicOrder(t *testing.T) {
	et := NewExampleTechnique()
	op := makeExampleOp(&spec.MediaType{
		Schema: &spec.Schema{Type: "object"},
		Examples: map[string]*spec.Example{
			"zebra":  {Value: map[string]any{"x": "z"}},
			"apple":  {Value: map[string]any{"x": "a"}},
			"mango":  {Value: map[string]any{"x": "m"}},
		},
	})
	cases, err := et.Generate(op)
	require.NoError(t, err)
	require.Len(t, cases, 3)
	// Named examples should be sorted alphabetically.
	assert.Contains(t, cases[0].Title, "apple")
	assert.Contains(t, cases[1].Title, "mango")
	assert.Contains(t, cases[2].Title, "zebra")
}

func TestExampleTechnique_SchemaExampleFallback(t *testing.T) {
	et := NewExampleTechnique()
	op := makeExampleOp(&spec.MediaType{
		Schema: &spec.Schema{
			Type:    "object",
			Example: map[string]any{"name": "FallbackItem"},
		},
	})
	cases, err := et.Generate(op)
	require.NoError(t, err)
	require.Len(t, cases, 1)
	assert.Contains(t, cases[0].Title, "schema")
}

func TestExampleTechnique_MediaTypeExampleTakesPrecedenceOverSchemaExample(t *testing.T) {
	et := NewExampleTechnique()
	op := makeExampleOp(&spec.MediaType{
		Schema: &spec.Schema{
			Type:    "object",
			Example: map[string]any{"name": "SchemaLevel"},
		},
		Example: map[string]any{"name": "MediaTypeLevel"},
	})
	cases, err := et.Generate(op)
	require.NoError(t, err)
	// mediaType.example is present, so schema.example fallback should NOT fire — 1 case only.
	require.Len(t, cases, 1)
	assert.Contains(t, cases[0].Title, "inline")
}
