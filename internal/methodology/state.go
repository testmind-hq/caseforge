// internal/methodology/state.go
package methodology

import (
	"fmt"

	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

type StateTechnique struct {
	gen *datagen.Generator
}

func NewStateTechnique() *StateTechnique {
	return &StateTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *StateTechnique) Name() string { return "state_transition" }

// Applies only when SemanticInfo indicates a state machine.
func (t *StateTechnique) Applies(op *spec.Operation) bool {
	return op.SemanticInfo != nil && op.SemanticInfo.HasStateMachine
}

func (t *StateTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	if op.SemanticInfo == nil || !op.SemanticInfo.HasStateMachine {
		return nil, nil
	}
	// Find the state field schema and its enum values
	stateField := op.SemanticInfo.StateField
	s := getJSONSchema(op.RequestBody)
	if s == nil {
		return nil, nil
	}
	stateSchema, ok := s.Properties[stateField]
	if !ok || len(stateSchema.Enum) == 0 {
		return nil, nil
	}

	var cases []schema.TestCase
	for _, stateVal := range stateSchema.Enum {
		body := map[string]any{}
		for name, fieldSchema := range s.Properties {
			if name == stateField {
				body[name] = stateVal
			} else {
				body[name] = t.gen.Generate(fieldSchema, name)
			}
		}
		tc := buildTestCase(op, body,
			fmt.Sprintf("transition %s to %v", stateField, stateVal),
			fmt.Sprintf("%s %s requestBody.properties.%s", op.Method, op.Path, stateField))
		tc.Priority = "P1"
		tc.Source = schema.CaseSource{
			Technique: "state_transition",
			SpecPath:  fmt.Sprintf("%s %s", op.Method, op.Path),
			Rationale: fmt.Sprintf("state transition: set %s to %v", stateField, stateVal),
		}
		cases = append(cases, tc)
	}
	return cases, nil
}
