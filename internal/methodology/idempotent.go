// internal/methodology/idempotent.go
package methodology

import (
	"fmt"

	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

type IdempotentTechnique struct {
	gen *datagen.Generator
}

func NewIdempotentTechnique() *IdempotentTechnique {
	return &IdempotentTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *IdempotentTechnique) Name() string { return "idempotency" }

// Applies to write methods (POST/PUT/DELETE).
func (t *IdempotentTechnique) Applies(op *spec.Operation) bool {
	m := op.Method
	return m == "POST" || m == "PUT" || m == "DELETE"
}

func (t *IdempotentTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	body := t.buildValidBody(op)
	tc := buildTestCase(op, body,
		"duplicate_request",
		"send identical request twice — second should not create duplicate",
		fmt.Sprintf("%s %s", op.Method, op.Path))
	tc.Priority = "P2"
	tc.Labels = map[string]string{"type": "idempotency"}
	tc.Source = schema.CaseSource{
		Technique: "idempotency",
		SpecPath:  fmt.Sprintf("%s %s", op.Method, op.Path),
		Rationale: fmt.Sprintf("%s is a write operation; test that repeat calls are safe", op.Method),
	}
	// Note: true idempotency testing needs a second step (Phase 2 chain support)
	// Phase 1: generate the single-request case with a comment in the title
	return []schema.TestCase{tc}, nil
}

func (t *IdempotentTechnique) buildValidBody(op *spec.Operation) map[string]any {
	s := getJSONSchema(op.RequestBody)
	if s == nil {
		return nil
	}
	body := map[string]any{}
	for name, fieldSchema := range s.Properties {
		body[name] = t.gen.Generate(fieldSchema, name)
	}
	return body
}
