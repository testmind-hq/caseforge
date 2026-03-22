// internal/methodology/engine_test.go
package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestEngineGeneratesFromParsedSpec(t *testing.T) {
	noop := &llm.NoopProvider{}
	engine := NewEngine(noop,
		NewEquivalenceTechnique(),
		NewBoundaryTechnique(),
		NewDecisionTechnique(),
		NewStateTechnique(),
		NewIdempotentTechnique(),
		NewPairwiseTechnique(),
	)

	ps := &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{
				OperationID: "createUser",
				Method:      "POST",
				Path:        "/users",
				Summary:     "Create a user",
				RequestBody: &spec.RequestBody{
					Content: map[string]*spec.MediaType{
						"application/json": {
							Schema: &spec.Schema{
								Type:     "object",
								Required: []string{"email"},
								Properties: map[string]*spec.Schema{
									"email": {Type: "string", Format: "email"},
									"age":   {Type: "integer", Minimum: floatPtr(0), Maximum: floatPtr(120)},
								},
							},
						},
					},
				},
				Responses: map[string]*spec.Response{"201": {Description: "Created"}},
			},
		},
	}

	cases, err := engine.Generate(ps)
	require.NoError(t, err)
	assert.NotEmpty(t, cases)

	// All cases must have traceability
	for _, tc := range cases {
		assert.NotEmpty(t, tc.ID)
		assert.NotEmpty(t, tc.Source.Technique)
		assert.NotEmpty(t, tc.Source.SpecPath)
		assert.NotEmpty(t, tc.Source.Rationale)
		assert.Equal(t, "1", tc.Version)
	}
}

func TestEngineCallsSpecTechnique(t *testing.T) {
	called := false
	var gotSpec *spec.ParsedSpec
	noop := &llm.NoopProvider{}
	engine := NewEngine(noop)
	engine.AddSpecTechnique(&mockSpecTechnique{onGenerate: func(s *spec.ParsedSpec) ([]schema.TestCase, error) {
		called = true
		gotSpec = s
		return []schema.TestCase{{ID: "chain-1"}}, nil
	}})

	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{Method: "GET", Path: "/x", Responses: map[string]*spec.Response{"200": {}}},
	}}
	cases, err := engine.Generate(ps)
	require.NoError(t, err)
	assert.True(t, called, "SpecTechnique should have been called")
	assert.NotNil(t, gotSpec, "SpecTechnique should receive the ParsedSpec")
	assert.Contains(t, cases, schema.TestCase{ID: "chain-1"}, "cases returned by SpecTechnique must be in output")
}

func TestEngineGeneratesChainCasesFromCRUDSpec(t *testing.T) {
	noop := &llm.NoopProvider{}
	engine := NewEngine(noop)
	engine.AddSpecTechnique(NewChainTechnique())

	ps := &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{
				OperationID: "createItem",
				Method:      "POST", Path: "/items",
				RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
					"application/json": {Schema: &spec.Schema{Type: "object",
						Properties: map[string]*spec.Schema{"name": {Type: "string"}}}},
				}},
				Responses: map[string]*spec.Response{"201": {Content: map[string]*spec.MediaType{
					"application/json": {Schema: &spec.Schema{Type: "object",
						Properties: map[string]*spec.Schema{"id": {Type: "integer"}}}},
				}}},
			},
			{
				OperationID: "getItem",
				Method:      "GET", Path: "/items/{itemId}",
				Parameters: []*spec.Parameter{{Name: "itemId", In: "path", Required: true,
					Schema: &spec.Schema{Type: "integer"}}},
				Responses: map[string]*spec.Response{"200": {Description: "OK"}},
			},
		},
	}

	cases, err := engine.Generate(ps)
	require.NoError(t, err)

	var chainCases []schema.TestCase
	for _, c := range cases {
		if c.Kind == "chain" {
			chainCases = append(chainCases, c)
		}
	}
	require.Len(t, chainCases, 1)
	assert.Equal(t, "chain_crud", chainCases[0].Source.Technique)
	assert.Len(t, chainCases[0].Steps, 2) // setup + test (no DELETE in spec)
}

// mockSpecTechnique is a test double for SpecTechnique.
type mockSpecTechnique struct {
	onGenerate func(*spec.ParsedSpec) ([]schema.TestCase, error)
}

func (m *mockSpecTechnique) Name() string { return "mock_spec" }
func (m *mockSpecTechnique) Generate(s *spec.ParsedSpec) ([]schema.TestCase, error) {
	return m.onGenerate(s)
}
