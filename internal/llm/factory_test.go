// internal/llm/factory_test.go
package llm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProviderWithConfigStruct_Noop(t *testing.T) {
	p := NewProviderWithConfig(ProviderConfig{Provider: "anthropic", APIKey: ""})
	assert.Equal(t, "noop", p.Name())
}

func TestNewProviderWithConfigStruct_Anthropic(t *testing.T) {
	p := NewProviderWithConfig(ProviderConfig{
		Provider: "anthropic",
		APIKey:   "sk-test-key",
		Model:    "claude-sonnet-4-6",
	})
	assert.Equal(t, "anthropic:claude-sonnet-4-6", p.Name())
}
