// internal/assert/schema.go
package assert

import (
	"fmt"

	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// SchemaAssertions returns field-existence assertions for a response schema.
// prefix is the JSON path prefix (e.g. "body").
func SchemaAssertions(prefix string, s *spec.Schema) []schema.Assertion {
	if s == nil || s.Type != "object" {
		return nil
	}
	var assertions []schema.Assertion
	for fieldName := range s.Properties {
		assertions = append(assertions, schema.Assertion{
			Target:   fmt.Sprintf("%s.%s", prefix, fieldName),
			Operator: "exists",
			Expected: true,
		})
	}
	return assertions
}
