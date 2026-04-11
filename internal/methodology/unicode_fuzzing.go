// internal/methodology/unicode_fuzzing.go
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

// UnicodeFuzzingTechnique generates test cases that inject problematic Unicode
// strings into string fields in the request body.
//
// Mutations per string field (inspired by CATS fuzzer):
//  1. control_char   — null byte embedded mid-string: "hello\u0000world"
//  2. zero_width     — zero-width space prefix: "\u200Bhello"
//  3. bidi_override  — right-to-left override: "\u202Ehello"
//  4. overlong       — >10k char string (tests implicit max-length)
//  5. zalgo          — combining diacriticals: zalgo text
//
// All cases: Priority P3, assert status_code eq 400, Scenario "UNICODE_INJECTION".
type UnicodeFuzzingTechnique struct {
	gen *datagen.Generator
}

// NewUnicodeFuzzingTechnique returns a new UnicodeFuzzingTechnique.
func NewUnicodeFuzzingTechnique() *UnicodeFuzzingTechnique {
	return &UnicodeFuzzingTechnique{gen: datagen.NewGenerator(nil)}
}

func (u *UnicodeFuzzingTechnique) Name() string { return "unicode_fuzzing" }

// Applies returns true if the operation has a JSON request body with at least
// one string-typed property.
func (u *UnicodeFuzzingTechnique) Applies(op *spec.Operation) bool {
	s := getJSONSchema(op.RequestBody)
	if s == nil {
		return false
	}
	for _, fieldSchema := range s.Properties {
		if fieldSchema != nil && fieldSchema.Type == "string" {
			return true
		}
	}
	return false
}

// unicodeMutation describes a single unicode fuzzing mutation.
type unicodeMutation struct {
	value string
	label string
}

// unicodeMutations returns the 5 unicode mutations to apply to a string field.
func unicodeMutations() []unicodeMutation {
	return []unicodeMutation{
		{value: "hello\u0000world", label: "control_char"},
		{value: "\u200Bhello", label: "zero_width"},
		{value: "\u202Ehello", label: "bidi_override"},
		{value: strings.Repeat("A", 10001), label: "overlong"},
		{value: "z\u0300\u0301\u0302\u0303\u0304\u0305\u0306\u0307a", label: "zalgo"},
	}
}

// Generate produces 5 test cases per string field (one per unicode mutation).
func (u *UnicodeFuzzingTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	s := getJSONSchema(op.RequestBody)
	if s == nil {
		return nil, nil
	}

	validBase := buildValidBody(u.gen, op)
	if validBase == nil {
		validBase = map[string]any{}
	}

	var cases []schema.TestCase

	for _, fieldName := range slices.Sorted(maps.Keys(s.Properties)) {
		fieldSchema := s.Properties[fieldName]
		if fieldSchema == nil || fieldSchema.Type != "string" {
			continue // skip non-string fields
		}

		for _, m := range unicodeMutations() {
			body := copyMap(validBase)
			body[fieldName] = m.value

			specPath := fmt.Sprintf("%s %s requestBody.properties.%s", op.Method, op.Path, fieldName)
			tc := buildTestCase(op, body,
				fmt.Sprintf("[unicode_fuzzing] %s %s", fieldName, m.label), specPath)
			tc.Priority = "P3"
			tc.Steps[0].Assertions = []schema.Assertion{
				{Target: "status_code", Operator: "eq", Expected: 400},
			}
			tc.Source = schema.CaseSource{
				Technique: "unicode_fuzzing",
				SpecPath:  specPath,
				Rationale: fmt.Sprintf("field %q receives unicode mutation %q — server must sanitize/reject with 400", fieldName, m.label),
				Scenario:  "UNICODE_INJECTION",
			}
			cases = append(cases, tc)
		}
	}

	return cases, nil
}
