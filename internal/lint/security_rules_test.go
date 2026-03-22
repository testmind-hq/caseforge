package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestL011_MissingSecurityScheme(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{Method: "POST", Path: "/users", Security: nil, Responses: map[string]*spec.Response{"201": {}}},
		{Method: "GET", Path: "/health", Security: nil, Responses: map[string]*spec.Response{"200": {}}},      // GET excluded
		{Method: "POST", Path: "/public/login", Security: nil, Responses: map[string]*spec.Response{"200": {}}}, // /public excluded
		{Method: "DELETE", Path: "/users/{id}", Security: []string{"bearerAuth"}, Responses: map[string]*spec.Response{"204": {}}}, // has security
	}}
	rule := &ruleL011{}
	issues := rule.Check(ps)
	assert.Len(t, issues, 1)
	assert.Equal(t, "L011", issues[0].RuleID)
	assert.Equal(t, "error", issues[0].Severity)
	assert.Contains(t, issues[0].Path, "POST /users")
}

func TestL011_ExcludedPaths(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{Method: "POST", Path: "/health/check", Responses: map[string]*spec.Response{"200": {}}},
		{Method: "POST", Path: "/auth/login", Responses: map[string]*spec.Response{"200": {}}},
		{Method: "POST", Path: "/users/register", Responses: map[string]*spec.Response{"201": {}}},
		{Method: "POST", Path: "/public/verify", Responses: map[string]*spec.Response{"200": {}}},
	}}
	rule := &ruleL011{}
	assert.Empty(t, rule.Check(ps))
}

func TestL012_SensitiveFieldExposed(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{Method: "GET", Path: "/users/{id}",
			Responses: map[string]*spec.Response{
				"200": {Content: map[string]*spec.MediaType{
					"application/json": {Schema: &spec.Schema{Properties: map[string]*spec.Schema{
						"id":           {Type: "integer"},
						"passwordHash": {Type: "string"},
						"name":         {Type: "string"},
					}}},
				}},
			},
		},
	}}
	rule := &ruleL012{}
	issues := rule.Check(ps)
	assert.Len(t, issues, 1)
	assert.Equal(t, "L012", issues[0].RuleID)
	assert.Equal(t, "error", issues[0].Severity)
	assert.Contains(t, issues[0].Message, "passwordHash")
}

func TestL012_NoSensitiveFields(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{Method: "GET", Path: "/users",
			Responses: map[string]*spec.Response{
				"200": {Content: map[string]*spec.MediaType{
					"application/json": {Schema: &spec.Schema{Properties: map[string]*spec.Schema{
						"id": {Type: "integer"}, "name": {Type: "string"},
					}}},
				}},
			},
		},
	}}
	rule := &ruleL012{}
	assert.Empty(t, rule.Check(ps))
}
