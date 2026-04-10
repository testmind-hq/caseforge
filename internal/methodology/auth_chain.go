// internal/methodology/auth_chain.go
package methodology

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	assertpkg "github.com/testmind-hq/caseforge/internal/assert"
	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// AuthChainTechnique detects operations that produce auth tokens and generates
// a 2-step chain for each secured operation:
//  1. Auth step: call the auth operation, capture the token
//  2. Test step: call the secured operation with Authorization: Bearer {{authToken}}
type AuthChainTechnique struct {
	gen *datagen.Generator
}

func NewAuthChainTechnique() *AuthChainTechnique {
	return &AuthChainTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *AuthChainTechnique) Name() string { return "auth_chain" }

func (t *AuthChainTechnique) Generate(s *spec.ParsedSpec) ([]schema.TestCase, error) {
	authOp := findAuthOperation(s.Operations)
	if authOp == nil {
		return nil, nil
	}
	tokenField := findTokenField(authOp)
	if tokenField == "" {
		return nil, nil
	}

	var secured []*spec.Operation
	for _, op := range s.Operations {
		if op == authOp {
			continue
		}
		if isSecuredOperation(op, s) {
			secured = append(secured, op)
		}
	}
	if len(secured) == 0 {
		return nil, nil
	}

	sort.Slice(secured, func(i, j int) bool {
		return secured[i].Method+secured[i].Path < secured[j].Method+secured[j].Path
	})

	authBody := buildValidBody(t.gen, authOp)
	var authBodyAny any
	if authBody != nil {
		authBodyAny = authBody
	}
	authHeaders := map[string]string{}
	if authBodyAny != nil {
		authHeaders["Content-Type"] = "application/json"
	}

	var cases []schema.TestCase
	for _, op := range secured {
		authStep := schema.Step{
			ID:      "step-auth",
			Title:   fmt.Sprintf("authenticate via %s %s", authOp.Method, authOp.Path),
			Type:    "setup",
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

		testHeaders := map[string]string{
			"Authorization": "Bearer {{authToken}}",
		}
		if op.RequestBody != nil {
			testHeaders["Content-Type"] = "application/json"
		}
		testStep := schema.Step{
			ID:         "step-test",
			Title:      fmt.Sprintf("%s %s with auth token", op.Method, op.Path),
			Type:       "test",
			Method:     op.Method,
			Path:       op.Path,
			Headers:    testHeaders,
			Assertions: assertpkg.BasicAssertions(op),
			DependsOn:  []string{"step-auth"},
		}

		tc := schema.TestCase{
			Schema:   schema.SchemaBaseURL,
			Version:  "1",
			ID:       fmt.Sprintf("TC-%s", uuid.New().String()[:8]),
			Title:    fmt.Sprintf("auth chain: %s %s", op.Method, op.Path),
			Kind:     "chain",
			Priority: "P1",
			Tags:     op.Tags,
			Source: schema.CaseSource{
				Technique: "auth_chain",
				SpecPath:  fmt.Sprintf("%s %s", op.Method, op.Path),
				Rationale: fmt.Sprintf("authenticate via %s then call secured endpoint %s %s",
					authOp.Path, op.Method, op.Path),
			},
			Steps:       []schema.Step{authStep, testStep},
			GeneratedAt: time.Now(),
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

// findAuthOperation locates the POST operation that produces an auth token.
func findAuthOperation(ops []*spec.Operation) *spec.Operation {
	sorted := make([]*spec.Operation, len(ops))
	copy(sorted, ops)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Method+sorted[i].Path < sorted[j].Method+sorted[j].Path
	})
	// First pass: response body has a token field
	for _, op := range sorted {
		if op.Method != "POST" {
			continue
		}
		if findTokenField(op) != "" {
			return op
		}
	}
	// Second pass: path heuristic
	for _, op := range sorted {
		if op.Method != "POST" {
			continue
		}
		lower := strings.ToLower(op.Path)
		for _, kw := range []string{"login", "signin", "auth", "token", "session"} {
			if strings.Contains(lower, kw) {
				return op
			}
		}
	}
	return nil
}

// findTokenField returns the name of a token field in the 2xx response body.
func findTokenField(op *spec.Operation) string {
	tokenNames := []string{"access_token", "token", "jwt", "id_token", "bearer", "auth_token"}
	codes := make([]string, 0, len(op.Responses))
	for code := range op.Responses {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	for _, code := range codes {
		resp := op.Responses[code]
		n := 0
		fmt.Sscanf(code, "%d", &n)
		if n < 200 || n >= 300 {
			continue
		}
		mt, ok := resp.Content["application/json"]
		if !ok || mt == nil || mt.Schema == nil {
			continue
		}
		for _, name := range tokenNames {
			if _, ok := mt.Schema.Properties[name]; ok {
				return name
			}
		}
	}
	return ""
}

// isSecuredOperation returns true when the operation requires authentication.
func isSecuredOperation(op *spec.Operation, s *spec.ParsedSpec) bool {
	if len(op.Security) > 0 {
		return true
	}
	return len(s.GlobalSecurity) > 0
}
