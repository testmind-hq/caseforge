// internal/llm/openai_test.go
package llm

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestOpenAIProviderName_Standard(t *testing.T) {
	p := &OpenAIProvider{model: "gpt-4o-mini", baseURL: ""}
	assert.Equal(t, "openai:gpt-4o-mini", p.Name())
	assert.True(t, p.IsAvailable())
}

func TestOpenAIProviderName_Compat(t *testing.T) {
	p := &OpenAIProvider{model: "deepseek-chat", baseURL: "https://api.deepseek.com/v1"}
	assert.Equal(t, "openai-compat:deepseek-chat", p.Name())
	assert.True(t, p.IsAvailable())
}

func TestNewProvider_OpenAIWithKey(t *testing.T) {
	p := NewProvider("sk-test", "openai", "gpt-4o-mini")
	assert.Equal(t, "openai:gpt-4o-mini", p.Name())
}

func TestNewProvider_OpenAICompatWithKey(t *testing.T) {
	p := NewProviderWithConfig(ProviderConfig{
		APIKey:   "sk-test",
		Provider: "openai-compat",
		Model:    "deepseek-chat",
		BaseURL:  "https://api.deepseek.com/v1",
	})
	assert.Equal(t, "openai-compat:deepseek-chat", p.Name())
}

func TestNewProvider_OpenAINoKey(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	p := NewProvider("", "openai", "gpt-4o-mini")
	assert.Equal(t, "noop", p.Name())
}
