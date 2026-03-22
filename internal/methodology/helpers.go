// internal/methodology/helpers.go
package methodology

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	assertpkg "github.com/testmind-hq/caseforge/internal/assert"
	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// buildValidBody generates a body map with valid values for all fields in op's JSON schema.
func buildValidBody(gen *datagen.Generator, op *spec.Operation) map[string]any {
	s := getJSONSchema(op.RequestBody)
	if s == nil {
		return nil
	}
	body := map[string]any{}
	for fieldName, fieldSchema := range s.Properties {
		body[fieldName] = gen.Generate(fieldSchema, fieldName)
	}
	return body
}

func buildTestCase(op *spec.Operation, body map[string]any, title, specPath string) schema.TestCase {
	id := fmt.Sprintf("TC-%s", uuid.New().String()[:8])
	headers := map[string]string{}
	if body != nil {
		headers["Content-Type"] = "application/json"
	}

	// Avoid nil-typed-interface issue: only set Body to a non-nil any when body is non-nil.
	// Assigning nil map[string]any to any creates a non-nil interface with a nil value,
	// which would cause the Hurl renderer to emit a spurious "null" body block.
	var bodyAny any
	if body != nil {
		bodyAny = body
	}

	tc := schema.TestCase{
		Schema:   schema.SchemaBaseURL,
		Version:  "1",
		ID:       id,
		Title:    fmt.Sprintf("%s %s - %s", op.Method, op.Path, title),
		Kind:     "single",
		Priority: "P1",
		Tags:     op.Tags,
		Steps: []schema.Step{
			{
				ID:     "step-main",
				Title:  title,
				Type:   "test",
				Method: op.Method,
				Path:   op.Path,
				Headers: headers,
				Body:    bodyAny,
				Assertions: assertpkg.BasicAssertions(op),
			},
		},
		GeneratedAt: time.Now(),
	}
	return tc
}

func getJSONSchema(rb *spec.RequestBody) *spec.Schema {
	if rb == nil {
		return nil
	}
	if mt, ok := rb.Content["application/json"]; ok {
		return mt.Schema
	}
	return nil
}

func responseSchemaAssertions(op *spec.Operation) []schema.Assertion {
	for code, resp := range op.Responses {
		n := 0
		fmt.Sscanf(code, "%d", &n)
		if n >= 200 && n < 300 {
			if mt, ok := resp.Content["application/json"]; ok {
				return assertpkg.SchemaAssertions("body", mt.Schema)
			}
		}
	}
	return nil
}
