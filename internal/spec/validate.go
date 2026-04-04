// internal/spec/validate.go
package spec

import (
	"fmt"
	"math"
)

// ValidateExample checks whether value conforms to the given schema.
// It performs structural validation:
//   - For object schemas: checks that all required fields are present.
//   - For all fields that exist: checks that the value type matches the schema type.
//
// Returns a slice of human-readable error strings. An empty slice means the
// example is valid against the schema.
func ValidateExample(value any, schema *Schema) []string {
	if schema == nil || value == nil {
		return nil
	}
	obj, ok := value.(map[string]any)
	if !ok {
		// Non-object examples (arrays, primitives) — check top-level type only.
		if errs := validateType(value, schema); len(errs) > 0 {
			return errs
		}
		return nil
	}

	var errs []string

	// Check required fields are present.
	for _, req := range schema.Required {
		if _, exists := obj[req]; !exists {
			errs = append(errs, fmt.Sprintf("missing required field %q", req))
		}
	}

	// Check field types for fields present in the example.
	for fieldName, fieldVal := range obj {
		fieldSchema, known := schema.Properties[fieldName]
		if !known || fieldSchema == nil {
			continue // unknown field — not an error (spec may have additionalProperties)
		}
		for _, e := range validateType(fieldVal, fieldSchema) {
			errs = append(errs, fmt.Sprintf("field %q: %s", fieldName, e))
		}
	}

	return errs
}

// validateType checks that v matches the expected schema type.
// Returns nil on success.
func validateType(v any, s *Schema) []string {
	if s.Type == "" || v == nil {
		return nil
	}
	switch s.Type {
	case "string":
		if _, ok := v.(string); !ok {
			return []string{fmt.Sprintf("expected string, got %T", v)}
		}
	case "integer":
		switch val := v.(type) {
		case int, int32, int64:
			// always integral
		case float64: // JSON numbers decode as float64
			if val != math.Trunc(val) {
				return []string{fmt.Sprintf("expected integer, got non-integral float64 %v", val)}
			}
		default:
			return []string{fmt.Sprintf("expected integer, got %T", v)}
		}
	case "number":
		switch v.(type) {
		case int, int32, int64, float32, float64:
		default:
			return []string{fmt.Sprintf("expected number, got %T", v)}
		}
	case "boolean":
		if _, ok := v.(bool); !ok {
			return []string{fmt.Sprintf("expected boolean, got %T", v)}
		}
	case "array":
		if _, ok := v.([]any); !ok {
			return []string{fmt.Sprintf("expected array, got %T", v)}
		}
	case "object":
		if _, ok := v.(map[string]any); !ok {
			return []string{fmt.Sprintf("expected object, got %T", v)}
		}
	}
	return nil
}
