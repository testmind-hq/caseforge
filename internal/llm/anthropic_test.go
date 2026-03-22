// internal/llm/anthropic_test.go
package llm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnthropicProviderName(t *testing.T) {
	p := &AnthropicProvider{model: "claude-sonnet-4-6"}
	assert.Equal(t, "anthropic:claude-sonnet-4-6", p.Name())
	assert.True(t, p.IsAvailable())
}

func TestNewProviderReturnsNoopWhenNoKey(t *testing.T) {
	p := NewProvider("", "anthropic", "claude-sonnet-4-6")
	assert.Equal(t, "noop", p.Name())
}

func TestNewProviderReturnsAnthropicWhenKeySet(t *testing.T) {
	p := NewProvider("sk-test-key", "anthropic", "claude-sonnet-4-6")
	assert.Equal(t, "anthropic:claude-sonnet-4-6", p.Name())
}
