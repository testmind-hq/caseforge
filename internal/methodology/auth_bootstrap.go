package methodology

import (
	"fmt"

	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// BootstrapAuth wraps every single-step case whose source operation is secured
// with an auth setup step, converting Kind "single" → "chain".
// Cases that are already chains are left unchanged.
// If no auth operation can be detected in s, all cases are returned unchanged.
func BootstrapAuth(cases []schema.TestCase, s *spec.ParsedSpec) []schema.TestCase {
	authOp := findAuthOperation(s.Operations)
	if authOp == nil {
		return cases
	}
	tokenField := findTokenField(authOp)
	if tokenField == "" {
		return cases
	}

	gen := datagen.NewGenerator(nil)
	authBody := buildValidBody(gen, authOp)
	var authBodyAny any
	if authBody != nil {
		authBodyAny = authBody
	}
	authHeaders := map[string]string{}
	if authBodyAny != nil {
		authHeaders["Content-Type"] = "application/json"
	}
	authStep := schema.Step{
		ID:    "step-auth",
		Title: fmt.Sprintf("authenticate via %s %s", authOp.Method, authOp.Path),
		Type:  "setup",
		Method:  authOp.Method,
		Path:    authOp.Path,
		Headers: authHeaders,
		Body:    authBodyAny,
		Assertions: []schema.Assertion{
			{Target: "status_code", Operator: schema.OperatorLt, Expected: 300},
		},
		Captures: []schema.Capture{
			{Name: "authToken", From: fmt.Sprintf("jsonpath $.%s", tokenField)},
		},
	}

	// Build a set of secured paths from the spec for fast lookup.
	securedPaths := make(map[string]bool)
	for _, op := range s.Operations {
		if isSecuredOperation(op, s) {
			securedPaths[op.Method+" "+op.Path] = true
		}
	}

	result := make([]schema.TestCase, len(cases))
	for i, tc := range cases {
		// Skip if already a chain
		if tc.Kind == "chain" {
			result[i] = tc
			continue
		}
		// Determine if the case targets a secured operation
		opKey := ""
		if len(tc.Steps) > 0 {
			opKey = tc.Steps[0].Method + " " + tc.Steps[0].Path
		} else if tc.Source.SpecPath != "" {
			opKey = tc.Source.SpecPath
		}
		if !securedPaths[opKey] {
			result[i] = tc
			continue
		}
		// Wrap: prepend auth step, add DependsOn to original first step
		newSteps := make([]schema.Step, len(tc.Steps)+1)
		newSteps[0] = authStep
		for j, step := range tc.Steps {
			step.DependsOn = []string{"step-auth"}
			if step.Headers == nil {
				step.Headers = map[string]string{}
			}
			step.Headers["Authorization"] = "Bearer {{authToken}}"
			newSteps[j+1] = step
		}
		tc.Kind = "chain"
		tc.Steps = newSteps
		result[i] = tc
	}
	return result
}
