// internal/sandbox/strategies.go
package sandbox

import (
	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// ResponseStrategy generates a response body for an operation and status code.
// Returns (body, true) when it can produce a value; (nil, false) to pass to next strategy.
type ResponseStrategy interface {
	Name() string
	Generate(op *spec.Operation, statusCode string) (any, bool)
}

// successResponse returns the first 2xx response found, preferring 200 → 201 → 202.
func successResponse(op *spec.Operation) *spec.Response {
	for _, code := range []string{"200", "201", "202"} {
		if resp, ok := op.Responses[code]; ok {
			return resp
		}
	}
	return nil
}

// --- ExampleStrategy ---

type exampleStrategy struct{}

func (e *exampleStrategy) Name() string { return "example" }

func (e *exampleStrategy) Generate(op *spec.Operation, statusCode string) (any, bool) {
	resp, ok := op.Responses[statusCode]
	if !ok {
		resp = successResponse(op)
		if resp == nil {
			return nil, false
		}
	}
	mt, ok := resp.Content["application/json"]
	if !ok {
		return nil, false
	}
	if mt.Example != nil {
		return mt.Example, true
	}
	for _, ex := range mt.Examples {
		if ex != nil && ex.Value != nil {
			return ex.Value, true
		}
	}
	return nil, false
}

// --- SchemaStrategy ---

type schemaStrategy struct{}

func (s *schemaStrategy) Name() string { return "schema" }

func (s *schemaStrategy) Generate(op *spec.Operation, statusCode string) (any, bool) {
	resp, ok := op.Responses[statusCode]
	if !ok {
		resp = successResponse(op)
		if resp == nil {
			return nil, false
		}
	}
	mt, ok := resp.Content["application/json"]
	if !ok || mt == nil || mt.Schema == nil {
		return nil, false
	}
	return schemaZeroValue(mt.Schema), true
}

func schemaZeroValue(s *spec.Schema) any {
	if s == nil {
		return nil
	}
	switch s.Type {
	case "object":
		obj := map[string]any{}
		for name, prop := range s.Properties {
			obj[name] = schemaZeroValue(prop)
		}
		return obj
	case "array":
		if s.Items != nil {
			return []any{schemaZeroValue(s.Items)}
		}
		return []any{}
	case "integer":
		return int64(0)
	case "number":
		return float64(0)
	case "boolean":
		return false
	default: // "string" and unknown
		return ""
	}
}

// --- FakerStrategy ---

type fakerStrategy struct {
	gen *datagen.Generator
}

func newFakerStrategy() *fakerStrategy {
	return &fakerStrategy{gen: datagen.NewGenerator(nil)}
}

func (f *fakerStrategy) Name() string { return "faker" }

func (f *fakerStrategy) Generate(op *spec.Operation, statusCode string) (any, bool) {
	resp, ok := op.Responses[statusCode]
	if !ok {
		resp = successResponse(op)
	}
	if resp == nil {
		return map[string]any{}, true
	}
	mt, ok := resp.Content["application/json"]
	if !ok || mt == nil {
		return map[string]any{}, true
	}
	return f.generateFromSchema(mt.Schema), true
}

func (f *fakerStrategy) generateFromSchema(s *spec.Schema) any {
	if s == nil {
		return map[string]any{}
	}
	if s.Type == "object" {
		obj := map[string]any{}
		for name, prop := range s.Properties {
			obj[name] = f.gen.Generate(prop, name)
		}
		return obj
	}
	if s.Type == "array" {
		if s.Items != nil {
			return []any{f.generateFromSchema(s.Items)}
		}
		return []any{}
	}
	return f.gen.Generate(s, "")
}
