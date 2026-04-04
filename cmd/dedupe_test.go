// cmd/dedupe_test.go
package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/output/schema"
)

func writeDedupeCaseFile(t *testing.T, dir, id string, steps []schema.Step) string {
	t.Helper()
	tc := schema.TestCase{ID: id, GeneratedAt: time.Now(), Steps: steps}
	data, err := json.Marshal(tc)
	require.NoError(t, err)
	fp := filepath.Join(dir, id+".json")
	require.NoError(t, os.WriteFile(fp, data, 0644))
	return fp
}

func makeTestStep(method, path string, status int, extras ...string) schema.Step {
	assertions := []schema.Assertion{
		{Target: "status_code", Operator: "eq", Expected: float64(status)},
	}
	for _, tgt := range extras {
		assertions = append(assertions, schema.Assertion{Target: tgt, Operator: "eq", Expected: "x"})
	}
	return schema.Step{
		ID: "s1", Type: "test",
		Method: method, Path: path,
		Assertions: assertions,
	}
}

func TestDedupeCommand_IsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "dedupe" {
			found = true
			break
		}
	}
	assert.True(t, found, "dedupe should be registered on rootCmd")
}

func TestDedupeCommand_HasAllFlags(t *testing.T) {
	var dc *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Name() == "dedupe" {
			dc = c
			break
		}
	}
	require.NotNil(t, dc)
	for _, flagName := range []string{"cases", "threshold", "merge", "dry-run", "format"} {
		assert.NotNil(t, dc.Flags().Lookup(flagName), "flag --%s must exist", flagName)
	}
}

func TestDedupeCommand_NonexistentDir_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"dedupe", "--cases", "/nonexistent/path/xyz"})
	err := rootCmd.Execute()
	rootCmd.SetArgs([]string{})
	assert.Error(t, err)
}

func TestDedupeCommand_NoDuplicates_ExitsZero(t *testing.T) {
	dir := t.TempDir()
	writeDedupeCaseFile(t, dir, "case-get", []schema.Step{makeTestStep("GET", "/users", 200)})
	writeDedupeCaseFile(t, dir, "case-post", []schema.Step{makeTestStep("POST", "/users", 201)})

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"dedupe", "--cases", dir})
	err := rootCmd.Execute()
	rootCmd.SetArgs([]string{})
	assert.NoError(t, err)
}

func TestDedupeCommand_ExactDuplicate_ReportsGroup(t *testing.T) {
	dir := t.TempDir()
	step := makeTestStep("POST", "/users", 201, "jsonpath $.id")
	writeDedupeCaseFile(t, dir, "case-a", []schema.Step{step})
	writeDedupeCaseFile(t, dir, "case-b", []schema.Step{step})

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"dedupe", "--cases", dir})
	_ = rootCmd.Execute()
	rootCmd.SetArgs([]string{})

	out := buf.String()
	assert.Contains(t, out, "Group 1")
	assert.Contains(t, out, "exact")
}

func TestDedupeCommand_DryRun_ExitsZeroAndKeepsFiles(t *testing.T) {
	dir := t.TempDir()
	step := makeTestStep("POST", "/users", 201, "jsonpath $.id")
	writeDedupeCaseFile(t, dir, "case-a", []schema.Step{step})
	bPath := writeDedupeCaseFile(t, dir, "case-b", []schema.Step{step})

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"dedupe", "--cases", dir, "--dry-run"})
	err := rootCmd.Execute()
	rootCmd.SetArgs([]string{})

	assert.NoError(t, err, "--dry-run must exit 0 even with duplicates")
	_, statErr := os.Stat(bPath)
	assert.NoError(t, statErr, "--dry-run must not delete files")
}

func TestDedupeCommand_JSONFormat_ValidJSON(t *testing.T) {
	dir := t.TempDir()
	step := makeTestStep("POST", "/users", 201, "jsonpath $.id")
	writeDedupeCaseFile(t, dir, "case-a", []schema.Step{step})
	writeDedupeCaseFile(t, dir, "case-b", []schema.Step{step})

	var stdout bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"dedupe", "--cases", dir, "--format", "json"})
	_ = rootCmd.Execute()
	rootCmd.SetArgs([]string{})

	out := stdout.String()
	assert.True(t, json.Valid([]byte(out)), "stdout must be valid JSON")
	assert.Contains(t, out, `"groups"`)
	assert.Contains(t, out, `"total_scanned"`)
}
