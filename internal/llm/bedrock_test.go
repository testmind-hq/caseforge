// internal/llm/bedrock_test.go
package llm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBedrockProviderName(t *testing.T) {
	p := &BedrockProvider{model: "us.anthropic.claude-sonnet-4-6"}
	assert.Equal(t, "bedrock:us.anthropic.claude-sonnet-4-6", p.Name())
	assert.True(t, p.IsAvailable())
}

func TestNewProviderWithConfig_BedrockNoRegion(t *testing.T) {
	p := NewProviderWithConfig(ProviderConfig{Provider: "bedrock", Region: ""})
	assert.Equal(t, "noop", p.Name())
}

func TestNewProviderWithConfig_BedrockWithRegion(t *testing.T) {
	p := NewProviderWithConfig(ProviderConfig{
		Provider: "bedrock",
		Region:   "us-east-1",
		Model:    "us.anthropic.claude-sonnet-4-6",
	})
	assert.Equal(t, "bedrock:us.anthropic.claude-sonnet-4-6", p.Name())
}
