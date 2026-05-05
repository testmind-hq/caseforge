// internal/config/config_test.go
package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadDefaults(t *testing.T) {
	viper.Reset()
	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "noop", cfg.AI.Provider)
	assert.Equal(t, "hurl", cfg.Output.DefaultFormat)
	assert.Equal(t, "error", cfg.Lint.FailOn)
}

func TestLoadFromViper(t *testing.T) {
	viper.Reset()
	viper.Set("ai.provider", "anthropic")
	viper.Set("ai.model", "claude-sonnet-4-6")
	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "anthropic", cfg.AI.Provider)
	assert.Equal(t, "claude-sonnet-4-6", cfg.AI.Model)
}

func TestLoadBaseURLFromViper(t *testing.T) {
	viper.Reset()
	viper.Set("ai.base_url", "https://api.deepseek.com/v1")
	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "https://api.deepseek.com/v1", cfg.AI.BaseURL)
}

func TestValidate_APIKeyLooksLikeURL(t *testing.T) {
	c := &AIConfig{Provider: "openai", APIKey: "https://open.bigmodel.cn/api/paas/v4"}
	err := c.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "api_key looks like a URL")
	assert.Contains(t, err.Error(), "base_url")
}

func TestValidate_BaseURLMissingScheme(t *testing.T) {
	c := &AIConfig{Provider: "openai-compat", BaseURL: "open.bigmodel.cn/api/paas/v4"}
	err := c.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no HTTP scheme")
}

func TestValidate_OpenAICompatMissingBaseURL(t *testing.T) {
	c := &AIConfig{Provider: "openai-compat", BaseURL: ""}
	err := c.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires ai.base_url")
}

func TestValidate_UnknownProvider(t *testing.T) {
	c := &AIConfig{Provider: "unknown-llm"}
	err := c.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown ai.provider")
}

func TestValidate_ValidOpenAICompat(t *testing.T) {
	c := &AIConfig{Provider: "openai-compat", BaseURL: "https://open.bigmodel.cn/api/paas/v4"}
	assert.NoError(t, c.Validate())
}

func TestValidate_ValidAnthropicNoBaseURL(t *testing.T) {
	c := &AIConfig{Provider: "anthropic", APIKey: "sk-ant-test"}
	assert.NoError(t, c.Validate())
}

func TestConfigBedrockRegionFromEnv(t *testing.T) {
	viper.Reset()
	viper.Set("ai.provider", "bedrock")
	t.Setenv("AWS_REGION", "us-west-2")

	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "us-west-2", cfg.AI.Region)

	viper.Reset()
}
