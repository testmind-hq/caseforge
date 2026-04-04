// internal/spec/validate_test.go
package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateExample_ValidObject(t *testing.T) {
	s := &Schema{
		Type:     "object",
		Required: []string{"name", "email"},
		Properties: map[string]*Schema{
			"name":  {Type: "string"},
			"email": {Type: "string"},
			"age":   {Type: "integer"},
		},
	}
	value := map[string]any{
		"name":  "Alice",
		"email": "alice@example.com",
		"age":   float64(30),
	}
	errs := ValidateExample(value, s)
	assert.Empty(t, errs, "valid object should have no errors")
}

func TestValidateExample_MissingRequiredField(t *testing.T) {
	s := &Schema{
		Type:     "object",
		Required: []string{"name", "email"},
		Properties: map[string]*Schema{
			"name":  {Type: "string"},
			"email": {Type: "string"},
		},
	}
	value := map[string]any{"name": "Alice"} // missing email
	errs := ValidateExample(value, s)
	assert.Len(t, errs, 1)
	assert.Contains(t, errs[0], "email")
}

func TestValidateExample_WrongType(t *testing.T) {
	s := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"age": {Type: "integer"},
		},
	}
	value := map[string]any{"age": "not-a-number"} // wrong type
	errs := ValidateExample(value, s)
	assert.Len(t, errs, 1)
	assert.Contains(t, errs[0], "age")
}

func TestValidateExample_NilSchema(t *testing.T) {
	errs := ValidateExample(map[string]any{"x": 1}, nil)
	assert.Empty(t, errs)
}

func TestValidateExample_NilValue(t *testing.T) {
	errs := ValidateExample(nil, &Schema{Type: "object"})
	assert.Empty(t, errs)
}

func TestValidateExample_UnknownFieldsAllowed(t *testing.T) {
	s := &Schema{
		Type:       "object",
		Properties: map[string]*Schema{"name": {Type: "string"}},
	}
	value := map[string]any{"name": "Alice", "extra": "unexpected"}
	errs := ValidateExample(value, s)
	assert.Empty(t, errs, "unknown fields should not be reported as errors")
}

func TestValidateExample_MultipleMissingRequired(t *testing.T) {
	s := &Schema{
		Type:     "object",
		Required: []string{"a", "b", "c"},
		Properties: map[string]*Schema{
			"a": {Type: "string"},
			"b": {Type: "string"},
			"c": {Type: "string"},
		},
	}
	value := map[string]any{} // all required fields missing
	errs := ValidateExample(value, s)
	assert.Len(t, errs, 3)
}
