// internal/assert/basic.go
package assert

import (
	"strconv"

	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// BasicAssertions returns status code and duration assertions for an operation.
func BasicAssertions(op *spec.Operation) []schema.Assertion {
	var assertions []schema.Assertion

	// Find the primary success status code
	for code := range op.Responses {
		n, err := strconv.Atoi(code)
		if err == nil && n >= 200 && n < 300 {
			assertions = append(assertions, schema.Assertion{
				Target:   "status_code",
				Operator: "eq",
				Expected: n,
			})
			break
		}
	}
	// Fallback to 200 if no explicit success response defined
	if len(assertions) == 0 {
		assertions = append(assertions, schema.Assertion{
			Target:   "status_code",
			Operator: "eq",
			Expected: 200,
		})
	}

	// Standard performance assertion
	assertions = append(assertions, schema.Assertion{
		Target:   "duration_ms",
		Operator: "lt",
		Expected: 2000,
	})

	return assertions
}
