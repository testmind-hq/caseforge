// cmd/rbt_test.go
package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRBTCommand_IsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "rbt" {
			found = true
			break
		}
	}
	assert.True(t, found, "rbt command should be registered on rootCmd")
}

func TestRBTCommand_HasRequiredFlags(t *testing.T) {
	var rbtCmd *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Name() == "rbt" {
			rbtCmd = c
			break
		}
	}
	require.NotNil(t, rbtCmd)
	assert.NotNil(t, rbtCmd.Flags().Lookup("spec"))
	assert.NotNil(t, rbtCmd.Flags().Lookup("cases"))
	assert.NotNil(t, rbtCmd.Flags().Lookup("format"))
	assert.NotNil(t, rbtCmd.Flags().Lookup("fail-on"))
	assert.NotNil(t, rbtCmd.Flags().Lookup("dry-run"))
	assert.NotNil(t, rbtCmd.Flags().Lookup("generate"))
	assert.NotNil(t, rbtCmd.Flags().Lookup("no-ai"))
	assert.NotNil(t, rbtCmd.Flags().Lookup("gen-format"))
}

func TestRBTCommand_MissingSpec_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"rbt"})
	err := rootCmd.Execute()
	rootCmd.SetArgs([]string{})
	assert.Error(t, err)
}

func TestRBTCommand_DryRunJSONFormat_ProducesValidJSON(t *testing.T) {
	specPath := filepath.Join("testdata", "petstore.yaml")
	dir := t.TempDir()

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"rbt", "--spec", specPath, "--format", "json",
		"--dry-run", "--output", dir})
	err := rootCmd.Execute()
	rootCmd.SetArgs([]string{})
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(dir, "rbt-report.json"))
	require.NoError(t, err)
	assert.Contains(t, string(data), `"diff_base"`)
}

func TestRBTCommand_DryRunHighFail_ExitsZero(t *testing.T) {
	specPath := filepath.Join("testdata", "petstore.yaml")
	dir := t.TempDir()

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"rbt", "--spec", specPath, "--dry-run",
		"--fail-on", "high", "--output", dir})
	err := rootCmd.Execute()
	rootCmd.SetArgs([]string{})
	assert.NoError(t, err)
}
