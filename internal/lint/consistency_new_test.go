// internal/lint/consistency_new_test.go
package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestL016_DuplicateOperationID(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{OperationID: "listUsers", Method: "GET", Path: "/users", Responses: map[string]*spec.Response{"200": {}}},
		{OperationID: "listUsers", Method: "GET", Path: "/admin/users", Responses: map[string]*spec.Response{"200": {}}},
	}}
	rule := &ruleL016{}
	issues := rule.Check(ps)
	assert.Len(t, issues, 1)
	assert.Equal(t, "L016", issues[0].RuleID)
	assert.Equal(t, "error", issues[0].Severity)
	assert.Contains(t, issues[0].Message, "listUsers")
}

func TestL016_UniqueOperationIDs_NoIssue(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{OperationID: "listUsers", Method: "GET", Path: "/users", Responses: map[string]*spec.Response{"200": {}}},
		{OperationID: "getUser", Method: "GET", Path: "/users/{id}", Responses: map[string]*spec.Response{"200": {}}},
	}}
	rule := &ruleL016{}
	assert.Empty(t, rule.Check(ps))
}

func TestL017_MixedPathVersions(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{OperationID: "listUsersV1", Method: "GET", Path: "/v1/users", Responses: map[string]*spec.Response{"200": {}}},
		{OperationID: "listUsersV2", Method: "GET", Path: "/v2/users", Responses: map[string]*spec.Response{"200": {}}},
	}}
	rule := &ruleL017{}
	issues := rule.Check(ps)
	assert.Len(t, issues, 1)
	assert.Equal(t, "L017", issues[0].RuleID)
	assert.Equal(t, "warning", issues[0].Severity)
}

func TestL017_SingleVersion_NoIssue(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{OperationID: "listUsers", Method: "GET", Path: "/v1/users", Responses: map[string]*spec.Response{"200": {}}},
		{OperationID: "getUser", Method: "GET", Path: "/v1/users/{id}", Responses: map[string]*spec.Response{"200": {}}},
	}}
	rule := &ruleL017{}
	assert.Empty(t, rule.Check(ps))
}

func TestL018_InconsistentContentType(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{Method: "GET", Path: "/users", Responses: map[string]*spec.Response{
			"200": {Content: map[string]*spec.MediaType{"application/json": {}}},
		}},
		{Method: "GET", Path: "/export", Responses: map[string]*spec.Response{
			"200": {Content: map[string]*spec.MediaType{"text/csv": {}}},
		}},
	}}
	rule := &ruleL018{}
	issues := rule.Check(ps)
	assert.Len(t, issues, 1)
	assert.Equal(t, "L018", issues[0].RuleID)
	assert.Equal(t, "warning", issues[0].Severity)
}

func TestL018_ConsistentContentType_NoIssue(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{Method: "GET", Path: "/users", Responses: map[string]*spec.Response{
			"200": {Content: map[string]*spec.MediaType{"application/json": {}}},
		}},
		{Method: "GET", Path: "/users/{id}", Responses: map[string]*spec.Response{
			"200": {Content: map[string]*spec.MediaType{"application/json": {}}},
		}},
	}}
	rule := &ruleL018{}
	assert.Empty(t, rule.Check(ps))
}
