// internal/rbt/parser_llm_test.go
package rbt

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/llm"
)

// llmStub is a controllable LLMProvider for rbt tests.
type llmStub struct {
	text string
	err  error
}

func (s *llmStub) Complete(_ context.Context, _ *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &llm.CompletionResponse{Text: s.text}, nil
}
func (s *llmStub) IsAvailable() bool { return true }
func (s *llmStub) Name() string      { return "stub" }

func TestLLMParser_InferRoutes_BareJSON(t *testing.T) {
	p := NewLLMParser(&llmStub{
		text: `[{"method":"GET","path":"/users","confidence":0.9}]`,
	}, "")
	routes, err := p.inferRoutes(context.Background(), "handler.go", "func listUsers() {}")
	require.NoError(t, err)
	require.Len(t, routes, 1)
	assert.Equal(t, "GET", routes[0].Method)
	assert.Equal(t, "/users", routes[0].Path)
	assert.InDelta(t, 0.9, routes[0].Confidence, 0.001)
}

func TestLLMParser_InferRoutes_FencedJSON(t *testing.T) {
	p := NewLLMParser(&llmStub{
		text: "```json\n[{\"method\":\"POST\",\"path\":\"/orders\",\"confidence\":0.8}]\n```",
	}, "")
	routes, err := p.inferRoutes(context.Background(), "order.go", "func createOrder() {}")
	require.NoError(t, err)
	require.Len(t, routes, 1)
	assert.Equal(t, "POST", routes[0].Method)
}

func TestLLMParser_InferRoutes_TextPreamble(t *testing.T) {
	p := NewLLMParser(&llmStub{
		text: "Here are the routes:\n[{\"method\":\"DELETE\",\"path\":\"/items/{id}\",\"confidence\":0.7}]",
	}, "")
	routes, err := p.inferRoutes(context.Background(), "item.go", "func deleteItem() {}")
	require.NoError(t, err)
	require.Len(t, routes, 1)
	assert.Equal(t, "DELETE", routes[0].Method)
}

func TestLLMParser_InferRoutes_EmptyArray(t *testing.T) {
	p := NewLLMParser(&llmStub{text: "[]"}, "")
	routes, err := p.inferRoutes(context.Background(), "util.go", "func helper() {}")
	require.NoError(t, err)
	assert.Empty(t, routes)
}

func TestLLMParser_InferRoutes_LLMError_ReturnsError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // immediately cancel so retry doesn't wait
	p := NewLLMParser(&llmStub{err: errors.New("network error")}, "")
	_, err := p.inferRoutes(ctx, "handler.go", "")
	require.Error(t, err)
}

func TestLLMParser_InferRoutes_BadJSON_ReturnsError(t *testing.T) {
	p := NewLLMParser(&llmStub{text: "not json"}, "")
	_, err := p.inferRoutes(context.Background(), "handler.go", "")
	require.Error(t, err)
}
