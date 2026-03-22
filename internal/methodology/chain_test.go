package methodology

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func crudSpec() *spec.ParsedSpec {
	return &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{
				OperationID: "createUser",
				Method:      "POST",
				Path:        "/users",
				Tags:        []string{"users"},
				RequestBody: &spec.RequestBody{
					Content: map[string]*spec.MediaType{
						"application/json": {Schema: &spec.Schema{
							Type: "object",
							Properties: map[string]*spec.Schema{
								"name":  {Type: "string"},
								"email": {Type: "string", Format: "email"},
							},
						}},
					},
				},
				Responses: map[string]*spec.Response{
					"201": {Content: map[string]*spec.MediaType{
						"application/json": {Schema: &spec.Schema{
							Type: "object",
							Properties: map[string]*spec.Schema{
								"id":    {Type: "integer"},
								"name":  {Type: "string"},
								"email": {Type: "string"},
							},
						}},
					}},
				},
			},
			{
				OperationID: "getUser",
				Method:      "GET",
				Path:        "/users/{userId}",
				Tags:        []string{"users"},
				Parameters: []*spec.Parameter{
					{Name: "userId", In: "path", Required: true, Schema: &spec.Schema{Type: "integer"}},
				},
				Responses: map[string]*spec.Response{"200": {Description: "OK"}},
			},
			{
				OperationID: "deleteUser",
				Method:      "DELETE",
				Path:        "/users/{userId}",
				Tags:        []string{"users"},
				Parameters: []*spec.Parameter{
					{Name: "userId", In: "path", Required: true, Schema: &spec.Schema{Type: "integer"}},
				},
				Responses: map[string]*spec.Response{"204": {Description: "No Content"}},
			},
		},
	}
}

func TestChainTechniqueApplies(t *testing.T) {
	ct := NewChainTechnique()
	assert.Equal(t, "chain_crud", ct.Name())

	cases, err := ct.Generate(crudSpec())
	require.NoError(t, err)
	assert.NotEmpty(t, cases, "should generate at least one chain case")
}

func TestChainCaseHasMultipleSteps(t *testing.T) {
	ct := NewChainTechnique()
	cases, err := ct.Generate(crudSpec())
	require.NoError(t, err)
	require.NotEmpty(t, cases)

	tc := cases[0]
	assert.Equal(t, "chain", tc.Kind)
	assert.GreaterOrEqual(t, len(tc.Steps), 2, "chain case must have at least 2 steps")
}

func TestChainCaseSetupStepHasCapture(t *testing.T) {
	ct := NewChainTechnique()
	cases, err := ct.Generate(crudSpec())
	require.NoError(t, err)
	require.NotEmpty(t, cases)

	setup := cases[0].Steps[0]
	assert.Equal(t, "setup", setup.Type)
	assert.Equal(t, "POST", setup.Method)
	assert.NotEmpty(t, setup.Captures, "setup step must capture the created resource ID")
	assert.Equal(t, "userId", setup.Captures[0].Name)
	assert.Contains(t, setup.Captures[0].From, "$.id")
}

func TestChainCaseReadStepUsesCapture(t *testing.T) {
	ct := NewChainTechnique()
	cases, err := ct.Generate(crudSpec())
	require.NoError(t, err)
	require.NotEmpty(t, cases)

	for _, s := range cases[0].Steps {
		if s.Method == "GET" {
			assert.Contains(t, s.Path, "{{userId}}", "GET path must use captured variable")
			assert.Equal(t, "test", s.Type)
		}
	}
}

func TestChainCaseTeardownStepPresent(t *testing.T) {
	ct := NewChainTechnique()
	cases, err := ct.Generate(crudSpec())
	require.NoError(t, err)
	require.NotEmpty(t, cases)

	var hasTeardown bool
	for _, s := range cases[0].Steps {
		if s.Type == "teardown" {
			hasTeardown = true
			assert.Equal(t, "DELETE", s.Method)
			assert.Contains(t, s.Path, "{{userId}}")
			assert.Equal(t, []string{"step-setup", "step-test"}, s.DependsOn, "teardown must depend on both setup and test")
		}
	}
	assert.True(t, hasTeardown, "should have a teardown step when DELETE exists")
}

func TestChainTechniqueNoGroupWhenNoCRUD(t *testing.T) {
	ps := &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{Method: "GET", Path: "/health", Responses: map[string]*spec.Response{"200": {}}},
		},
	}
	ct := NewChainTechnique()
	cases, err := ct.Generate(ps)
	require.NoError(t, err)
	assert.Empty(t, cases, "no chain case when no POST+GET pair exists")
}

func TestChainTechniqueSourceAnnotation(t *testing.T) {
	ct := NewChainTechnique()
	cases, err := ct.Generate(crudSpec())
	require.NoError(t, err)
	require.NotEmpty(t, cases)

	assert.Equal(t, "chain_crud", cases[0].Source.Technique)
	assert.True(t, strings.HasPrefix(cases[0].Source.SpecPath, "/users"))
}
