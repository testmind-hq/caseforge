package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
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
	t.Cleanup(func() { rootCmd.SetArgs(nil) })
	rootCmd.SetArgs([]string{"ask"})
	err := rootCmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "arg")
}

func TestAskCommand_FailsWhenNoProvider(t *testing.T) {
	// Reset viper and set noop provider — prior tests may have loaded ~/.caseforge.yaml
	// via rootCmd.Execute(), leaving a real API key in viper's state.
	viper.Reset()
	t.Cleanup(func() { viper.Reset() })
	viper.Set("ai.provider", "noop")

	err := runAsk(askCmd, []string{"POST /users - create user"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "AI provider")
}

func TestAskCommand_Integration_NoProvider(t *testing.T) {
	// Use a temp config with provider=noop so initConfig() doesn't pick up
	// the real ~/.caseforge.yaml (which may have a real API key).
	tmpCfg := filepath.Join(t.TempDir(), ".caseforge.yaml")
	require.NoError(t, os.WriteFile(tmpCfg, []byte("ai:\n  provider: noop\n"), 0600))

	viper.Reset()
	t.Cleanup(func() {
		viper.Reset()
		cfgFile = "" // prevent stale temp path from leaking into subsequent rootCmd.Execute() calls
	})

	outDir := t.TempDir()
	t.Cleanup(func() { rootCmd.SetArgs(nil) })
	rootCmd.SetArgs([]string{"ask", "--config", tmpCfg, "--output", outDir, "POST /users create user"})
	err := rootCmd.Execute()
	// With noop provider, Generator.Generate returns "AI provider" error.
	require.Error(t, err)
	assert.Contains(t, err.Error(), "AI provider")
}
