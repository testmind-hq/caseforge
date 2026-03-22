package diff

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func v1Spec() *spec.ParsedSpec {
	return &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{Method: "GET", Path: "/users", Parameters: []*spec.Parameter{{Name: "limit", In: "query", Required: false, Schema: &spec.Schema{Type: "integer"}}}, Responses: map[string]*spec.Response{"200": {Content: map[string]*spec.MediaType{"application/json": {Schema: &spec.Schema{Properties: map[string]*spec.Schema{"id": {Type: "integer"}, "email": {Type: "string"}}}}}}}},
			{Method: "DELETE", Path: "/users/{id}", Responses: map[string]*spec.Response{"204": {}}},
			{Method: "POST", Path: "/orders", RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{"application/json": {Schema: &spec.Schema{Properties: map[string]*spec.Schema{"customerId": {Type: "integer"}}}}}}, Responses: map[string]*spec.Response{"201": {}}},
		},
	}
}

func v2Spec() *spec.ParsedSpec {
	return &spec.ParsedSpec{
		Operations: []*spec.Operation{
			// GET /users: response field "email" removed, new optional param added
			{Method: "GET", Path: "/users", Parameters: []*spec.Parameter{{Name: "limit", In: "query", Required: false, Schema: &spec.Schema{Type: "integer"}}, {Name: "includeDeleted", In: "query", Required: false, Schema: &spec.Schema{Type: "boolean"}}}, Responses: map[string]*spec.Response{"200": {Content: map[string]*spec.MediaType{"application/json": {Schema: &spec.Schema{Properties: map[string]*spec.Schema{"id": {Type: "integer"}}}}}}}},
			// DELETE /users/{id} removed
			// POST /orders: customerId type changed
			{Method: "POST", Path: "/orders", RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{"application/json": {Schema: &spec.Schema{Properties: map[string]*spec.Schema{"customerId": {Type: "string"}}}}}}, Responses: map[string]*spec.Response{"201": {}}},
			// New endpoint
			{Method: "POST", Path: "/payments", Responses: map[string]*spec.Response{"201": {}}},
		},
	}
}

func TestDiffEndpointRemoved(t *testing.T) {
	result := Diff(v1Spec(), v2Spec())
	var found bool
	for _, c := range result.Changes {
		if c.Kind == Breaking && c.Method == "DELETE" && c.Path == "/users/{id}" {
			found = true
			assert.Contains(t, c.Description, "endpoint removed")
		}
	}
	assert.True(t, found, "should detect endpoint removal as BREAKING")
}

func TestDiffResponseFieldRemoved(t *testing.T) {
	result := Diff(v1Spec(), v2Spec())
	var found bool
	for _, c := range result.Changes {
		if c.Kind == Breaking && c.Path == "/users" && strings.Contains(c.Description, "email") {
			found = true
		}
	}
	assert.True(t, found, "response field deletion should be BREAKING")
}

func TestDiffParamTypeChanged(t *testing.T) {
	result := Diff(v1Spec(), v2Spec())
	var found bool
	for _, c := range result.Changes {
		if c.Kind == Breaking && c.Path == "/orders" && strings.Contains(c.Description, "customerId") {
			found = true
		}
	}
	assert.True(t, found, "param type change should be BREAKING")
}

func TestDiffNewEndpoint(t *testing.T) {
	result := Diff(v1Spec(), v2Spec())
	var found bool
	for _, c := range result.Changes {
		if c.Kind == NonBreaking && c.Path == "/payments" {
			found = true
		}
	}
	assert.True(t, found, "new endpoint should be NON_BREAKING")
}

func TestDiffNewOptionalParam(t *testing.T) {
	result := Diff(v1Spec(), v2Spec())
	var found bool
	for _, c := range result.Changes {
		if c.Kind == NonBreaking && c.Path == "/users" && strings.Contains(c.Description, "includeDeleted") {
			found = true
		}
	}
	assert.True(t, found, "new optional param should be NON_BREAKING")
}

// ensure require import is used
var _ = require.New
