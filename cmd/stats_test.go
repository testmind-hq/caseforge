// cmd/stats_test.go
package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTestIndex(t *testing.T, dir string) {
	t.Helper()
	index := map[string]any{
		"$schema":      "https://caseforge.dev/schema/v1/index.json",
		"version":      "1",
		"generated_at": time.Now().Format(time.RFC3339),
		"meta": map[string]any{
			"caseforge_version": "v0.8.0",
			"by_technique": map[string]int{
				"equivalence_partitioning": 10,
				"boundary_value":           5,
				"owasp":                    8,
			},
			"by_priority": map[string]int{"P0": 6, "P1": 12, "P2": 5},
			"by_kind":     map[string]int{"single": 20, "chain": 3},
		},
		"test_cases": []map[string]any{},
	}
	data, _ := json.MarshalIndent(index, "", "  ")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "index.json"), data, 0644))
}

func TestStatsCommand_IsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "stats" {
			found = true
			break
		}
	}
	assert.True(t, found, "stats command should be registered")
}

func TestStats_TerminalOutput(t *testing.T) {
	dir := t.TempDir()
	writeTestIndex(t, dir)

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetArgs([]string{"stats", "--cases", dir})
	err := rootCmd.Execute()
	require.NoError(t, err)

	output := out.String()
	assert.Contains(t, output, "Total cases:")
	assert.Contains(t, output, "equivalence_partitioning")
	assert.Contains(t, output, "P0")
}

func TestStats_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	writeTestIndex(t, dir)

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetArgs([]string{"stats", "--cases", dir, "--format", "json"})
	err := rootCmd.Execute()
	require.NoError(t, err)

	var result map[string]any
	require.NoError(t, json.Unmarshal(out.Bytes(), &result))
	assert.Equal(t, float64(0), result["total"]) // no test_cases in fixture
	assert.NotEmpty(t, result["generated_at"])
}

func TestStats_MissingIndexJSON(t *testing.T) {
	dir := t.TempDir()
	rootCmd.SetArgs([]string{"stats", "--cases", dir})
	err := rootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "index.json")
}

func TestStats_InvalidFormat(t *testing.T) {
	dir := t.TempDir()
	writeTestIndex(t, dir)
	rootCmd.SetArgs([]string{"stats", "--cases", dir, "--format", "xml"})
	err := rootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "xml")
}

func TestStats_TerminalHasTechniqueBar(t *testing.T) {
	dir := t.TempDir()
	writeTestIndex(t, dir)

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetArgs([]string{"stats", "--cases", dir, "--format", "terminal"})
	require.NoError(t, rootCmd.Execute())

	output := out.String()
	// Should show bar characters for non-zero technique counts
	assert.True(t, strings.Contains(output, "█") || strings.Contains(output, "Technique distribution:"),
		"terminal output should contain technique distribution")
}
