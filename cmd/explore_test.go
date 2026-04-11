// cmd/explore_test.go
package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExploreCommand_IsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "explore" {
			found = true
			break
		}
	}
	assert.True(t, found, "explore command must be registered on rootCmd")
}

func TestExploreCommand_HasRequiredFlags(t *testing.T) {
	var exploreCmd *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Use == "explore" {
			exploreCmd = c
			break
		}
	}
	require.NotNil(t, exploreCmd)

	assert.NotNil(t, exploreCmd.Flags().Lookup("spec"),   "--spec flag required")
	assert.NotNil(t, exploreCmd.Flags().Lookup("target"), "--target flag required")
}

func TestExploreCommand_HasOptionalFlags(t *testing.T) {
	var exploreCmd *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Use == "explore" {
			exploreCmd = c
			break
		}
	}
	require.NotNil(t, exploreCmd)

	assert.NotNil(t, exploreCmd.Flags().Lookup("max-probes"), "--max-probes flag must exist")
	assert.NotNil(t, exploreCmd.Flags().Lookup("output"),     "--output flag must exist")
	assert.NotNil(t, exploreCmd.Flags().Lookup("dry-run"),    "--dry-run flag must exist")
}

func TestExploreCommand_MissingSpecReturnsError(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("spec", "", "")
	cmd.Flags().String("target", "", "")
	cmd.Flags().Int("max-probes", 50, "")
	cmd.Flags().String("output", "./reports", "")
	cmd.Flags().Bool("dry-run", false, "")
	err := runExplore(cmd, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "spec")
}

func TestExploreCommand_MissingTargetWithoutDryRunReturnsError(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("spec", "testdata/petstore.yaml", "")
	cmd.Flags().String("target", "", "")
	cmd.Flags().Int("max-probes", 50, "")
	cmd.Flags().String("output", "./reports", "")
	cmd.Flags().Bool("dry-run", false, "")
	err := runExplore(cmd, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "target")
}

func TestExploreCommand_ExportPool_DryRun(t *testing.T) {
	// explore --dry-run should still write an (empty) pool file when --export-pool is set
	tmp := t.TempDir()
	specFile := filepath.Join(tmp, "spec.yaml")
	const specYAML = `
openapi: "3.0.0"
info: {title: T, version: "1"}
paths:
  /items:
    post:
      responses:
        "201": {description: created}
`
	require.NoError(t, os.WriteFile(specFile, []byte(specYAML), 0644))
	poolFile := filepath.Join(tmp, "pool.json")

	t.Cleanup(func() { _ = exploreCmd.Flags().Set("export-pool", "") })
	rootCmd.SetArgs([]string{
		"explore", "--spec", specFile, "--dry-run", "--export-pool", poolFile,
	})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("explore: %v", err)
	}
	if _, err := os.Stat(poolFile); err != nil {
		t.Errorf("pool file not written: %v", err)
	}
}

func TestExploreCommand_PrioritizeUncoveredFlag(t *testing.T) {
	tmp := t.TempDir()
	const specYAML = `
openapi: "3.0.0"
info: {title: T, version: "1"}
paths:
  /items:
    post:
      responses:
        "201": {description: created}
`
	specFile := filepath.Join(tmp, "spec.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte(specYAML), 0644))
	outDir := filepath.Join(tmp, "reports")

	rootCmd.SetArgs([]string{
		"explore", "--spec", specFile, "--dry-run",
		"--prioritize-uncovered", "--output", outDir,
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("explore with --prioritize-uncovered: %v", err)
	}
}
