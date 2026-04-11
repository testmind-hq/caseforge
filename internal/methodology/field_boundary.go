// internal/methodology/field_boundary.go
package methodology

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// FieldBoundaryTechnique performs a depth-first walk of the entire request body
// schema (including nested objects and arrays) and generates boundary test cases
// for every field that declares Minimum, Maximum, MinLength, or MaxLength.
type FieldBoundaryTechnique struct {
	gen *datagen.Generator
}

func NewFieldBoundaryTechnique() *FieldBoundaryTechnique {
	return &FieldBoundaryTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *FieldBoundaryTechnique) Name() string { return "field_boundary" }

// fieldEntry holds a dot-notation path to a schema field and the field's schema.
type fieldEntry struct {
	dotPath string // e.g. "address.zip" or "items[0].count"
	schema  *spec.Schema
}

// walkSchemaFields does a depth-first traversal of s.Properties, producing
// one fieldEntry per leaf/intermediate field (in sorted order for determinism).
func walkSchemaFields(s *spec.Schema, prefix string) []fieldEntry {
	if s == nil {
		return nil
	}
	var result []fieldEntry
	names := slices.Sorted(maps.Keys(s.Properties))
	for _, name := range names {
		prop := s.Properties[name]
		if prop == nil {
			continue
		}
		path := name
		if prefix != "" {
			path = prefix + "." + name
		}
		result = append(result, fieldEntry{path, prop})
		if prop.Type == "object" && len(prop.Properties) > 0 {
			result = append(result, walkSchemaFields(prop, path)...)
		}
		if prop.Type == "array" && prop.Items != nil && prop.Items.Type == "object" {
			result = append(result, walkSchemaFields(prop.Items, path+"[0]")...)
		}
	}
	return result
}

// anyHasBoundary reports whether any of the fields has at least one boundary constraint.
func anyHasBoundary(fields []fieldEntry) bool {
	for _, f := range fields {
		if f.schema.Minimum != nil || f.schema.Maximum != nil ||
			f.schema.MinLength != nil || f.schema.MaxLength != nil {
			return true
		}
	}
	return false
}

func (t *FieldBoundaryTechnique) Applies(op *spec.Operation) bool {
	s := getJSONSchema(op.RequestBody)
	if s == nil {
		return false
	}
	fields := walkSchemaFields(s, "")
	return len(fields) > 0 && anyHasBoundary(fields)
}

func (t *FieldBoundaryTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	s := getJSONSchema(op.RequestBody)
	if s == nil {
		return nil, nil
	}
	base := buildValidBody(t.gen, op)
	if base == nil {
		base = map[string]any{}
	}

	var cases []schema.TestCase
	for _, entry := range walkSchemaFields(s, "") {
		fs := entry.schema
		if fs == nil {
			continue
		}

		type boundCase struct {
			label    string
			value    any
			valid    bool
			scenario string
		}
		var bcs []boundCase

		switch fs.Type {
		case "integer", "number":
			if fs.Minimum != nil {
				bcs = append(bcs,
					boundCase{"valid_min", t.gen.GenerateBoundary(fs, datagen.BoundaryMin), true, "FIELD_BOUNDARY_VALID"},
					boundCase{"invalid_below_min", t.gen.GenerateBoundary(fs, datagen.BoundaryMinMinusOne), false, "FIELD_BOUNDARY_INVALID"},
				)
			}
			if fs.Maximum != nil {
				bcs = append(bcs,
					boundCase{"valid_max", t.gen.GenerateBoundary(fs, datagen.BoundaryMax), true, "FIELD_BOUNDARY_VALID"},
					boundCase{"invalid_above_max", t.gen.GenerateBoundary(fs, datagen.BoundaryMaxPlusOne), false, "FIELD_BOUNDARY_INVALID"},
				)
			}
		case "string":
			if fs.MinLength != nil {
				bcs = append(bcs,
					boundCase{"valid_min", strings.Repeat("a", int(*fs.MinLength)), true, "FIELD_BOUNDARY_VALID"},
					boundCase{"invalid_below_min", strings.Repeat("a", max(0, int(*fs.MinLength)-1)), false, "FIELD_BOUNDARY_INVALID"},
				)
			}
			if fs.MaxLength != nil {
				bcs = append(bcs,
					boundCase{"valid_max", strings.Repeat("a", int(*fs.MaxLength)), true, "FIELD_BOUNDARY_VALID"},
					boundCase{"invalid_above_max", strings.Repeat("a", int(*fs.MaxLength)+1), false, "FIELD_BOUNDARY_INVALID"},
				)
			}
		}

		for _, bc := range bcs {
			body := maps.Clone(base)
			setAtPath(body, entry.dotPath, bc.value)

			specPath := fmt.Sprintf("%s %s requestBody.%s", op.Method, op.Path, entry.dotPath)
			title := fmt.Sprintf("[field_boundary] %s %s", entry.dotPath, bc.label)
			tc := buildTestCase(op, body, title, specPath)

			if bc.valid {
				tc.Priority = "P1"
				tc.Steps[0].Assertions = []schema.Assertion{
					{Target: "status_code", Operator: "gte", Expected: 200},
					{Target: "status_code", Operator: "lt", Expected: 300},
				}
			} else {
				tc.Priority = "P2"
				tc.Steps[0].Assertions = []schema.Assertion{
					{Target: "status_code", Operator: "gte", Expected: 400},
					{Target: "status_code", Operator: "lt", Expected: 500},
				}
			}
			tc.Source = schema.CaseSource{
				Technique: "field_boundary",
				SpecPath:  specPath,
				Rationale: fmt.Sprintf("field %q boundary test: %s", entry.dotPath, bc.label),
				Scenario:  bc.scenario,
			}
			cases = append(cases, tc)
		}
	}
	return cases, nil
}

// setAtPath sets a value at the given dot-notation path inside body.
// Segments ending with "[0]" set an array with one element.
// Intermediate maps are created as needed.
func setAtPath(body map[string]any, dotPath string, value any) {
	segments := strings.Split(dotPath, ".")
	current := body
	for i, seg := range segments {
		isLast := i == len(segments)-1
		if strings.HasSuffix(seg, "[0]") {
			key := seg[:len(seg)-3]
			if isLast {
				// Set the array with one element
				current[key] = []any{value}
			} else {
				// Ensure an array with a map exists for further traversal
				existing, ok := current[key]
				var inner map[string]any
				if ok {
					if arr, ok2 := existing.([]any); ok2 && len(arr) > 0 {
						if m, ok3 := arr[0].(map[string]any); ok3 {
							inner = m
						}
					}
				}
				if inner == nil {
					inner = map[string]any{}
					current[key] = []any{inner}
				}
				current = inner
			}
		} else {
			if isLast {
				current[seg] = value
			} else {
				existing, ok := current[seg]
				var next map[string]any
				if ok {
					next, _ = existing.(map[string]any)
				}
				if next == nil {
					next = map[string]any{}
					current[seg] = next
				}
				current = next
			}
		}
	}
}

