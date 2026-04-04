// cmd/export_test.go
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

func TestExportCmd_IsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "export" {
			found = true
			break
		}
	}
	assert.True(t, found, "export command must be registered on rootCmd")
}

func TestExportCmd_RequiredFlags(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "export" {
			assert.NotNil(t, c.Flags().Lookup("cases"), "--cases flag must exist")
			assert.NotNil(t, c.Flags().Lookup("format"), "--format flag must exist")
			assert.NotNil(t, c.Flags().Lookup("output"), "--output flag must exist")
			return
		}
	}
	t.Fatal("export command not found")
}

func TestRunExport_AllureWritesOutputAndPrintsSuccess(t *testing.T) {
	dir := t.TempDir()
	casesDir := filepath.Join(dir, "cases")
	require.NoError(t, os.MkdirAll(casesDir, 0o755))

	indexJSON := `{
  "$schema": "https://caseforge.dev/schema/v1/index.json",
  "version": "1",
  "generated_at": "2026-01-01T00:00:00Z",
  "meta": {},
  "test_cases": [
    {
      "id": "TC-001", "title": "GET /items", "kind": "single",
      "priority": "P1", "tags": [],
      "source": {"technique": "equivalence_partitioning", "spec_path": "GET /items"},
      "steps": []
    }
  ]
}`
	require.NoError(t, os.WriteFile(filepath.Join(casesDir, "index.json"), []byte(indexJSON), 0o644))

	outDir := filepath.Join(dir, "out")
	exportCases = casesDir
	exportFormat = "allure"
	exportOutput = outDir

	var buf bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)

	require.NoError(t, runExport(cmd, nil))
	assert.Contains(t, buf.String(), "Exported 1 test case")
	entries, err := os.ReadDir(filepath.Join(outDir, "allure"))
	require.NoError(t, err)
	assert.NotEmpty(t, entries)
}

func TestRunExport_InvalidFormat_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	casesDir := filepath.Join(dir, "cases")
	require.NoError(t, os.MkdirAll(casesDir, 0o755))

	indexJSON := `{"$schema":"","version":"1","generated_at":"2026-01-01T00:00:00Z","meta":{},"test_cases":[]}`
	require.NoError(t, os.WriteFile(filepath.Join(casesDir, "index.json"), []byte(indexJSON), 0o644))

	exportCases = casesDir
	exportFormat = "badformat"
	exportOutput = filepath.Join(dir, "out")

	err := runExport(&cobra.Command{}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "badformat")
}

func TestRunExport_MissingIndexJSON_ReturnsError(t *testing.T) {
	exportCases = t.TempDir() // no index.json inside
	exportFormat = "allure"
	exportOutput = t.TempDir()

	err := runExport(&cobra.Command{}, nil)
	assert.Error(t, err)
}
