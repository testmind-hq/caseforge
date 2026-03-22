// internal/assert/basic_test.go
package assert

import (
	"testing"

	"github.com/stretchr/testify/assert"
	specpkg "github.com/testmind-hq/caseforge/internal/spec"
)

func TestStatusCodeAssertion(t *testing.T) {
	op := &specpkg.Operation{
		Responses: map[string]*specpkg.Response{
			"200": {Description: "OK"},
		},
	}
	assertions := BasicAssertions(op)
	var found bool
	for _, a := range assertions {
		if a.Target == "status_code" && a.Operator == "eq" && a.Expected == 200 {
			found = true
		}
	}
	assert.True(t, found, "expected status_code eq 200 assertion")
}

func TestDurationAssertion(t *testing.T) {
	op := &specpkg.Operation{}
	assertions := BasicAssertions(op)
	var found bool
	for _, a := range assertions {
		if a.Target == "duration_ms" && a.Operator == "lt" {
			found = true
		}
	}
	assert.True(t, found, "expected duration_ms lt assertion")
}

func TestSchemaAssertionsForObjectResponse(t *testing.T) {
	props := map[string]*specpkg.Schema{
		"id":   {Type: "string", Format: "uuid"},
		"name": {Type: "string"},
	}
	s := &specpkg.Schema{Type: "object", Properties: props}
	assertions := SchemaAssertions("body", s)
	targets := make(map[string]bool)
	for _, a := range assertions {
		targets[a.Target] = true
	}
	assert.True(t, targets["body.id"])
	assert.True(t, targets["body.name"])
}
