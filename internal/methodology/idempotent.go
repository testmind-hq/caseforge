// internal/methodology/idempotent.go
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

type IdempotentTechnique struct {
	gen *datagen.Generator
}

func NewIdempotentTechnique() *IdempotentTechnique {
	return &IdempotentTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *IdempotentTechnique) Name() string { return "idempotency" }

// Applies returns true for POST, PUT, and DELETE.
// PATCH is excluded because its idempotency is API-contract-dependent (RFC 5789 does not guarantee it),
// which would cause false-positive idempotency test generation for non-idempotent PATCH endpoints.
func (t *IdempotentTechnique) Applies(op *spec.Operation) bool {
	m := op.Method
	return m == "POST" || m == "PUT" || m == "DELETE"
}

func (t *IdempotentTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	body := buildValidBody(t.gen, op)

	var bodyAny any
	if body != nil {
		bodyAny = body
	}
	headers := map[string]string{}
	if bodyAny != nil {
		headers["Content-Type"] = "application/json"
	}

	assertions := assertpkg.BasicAssertions(op)

	steps := []schema.Step{
		{
			ID:         "step-setup",
			Title:      fmt.Sprintf("%s %s — first call", op.Method, op.Path),
			Type:       "setup",
			Method:     op.Method,
			Path:       op.Path,
			Headers:    headers,
			Body:       bodyAny,
			Assertions: assertions,
		},
		{
			ID:         "step-test",
			Title:      fmt.Sprintf("%s %s — identical second call must be safe", op.Method, op.Path),
			Type:       "test",
			Method:     op.Method,
			Path:       op.Path,
			Headers:    headers,
			Body:       bodyAny,
			Assertions: assertions,
			DependsOn:  []string{"step-setup"},
		},
	}

	tc := schema.TestCase{
		Schema:   schema.SchemaBaseURL,
		Version:  "1",
		ID:       fmt.Sprintf("TC-%s", uuid.New().String()[:8]),
		Title:    fmt.Sprintf("%s %s - idempotent: second call must be safe", op.Method, op.Path),
		Kind:     "chain",
		Priority: "P2",
		Tags:     op.Tags,
		Labels:   map[string]string{"type": "idempotency"},
		Source: schema.CaseSource{
			Technique: "idempotency",
			SpecPath:  fmt.Sprintf("%s %s", op.Method, op.Path),
			Rationale: fmt.Sprintf("%s is a write operation; test that repeat calls are safe", op.Method),
		},
		Steps:       steps,
		GeneratedAt: time.Now(),
	}
	return []schema.TestCase{tc}, nil
}
