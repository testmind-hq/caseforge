// internal/datagen/generator.go
package datagen

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/testmind-hq/caseforge/internal/spec"
)

type BoundaryKind int

const (
	BoundaryMin        BoundaryKind = iota
	BoundaryMinMinusOne             // min - 1 (invalid)
	BoundaryMax
	BoundaryMaxPlusOne // max + 1 (invalid)
)

type Generator struct {
	// llm reserved for Phase 2 semantic data generation
	llm any
}

func NewGenerator(llm any) *Generator {
	return &Generator{llm: llm}
}

// Generate produces a valid fake value for the given schema and field name.
// Three-tier fallback: format → field name → random by type.
func (g *Generator) Generate(s *spec.Schema, fieldName string) any {
	if s == nil {
		return gofakeit.Word()
	}

	// Tier 0: enum — always pick from enum
	if len(s.Enum) > 0 {
		return s.Enum[gofakeit.Number(0, len(s.Enum)-1)]
	}

	// Tier 1: format-aware
	if s.Format != "" {
		if val, ok := generateByFormat(s.Format); ok {
			return val
		}
	}

	// Tier 2: field name semantic
	if fieldName != "" {
		if val, ok := generateByFieldName(fieldName); ok {
			return val
		}
	}

	// Tier 3: fallback by type
	return generateByType(s)
}

func generateByType(s *spec.Schema) any {
	switch s.Type {
	case "string":
		return gofakeit.Word()
	case "integer":
		min, max := int64(0), int64(1000)
		if s.Minimum != nil {
			min = int64(*s.Minimum)
		}
		if s.Maximum != nil {
			max = int64(*s.Maximum)
		}
		return int64(gofakeit.Number(int(min), int(max)))
	case "number":
		return gofakeit.Float64Range(0, 1000)
	case "boolean":
		return gofakeit.Bool()
	case "array":
		// Return a single-element array
		if s.Items != nil {
			return []any{generateByType(s.Items)}
		}
		return []any{}
	case "object":
		result := map[string]any{}
		for name, prop := range s.Properties {
			result[name] = generateByType(prop)
		}
		return result
	default:
		return gofakeit.Word()
	}
}

// GenerateBoundary produces a boundary value for numeric/string schemas.
func (g *Generator) GenerateBoundary(s *spec.Schema, kind BoundaryKind) any {
	if s == nil {
		return nil
	}
	switch s.Type {
	case "integer":
		switch kind {
		case BoundaryMin:
			if s.Minimum != nil {
				return int64(*s.Minimum)
			}
			return int64(0)
		case BoundaryMinMinusOne:
			if s.Minimum != nil {
				return int64(*s.Minimum) - 1
			}
			return int64(-1)
		case BoundaryMax:
			if s.Maximum != nil {
				return int64(*s.Maximum)
			}
			return int64(1000)
		case BoundaryMaxPlusOne:
			if s.Maximum != nil {
				return int64(*s.Maximum) + 1
			}
			return int64(1001)
		}
	case "string":
		switch kind {
		case BoundaryMin:
			if s.MinLength != nil {
				return gofakeit.LetterN(uint(*s.MinLength))
			}
			return ""
		case BoundaryMinMinusOne:
			if s.MinLength != nil && *s.MinLength > 0 {
				return gofakeit.LetterN(uint(*s.MinLength - 1))
			}
			return ""
		case BoundaryMax:
			if s.MaxLength != nil {
				return gofakeit.LetterN(uint(*s.MaxLength))
			}
			return gofakeit.LetterN(255)
		case BoundaryMaxPlusOne:
			if s.MaxLength != nil {
				return gofakeit.LetterN(uint(*s.MaxLength + 1))
			}
			return gofakeit.LetterN(256)
		}
	}
	return nil
}
