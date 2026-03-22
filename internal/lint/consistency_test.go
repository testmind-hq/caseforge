package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestL007_VerbInPath(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{Method: "POST", Path: "/createUser", Responses: map[string]*spec.Response{"200": {}}},
		{Method: "GET", Path: "/users", Responses: map[string]*spec.Response{"200": {}}},
	}}
	rule := &ruleL007{}
	issues := rule.Check(ps)
	assert.Len(t, issues, 1)
	assert.Equal(t, "L007", issues[0].RuleID)
	assert.Equal(t, "warning", issues[0].Severity)
	assert.Contains(t, issues[0].Message, "create")
}

func TestL008_InconsistentNaming(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{Method: "POST", Path: "/users",
			RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
				"application/json": {Schema: &spec.Schema{Properties: map[string]*spec.Schema{
					"userId": {Type: "string"},
				}}},
			}},
			Parameters: []*spec.Parameter{{Name: "user_id", In: "query", Schema: &spec.Schema{Type: "string"}}},
			Responses:  map[string]*spec.Response{"200": {}},
		},
	}}
	rule := &ruleL008{}
	issues := rule.Check(ps)
	assert.Len(t, issues, 1)
	assert.Equal(t, "L008", issues[0].RuleID)
}

func TestL008_ConsistentNaming_NoIssue(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{Method: "GET", Path: "/users",
			Parameters: []*spec.Parameter{
				{Name: "userId", In: "query", Schema: &spec.Schema{Type: "string"}},
				{Name: "orderId", In: "query", Schema: &spec.Schema{Type: "string"}},
			},
			Responses: map[string]*spec.Response{"200": {}},
		},
	}}
	rule := &ruleL008{}
	assert.Empty(t, rule.Check(ps))
}

func TestL009_InconsistentPagination(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{Method: "GET", Path: "/users",
			Parameters: []*spec.Parameter{{Name: "page", In: "query", Schema: &spec.Schema{Type: "integer"}}},
			Responses:  map[string]*spec.Response{"200": {}},
		},
		{Method: "GET", Path: "/orders",
			Parameters: []*spec.Parameter{{Name: "offset", In: "query", Schema: &spec.Schema{Type: "integer"}}},
			Responses:  map[string]*spec.Response{"200": {}},
		},
	}}
	rule := &ruleL009{}
	issues := rule.Check(ps)
	assert.Len(t, issues, 1)
	assert.Equal(t, "L009", issues[0].RuleID)
}

func TestL010_InconsistentErrorSchema(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{Method: "POST", Path: "/users", Responses: map[string]*spec.Response{
			"400": {Content: map[string]*spec.MediaType{
				"application/json": {Schema: &spec.Schema{Properties: map[string]*spec.Schema{"error": {Type: "string"}}}},
			}},
		}},
		{Method: "POST", Path: "/orders", Responses: map[string]*spec.Response{
			"400": {Content: map[string]*spec.MediaType{
				"application/json": {Schema: &spec.Schema{Properties: map[string]*spec.Schema{"message": {Type: "string"}}}},
			}},
		}},
	}}
	rule := &ruleL010{}
	issues := rule.Check(ps)
	assert.Len(t, issues, 1)
	assert.Equal(t, "L010", issues[0].RuleID)
}
