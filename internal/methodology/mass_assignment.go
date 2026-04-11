// internal/methodology/mass_assignment.go
package methodology

import (
	"fmt"
	"maps"

	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/score"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// MassAssignmentTechnique generates test cases that inject extra fields not
// declared in the OpenAPI schema to detect mass assignment / object property
// injection vulnerabilities. Inspired by the CATS fuzzer.
//
// For each operation with a JSON request body, four cases are generated — one
// per probe category (privilege, status, financial, identity). Each case starts
// with a valid base body and adds all probe fields for that category.
type MassAssignmentTechnique struct {
	gen *datagen.Generator
}

func NewMassAssignmentTechnique() *MassAssignmentTechnique {
	return &MassAssignmentTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *MassAssignmentTechnique) Name() string { return "mass_assignment" }

func (t *MassAssignmentTechnique) Applies(op *spec.Operation) bool {
	return op.RequestBody != nil
}

type probeCategory struct {
	name     string
	scenario string
	fields   map[string]any
}

func probeCategories() []probeCategory {
	return []probeCategory{
		{
			name:     "privilege",
			scenario: string(score.ScenarioMassAssignmentPrivilege),
			fields: map[string]any{
				"role":     "__probe__",
				"admin":    true,
				"isAdmin":  true,
				"is_admin": true,
			},
		},
		{
			name:     "status",
			scenario: string(score.ScenarioMassAssignmentStatus),
			fields: map[string]any{
				"verified": true,
				"approved": true,
				"banned":   false,
				"disabled": false,
			},
		},
		{
			name:     "financial",
			scenario: string(score.ScenarioMassAssignmentFinancial),
			fields: map[string]any{
				"balance":  float64(1),
				"credits":  float64(1),
				"discount": float64(0),
				"price":    float64(1),
			},
		},
		{
			name:     "identity",
			scenario: string(score.ScenarioMassAssignmentIdentity),
			fields: map[string]any{
				"userId":    "__probe__",
				"user_id":   "__probe__",
				"ownerId":   "__probe__",
				"createdBy": "__probe__",
			},
		},
	}
}

func (t *MassAssignmentTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	validBase := buildValidBody(t.gen, op)
	if validBase == nil {
		validBase = map[string]any{}
	}

	var cases []schema.TestCase
	for _, cat := range probeCategories() {
		body := maps.Clone(validBase)
		for k, v := range cat.fields {
			body[k] = v
		}
		specPath := fmt.Sprintf("%s %s requestBody", op.Method, op.Path)
		tc := buildTestCase(op, body,
			fmt.Sprintf("[mass_assignment] %s probe", cat.name), specPath)
		tc.Priority = "P2"
		tc.Steps[0].Assertions = []schema.Assertion{
			{Target: "status_code", Operator: "eq", Expected: 400},
		}
		tc.Source = schema.CaseSource{
			Technique: "mass_assignment",
			SpecPath:  specPath,
			Rationale: fmt.Sprintf("inject %s probe fields not declared in schema to detect mass assignment vulnerability", cat.name),
			Scenario:  cat.scenario,
		}
		cases = append(cases, tc)
	}
	return cases, nil
}
