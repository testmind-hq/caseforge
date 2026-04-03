// cmd/explore_test.go
package cmd

import (
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
