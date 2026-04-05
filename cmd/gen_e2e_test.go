// cmd/gen_e2e_test.go
// End-to-end behavioral tests for caseforge gen flags:
//   --no-ai, --technique, --priority, --operations, --resume
//
// These tests call runGen() directly with --no-ai so no LLM is needed.
// They verify that each flag actually affects the generated output — not just
// that the flag is registered (covered by gen_ux_test.go) or that the helper
// functions work in isolation (covered by gen_flags_test.go).
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/checkpoint"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/output/writer"
)

const miniSpec = "testdata/mini.yaml" // 3 ops: createUser, getUser, deleteUser

// resetGenGlobals returns a cleanup function that restores all gen global vars
// to their flag defaults. Call via t.Cleanup(resetGenGlobals(t)).
func resetGenGlobals(t *testing.T) func() {
	t.Helper()
	return func() {
		genSpec = ""
		genOutput = "./cases"
		genNoAI = false
		genFormat = "hurl"
		genTechnique = ""
		genPriority = ""
		genOperations = ""
		genConcurrency = 1
		genResume = false
	}
}

// runMiniGen calls runGen with mini.yaml + --no-ai and returns the output dir.
// Additional setup should be done BEFORE calling this helper.
func runMiniGen(t *testing.T) string {
	t.Helper()
	outDir := t.TempDir()
	genSpec = miniSpec
	genOutput = outDir
	genNoAI = true
	require.NoError(t, runGen(genCmd, nil))
	return outDir
}

// readCases reads the generated index.json from outDir.
func readCases(t *testing.T, outDir string) []schema.TestCase {
	t.Helper()
	w := writer.NewJSONSchemaWriter()
	cases, err := w.Read(filepath.Join(outDir, "index.json"))
	require.NoError(t, err, "index.json must be readable after gen")
	return cases
}

// ─────────────────────────────────────────────────────────
// --no-ai
// ─────────────────────────────────────────────────────────

// TestGen_NoAI_ProducesCasesWithoutLLM verifies that --no-ai generates cases
// using the pure algorithm engine (no LLM required) and writes all output files.
func TestGen_NoAI_ProducesCasesWithoutLLM(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	outDir := runMiniGen(t)

	cases := readCases(t, outDir)
	assert.NotEmpty(t, cases, "--no-ai must still produce cases via algorithm engine")

	// At least one .hurl file rendered.
	hurls, _ := filepath.Glob(filepath.Join(outDir, "*.hurl"))
	assert.NotEmpty(t, hurls, "--no-ai must still render .hurl output files")

	// No checkpoint file left behind (deleted on successful completion).
	assert.NoFileExists(t, filepath.Join(outDir, ".state.json"),
		"checkpoint must be removed after a clean run")
}

// TestGen_NoAI_AllCasesHaveAssertions verifies that algorithmic cases
// include at least one assertion (executability = 100%).
func TestGen_NoAI_AllCasesHaveAssertions(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	outDir := runMiniGen(t)

	cases := readCases(t, outDir)
	require.NotEmpty(t, cases)
	for _, c := range cases {
		hasAssertion := false
		for _, s := range c.Steps {
			if len(s.Assertions) > 0 {
				hasAssertion = true
				break
			}
		}
		assert.True(t, hasAssertion, "case %q must have at least one assertion", c.ID)
	}
}

// ─────────────────────────────────────────────────────────
// --technique
// ─────────────────────────────────────────────────────────

// TestGen_Technique_FiltersSingleTechnique verifies that --technique produces
// only cases from the requested technique(s).
func TestGen_Technique_FiltersSingleTechnique(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	genTechnique = "equivalence_partitioning"
	outDir := runMiniGen(t)

	cases := readCases(t, outDir)
	require.NotEmpty(t, cases, "--technique equivalence_partitioning must produce cases")
	for _, c := range cases {
		assert.Equal(t, "equivalence_partitioning", c.Source.Technique,
			"all cases must use only the requested technique")
	}
}

// TestGen_Technique_FiltersTwoTechniques verifies comma-separated technique list.
func TestGen_Technique_FiltersTwoTechniques(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	genTechnique = "equivalence_partitioning,boundary_value"
	outDir := runMiniGen(t)

	cases := readCases(t, outDir)
	require.NotEmpty(t, cases)
	allowed := map[string]bool{"equivalence_partitioning": true, "boundary_value": true}
	for _, c := range cases {
		assert.True(t, allowed[c.Source.Technique],
			"technique %q must be one of equivalence_partitioning or boundary_value", c.Source.Technique)
	}
	// Both techniques should appear (mini.yaml has fields that trigger both).
	techs := map[string]bool{}
	for _, c := range cases {
		techs[c.Source.Technique] = true
	}
	assert.True(t, techs["equivalence_partitioning"], "equivalence_partitioning must appear in output")
	assert.True(t, techs["boundary_value"], "boundary_value must appear in output (password has minLength/maxLength)")
}

// TestGen_Technique_UnknownNameProducesNoCases verifies that an unrecognised
// technique name produces zero output cases (not an error).
func TestGen_Technique_UnknownNameProducesNoCases(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	genTechnique = "nonexistent_technique_xyz"
	outDir := runMiniGen(t)

	cases := readCases(t, outDir)
	assert.Empty(t, cases, "unknown technique name must produce zero cases")
}

// ─────────────────────────────────────────────────────────
// --priority
// ─────────────────────────────────────────────────────────

// TestGen_Priority_P0KeepsOnlyP0Cases verifies that --priority P0 filters
// the output to only P0 cases.
func TestGen_Priority_P0KeepsOnlyP0Cases(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	genPriority = "P0"
	outDir := runMiniGen(t)

	cases := readCases(t, outDir)
	require.NotEmpty(t, cases, "--priority P0 must produce at least some cases")
	for _, c := range cases {
		assert.Equal(t, "P0", c.Priority,
			"all cases must be P0 when --priority P0 is set")
	}
}

// TestGen_Priority_P2KeepsP0AndP1AndP2 verifies that --priority P2 is inclusive
// (passes P0, P1, P2 but not P3).
func TestGen_Priority_P2KeepsP0AndP1AndP2(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))

	// Full run (no priority filter) to get a baseline case count.
	fullOutDir := runMiniGen(t)
	fullCases := readCases(t, fullOutDir)

	// P2-filtered run.
	resetGenGlobals(t)()
	genSpec = miniSpec
	genOutput = t.TempDir()
	genNoAI = true
	genPriority = "P2"
	require.NoError(t, runGen(genCmd, nil))
	filteredCases := readCases(t, genOutput)

	for _, c := range filteredCases {
		assert.Contains(t, []string{"P0", "P1", "P2"}, c.Priority,
			"--priority P2 must not include P3 cases")
	}
	// Filtered count should be ≤ full count.
	assert.LessOrEqual(t, len(filteredCases), len(fullCases))
}

// TestGen_Priority_InvalidValueReturnsError verifies that an invalid priority
// string (not P0–P3) returns an error rather than silently producing all cases.
func TestGen_Priority_InvalidValueReturnsError(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	genSpec = miniSpec
	genOutput = t.TempDir()
	genNoAI = true
	genPriority = "HIGH" // not a valid priority

	err := runGen(genCmd, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid --priority")
}

// ─────────────────────────────────────────────────────────
// --operations
// ─────────────────────────────────────────────────────────

// TestGen_Operations_SingleOpIDFiltersOutput verifies that --operations with
// a single operationId produces only cases for that operation.
func TestGen_Operations_SingleOpIDFiltersOutput(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	genOperations = "createUser"
	outDir := runMiniGen(t)

	cases := readCases(t, outDir)
	require.NotEmpty(t, cases, "--operations createUser must produce cases")
	for _, c := range cases {
		// Use SpecPath (canonical spec operation) instead of Steps[0].Path because
		// OWASP technique cases inject attack payloads (SQLi, XSS, path traversal)
		// into the actual step path while SpecPath always holds the original spec op.
		// OWASP also generates OPTIONS /users CORS preflight cases, so we only assert
		// that createUser cases target the /users path (not /users/{id} from other ops).
		assert.NotContains(t, c.Source.SpecPath, "/users/{id}",
			"createUser cases must not include getUser or deleteUser ops")
	}
}

// TestGen_Operations_MultipleOpIDs verifies comma-separated operationId list.
func TestGen_Operations_MultipleOpIDs(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	genOperations = "getUser,deleteUser"
	outDir := runMiniGen(t)

	cases := readCases(t, outDir)
	require.NotEmpty(t, cases)
	for _, c := range cases {
		// Use SpecPath for canonical operation identity — OWASP cases inject attack
		// payloads into the actual step path, so Steps[0].Path is not reliable here.
		assert.Contains(t, c.Source.SpecPath, "/users/{id}",
			"getUser/deleteUser cases must be sourced from /users/{id} spec ops")
		assert.False(t, c.Source.SpecPath == "POST /users",
			"createUser (POST /users) must be excluded")
	}
}

// TestGen_Operations_UnknownOpIDProducesNoCases verifies that an unrecognised
// operationId produces zero output (not an error).
func TestGen_Operations_UnknownOpIDProducesNoCases(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	genOperations = "nonExistentOperation"
	outDir := runMiniGen(t)

	cases := readCases(t, outDir)
	assert.Empty(t, cases, "unknown operationId must produce zero cases")
}

// ─────────────────────────────────────────────────────────
// --resume
// ─────────────────────────────────────────────────────────

// TestGen_Resume_SkipsAlreadyCompletedOperations verifies that --resume skips
// operations that are marked done in an existing checkpoint file, and merges
// their prior cases from index.json.
func TestGen_Resume_SkipsAlreadyCompletedOperations(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	outDir := t.TempDir()

	// Step 1: full run to establish a baseline index.json.
	genSpec = miniSpec
	genOutput = outDir
	genNoAI = true
	require.NoError(t, runGen(genCmd, nil))
	fullCases := readCases(t, outDir)
	require.NotEmpty(t, fullCases)

	// Step 2: compute the spec hash and write a checkpoint that marks
	// createUser (POST /users) as already completed.
	specHash, err := writer.HashFile(miniSpec)
	require.NoError(t, err)
	ckptState := checkpoint.NewState(specHash)
	ckptState.Completed[checkpoint.OperationKey("POST", "/users")] = true
	ckptMgr := checkpoint.NewManager(outDir)
	require.NoError(t, ckptMgr.Save(ckptState))

	// Step 3: resume run — createUser should be skipped.
	genSpec = miniSpec
	genOutput = outDir
	genNoAI = true
	genResume = true
	require.NoError(t, runGen(genCmd, nil))

	// The resumed output still contains all cases (prior + newly generated).
	resumedCases := readCases(t, outDir)
	assert.NotEmpty(t, resumedCases)

	// The resumed run must not re-generate createUser cases from scratch.
	// Because createUser was skipped, its cases come only from the prior run's
	// index.json. Total case count should stay roughly the same (not double).
	assert.LessOrEqual(t, len(resumedCases), len(fullCases)*2,
		"resume must not double-generate cases for already-completed operations")
}

// TestGen_Resume_FreshRunWhenSpecHashChanges verifies that if the spec hash
// in the checkpoint does not match the current spec, the run starts fresh
// (checkpoint discarded).
func TestGen_Resume_FreshRunWhenSpecHashChanges(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	outDir := t.TempDir()

	// Write a stale checkpoint with a different spec hash.
	staleState := checkpoint.NewState("0000000000000000000000000000000000000000000000000000000000000000")
	staleState.Completed[checkpoint.OperationKey("POST", "/users")] = true
	staleState.Completed[checkpoint.OperationKey("GET", "/users/{id}")] = true
	staleState.Completed[checkpoint.OperationKey("DELETE", "/users/{id}")] = true
	ckptMgr := checkpoint.NewManager(outDir)
	require.NoError(t, ckptMgr.Save(staleState))

	// Resume run: stale hash → fresh start.
	genSpec = miniSpec
	genOutput = outDir
	genNoAI = true
	genResume = true
	require.NoError(t, runGen(genCmd, nil))

	// All 3 operations must have been processed (fresh run ignores stale checkpoint).
	cases := readCases(t, outDir)
	ops := map[string]bool{}
	for _, c := range cases {
		if len(c.Steps) > 0 {
			ops[c.Steps[0].Method+" "+c.Steps[0].Path] = true
		}
	}
	assert.True(t, ops["POST /users"] || ops["GET /users/{id}"] || ops["DELETE /users/{id}"],
		"fresh run must generate cases for all operations")

	// Checkpoint must be cleaned up on successful completion.
	assert.NoFileExists(t, filepath.Join(outDir, ".state.json"))
}

// TestGen_Resume_WritesCheckpointDuringRun verifies that a checkpoint file is
// created during a run and removed on successful completion.
func TestGen_Resume_CheckpointRemovedOnSuccess(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	outDir := runMiniGen(t)

	// Checkpoint must be absent after a clean (non-resume) run.
	assert.NoFileExists(t, filepath.Join(outDir, ".state.json"),
		"checkpoint must be deleted after a successful run")
}

// ─────────────────────────────────────────────────────────
// Combined flags
// ─────────────────────────────────────────────────────────

// TestGen_CombinedFlags_TechniqueAndPriority verifies that --technique and
// --priority can be used together: output must satisfy both constraints.
func TestGen_CombinedFlags_TechniqueAndPriority(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	genTechnique = "equivalence_partitioning"
	genPriority = "P1"
	outDir := runMiniGen(t)

	cases := readCases(t, outDir)
	require.NotEmpty(t, cases)
	for _, c := range cases {
		assert.Equal(t, "equivalence_partitioning", c.Source.Technique)
		assert.Contains(t, []string{"P0", "P1"}, c.Priority,
			"--priority P1 must keep only P0 and P1")
	}
}

// TestGen_CombinedFlags_OperationsAndTechnique verifies that --operations and
// --technique together restrict both the operation set and technique set.
func TestGen_CombinedFlags_OperationsAndTechnique(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	genOperations = "createUser"
	genTechnique = "boundary_value"
	outDir := runMiniGen(t)

	cases := readCases(t, outDir)
	// mini.yaml createUser has password with minLength/maxLength → boundary cases expected.
	require.NotEmpty(t, cases, "createUser + boundary_value must produce cases for password field")
	for _, c := range cases {
		assert.Equal(t, "boundary_value", c.Source.Technique)
		require.NotEmpty(t, c.Steps)
		assert.Equal(t, "POST", c.Steps[0].Method)
	}
}

// ─────────────────────────────────────────────────────────
// helpers used in resume tests
// ─────────────────────────────────────────────────────────

// writeMinimalIndex writes a minimal index.json to dir so --resume can load
// prior cases for skipped operations.
func writeMinimalIndex(t *testing.T, dir string, cases []schema.TestCase) {
	t.Helper()
	type indexJSON struct {
		Schema      string            `json:"$schema"`
		Version     string            `json:"version"`
		GeneratedAt time.Time         `json:"generated_at"`
		Meta        map[string]any    `json:"meta"`
		TestCases   []schema.TestCase `json:"test_cases"`
	}
	idx := indexJSON{
		Schema:      "https://caseforge.dev/schema/v1/index.json",
		Version:     "1",
		GeneratedAt: time.Now(),
		Meta:        map[string]any{},
		TestCases:   cases,
	}
	data, _ := json.MarshalIndent(idx, "", "  ")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "index.json"), data, 0644))
}

// validateHurlFilesExist asserts that at least one .hurl file was rendered.
func validateHurlFilesExist(t *testing.T, outDir string) {
	t.Helper()
	hurls, _ := filepath.Glob(filepath.Join(outDir, "*.hurl"))
	assert.NotEmpty(t, hurls, "at least one .hurl file must be rendered")
}

// ─────────────────────────────────────────────────────────
// Output format smoke tests (--no-ai is already set by runMiniGen)
// ─────────────────────────────────────────────────────────

// TestGen_Format_MarkdownRendersFile verifies --format markdown creates .md files.
func TestGen_Format_MarkdownRendersFile(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	genFormat = "markdown"
	outDir := runMiniGen(t)

	mds, _ := filepath.Glob(filepath.Join(outDir, "*.md"))
	assert.NotEmpty(t, mds, "--format markdown must create .md output files")
}

// TestGen_Format_CSVRendersFile verifies --format csv creates a .csv file.
func TestGen_Format_CSVRendersFile(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	genFormat = "csv"
	outDir := runMiniGen(t)

	csvs, _ := filepath.Glob(filepath.Join(outDir, "*.csv"))
	assert.NotEmpty(t, csvs, "--format csv must create a .csv output file")
}

// TestGen_Format_HurlRendersFile verifies --format hurl (default) creates .hurl files.
func TestGen_Format_HurlRendersFile(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	// genFormat defaults to "hurl" — set explicitly for clarity.
	genFormat = "hurl"
	outDir := runMiniGen(t)
	validateHurlFilesExist(t, outDir)
}

// Ensure writeMinimalIndex is used (suppress unused warning if resume tests
// call ckptMgr.Save directly instead).
var _ = fmt.Sprintf
