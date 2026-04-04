// cmd/watch_test.go
package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWatchCommand_IsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "watch" {
			found = true
			break
		}
	}
	assert.True(t, found, "watch command should be registered")
}

func TestWatchCommand_HasRequiredSpecFlag(t *testing.T) {
	var watchC *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Name() == "watch" {
			watchC = c
			break
		}
	}
	require.NotNil(t, watchC)
	assert.NotNil(t, watchC.Flags().Lookup("spec"))
	assert.NotNil(t, watchC.Flags().Lookup("output"))
	assert.NotNil(t, watchC.Flags().Lookup("no-ai"))
	assert.NotNil(t, watchC.Flags().Lookup("format"))
}

func TestWatch_RejectsURL(t *testing.T) {
	rootCmd.SetArgs([]string{"watch", "--spec", "https://example.com/openapi.yaml"})
	err := rootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "local files")
}

func TestWatch_RejectsMissingFile(t *testing.T) {
	rootCmd.SetArgs([]string{"watch", "--spec", "/nonexistent/openapi.yaml"})
	err := rootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
