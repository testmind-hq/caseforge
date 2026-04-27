// cmd/gen_ux_test.go
// Tests for gen UX enhancements: --resume flag, dynamic flag completion.
package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenResume_FlagExists verifies that --resume is registered on genCmd.
func TestGenResume_FlagExists(t *testing.T) {
	f := genCmd.Flags().Lookup("resume")
	require.NotNil(t, f, "--resume flag must be registered")
	assert.Equal(t, "bool", f.Value.Type())
}

// TestGenCompletion_FlagsHaveCompletionFuncs verifies that --operations,
// --technique, --format, and --priority have completion functions registered.
func TestGenCompletion_FlagsHaveCompletionFuncs(t *testing.T) {
	for _, flagName := range []string{"technique", "format", "priority"} {
		fn, ok := genCmd.GetFlagCompletionFunc(flagName)
		require.True(t, ok, "flag --%s must have a completion function registered", flagName)
		completions, directive := fn(genCmd, nil, "")
		assert.NotEqual(t, cobra.ShellCompDirectiveError, directive,
			"flag --%s completion must not return error directive", flagName)
		assert.NotEmpty(t, completions, "flag --%s must return at least one completion", flagName)
	}
	// --operations also has a completion func (returns nil list without --spec, but func exists).
	_, ok := genCmd.GetFlagCompletionFunc("operations")
	assert.True(t, ok, "flag --operations must have a completion function registered")
}

// TestGenCompletion_TechniqueCompletionContainsKnownNames verifies the
// allTechniqueNames slice is non-empty and contains expected entries.
func TestGenCompletion_TechniqueCompletionContainsKnownNames(t *testing.T) {
	require.NotEmpty(t, allTechniqueNames)
	assert.Contains(t, allTechniqueNames, "equivalence_partitioning")
	assert.Contains(t, allTechniqueNames, "boundary_value")
	assert.Contains(t, allTechniqueNames, "pairwise")
	assert.Contains(t, allTechniqueNames, "owasp_api_top10")
	assert.Contains(t, allTechniqueNames, "classification_tree")
	assert.Contains(t, allTechniqueNames, "orthogonal_array")
}

// TestGenCompletion_TechniqueCompletionCoversAllRegistered verifies that every
// technique registered in the engine is covered in allTechniqueNames.
func TestGenCompletion_TechniqueCompletionCoversAllRegistered(t *testing.T) {
	registered := []string{
		"equivalence_partitioning",
		"boundary_value",
		"decision_table",
		"state_transition",
		"idempotency",
		"pairwise",
		"classification_tree",
		"orthogonal_array",
		"owasp_api_top10",
		"example_extraction",
		"chain_crud",
		"owasp_api_top10_spec",
	}
	nameSet := make(map[string]bool, len(allTechniqueNames))
	for _, n := range allTechniqueNames {
		nameSet[n] = true
	}
	for _, name := range registered {
		assert.True(t, nameSet[name], "allTechniqueNames must include %q", name)
	}
}
