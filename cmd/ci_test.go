// cmd/ci_test.go
package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCICommand_IsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "ci" {
			found = true
			break
		}
	}
	assert.True(t, found, "ci command should be registered")
}

func TestCIInit_IsRegistered(t *testing.T) {
	var ciC *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Name() == "ci" {
			ciC = c
			break
		}
	}
	require.NotNil(t, ciC)
	found := false
	for _, sub := range ciC.Commands() {
		if sub.Name() == "init" {
			found = true
			break
		}
	}
	assert.True(t, found, "ci init subcommand should be registered")
}

func TestCIInit_GitHubActions(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "api-test.yml")

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"ci", "init", "--platform", "github-actions", "--output", out})
	err := rootCmd.Execute()
	require.NoError(t, err)

	content, err := os.ReadFile(out)
	require.NoError(t, err)
	body := string(content)

	assert.Contains(t, body, "caseforge lint")
	assert.Contains(t, body, "caseforge gen")
	assert.Contains(t, body, "caseforge run")
	assert.Contains(t, body, "ANTHROPIC_API_KEY")
	assert.Contains(t, body, "actions/checkout")
}

func TestCIInit_GitLabCI(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".gitlab-ci.yml")

	rootCmd.SetArgs([]string{"ci", "init", "--platform", "gitlab-ci", "--output", out})
	require.NoError(t, rootCmd.Execute())

	content, _ := os.ReadFile(out)
	body := string(content)
	assert.Contains(t, body, "stages:")
	assert.Contains(t, body, "caseforge lint")
	assert.Contains(t, body, "caseforge gen")
}

func TestCIInit_Jenkins(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "Jenkinsfile")

	rootCmd.SetArgs([]string{"ci", "init", "--platform", "jenkins", "--output", out})
	require.NoError(t, rootCmd.Execute())

	content, _ := os.ReadFile(out)
	body := string(content)
	assert.Contains(t, body, "pipeline")
	assert.Contains(t, body, "caseforge lint")
}

func TestCIInit_Shell(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "api-test.sh")

	rootCmd.SetArgs([]string{"ci", "init", "--platform", "shell", "--output", out})
	require.NoError(t, rootCmd.Execute())

	content, _ := os.ReadFile(out)
	body := string(content)
	assert.True(t, strings.HasPrefix(body, "#!/usr/bin/env bash"))
	assert.Contains(t, body, "caseforge lint")
	assert.Contains(t, body, "caseforge gen")
}

func TestCIInit_DefaultOutputPaths(t *testing.T) {
	platforms := map[string]string{
		"github-actions": ".github/workflows/api-test.yml",
		"gitlab-ci":      ".gitlab-ci.yml",
		"jenkins":        "Jenkinsfile",
		"shell":          "scripts/api-test.sh",
	}
	for platform, expectedPath := range platforms {
		assert.Equal(t, expectedPath, ciDefaultPath[platform],
			"platform %s default path mismatch", platform)
	}
}

func TestCIInit_InvalidPlatform(t *testing.T) {
	rootCmd.SetArgs([]string{"ci", "init", "--platform", "circleci"})
	err := rootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circleci")
}

func TestCIInit_CustomSpecPath(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "workflow.yml")

	rootCmd.SetArgs([]string{"ci", "init", "--platform", "github-actions",
		"--spec", "api/v2/openapi.yaml", "--output", out})
	require.NoError(t, rootCmd.Execute())

	content, _ := os.ReadFile(out)
	assert.Contains(t, string(content), "api/v2/openapi.yaml")
}

