package ask_test

import (
	"context"
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
}

func (s *stubProvider) Complete(_ context.Context, _ *llm.CompletionRequest) (*llm.CompletionResponse, error) {
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
	assert.Equal(t, "POST /users - create user", tc.Source.SpecPath)
	assert.False(t, tc.GeneratedAt.IsZero())
	assert.NotEmpty(t, tc.Steps[0].ID)
}
