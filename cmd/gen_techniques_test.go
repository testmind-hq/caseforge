// cmd/gen_techniques_test.go
// Behavioral acceptance tests for the six Tcases-inspired generation techniques.
// Each test calls runGen() directly with --no-ai so no LLM is needed.
package cmd

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGen_IsolatedNegative_GeneratesCases(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	outDir := t.TempDir()
	genSpec = miniSpec
	genOutput = outDir
	genNoAI = true
	genTechnique = "isolated_negative"
	require.NoError(t, runGen(genCmd, nil))

	cases := readCases(t, outDir)
	require.NotEmpty(t, cases)
	for _, tc := range cases {
		assert.Equal(t, "isolated_negative", tc.Source.Technique)
	}
}

func TestGen_SchemaViolation_GeneratesCases(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	outDir := t.TempDir()
	genSpec = miniSpec
	genOutput = outDir
	genNoAI = true
	genTechnique = "schema_violation"
	require.NoError(t, runGen(genCmd, nil))

	cases := readCases(t, outDir)
	require.NotEmpty(t, cases)
	for _, tc := range cases {
		assert.Equal(t, "schema_violation", tc.Source.Technique)
		require.NotEmpty(t, tc.Steps[0].Assertions)
		// JSON round-trip deserializes integers as float64
		assert.EqualValues(t, 422, tc.Steps[0].Assertions[0].Expected)
	}
}

func TestGen_VariableIrrelevance_RequiresDependencyParams(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	outDir := t.TempDir()
	// mini.yaml has no sort/filter dependency groups → technique applies=false → 0 cases
	genSpec = miniSpec
	genOutput = outDir
	genNoAI = true
	genTechnique = "variable_irrelevance"
	require.NoError(t, runGen(genCmd, nil))
	// No error even when no cases generated — technique just doesn't apply
}

func TestGen_TupleLevel3_FlagAccepted(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	outDir := t.TempDir()
	genSpec = miniSpec
	genOutput = outDir
	genNoAI = true
	genTechnique = "pairwise"
	genTupleLevel = 3
	require.NoError(t, runGen(genCmd, nil))
}

func TestGen_Seed_DeterministicOutput(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	dir1 := t.TempDir()
	genSpec = miniSpec
	genOutput = dir1
	genNoAI = true
	genSeed = 42
	require.NoError(t, runGen(genCmd, nil))
	cases1 := readCases(t, dir1)

	resetGenGlobals(t)()
	dir2 := t.TempDir()
	genSpec = miniSpec
	genOutput = dir2
	genNoAI = true
	genSeed = 42
	require.NoError(t, runGen(genCmd, nil))
	cases2 := readCases(t, dir2)

	require.Equal(t, len(cases1), len(cases2), "same seed must produce same number of cases")
	// Collect the set of titles from both runs; operation processing order can vary
	// between runs (Go map iteration over spec paths), so we compare as sorted sets.
	titles1 := make([]string, len(cases1))
	titles2 := make([]string, len(cases2))
	for i, tc := range cases1 {
		titles1[i] = tc.Title
	}
	for i, tc := range cases2 {
		titles2[i] = tc.Title
	}
	sort.Strings(titles1)
	sort.Strings(titles2)
	assert.Equal(t, titles1, titles2, "same seed must produce the same set of case titles")
}

func TestGen_TupleLevelInvalid_ReturnsError(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	genSpec = miniSpec
	genOutput = t.TempDir()
	genNoAI = true
	genTupleLevel = 5
	err := runGen(genCmd, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "tuple-level")
}
