package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testmind-hq/caseforge/internal/output/schema"
)

func TestSuggest_SpecPathMatch(t *testing.T) {
	result := DiffResult{Changes: []Change{
		{Kind: Breaking, Method: "DELETE", Path: "/users/{id}", Description: "endpoint removed"},
	}}
	cases := []schema.TestCase{
		{ID: "TC-0001", Title: "delete user", Source: schema.CaseSource{SpecPath: "DELETE /users/{id}"}},
		{ID: "TC-0002", Title: "create user", Source: schema.CaseSource{SpecPath: "POST /users"}},
	}
	affected := Suggest(result, cases)
	assert.Len(t, affected, 1)
	assert.Equal(t, "TC-0001", affected[0].ID)
	assert.NotEmpty(t, affected[0].Reason)
}

func TestSuggest_StepPathMatch(t *testing.T) {
	result := DiffResult{Changes: []Change{
		{Kind: Breaking, Method: "GET", Path: "/users/{id}", Description: "endpoint removed"},
	}}
	cases := []schema.TestCase{
		{ID: "TC-0003", Title: "chain case", Source: schema.CaseSource{SpecPath: "POST /users"},
			Steps: []schema.Step{
				{Path: "/users/{{userId}}"},
			},
		},
	}
	affected := Suggest(result, cases)
	assert.Len(t, affected, 1)
	assert.Equal(t, "TC-0003", affected[0].ID)
}

func TestSuggest_NonBreakingNotIncluded(t *testing.T) {
	result := DiffResult{Changes: []Change{
		{Kind: NonBreaking, Method: "GET", Path: "/users", Description: "new optional param"},
	}}
	cases := []schema.TestCase{
		{ID: "TC-0001", Source: schema.CaseSource{SpecPath: "GET /users"}},
	}
	affected := Suggest(result, cases)
	assert.Empty(t, affected)
}

func TestSuggest_Empty(t *testing.T) {
	result := DiffResult{}
	affected := Suggest(result, nil)
	assert.Empty(t, affected)
}
