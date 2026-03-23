// internal/llm/mcp_test.go
package llm

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMCPSamplingProviderComplete(t *testing.T) {
	called := false
	provider := NewMCPSamplingProvider(func(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
		called = true
		assert.Equal(t, "be concise", req.System)
		assert.Equal(t, "hello", req.Messages[0].Content)
		return &CompletionResponse{Text: "world", Tokens: 5}, nil
	})

	resp, err := provider.Complete(context.Background(), &CompletionRequest{
		System:    "be concise",
		Messages:  []Message{{Role: "user", Content: "hello"}},
		MaxTokens: 100,
	})
	require.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, "world", resp.Text)
	assert.Equal(t, 5, resp.Tokens)
}

func TestMCPSamplingProviderPropagatesError(t *testing.T) {
	provider := NewMCPSamplingProvider(func(_ context.Context, _ *CompletionRequest) (*CompletionResponse, error) {
		return nil, errors.New("sampling timeout")
	})
	_, err := provider.Complete(context.Background(), &CompletionRequest{})
	assert.ErrorContains(t, err, "sampling timeout")
}

func TestMCPSamplingProviderIsAvailableAndName(t *testing.T) {
	provider := NewMCPSamplingProvider(func(_ context.Context, _ *CompletionRequest) (*CompletionResponse, error) {
		return nil, nil
	})
	assert.True(t, provider.IsAvailable())
	assert.Equal(t, "mcp-sampling", provider.Name())
}
