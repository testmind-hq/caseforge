// internal/llm/noop_test.go
package llm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoopProviderAlwaysReturnsEmpty(t *testing.T) {
	p := &NoopProvider{}
	resp, err := p.Complete(context.Background(), &CompletionRequest{
		System:    "you are a test expert",
		Messages:  []Message{{Role: "user", Content: "analyze this"}},
		MaxTokens: 100,
	})
	require.NoError(t, err)
	assert.Equal(t, "", resp.Text)
	assert.Equal(t, 0, resp.Tokens)
}

func TestNoopProviderIsNotAvailable(t *testing.T) {
	p := &NoopProvider{}
	assert.False(t, p.IsAvailable())
	assert.Equal(t, "noop", p.Name())
}
