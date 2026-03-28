package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompletionCommand_IsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "completion [bash|zsh|fish|powershell]" {
			found = true
			break
		}
	}
	assert.True(t, found, "completion command must be registered on rootCmd")
}

func TestCompletionCommand_ValidArgs(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "completion [bash|zsh|fish|powershell]" {
			assert.ElementsMatch(t, []string{"bash", "zsh", "fish", "powershell"}, c.ValidArgs)
			return
		}
	}
	t.Fatal("completion command not found")
}

func TestCompletion_Bash(t *testing.T) {
	var buf bytes.Buffer
	completionCmd.SetOut(&buf)
	require.NoError(t, runCompletion(completionCmd, []string{"bash"}))
	assert.NotEmpty(t, buf.String())
}

func TestCompletion_Zsh(t *testing.T) {
	var buf bytes.Buffer
	completionCmd.SetOut(&buf)
	require.NoError(t, runCompletion(completionCmd, []string{"zsh"}))
	assert.NotEmpty(t, buf.String())
}

func TestCompletion_Fish(t *testing.T) {
	var buf bytes.Buffer
	completionCmd.SetOut(&buf)
	require.NoError(t, runCompletion(completionCmd, []string{"fish"}))
	assert.NotEmpty(t, buf.String())
}

func TestCompletion_PowerShell(t *testing.T) {
	var buf bytes.Buffer
	completionCmd.SetOut(&buf)
	require.NoError(t, runCompletion(completionCmd, []string{"powershell"}))
	assert.NotEmpty(t, buf.String())
}

func TestCompletion_InvalidShell_ReturnsError(t *testing.T) {
	err := runCompletion(completionCmd, []string{"tcsh"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported shell")
}
