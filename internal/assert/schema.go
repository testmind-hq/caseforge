// internal/assert/schema.go
package assert

import (
	"fmt"

	"github.com/testmind-hq/caseforge/internal/output/schema"
	specpkg "github.com/testmind-hq/caseforge/internal/spec"
)

// SchemaAssertions returns assertions for each property of a response object schema.
// It selects the operator based on the field's JSON Schema format:
//   - format "uuid"      → is_uuid
//   - format "date-time" → is_iso8601
//   - otherwise          → exists
//
// prefix is the JSON path prefix (e.g. "body").
func SchemaAssertions(prefix string, s *specpkg.Schema) []schema.Assertion {
	if s == nil || s.Type != "object" {
		return nil
	}
	var assertions []schema.Assertion
	for fieldName, fieldSchema := range s.Properties {
		op := operatorForFormat(fieldSchema)
		a := schema.Assertion{
			Target:   fmt.Sprintf("%s.%s", prefix, fieldName),
			Operator: op,
		}
		assertions = append(assertions, a)
	}
	return assertions
}

// operatorForFormat maps a JSON Schema format string to the appropriate assertion operator.
func operatorForFormat(s *specpkg.Schema) string {
	if s == nil {
		return schema.OperatorExists
	}
	switch s.Format {
	case "uuid":
		return schema.OperatorIsUUID
	case "date-time", "date", "time":
		return schema.OperatorIsISO8601
	default:
		return schema.OperatorExists
	}
}
