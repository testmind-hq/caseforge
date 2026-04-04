// cmd/diff_test.go
package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/output/writer"
)

func resetDiffFlags(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		diffOld = ""
		diffNew = ""
		diffFormat = "text"
		diffCases = ""
		diffGenCases = ""
		diffCmd.SetOut(nil)
		diffCmd.SetErr(nil)
	})
}

func TestDiffCommand_BasicText(t *testing.T) {
	resetDiffFlags(t)
	diffOld = "../testdata/petstore_v1.yaml"
	diffNew = "../testdata/petstore_v2.yaml"
	diffFormat = "text"

	buf := &bytes.Buffer{}
	diffCmd.SetOut(buf)
	diffCmd.SetErr(buf)
	err := runDiff(diffCmd, nil)
	require.ErrorIs(t, err, errBreakingChanges)
	output := buf.String()
	assert.Contains(t, output, "BREAKING")
	assert.Contains(t, output, "/pets/{petId}")
}

func TestDiffCommand_JSONFormat(t *testing.T) {
	resetDiffFlags(t)
	diffOld = "../testdata/petstore_v1.yaml"
	diffNew = "../testdata/petstore_v2.yaml"
	diffFormat = "json"

	buf := &bytes.Buffer{}
	diffCmd.SetOut(buf)
	diffCmd.SetErr(buf)
	err := runDiff(diffCmd, nil)
	require.ErrorIs(t, err, errBreakingChanges)
	output := buf.String()
	assert.Contains(t, output, `"kind"`)
	assert.Contains(t, output, "BREAKING")
}

func TestDiffCommand_HasGenCasesFlag(t *testing.T) {
	assert.NotNil(t, diffCmd.Flags().Lookup("gen-cases"), "--gen-cases flag must exist")
}

func TestDiffCommand_GenCases_WritesIndexJSON(t *testing.T) {
	resetDiffFlags(t)
	diffOld = "../testdata/petstore_v1.yaml"
	diffNew = "../testdata/petstore_v2.yaml"
	diffFormat = "text"
	diffGenCases = t.TempDir()

	buf := &bytes.Buffer{}
	diffCmd.SetOut(buf)
	diffCmd.SetErr(buf)

	err := runDiff(diffCmd, nil)
	// Breaking changes are detected → errBreakingChanges expected.
	require.ErrorIs(t, err, errBreakingChanges)

	// index.json must have been written to diffGenCases.
	indexPath := filepath.Join(diffGenCases, "index.json")
	_, statErr := os.Stat(indexPath)
	require.NoError(t, statErr, "index.json must exist in --gen-cases dir")

	// Must be readable as a proper index file with at least one test case.
	cases, readErr := writer.NewJSONSchemaWriter().Read(indexPath)
	require.NoError(t, readErr)
	assert.NotEmpty(t, cases, "generated index.json must contain at least one test case")

	// Output must mention "Generated".
	assert.Contains(t, buf.String(), "Generated")
}

func TestDiffCommand_GenCases_NoBreakingChanges_NoOutput(t *testing.T) {
	resetDiffFlags(t)
	// identical specs → no breaking changes
	diffOld = "../testdata/petstore_v1.yaml"
	diffNew = "../testdata/petstore_v1.yaml"
	diffFormat = "text"
	diffGenCases = t.TempDir()

	buf := &bytes.Buffer{}
	diffCmd.SetOut(buf)
	diffCmd.SetErr(buf)

	err := runDiff(diffCmd, nil)
	require.NoError(t, err)

	// No index.json should be written when there are no breaking changes.
	_, statErr := os.Stat(filepath.Join(diffGenCases, "index.json"))
	assert.True(t, os.IsNotExist(statErr), "index.json must not exist when there are no breaking changes")
}

func TestDiffCommand_GenCases_JSONIsValid(t *testing.T) {
	resetDiffFlags(t)
	diffOld = "../testdata/petstore_v1.yaml"
	diffNew = "../testdata/petstore_v2.yaml"
	diffFormat = "text"
	diffGenCases = t.TempDir()

	buf := &bytes.Buffer{}
	diffCmd.SetOut(buf)

	_ = runDiff(diffCmd, nil)

	data, err := os.ReadFile(filepath.Join(diffGenCases, "index.json"))
	require.NoError(t, err)

	var raw map[string]any
	require.NoError(t, json.Unmarshal(data, &raw), "index.json must be valid JSON")
	_, hasTestCases := raw["test_cases"]
	assert.True(t, hasTestCases, "index.json must have 'test_cases' key")
}

func TestDiffCommand_Cases_ReadsProperIndexJSON(t *testing.T) {
	resetDiffFlags(t)
	dir := t.TempDir()

	// Write a proper index.json (wrapper format) with a case for an affected endpoint.
	indexJSON := `{
  "$schema": "https://caseforge.dev/schema/v1/index.json",
  "version": "1",
  "generated_at": "2026-01-01T00:00:00Z",
  "meta": {},
  "test_cases": [
    {
      "id": "TC-001", "title": "DELETE /pets/{petId}", "kind": "single",
      "priority": "P1", "tags": [],
      "source": {"technique": "equivalence_partitioning", "spec_path": "DELETE /pets/{petId}"},
      "steps": [{"id": "s1", "title": "t", "type": "test", "method": "DELETE", "path": "/pets/{petId}", "assertions": []}]
    }
  ]
}`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "index.json"), []byte(indexJSON), 0o644))

	diffOld = "../testdata/petstore_v1.yaml"
	diffNew = "../testdata/petstore_v2.yaml"
	diffFormat = "text"
	diffCases = dir

	buf := &bytes.Buffer{}
	diffCmd.SetOut(buf)
	diffCmd.SetErr(buf)

	err := runDiff(diffCmd, nil)
	require.ErrorIs(t, err, errBreakingChanges)

	// The case for the removed endpoint must appear in Affected test cases.
	assert.Contains(t, buf.String(), "Affected test cases")
	assert.Contains(t, buf.String(), "TC-001")
}
