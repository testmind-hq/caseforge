// internal/assert/schema.go
package assert

import (
	"fmt"

	"github.com/testmind-hq/caseforge/internal/output/schema"
	specpkg "github.com/testmind-hq/caseforge/internal/spec"
)

// SchemaAssertions returns assertions for each property of a response object schema.
// The assertion operator is chosen based on the field's JSON Schema format:
//
//   - format "uuid"                    → is_uuid
//   - format "date-time"/"date"/"time" → is_iso8601
//   - format "email"                   → matches (RFC 5322 simplified pattern)
//   - format "uri"/"url"/"uri-reference" → matches (http/https URL pattern)
//   - format "ipv4"                    → matches (dotted-decimal IPv4)
//   - format "ipv6"                    → matches (hex-colon IPv6)
//   - otherwise                        → exists
//
// prefix is the JSON path prefix (e.g. "body").
func SchemaAssertions(prefix string, s *specpkg.Schema) []schema.Assertion {
	if s == nil || s.Type != "object" {
		return nil
	}
	var assertions []schema.Assertion
	for fieldName, fieldSchema := range s.Properties {
		target := fmt.Sprintf("%s.%s", prefix, fieldName)
		assertions = append(assertions, assertionForSchema(fieldSchema, target))
	}
	return assertions
}

// RangeAssertions returns gte/lte assertions for numeric properties that declare
// minimum or maximum constraints in the JSON Schema.
func RangeAssertions(prefix string, s *specpkg.Schema) []schema.Assertion {
	if s == nil || s.Type != "object" {
		return nil
	}
	var out []schema.Assertion
	for fieldName, fieldSchema := range s.Properties {
		if fieldSchema == nil {
			continue
		}
		target := fmt.Sprintf("%s.%s", prefix, fieldName)
		if fieldSchema.Minimum != nil {
			out = append(out, schema.Assertion{
				Target:   target,
				Operator: schema.OperatorGte,
				Expected: *fieldSchema.Minimum,
			})
		}
		if fieldSchema.Maximum != nil {
			out = append(out, schema.Assertion{
				Target:   target,
				Operator: schema.OperatorLte,
				Expected: *fieldSchema.Maximum,
			})
		}
	}
	return out
}

// assertionForSchema builds a single assertion for a field, selecting the
// appropriate operator and expected value based on the field's format.
func assertionForSchema(s *specpkg.Schema, target string) schema.Assertion {
	if s == nil {
		return schema.Assertion{Target: target, Operator: schema.OperatorExists}
	}
	switch s.Format {
	case "uuid":
		return schema.Assertion{Target: target, Operator: schema.OperatorIsUUID}
	case "date-time", "date", "time":
		return schema.Assertion{Target: target, Operator: schema.OperatorIsISO8601}
	case "email":
		return schema.Assertion{Target: target, Operator: schema.OperatorMatches, Expected: `^.+@.+\..+$`}
	case "uri", "url", "uri-reference":
		return schema.Assertion{Target: target, Operator: schema.OperatorMatches, Expected: `^https?://.+`}
	case "ipv4":
		return schema.Assertion{Target: target, Operator: schema.OperatorMatches, Expected: `^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`}
	case "ipv6":
		return schema.Assertion{Target: target, Operator: schema.OperatorMatches, Expected: `^[0-9a-fA-F:]+$`}
	default:
		return schema.Assertion{Target: target, Operator: schema.OperatorExists}
	}
}
