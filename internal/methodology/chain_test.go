package methodology

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/output/schema"
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
			assert.Equal(t, []string{"step-test"}, s.DependsOn, "teardown must depend on test step")
		}
	}
	assert.True(t, hasTeardown, "should have a teardown step when DELETE exists")
}

func TestChainCaseTeardownWithMismatchedParamName(t *testing.T) {
	// Verify that teardown substitution works even when DELETE uses a different
	// param name than GET (e.g., GET /users/{userId} but DELETE /users/{id}).
	ps := &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{
				OperationID: "createUser", Method: "POST", Path: "/users",
				Responses: map[string]*spec.Response{"201": {Content: map[string]*spec.MediaType{
					"application/json": {Schema: &spec.Schema{Type: "object",
						Properties: map[string]*spec.Schema{"id": {Type: "integer"}}}},
				}}},
			},
			{
				OperationID: "getUser", Method: "GET", Path: "/users/{userId}",
				Parameters: []*spec.Parameter{{Name: "userId", In: "path", Required: true, Schema: &spec.Schema{Type: "integer"}}},
				Responses:  map[string]*spec.Response{"200": {}},
			},
			{
				OperationID: "deleteUser", Method: "DELETE", Path: "/users/{id}",
				Parameters: []*spec.Parameter{{Name: "id", In: "path", Required: true, Schema: &spec.Schema{Type: "integer"}}},
				Responses:  map[string]*spec.Response{"204": {}},
			},
		},
	}
	ct := NewChainTechnique()
	cases, err := ct.Generate(ps)
	require.NoError(t, err)
	require.NotEmpty(t, cases)

	for _, s := range cases[0].Steps {
		if s.Type == "teardown" {
			assert.Contains(t, s.Path, "{{userId}}", "teardown path must use captured variable, even when DELETE param name differs")
			assert.NotContains(t, s.Path, "{id}", "raw OpenAPI param must not appear in rendered path")
		}
	}
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

func TestChainTechnique_Source_ScenarioCRUDFlow(t *testing.T) {
	ct := NewChainTechnique()
	cases, err := ct.Generate(crudSpec())
	require.NoError(t, err)
	require.NotEmpty(t, cases)

	assert.Equal(t, "CRUD_FLOW", cases[0].Source.Scenario,
		"chain_crud source scenario must be CRUD_FLOW")
}

func crudSpecWithUpdate() *spec.ParsedSpec {
	return &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{
				OperationID: "createUser", Method: "POST", Path: "/users",
				Tags: []string{"users"},
				RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
					"application/json": {Schema: &spec.Schema{
						Type: "object",
						Properties: map[string]*spec.Schema{
							"name":  {Type: "string"},
							"email": {Type: "string"},
						},
					}},
				}},
				Responses: map[string]*spec.Response{"201": {
					Headers: map[string]string{},
					Content: map[string]*spec.MediaType{"application/json": {Schema: &spec.Schema{
						Type: "object",
						Properties: map[string]*spec.Schema{"id": {Type: "integer"}},
					}}},
				}},
			},
			{
				OperationID: "updateUser", Method: "PUT", Path: "/users/{userId}",
				Parameters: []*spec.Parameter{{Name: "userId", In: "path", Required: true}},
				RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
					"application/json": {Schema: &spec.Schema{
						Type: "object",
						Properties: map[string]*spec.Schema{"name": {Type: "string"}},
					}},
				}},
				Responses: map[string]*spec.Response{"200": {}},
			},
			{
				OperationID: "getUser", Method: "GET", Path: "/users/{userId}",
				Parameters: []*spec.Parameter{{Name: "userId", In: "path", Required: true}},
				Responses:  map[string]*spec.Response{"200": {}},
			},
			{
				OperationID: "deleteUser", Method: "DELETE", Path: "/users/{userId}",
				Parameters: []*spec.Parameter{{Name: "userId", In: "path", Required: true}},
				Responses:  map[string]*spec.Response{"204": {}},
			},
		},
	}
}

func TestChainTechnique_NStepChain_WithUpdate(t *testing.T) {
	ct := NewChainTechnique()
	cases, err := ct.Generate(crudSpecWithUpdate())
	require.NoError(t, err)
	require.NotEmpty(t, cases)

	tc := cases[0]
	assert.Equal(t, "chain", tc.Kind)
	// Expect: setup(POST) + update(PUT) + test(GET) + teardown(DELETE) = 4 steps
	assert.Equal(t, 4, len(tc.Steps))

	types := make([]string, len(tc.Steps))
	for i, s := range tc.Steps {
		types[i] = s.Type
	}
	assert.Equal(t, []string{"setup", "update", "test", "teardown"}, types)
}

func TestChainTechnique_UpdateStepUsesCapture(t *testing.T) {
	ct := NewChainTechnique()
	cases, err := ct.Generate(crudSpecWithUpdate())
	require.NoError(t, err)
	require.NotEmpty(t, cases)

	var updateStep *schema.Step
	for i := range cases[0].Steps {
		if cases[0].Steps[i].Type == "update" {
			updateStep = &cases[0].Steps[i]
		}
	}
	require.NotNil(t, updateStep, "should have an update step")
	assert.Contains(t, updateStep.Path, "{{userId}}", "update path must use captured variable")
	assert.Contains(t, updateStep.DependsOn, "step-setup", "update step must depend on setup")
}

func TestChainTechnique_LocationHeaderCapture(t *testing.T) {
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{
			OperationID: "createItem", Method: "POST", Path: "/items",
			Responses: map[string]*spec.Response{"201": {
				Headers: map[string]string{"Location": "string"},
			}},
		},
		{
			OperationID: "getItem", Method: "GET", Path: "/items/{itemId}",
			Responses: map[string]*spec.Response{"200": {}},
		},
	}}
	ct := NewChainTechnique()
	cases, err := ct.Generate(ps)
	require.NoError(t, err)
	require.NotEmpty(t, cases)

	setup := cases[0].Steps[0]
	require.NotEmpty(t, setup.Captures)
	assert.Equal(t, "header Location", setup.Captures[0].From,
		"should prefer Location header capture when documented in spec")
}
