// internal/lint/completeness_test.go
package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestL013_ParameterMissingType(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{
			OperationID: "listUsers",
			Method:      "GET",
			Path:        "/users",
			Parameters:  []*spec.Parameter{{Name: "limit", In: "query", Schema: &spec.Schema{}}}, // no type
			Responses:   map[string]*spec.Response{"200": {}},
		},
	}}
	rule := &ruleL013{}
	issues := rule.Check(ps)
	assert.Len(t, issues, 1)
	assert.Equal(t, "L013", issues[0].RuleID)
	assert.Equal(t, "warning", issues[0].Severity)
	assert.Contains(t, issues[0].Message, "limit")
}

func TestL013_ParameterWithType_NoIssue(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{
			OperationID: "listUsers",
			Method:      "GET",
			Path:        "/users",
			Parameters:  []*spec.Parameter{{Name: "limit", In: "query", Schema: &spec.Schema{Type: "integer"}}},
			Responses:   map[string]*spec.Response{"200": {}},
		},
	}}
	rule := &ruleL013{}
	assert.Empty(t, rule.Check(ps))
}

func TestL014_Missing4xxResponse(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{
			OperationID: "createUser",
			Method:      "POST",
			Path:        "/users",
			Responses:   map[string]*spec.Response{"201": {}}, // no 4xx
		},
	}}
	rule := &ruleL014{}
	issues := rule.Check(ps)
	assert.Len(t, issues, 1)
	assert.Equal(t, "L014", issues[0].RuleID)
	assert.Equal(t, "warning", issues[0].Severity)
}

func TestL014_Has4xxResponse_NoIssue(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{
			OperationID: "createUser",
			Method:      "POST",
			Path:        "/users",
			Responses:   map[string]*spec.Response{"201": {}, "400": {}},
		},
	}}
	rule := &ruleL014{}
	assert.Empty(t, rule.Check(ps))
}

func TestL015_ResponseFieldMissingExample(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{
			OperationID: "getUser",
			Method:      "GET",
			Path:        "/users/{id}",
			Responses: map[string]*spec.Response{
				"200": {Content: map[string]*spec.MediaType{
					"application/json": {Schema: &spec.Schema{
						Properties: map[string]*spec.Schema{
							"id":   {Type: "integer"}, // no Example
							"name": {Type: "string"},  // no Example
						},
					}},
				}},
			},
		},
	}}
	rule := &ruleL015{}
	issues := rule.Check(ps)
	assert.Len(t, issues, 2)
	for _, iss := range issues {
		assert.Equal(t, "L015", iss.RuleID)
		assert.Equal(t, "warning", iss.Severity)
	}
}

func TestL015_ResponseFieldHasExample_NoIssue(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{
			OperationID: "getUser",
			Method:      "GET",
			Path:        "/users/{id}",
			Responses: map[string]*spec.Response{
				"200": {Content: map[string]*spec.MediaType{
					"application/json": {Schema: &spec.Schema{
						Properties: map[string]*spec.Schema{
							"id":   {Type: "integer", Example: 42},
							"name": {Type: "string", Example: "Alice"},
						},
					}},
				}},
			},
		},
	}}
	rule := &ruleL015{}
	assert.Empty(t, rule.Check(ps))
}
