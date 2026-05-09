// internal/sandbox/strategies_test.go
package sandbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func opWithExample() *spec.Operation {
	return &spec.Operation{
		Method: "GET",
		Path:   "/pets",
		Responses: map[string]*spec.Response{
			"200": {
				Content: map[string]*spec.MediaType{
					"application/json": {
						Example: map[string]any{"id": "abc", "name": "Fido"},
					},
				},
			},
		},
	}
}

func opWithSchema() *spec.Operation {
	return &spec.Operation{
		Method: "GET",
		Path:   "/pets/1",
		Responses: map[string]*spec.Response{
			"200": {
				Content: map[string]*spec.MediaType{
					"application/json": {
						Schema: &spec.Schema{
							Type: "object",
							Properties: map[string]*spec.Schema{
								"id":   {Type: "string", Format: "uuid"},
								"name": {Type: "string"},
							},
						},
					},
				},
			},
		},
	}
}

func opNoResponse() *spec.Operation {
	return &spec.Operation{
		Method:    "DELETE",
		Path:      "/pets/1",
		Responses: map[string]*spec.Response{"204": {Description: "Deleted"}},
	}
}

func TestExampleStrategy_UsesExample(t *testing.T) {
	s := &exampleStrategy{}
	out, ok := s.Generate(opWithExample(), "200")
	assert.True(t, ok)
	m, _ := out.(map[string]any)
	assert.Equal(t, "Fido", m["name"])
}

func TestExampleStrategy_NoExample(t *testing.T) {
	s := &exampleStrategy{}
	_, ok := s.Generate(opWithSchema(), "200")
	assert.False(t, ok)
}

func TestSchemaStrategy_ProducesObject(t *testing.T) {
	s := &schemaStrategy{}
	out, ok := s.Generate(opWithSchema(), "200")
	assert.True(t, ok)
	m, _ := out.(map[string]any)
	assert.Contains(t, m, "id")
	assert.Contains(t, m, "name")
}

func TestSchemaStrategy_NoSchema(t *testing.T) {
	s := &schemaStrategy{}
	_, ok := s.Generate(opNoResponse(), "204")
	assert.False(t, ok)
}

func TestFakerStrategy_AlwaysProduces(t *testing.T) {
	s := newFakerStrategy()
	out, ok := s.Generate(opWithSchema(), "200")
	assert.True(t, ok)
	assert.NotNil(t, out)
}

func TestFakerStrategy_NoSchema(t *testing.T) {
	s := newFakerStrategy()
	out, ok := s.Generate(opNoResponse(), "204")
	assert.True(t, ok) // faker always returns true
	assert.NotNil(t, out)
}
