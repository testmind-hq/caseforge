package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

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

func TestExportCmd_Integration(t *testing.T) {
	// Write a minimal index.json fixture.
	casesDir := t.TempDir()
	index := map[string]any{
		"$schema":      "https://caseforge.dev/schema/v1/index.json",
		"version":      "1",
		"generated_at": "2026-03-23T00:00:00Z",
		"test_cases": []map[string]any{
			{
				"id":           "TC-0001",
				"title":        "GET /ping",
				"kind":         "single",
				"priority":     "P1",
				"tags":         []string{"smoke"},
				"source":       map[string]any{"technique": "boundary_value", "spec_path": "GET /ping", "rationale": ""},
				"steps":        []any{},
				"generated_at": "2026-03-23T00:00:00Z",
			},
		},
	}
	data, err := json.Marshal(index)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(casesDir, "index.json"), data, 0o644))

	outDir := t.TempDir()

	// Run the export command end-to-end for each format.
	for _, format := range []string{"allure", "xray", "testrail"} {
		t.Run(format, func(t *testing.T) {
			var buf bytes.Buffer
			exportCmd.SetOut(&buf)
			// exportCmd.Execute() delegates to rootCmd, so pass "export" as first arg.
			rootCmd.SetArgs([]string{
				"export",
				"--cases", casesDir,
				"--format", format,
				"--output", outDir,
			})
			require.NoError(t, rootCmd.Execute())
			assert.Contains(t, buf.String(), "Exported 1 test cases")

			// Verify output directory was created.
			entries, err := os.ReadDir(filepath.Join(outDir, format))
			require.NoError(t, err)
			assert.NotEmpty(t, entries, "expected output files in %s/", format)
		})
	}
}

func TestExportCmd_InvalidFormat(t *testing.T) {
	casesDir := t.TempDir()
	index := map[string]any{
		"$schema": "x", "version": "1", "generated_at": "2026-03-23T00:00:00Z", "test_cases": []any{},
	}
	data, _ := json.Marshal(index)
	_ = os.WriteFile(filepath.Join(casesDir, "index.json"), data, 0o644)

	rootCmd.SetArgs([]string{"export", "--cases", casesDir, "--format", "bogus", "--output", t.TempDir()})
	err := rootCmd.Execute()
	assert.ErrorContains(t, err, "unknown export format")
}
