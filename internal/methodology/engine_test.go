// internal/methodology/engine_test.go
package methodology

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/event"
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

func TestEngineEmitsEventsToSink(t *testing.T) {
	type emitted struct{ typ event.EventType }
	var got []emitted
	mu := sync.Mutex{}

	sink := event.SinkFunc(func(e event.Event) {
		mu.Lock()
		got = append(got, emitted{e.Type})
		mu.Unlock()
	})

	noop := &llm.NoopProvider{}
	engine := NewEngine(noop, NewEquivalenceTechnique())
	engine.SetSink(sink)

	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{
			OperationID: "createUser",
			Method:      "POST", Path: "/users",
			RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
				"application/json": {Schema: &spec.Schema{
					Type: "object",
					Properties: map[string]*spec.Schema{
						"email": {Type: "string", Format: "email"},
					},
				}},
			}},
			Responses: map[string]*spec.Response{"201": {}},
		},
	}}
	_, err := engine.Generate(ps)
	require.NoError(t, err)

	types := make([]event.EventType, len(got))
	for i, g := range got {
		types[i] = g.typ
	}
	assert.Contains(t, types, event.EventOperationDone)
	assert.Contains(t, types, event.EventCaseGenerated)
	// NoopProvider skips annotation — EventOperationAnnotating must NOT be emitted.
	assert.NotContains(t, types, event.EventOperationAnnotating)
}

func TestEngineEmitsAnnotatingEventsWhenLLMAvailable(t *testing.T) {
	var got []event.EventType
	mu := sync.Mutex{}
	sink := event.SinkFunc(func(e event.Event) {
		mu.Lock()
		got = append(got, e.Type)
		mu.Unlock()
	})

	// stubLLM returns an available provider that always succeeds with empty JSON.
	stub := &stubLLMProvider{}
	engine := NewEngine(stub, NewEquivalenceTechnique())
	engine.SetSink(sink)

	ps := &spec.ParsedSpec{Operations: []*spec.Operation{
		{OperationID: "op1", Method: "GET", Path: "/a",
			Responses: map[string]*spec.Response{"200": {}}},
		{OperationID: "op2", Method: "POST", Path: "/b",
			Responses: map[string]*spec.Response{"201": {}}},
	}}
	_, err := engine.Generate(ps)
	require.NoError(t, err)

	var annotatingCount int
	for _, typ := range got {
		if typ == event.EventOperationAnnotating {
			annotatingCount++
		}
	}
	assert.Equal(t, 2, annotatingCount, "one EventOperationAnnotating per operation")
}

// stubLLMProvider is an available provider that returns empty JSON for any call.
type stubLLMProvider struct{}

func (s *stubLLMProvider) IsAvailable() bool                      { return true }
func (s *stubLLMProvider) Name() string                            { return "stub" }
func (s *stubLLMProvider) Complete(_ context.Context, _ *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	return &llm.CompletionResponse{Text: "{}"}, nil
}

func TestEngineConcurrentProducesSameResults(t *testing.T) {
	// Build a spec with several operations so the worker pool is exercised.
	noop := &llm.NoopProvider{}
	ps := &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{OperationID: "op1", Method: "POST", Path: "/a",
				RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
					"application/json": {Schema: &spec.Schema{Type: "object",
						Properties: map[string]*spec.Schema{"x": {Type: "string"}}}},
				}},
				Responses: map[string]*spec.Response{"201": {}}},
			{OperationID: "op2", Method: "GET", Path: "/b",
				Responses: map[string]*spec.Response{"200": {}}},
			{OperationID: "op3", Method: "PUT", Path: "/c",
				RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
					"application/json": {Schema: &spec.Schema{Type: "object",
						Properties: map[string]*spec.Schema{"y": {Type: "integer"}}}},
				}},
				Responses: map[string]*spec.Response{"200": {}}},
		},
	}

	newEngine := func(concurrency int) *Engine {
		e := NewEngine(noop, NewEquivalenceTechnique(), NewBoundaryTechnique())
		e.SetConcurrency(concurrency)
		return e
	}

	serialCases, err := newEngine(1).Generate(ps)
	require.NoError(t, err)

	parallelCases, err := newEngine(3).Generate(ps)
	require.NoError(t, err)

	// IDs are random UUIDs so we compare by (title, technique, path) fingerprint.
	type fingerprint struct{ title, technique, path string }
	toSet := func(cases []schema.TestCase) map[fingerprint]int {
		m := make(map[fingerprint]int)
		for _, c := range cases {
			fp := fingerprint{c.Title, c.Source.Technique, c.Source.SpecPath}
			m[fp]++
		}
		return m
	}

	serialSet := toSet(serialCases)
	parallelSet := toSet(parallelCases)
	assert.Equal(t, serialSet, parallelSet, "concurrent and serial must produce equivalent cases")
}

func TestEngine_SeedProducesIdenticalOutput(t *testing.T) {
	op := &spec.Operation{
		OperationID: "listItems",
		Method:      "GET",
		Path:        "/items",
		Parameters: []*spec.Parameter{
			{Name: "sort", In: "query", Schema: &spec.Schema{Type: "boolean"}},
			{Name: "filter", In: "query", Schema: &spec.Schema{Enum: []any{"all", "active", "archived"}}},
			{Name: "format", In: "query", Schema: &spec.Schema{Enum: []any{"json", "csv"}}},
			{Name: "lang", In: "query", Schema: &spec.Schema{Enum: []any{"en", "zh"}}},
		},
		Responses: map[string]*spec.Response{"200": {}},
	}
	ps := &spec.ParsedSpec{Operations: []*spec.Operation{op}}

	run := func(seed int64) []schema.TestCase {
		e := NewEngine(&llm.NoopProvider{}, NewPairwiseTechnique())
		e.SetSeed(seed)
		cases, err := e.Generate(ps)
		require.NoError(t, err)
		return cases
	}

	cases1 := run(42)
	cases2 := run(42)
	cases3 := run(99) // different seed

	require.Equal(t, len(cases1), len(cases2), "same seed must produce same number of cases")
	for i := range cases1 {
		assert.Equal(t, cases1[i].Steps[0].Path, cases2[i].Steps[0].Path,
			"same seed must produce identical paths at index %d", i)
	}

	// Different seed may (not guaranteed) produce different order
	_ = cases3 // just ensure it doesn't panic
}

// mockSpecTechnique is a test double for SpecTechnique.
type mockSpecTechnique struct {
	onGenerate func(*spec.ParsedSpec) ([]schema.TestCase, error)
}

func (m *mockSpecTechnique) Name() string { return "mock_spec" }
func (m *mockSpecTechnique) Generate(s *spec.ParsedSpec) ([]schema.TestCase, error) {
	return m.onGenerate(s)
}

func TestEngine_MaxCasesPerOp_TruncatesByPriority(t *testing.T) {
	engine := NewEngine(&llm.NoopProvider{},
		NewEquivalenceTechnique(),
		NewBoundaryTechnique(),
		NewIsolatedNegativeTechnique(),
		NewSchemaViolationTechnique(),
	)
	engine.SetMaxCasesPerOp(2)

	ps := &spec.ParsedSpec{
		Operations: []*spec.Operation{{
			OperationID: "createUser",
			Method:      "POST",
			Path:        "/users",
			RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
				"application/json": {Schema: &spec.Schema{
					Type:     "object",
					Required: []string{"email"},
					Properties: map[string]*spec.Schema{
						"email": {Type: "string", Format: "email"},
						"age": {Type: "integer",
							Minimum: func() *float64 { v := float64(18); return &v }(),
							Maximum: func() *float64 { v := float64(120); return &v }(),
						},
					},
				}},
			}},
			Responses: map[string]*spec.Response{"201": {Description: "Created"}},
		}},
	}

	cases, err := engine.Generate(ps)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(cases), 2,
		"engine must not produce more than maxCasesPerOp cases for a single operation")
}
