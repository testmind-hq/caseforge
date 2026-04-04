// cmd/suite_test.go
package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	suitepkg "github.com/testmind-hq/caseforge/internal/suite"
)

// ── suite create ─────────────────────────────────────────────────────────────

func TestSuiteCommand_IsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "suite" {
			found = true
			break
		}
	}
	assert.True(t, found, "suite command must be registered on rootCmd")
}

func TestSuiteCreateCmd_HasRequiredFlags(t *testing.T) {
	assert.NotNil(t, suiteCreateCmd.Flags().Lookup("id"), "--id flag must exist")
	assert.NotNil(t, suiteCreateCmd.Flags().Lookup("title"), "--title flag must exist")
	assert.NotNil(t, suiteCreateCmd.Flags().Lookup("kind"), "--kind flag must exist")
	assert.NotNil(t, suiteCreateCmd.Flags().Lookup("cases"), "--cases flag must exist")
	assert.NotNil(t, suiteCreateCmd.Flags().Lookup("output"), "--output flag must exist")
}

func TestSuiteCreate_WritesValidSuiteFile(t *testing.T) {
	dir := t.TempDir()
	outFile := filepath.Join(dir, "suite.json")

	suiteCreateID = "SUITE-001"
	suiteCreateTitle = "E2E Flow"
	suiteCreateKind = "chain"
	suiteCreateCases = "TC-CREATE,TC-READ,TC-DELETE"
	suiteCreateOutput = outFile
	t.Cleanup(func() {
		suiteCreateID = ""
		suiteCreateTitle = ""
		suiteCreateKind = "sequential"
		suiteCreateCases = ""
		suiteCreateOutput = "suite.json"
	})

	var buf bytes.Buffer
	suiteCreateCmd.SetOut(&buf)

	require.NoError(t, runSuiteCreate(suiteCreateCmd, nil))

	// File must exist.
	_, err := os.Stat(outFile)
	require.NoError(t, err)

	// File must be valid JSON with required fields.
	loaded, err := suitepkg.LoadSuiteFile(outFile)
	require.NoError(t, err)
	assert.Equal(t, "SUITE-001", loaded.ID)
	assert.Equal(t, "chain", loaded.Kind)
	assert.Len(t, loaded.Cases, 3)
	assert.Equal(t, schema.SuiteSchemaURL, loaded.Schema)

	// Output must mention the suite ID.
	assert.Contains(t, buf.String(), "SUITE-001")
}

func TestSuiteCreate_NoCases_WritesEmptyCasesList(t *testing.T) {
	dir := t.TempDir()
	outFile := filepath.Join(dir, "suite.json")

	suiteCreateID = "SUITE-EMPTY"
	suiteCreateTitle = "Empty Suite"
	suiteCreateKind = "sequential"
	suiteCreateCases = ""
	suiteCreateOutput = outFile
	t.Cleanup(func() {
		suiteCreateID = ""
		suiteCreateTitle = ""
		suiteCreateKind = "sequential"
		suiteCreateCases = ""
		suiteCreateOutput = "suite.json"
	})

	suiteCreateCmd.SetOut(&bytes.Buffer{})
	require.NoError(t, runSuiteCreate(suiteCreateCmd, nil))

	loaded, err := suitepkg.LoadSuiteFile(outFile)
	require.NoError(t, err)
	assert.Empty(t, loaded.Cases)
}

func TestSuiteCreate_InvalidKind_ReturnsError(t *testing.T) {
	suiteCreateID = "S"
	suiteCreateTitle = "T"
	suiteCreateKind = "invalid"
	suiteCreateOutput = filepath.Join(t.TempDir(), "suite.json")
	t.Cleanup(func() {
		suiteCreateID = ""
		suiteCreateTitle = ""
		suiteCreateKind = "sequential"
		suiteCreateOutput = "suite.json"
	})

	err := runSuiteCreate(suiteCreateCmd, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")
}

// ── suite validate ────────────────────────────────────────────────────────────

func TestSuiteValidateCmd_HasRequiredFlags(t *testing.T) {
	assert.NotNil(t, suiteValidateCmd.Flags().Lookup("suite"), "--suite flag must exist")
	assert.NotNil(t, suiteValidateCmd.Flags().Lookup("cases"), "--cases flag must exist")
}

func TestSuiteValidate_ValidSuiteNoIndex_Passes(t *testing.T) {
	dir := t.TempDir()
	suiteFile := filepath.Join(dir, "suite.json")

	s := &schema.TestSuite{
		ID: "S", Title: "T", Kind: "sequential",
		Cases: []schema.SuiteCase{
			{CaseID: "TC-A"},
			{CaseID: "TC-B", DependsOn: []string{"TC-A"}},
		},
	}
	require.NoError(t, suitepkg.WriteSuiteFile(s, suiteFile))

	suiteValidateSuite = suiteFile
	suiteValidateCases = ""
	t.Cleanup(func() { suiteValidateSuite = ""; suiteValidateCases = "" })

	var buf bytes.Buffer
	suiteValidateCmd.SetOut(&buf)

	require.NoError(t, runSuiteValidate(suiteValidateCmd, nil))
	assert.Contains(t, buf.String(), "valid")
	assert.Contains(t, buf.String(), "TC-A → TC-B")
}

func TestSuiteValidate_WithMatchingIndex_Passes(t *testing.T) {
	dir := t.TempDir()

	// Write index.json with matching case IDs.
	indexJSON := `{
  "$schema": "https://caseforge.dev/schema/v1/index.json",
  "version": "1",
  "generated_at": "2026-01-01T00:00:00Z",
  "meta": {},
  "test_cases": [
    {"id": "TC-CREATE", "title": "create", "kind": "single", "priority": "P1", "tags": [], "source": {"technique": "eq", "spec_path": "POST /users"}, "steps": []},
    {"id": "TC-READ",   "title": "read",   "kind": "single", "priority": "P1", "tags": [], "source": {"technique": "eq", "spec_path": "GET /users"}, "steps": []}
  ]
}`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "index.json"), []byte(indexJSON), 0o644))

	suiteFile := filepath.Join(dir, "suite.json")
	s := &schema.TestSuite{
		ID: "S", Title: "T", Kind: "chain",
		Cases: []schema.SuiteCase{
			{CaseID: "TC-CREATE", Exports: []string{"user_id"}},
			{CaseID: "TC-READ", DependsOn: []string{"TC-CREATE"}},
		},
	}
	require.NoError(t, suitepkg.WriteSuiteFile(s, suiteFile))

	suiteValidateSuite = suiteFile
	suiteValidateCases = dir
	t.Cleanup(func() { suiteValidateSuite = ""; suiteValidateCases = "" })

	var buf bytes.Buffer
	suiteValidateCmd.SetOut(&buf)

	require.NoError(t, runSuiteValidate(suiteValidateCmd, nil))
	assert.Contains(t, buf.String(), "valid")
}

func TestSuiteValidate_CyclicSuite_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	suiteFile := filepath.Join(dir, "suite.json")

	// Write a suite with a cycle directly as JSON (bypassing WriteSuiteFile validation).
	cycleJSON := `{
  "$schema": "https://caseforge.dev/schema/v1/suite.json",
  "id": "S", "title": "T", "kind": "chain",
  "cases": [
    {"case_id": "TC-A", "depends_on": ["TC-B"]},
    {"case_id": "TC-B", "depends_on": ["TC-A"]}
  ]
}`
	require.NoError(t, os.WriteFile(suiteFile, []byte(cycleJSON), 0o644))

	suiteValidateSuite = suiteFile
	suiteValidateCases = ""
	t.Cleanup(func() { suiteValidateSuite = ""; suiteValidateCases = "" })

	var buf bytes.Buffer
	suiteValidateCmd.SetOut(&buf)

	err := runSuiteValidate(suiteValidateCmd, nil)
	require.Error(t, err)
}

func TestSuiteValidate_MissingCaseInIndex_ReturnsError(t *testing.T) {
	dir := t.TempDir()

	// index.json with only TC-CREATE
	indexJSON := `{
  "$schema": "https://caseforge.dev/schema/v1/index.json",
  "version": "1", "generated_at": "2026-01-01T00:00:00Z", "meta": {},
  "test_cases": [{"id": "TC-CREATE", "title": "c", "kind": "single", "priority": "P1", "tags": [], "source": {"technique": "eq", "spec_path": ""}, "steps": []}]
}`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "index.json"), []byte(indexJSON), 0o644))

	suiteFile := filepath.Join(dir, "suite.json")
	s := &schema.TestSuite{
		ID: "S", Title: "T", Kind: "sequential",
		Cases: []schema.SuiteCase{
			{CaseID: "TC-CREATE"},
			{CaseID: "TC-MISSING"}, // not in index
		},
	}
	require.NoError(t, suitepkg.WriteSuiteFile(s, suiteFile))

	suiteValidateSuite = suiteFile
	suiteValidateCases = dir
	t.Cleanup(func() { suiteValidateSuite = ""; suiteValidateCases = "" })

	var buf bytes.Buffer
	suiteValidateCmd.SetOut(&buf)

	err := runSuiteValidate(suiteValidateCmd, nil)
	require.Error(t, err)
	assert.Contains(t, buf.String(), "TC-MISSING")
}

// ── index.json suites field round-trip ───────────────────────────────────────

func TestIndexJSON_SuitesFieldRoundTrip(t *testing.T) {
	dir := t.TempDir()
	suiteFile := filepath.Join(dir, "suite.json")

	s := &schema.TestSuite{
		ID: "SUITE-RT", Title: "Round-trip", Kind: "sequential",
		Cases: []schema.SuiteCase{{CaseID: "TC-1"}},
	}
	require.NoError(t, suitepkg.WriteSuiteFile(s, suiteFile))

	// Verify the suite.json is valid JSON with the $schema key.
	data, err := os.ReadFile(suiteFile)
	require.NoError(t, err)
	var raw map[string]any
	require.NoError(t, json.Unmarshal(data, &raw))
	assert.Equal(t, schema.SuiteSchemaURL, raw["$schema"])
	assert.Equal(t, "SUITE-RT", raw["id"])
}
