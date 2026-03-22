// internal/lint/lint_test.go
package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestL004MissingSuccessResponse(t *testing.T) {
	ps := &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{
				OperationID: "listUsers",
				Method:      "GET",
				Path:        "/users",
				Responses: map[string]*spec.Response{
					"400": {Description: "Bad request"},
				},
			},
		},
	}
	issues := RunAll(ps)
	found := false
	for _, iss := range issues {
		if iss.RuleID == "L004" {
			found = true
		}
	}
	assert.True(t, found, "expected L004 for missing 2xx response")
}

func TestL001MissingOperationID(t *testing.T) {
	ps := &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{Method: "GET", Path: "/users"},
		},
	}
	issues := RunAll(ps)
	found := false
	for _, iss := range issues {
		if iss.RuleID == "L001" {
			found = true
		}
	}
	assert.True(t, found)
}

func TestL006UndeclaredPathParam(t *testing.T) {
	ps := &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{
				OperationID: "getUser",
				Method:      "GET",
				Path:        "/users/{id}",
				Parameters:  []*spec.Parameter{}, // {id} not declared
				Responses:   map[string]*spec.Response{"200": {Description: "OK"}},
			},
		},
	}
	issues := RunAll(ps)
	found := false
	for _, iss := range issues {
		if iss.RuleID == "L006" {
			found = true
		}
	}
	assert.True(t, found, "expected L006 for undeclared path param {id}")
}

func TestL006PassesWhenParamDeclared(t *testing.T) {
	ps := &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{
				OperationID: "getUser",
				Method:      "GET",
				Path:        "/users/{id}",
				Parameters: []*spec.Parameter{
					{Name: "id", In: "path", Required: true},
				},
				Responses: map[string]*spec.Response{"200": {Description: "OK"}},
			},
		},
	}
	issues := RunAll(ps)
	for _, iss := range issues {
		assert.NotEqual(t, "L006", iss.RuleID, "L006 should not fire when path param is declared")
	}
}

func TestNoIssuesForCleanSpec(t *testing.T) {
	ps := &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{
				OperationID: "listUsers",
				Method:      "GET",
				Path:        "/users",
				Summary:     "List all users",
				Responses:   map[string]*spec.Response{"200": {Description: "OK"}},
			},
		},
	}
	issues := RunAll(ps)
	for _, iss := range issues {
		if iss.Severity == "error" {
			t.Errorf("unexpected error: %s %s", iss.RuleID, iss.Message)
		}
	}
}
