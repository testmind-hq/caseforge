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
