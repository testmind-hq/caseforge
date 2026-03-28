package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigCommand_IsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "config" {
			found = true
			break
		}
	}
	assert.True(t, found, "config command must be registered on rootCmd")
}

func TestConfigShowCommand_IsRegistered(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "config" {
			for _, sub := range c.Commands() {
				if sub.Use == "show" {
					return // found
				}
			}
			t.Fatal("config show subcommand not found")
		}
	}
	t.Fatal("config command not found")
}

func TestConfigShow_PrintsYAML(t *testing.T) {
	t.Cleanup(func() { viper.Reset() })
	var buf bytes.Buffer
	configShowCmd.SetOut(&buf)
	require.NoError(t, runConfigShow(configShowCmd, nil))
	out := buf.String()
	assert.Contains(t, out, "provider:")
	assert.Contains(t, out, "default_format:")
}

func TestConfigShow_MasksAPIKey_LongKey(t *testing.T) {
	t.Cleanup(func() { viper.Reset() })
	viper.Set("ai.api_key", "sk-abcdefghijklmnop")
	var buf bytes.Buffer
	configShowCmd.SetOut(&buf)
	require.NoError(t, runConfigShow(configShowCmd, nil))
	out := buf.String()
	assert.Contains(t, out, "sk-abc...")
	assert.NotContains(t, out, "sk-abcdefghijklmnop")
}

func TestConfigShow_MasksAPIKey_ShortKey(t *testing.T) {
	t.Cleanup(func() { viper.Reset() })
	viper.Set("ai.api_key", "short")
	var buf bytes.Buffer
	configShowCmd.SetOut(&buf)
	require.NoError(t, runConfigShow(configShowCmd, nil))
	out := buf.String()
	assert.Contains(t, out, "short")
}

func TestConfigShow_EmptyAPIKey(t *testing.T) {
	t.Cleanup(func() { viper.Reset() })
	var buf bytes.Buffer
	configShowCmd.SetOut(&buf)
	require.NoError(t, runConfigShow(configShowCmd, nil))
	// Just check it doesn't error; empty key is fine
	assert.NotEmpty(t, buf.String())
}
