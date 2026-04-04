// internal/lint/security_new_test.go
package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestL019_GET_MissingSecurityScheme(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{Method: "GET", Path: "/users", Security: nil, Responses: map[string]*spec.Response{"200": {}}},
		{Method: "GET", Path: "/public/health", Security: nil, Responses: map[string]*spec.Response{"200": {}}}, // excluded
		{Method: "GET", Path: "/users/{id}", Security: []string{"bearerAuth"}, Responses: map[string]*spec.Response{"200": {}}},
	}}
	rule := &ruleL019{}
	issues := rule.Check(ps)
	assert.Len(t, issues, 1)
	assert.Equal(t, "L019", issues[0].RuleID)
	assert.Equal(t, "warning", issues[0].Severity)
	assert.Contains(t, issues[0].Path, "GET /users")
}

func TestL019_GET_WithSecurity_NoIssue(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{Method: "GET", Path: "/users", Security: []string{"bearerAuth"}, Responses: map[string]*spec.Response{"200": {}}},
	}}
	rule := &ruleL019{}
	assert.Empty(t, rule.Check(ps))
}

func TestL020_SensitiveQueryParam(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{
			Method: "GET", Path: "/users",
			Parameters: []*spec.Parameter{
				{Name: "token", In: "query", Schema: &spec.Schema{Type: "string"}},
				{Name: "page", In: "query", Schema: &spec.Schema{Type: "integer"}},
			},
			Responses: map[string]*spec.Response{"200": {}},
		},
	}}
	rule := &ruleL020{}
	issues := rule.Check(ps)
	assert.Len(t, issues, 1)
	assert.Equal(t, "L020", issues[0].RuleID)
	assert.Equal(t, "error", issues[0].Severity)
	assert.Contains(t, issues[0].Message, "token")
}

func TestL020_NoSensitiveQueryParam_NoIssue(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{
			Method: "GET", Path: "/users",
			Parameters: []*spec.Parameter{
				{Name: "page", In: "query", Schema: &spec.Schema{Type: "integer"}},
				{Name: "limit", In: "query", Schema: &spec.Schema{Type: "integer"}},
			},
			Responses: map[string]*spec.Response{"200": {}},
		},
	}}
	rule := &ruleL020{}
	assert.Empty(t, rule.Check(ps))
}

func TestL021_NoSecuritySchemes(t *testing.T) {
	ps := &spec.ParsedSpec{
		SecuritySchemes: nil,
		Operations: []*spec.Operation{
			{OperationID: "listUsers", Method: "GET", Path: "/users", Responses: map[string]*spec.Response{"200": {}}},
		},
	}
	rule := &ruleL021{}
	issues := rule.Check(ps)
	assert.Len(t, issues, 1)
	assert.Equal(t, "L021", issues[0].RuleID)
	assert.Equal(t, "warning", issues[0].Severity)
}

func TestL021_HasSecuritySchemes_NoIssue(t *testing.T) {
	ps := &spec.ParsedSpec{
		SecuritySchemes: []string{"bearerAuth"},
		Operations: []*spec.Operation{
			{OperationID: "listUsers", Method: "GET", Path: "/users", Responses: map[string]*spec.Response{"200": {}}},
		},
	}
	rule := &ruleL021{}
	assert.Empty(t, rule.Check(ps))
}
