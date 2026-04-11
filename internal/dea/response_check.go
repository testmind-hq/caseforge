// internal/dea/response_check.go
package dea

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/testmind-hq/caseforge/internal/spec"
)

// validateProbeResponse checks whether ev.ActualBody conforms to the response
// schema declared in op for ev.ActualStatus. Returns nil when:
//   - no JSON schema is declared for that status code (or any 2xx fallback)
//   - the body is empty
//   - the body passes all declared constraints
func validateProbeResponse(op *spec.Operation, probe Probe, ev *Evidence) *DiscoveredRule {
	if ev.ActualBody == "" {
		return nil
	}
	respSchema := findResponseSchema(op, ev.ActualStatus)
	if respSchema == nil {
		return nil
	}
	violations := checkResponseBody(ev.ActualBody, respSchema)
	if len(violations) == 0 {
		return nil
	}
	return &DiscoveredRule{
		ID:        fmt.Sprintf("RULE-%s", strings.ToUpper(fmt.Sprintf("%x", hashID(op.Path+probe.Method))[:6])),
		Operation: fmt.Sprintf("%s %s", probe.Method, probe.Path),
		Category:  CategoryResponseSchemaMismatch,
		Description: fmt.Sprintf("Response body does not match declared schema: %s",
			strings.Join(violations, "; ")),
		Implicit:   false,
		Confidence: ConfidenceMedium,
		Evidence: []RuleEvidence{{
			ProbeDescription: "response schema validation",
			ActualStatus:     ev.ActualStatus,
			ActualBody:       ev.ActualBody,
		}},
	}
}

// hashID produces a simple uint32 hash for deterministic ID generation.
func hashID(s string) uint32 {
	h := uint32(2166136261)
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}

// findResponseSchema returns the application/json schema for the given status code.
// Falls back to any 2xx response schema when the exact code is not declared.
func findResponseSchema(op *spec.Operation, statusCode int) *spec.Schema {
	key := fmt.Sprintf("%d", statusCode)
	if r, ok := op.Responses[key]; ok {
		if mt, ok2 := r.Content["application/json"]; ok2 && mt.Schema != nil {
			return mt.Schema
		}
	}
	// Fallback: any 2xx response with a JSON schema
	for code, r := range op.Responses {
		var n int
		fmt.Sscanf(code, "%d", &n)
		if n >= 200 && n < 300 {
			if mt, ok := r.Content["application/json"]; ok && mt.Schema != nil {
				return mt.Schema
			}
		}
	}
	return nil
}

// checkResponseBody parses body as JSON and validates it against schema.
// Returns human-readable violation messages; nil when body is empty or not JSON.
func checkResponseBody(body string, s *spec.Schema) []string {
	if body == "" {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(body), &m); err != nil {
		// Non-JSON body — skip structural validation (binary, HTML error pages, etc.)
		return nil
	}
	return checkObjectSchema(m, s)
}

// checkObjectSchema validates a parsed JSON object against s.
// Only validates "object" schemas; returns nil for non-object types.
func checkObjectSchema(obj map[string]any, s *spec.Schema) []string {
	if s == nil || s.Type != "object" {
		return nil
	}
	var violations []string
	for _, req := range s.Required {
		if _, ok := obj[req]; !ok {
			violations = append(violations, fmt.Sprintf("required field %q absent in response", req))
		}
	}
	for fieldName, fieldSchema := range s.Properties {
		val, ok := obj[fieldName]
		if !ok || val == nil {
			continue
		}
		if v := checkFieldType(val, fieldSchema); v != "" {
			violations = append(violations, fmt.Sprintf("field %q: %s", fieldName, v))
		}
	}
	return violations
}

// checkFieldType returns a violation message if val does not match s.Type, or "".
func checkFieldType(val any, s *spec.Schema) string {
	if s == nil {
		return ""
	}
	switch s.Type {
	case "string":
		if _, ok := val.(string); !ok {
			return fmt.Sprintf("expected string, got %T", val)
		}
	case "integer", "number":
		if _, ok := val.(float64); !ok {
			return fmt.Sprintf("expected number, got %T", val)
		}
	case "boolean":
		if _, ok := val.(bool); !ok {
			return fmt.Sprintf("expected boolean, got %T", val)
		}
	case "array":
		if _, ok := val.([]any); !ok {
			return fmt.Sprintf("expected array, got %T", val)
		}
	case "object":
		if _, ok := val.(map[string]any); !ok {
			return fmt.Sprintf("expected object, got %T", val)
		}
	}
	return ""
}
