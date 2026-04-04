// cmd/rbt_index_test.go
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

func TestRBTIndexCommand_IsRegistered(t *testing.T) {
	var rbtCmd *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Name() == "rbt" {
			rbtCmd = c
			break
		}
	}
	require.NotNil(t, rbtCmd)

	found := false
	for _, sub := range rbtCmd.Commands() {
		if sub.Name() == "index" {
			found = true
			break
		}
	}
	assert.True(t, found, "index subcommand should be registered on rbt")
}

func TestRBTIndexCommand_HasFlags(t *testing.T) {
	var indexCmd *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Name() == "rbt" {
			for _, sub := range c.Commands() {
				if sub.Name() == "index" {
					indexCmd = sub
					break
				}
			}
		}
	}
	require.NotNil(t, indexCmd)
	assert.NotNil(t, indexCmd.Flags().Lookup("spec"))
	assert.NotNil(t, indexCmd.Flags().Lookup("strategy"))
	assert.NotNil(t, indexCmd.Flags().Lookup("out"))
	assert.NotNil(t, indexCmd.Flags().Lookup("overwrite"))
}

func TestRBTIndexCommand_MissingSpec_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"rbt", "index"})
	err := rootCmd.Execute()
	rootCmd.SetArgs([]string{})
	assert.Error(t, err)
}

func TestRBTIndexCommand_LLMStrategy_WritesMapFile(t *testing.T) {
	dir := t.TempDir()
	specFile := filepath.Join(dir, "openapi.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte(`openapi: "3.0.0"
info:
  title: Test
  version: "1"
paths:
  /users:
    post:
      operationId: createUser
`), 0644))

	outFile := filepath.Join(dir, "map.yaml")

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"rbt", "index",
		"--spec", specFile,
		"--src", dir,
		"--strategy", "llm",
		"--out", outFile,
	})
	err := rootCmd.Execute()
	rootCmd.SetArgs([]string{})
	require.NoError(t, err)

	_, statErr := os.Stat(outFile)
	assert.NoError(t, statErr, "map.yaml should be created")
}
