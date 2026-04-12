package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConformanceCommand_IsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "conformance" {
			found = true
		}
	}
	assert.True(t, found, "conformance command should be registered on rootCmd")
}

func TestConformanceCommand_HasRequiredFlags(t *testing.T) {
	var conformanceCmd *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Use == "conformance" {
			conformanceCmd = c
		}
	}
	require.NotNil(t, conformanceCmd)
	assert.NotNil(t, conformanceCmd.Flags().Lookup("spec"))
	assert.NotNil(t, conformanceCmd.Flags().Lookup("target"))
	assert.NotNil(t, conformanceCmd.Flags().Lookup("output"))
}

func TestConformanceCommand_RequiresSpec(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"conformance", "--target", "http://localhost:8080"})
	err := rootCmd.Execute()
	assert.Error(t, err)
	t.Cleanup(func() { rootCmd.SetArgs(nil) })
}
