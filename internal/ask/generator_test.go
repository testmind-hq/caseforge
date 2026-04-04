package ask_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/ask"
	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/output/schema"
)

// stubProvider is a controllable LLM provider for tests.
type stubProvider struct {
	text      string
	available bool
	err       error
	callFn    func() (*llm.CompletionResponse, error)
}

func (s *stubProvider) Complete(_ context.Context, _ *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	if s.callFn != nil {
		return s.callFn()
	}
	if s.err != nil {
		return nil, s.err
	}
	return &llm.CompletionResponse{Text: s.text}, nil
}
func (s *stubProvider) IsAvailable() bool { return s.available }
func (s *stubProvider) Name() string      { return "stub" }

const sampleJSON = `[
  {
    "title": "POST /users - valid email",
    "kind": "single",
    "priority": "P1",
    "tags": ["users"],
    "steps": [
      {
        "title": "create user",
        "method": "POST",
        "path": "/users",
        "body": {"email": "test@example.com"},
        "assertions": [{"target": "status_code", "operator": "eq", "expected": 201}]
      }
    ]
  }
]`

func TestGenerator_UnavailableProvider_ReturnsError(t *testing.T) {
	gen := ask.NewGenerator(&llm.NoopProvider{})
	_, err := gen.Generate(context.Background(), "POST /users")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "AI provider")
}

func TestGenerator_ParsesLLMJSON_ReturnsCases(t *testing.T) {
	gen := ask.NewGenerator(&stubProvider{text: sampleJSON, available: true})
	cases, err := gen.Generate(context.Background(), "POST /users")
	require.NoError(t, err)
	require.Len(t, cases, 1)

	tc := cases[0]
	assert.Equal(t, "POST /users - valid email", tc.Title)
	assert.Equal(t, "single", tc.Kind)
	assert.Equal(t, "P1", tc.Priority)
	assert.Equal(t, []string{"users"}, tc.Tags)
	require.Len(t, tc.Steps, 1)
	assert.Equal(t, "test", tc.Steps[0].Type)
}

func TestGenerator_ParsesMarkdownFencedJSON(t *testing.T) {
	fenced := "```json\n" + sampleJSON + "\n```"
	gen := ask.NewGenerator(&stubProvider{text: fenced, available: true})
	cases, err := gen.Generate(context.Background(), "POST /users")
	require.NoError(t, err)
	assert.Len(t, cases, 1)
}

func TestGenerator_ParsesPlainFencedJSON(t *testing.T) {
	fenced := "```\n" + sampleJSON + "\n```"
	gen := ask.NewGenerator(&stubProvider{text: fenced, available: true})
	cases, err := gen.Generate(context.Background(), "POST /users")
	require.NoError(t, err)
	assert.Len(t, cases, 1)
}

func TestGenerator_EmptyJSONArray_ReturnsEmptyCases(t *testing.T) {
	gen := ask.NewGenerator(&stubProvider{text: "[]", available: true})
	cases, err := gen.Generate(context.Background(), "GET /health")
	require.NoError(t, err)
	assert.Empty(t, cases)
}

func TestGenerator_InvalidJSON_ReturnsError(t *testing.T) {
	gen := ask.NewGenerator(&stubProvider{text: "{not json", available: true})
	_, err := gen.Generate(context.Background(), "GET /health")
	require.Error(t, err)
}

func TestGenerator_LLMError_ReturnsError(t *testing.T) {
	gen := ask.NewGenerator(&stubProvider{available: true, err: fmt.Errorf("network timeout")})
	_, err := gen.Generate(context.Background(), "POST /users")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "LLM completion failed")
}

func TestGenerator_RetriesOnTransientError(t *testing.T) {
	calls := 0
	provider := &stubProvider{
		available: true,
		callFn: func() (*llm.CompletionResponse, error) {
			calls++
			if calls < 2 {
				return nil, fmt.Errorf("transient")
			}
			return &llm.CompletionResponse{Text: sampleJSON}, nil
		},
	}
	gen := ask.NewGenerator(provider)
	cases, err := gen.Generate(context.Background(), "POST /users")
	require.NoError(t, err)
	assert.Len(t, cases, 1)
	assert.Equal(t, 2, calls)
}

func TestGenerator_NoNewlineAfterFence_Passthrough(t *testing.T) {
	// A fence with no newline after the opening (e.g. just "```") is passed through
	// unchanged. JSON parsing will then fail as expected.
	gen := ask.NewGenerator(&stubProvider{text: "```", available: true})
	_, err := gen.Generate(context.Background(), "GET /health")
	require.Error(t, err) // fails JSON parsing, not fence stripping
}

func TestGenerator_FillsProgrammaticFields(t *testing.T) {
	gen := ask.NewGenerator(&stubProvider{text: sampleJSON, available: true})
	cases, err := gen.Generate(context.Background(), "POST /users - create user")
	require.NoError(t, err)
	require.Len(t, cases, 1)

	tc := cases[0]
	assert.Equal(t, schema.SchemaBaseURL, tc.Schema)
	assert.Equal(t, "1", tc.Version)
	assert.True(t, len(tc.ID) > 0 && tc.ID[:3] == "TC-", "ID must start with TC-")
	assert.Equal(t, "ask", tc.Source.Technique)
	assert.Equal(t, "POST /users - create user", tc.Source.Rationale)
	assert.Empty(t, tc.Source.SpecPath)
	assert.False(t, tc.GeneratedAt.IsZero())
	assert.NotEmpty(t, tc.Steps[0].ID)
}

// --- LLM intent parsing hardening tests ---

const missingKindPriorityJSON = `[{
  "title": "bare case",
  "steps": [{"title": "s", "method": "GET", "path": "/health"}]
}]`

func TestGenerator_DefaultsEmptyKind_ToSingle(t *testing.T) {
	gen := ask.NewGenerator(&stubProvider{text: missingKindPriorityJSON, available: true})
	cases, err := gen.Generate(context.Background(), "GET /health")
	require.NoError(t, err)
	require.Len(t, cases, 1)
	assert.Equal(t, "single", cases[0].Kind, "empty kind must default to 'single'")
}

func TestGenerator_DefaultsEmptyPriority_ToP1(t *testing.T) {
	gen := ask.NewGenerator(&stubProvider{text: missingKindPriorityJSON, available: true})
	cases, err := gen.Generate(context.Background(), "GET /health")
	require.NoError(t, err)
	require.Len(t, cases, 1)
	assert.Equal(t, "P1", cases[0].Priority, "empty priority must default to 'P1'")
}

const threeCasesJSON = `[
  {"title":"A","kind":"single","priority":"P0","tags":[],"steps":[{"title":"s","method":"GET","path":"/a","assertions":[]}]},
  {"title":"B","kind":"single","priority":"P1","tags":[],"steps":[{"title":"s","method":"POST","path":"/b","assertions":[]}]},
  {"title":"C","kind":"single","priority":"P2","tags":[],"steps":[{"title":"s","method":"DELETE","path":"/c","assertions":[]}]}
]`

func TestGenerator_MultipleTestCases_AllParsed(t *testing.T) {
	gen := ask.NewGenerator(&stubProvider{text: threeCasesJSON, available: true})
	cases, err := gen.Generate(context.Background(), "multi case test")
	require.NoError(t, err)
	require.Len(t, cases, 3)
	assert.Equal(t, "A", cases[0].Title)
	assert.Equal(t, "B", cases[1].Title)
	assert.Equal(t, "C", cases[2].Title)
}

const multiStepJSON = `[{
  "title": "chain case",
  "kind": "chain",
  "priority": "P1",
  "tags": [],
  "steps": [
    {"title": "create", "method": "POST", "path": "/users", "assertions": []},
    {"title": "read",   "method": "GET",  "path": "/users/1", "assertions": []},
    {"title": "delete", "method": "DELETE","path": "/users/1", "assertions": []}
  ]
}]`

func TestGenerator_MultiStepCase_StepIDsSequential(t *testing.T) {
	gen := ask.NewGenerator(&stubProvider{text: multiStepJSON, available: true})
	cases, err := gen.Generate(context.Background(), "chain test")
	require.NoError(t, err)
	require.Len(t, cases, 1)
	steps := cases[0].Steps
	require.Len(t, steps, 3)
	assert.Equal(t, "step-1", steps[0].ID)
	assert.Equal(t, "step-2", steps[1].ID)
	assert.Equal(t, "step-3", steps[2].ID)
}

func TestGenerator_ChainKind_Preserved(t *testing.T) {
	gen := ask.NewGenerator(&stubProvider{text: multiStepJSON, available: true})
	cases, err := gen.Generate(context.Background(), "chain test")
	require.NoError(t, err)
	require.Len(t, cases, 1)
	assert.Equal(t, "chain", cases[0].Kind)
	assert.Len(t, cases[0].Steps, 3)
}

const headersJSON = `[{
  "title": "auth request",
  "kind": "single",
  "priority": "P1",
  "tags": [],
  "steps": [{
    "title": "authenticated call",
    "method": "GET",
    "path": "/profile",
    "headers": {"Authorization": "Bearer tok", "X-Tenant": "acme"},
    "assertions": []
  }]
}]`

func TestGenerator_HeadersPassThrough(t *testing.T) {
	gen := ask.NewGenerator(&stubProvider{text: headersJSON, available: true})
	cases, err := gen.Generate(context.Background(), "GET /profile with auth")
	require.NoError(t, err)
	require.Len(t, cases, 1)
	step := cases[0].Steps[0]
	assert.Equal(t, "Bearer tok", step.Headers["Authorization"])
	assert.Equal(t, "acme", step.Headers["X-Tenant"])
}

const multiOperatorJSON = `[{
  "title": "assertion operators",
  "kind": "single",
  "priority": "P1",
  "tags": [],
  "steps": [{
    "title": "check response",
    "method": "GET",
    "path": "/users",
    "assertions": [
      {"target": "status_code", "operator": "eq",       "expected": 200},
      {"target": "body",        "operator": "contains",  "expected": "id"},
      {"target": "header",      "operator": "ne",        "expected": ""},
      {"target": "jsonpath",    "operator": "matches",   "expected": "^[0-9]+$"},
      {"target": "jsonpath",    "operator": "exists"},
      {"target": "jsonpath",    "operator": "gt",        "expected": 0}
    ]
  }]
}]`

func TestGenerator_AssertionOperatorsPassThrough(t *testing.T) {
	gen := ask.NewGenerator(&stubProvider{text: multiOperatorJSON, available: true})
	cases, err := gen.Generate(context.Background(), "GET /users assertions")
	require.NoError(t, err)
	require.Len(t, cases, 1)
	assertions := cases[0].Steps[0].Assertions
	require.Len(t, assertions, 6)
	ops := make([]string, len(assertions))
	for i, a := range assertions {
		ops[i] = a.Operator
	}
	assert.Equal(t, []string{"eq", "contains", "ne", "matches", "exists", "gt"}, ops)
}

func TestGenerator_AllRetriesExhausted_ReturnsError(t *testing.T) {
	calls := 0
	provider := &stubProvider{
		available: true,
		callFn: func() (*llm.CompletionResponse, error) {
			calls++
			return nil, fmt.Errorf("service unavailable")
		},
	}
	gen := ask.NewGenerator(provider)
	_, err := gen.Generate(context.Background(), "POST /orders")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "LLM completion failed")
	assert.Equal(t, 3, calls, "should attempt exactly 3 times before giving up")
}

const prioritiesJSON = `[
  {"title":"P0 case","kind":"single","priority":"P0","tags":[],"steps":[{"title":"s","method":"GET","path":"/","assertions":[]}]},
  {"title":"P2 case","kind":"single","priority":"P2","tags":[],"steps":[{"title":"s","method":"GET","path":"/","assertions":[]}]},
  {"title":"P3 case","kind":"single","priority":"P3","tags":[],"steps":[{"title":"s","method":"GET","path":"/","assertions":[]}]}
]`

func TestGenerator_AllPrioritiesPreserved(t *testing.T) {
	gen := ask.NewGenerator(&stubProvider{text: prioritiesJSON, available: true})
	cases, err := gen.Generate(context.Background(), "priority round-trip")
	require.NoError(t, err)
	require.Len(t, cases, 3)
	assert.Equal(t, "P0", cases[0].Priority)
	assert.Equal(t, "P2", cases[1].Priority)
	assert.Equal(t, "P3", cases[2].Priority)
}

func TestGenerator_ProseWrappedJSON(t *testing.T) {
	prose := "Here are some test cases for the API:\n\n" + sampleJSON + "\n\nI hope this helps!"
	gen := ask.NewGenerator(&stubProvider{text: prose, available: true})
	cases, err := gen.Generate(context.Background(), "POST /users")
	require.NoError(t, err)
	assert.Len(t, cases, 1, "should extract JSON array from surrounding prose")
}

func TestGenerator_JSONObjectInsteadOfArray_ReturnsError(t *testing.T) {
	gen := ask.NewGenerator(&stubProvider{text: `{"title":"oops","steps":[]}`, available: true})
	_, err := gen.Generate(context.Background(), "GET /")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing LLM response")
}

const bodyPassJSON = `[{
  "title": "create user",
  "kind": "single",
  "priority": "P1",
  "tags": [],
  "steps": [{
    "title": "POST create",
    "method": "POST",
    "path": "/users",
    "body": {"email": "a@b.com", "role": "admin"},
    "assertions": []
  }]
}]`

func TestGenerator_BodyPassThrough(t *testing.T) {
	gen := ask.NewGenerator(&stubProvider{text: bodyPassJSON, available: true})
	cases, err := gen.Generate(context.Background(), "POST /users create")
	require.NoError(t, err)
	require.Len(t, cases, 1)
	body, ok := cases[0].Steps[0].Body.(map[string]any)
	require.True(t, ok, "step body should be a map")
	assert.Equal(t, "a@b.com", body["email"])
	assert.Equal(t, "admin", body["role"])
}

const allMethodsJSON = `[
  {"title":"GET","kind":"single","priority":"P1","tags":[],"steps":[{"title":"s","method":"GET","path":"/r","assertions":[]}]},
  {"title":"PUT","kind":"single","priority":"P1","tags":[],"steps":[{"title":"s","method":"PUT","path":"/r/1","assertions":[]}]},
  {"title":"PATCH","kind":"single","priority":"P1","tags":[],"steps":[{"title":"s","method":"PATCH","path":"/r/1","assertions":[]}]},
  {"title":"DELETE","kind":"single","priority":"P1","tags":[],"steps":[{"title":"s","method":"DELETE","path":"/r/1","assertions":[]}]}
]`

func TestGenerator_HTTPMethodsPreserved(t *testing.T) {
	gen := ask.NewGenerator(&stubProvider{text: allMethodsJSON, available: true})
	cases, err := gen.Generate(context.Background(), "all HTTP methods")
	require.NoError(t, err)
	require.Len(t, cases, 4)
	assert.Equal(t, "GET", cases[0].Steps[0].Method)
	assert.Equal(t, "PUT", cases[1].Steps[0].Method)
	assert.Equal(t, "PATCH", cases[2].Steps[0].Method)
	assert.Equal(t, "DELETE", cases[3].Steps[0].Method)
}
