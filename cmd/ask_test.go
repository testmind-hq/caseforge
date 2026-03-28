package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAskCommand_IsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "ask <description>" {
			found = true
			break
		}
	}
	assert.True(t, found, "ask command must be registered on rootCmd")
}

func TestAskCommand_HasFlags(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "ask <description>" {
			assert.NotNil(t, c.Flags().Lookup("output"), "--output flag must exist")
			assert.NotNil(t, c.Flags().Lookup("format"), "--format flag must exist")
			assert.Equal(t, "./cases", c.Flags().Lookup("output").DefValue)
			assert.Equal(t, "hurl", c.Flags().Lookup("format").DefValue)
			return
		}
	}
	t.Fatal("ask command not found")
}

func TestAskCommand_RequiresDescription(t *testing.T) {
	err := runAsk(askCmd, []string{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "description")
}

func TestAskCommand_FailsWhenNoProvider(t *testing.T) {
	// Clear all API key env vars so config.Load() returns a noop provider.
	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("GEMINI_API_KEY", "")
	// NoopProvider.IsAvailable() returns false → Generator.Generate returns error.
	err := runAsk(askCmd, []string{"POST /users - create user"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "AI provider")
}
