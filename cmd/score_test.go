// cmd/score_test.go
package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeScoreTestIndex writes an index.json with representative test cases.
func writeScoreTestIndex(t *testing.T, dir string) {
	t.Helper()
	index := map[string]any{
		"$schema":      "https://caseforge.dev/schema/v1/index.json",
		"version":      "1",
		"generated_at": time.Now().Format(time.RFC3339),
		"meta":         map[string]any{},
		"test_cases": []map[string]any{
			{
				"id": "tc-1", "title": "GET /pets eq", "kind": "single", "priority": "P1",
				"tags": []string{},
				"source": map[string]any{
					"technique": "equivalence_partitioning",
					"spec_path": "GET /pets",
					"rationale": "eq test",
				},
				"steps": []map[string]any{
					{
						"id": "s1", "title": "step", "type": "test",
						"method": "GET", "path": "/pets",
						"assertions": []map[string]any{
							{"target": "status_code", "operator": "eq", "expected": 200},
						},
					},
				},
				"generated_at": time.Now().Format(time.RFC3339),
			},
			{
				"id": "tc-2", "title": "GET /pets bv", "kind": "single", "priority": "P1",
				"tags": []string{},
				"source": map[string]any{
					"technique": "boundary_value",
					"spec_path": "GET /pets",
					"rationale": "bv test",
				},
				"steps": []map[string]any{
					{
						"id": "s1", "title": "step", "type": "test",
						"method": "GET", "path": "/pets",
						"assertions": []map[string]any{
							{"target": "status_code", "operator": "eq", "expected": 400},
						},
					},
				},
				"generated_at": time.Now().Format(time.RFC3339),
			},
			{
				"id": "tc-3", "title": "POST /pets owasp", "kind": "single", "priority": "P0",
				"tags": []string{},
				"source": map[string]any{
					"technique": "owasp_api_top10",
					"spec_path": "POST /pets",
					"rationale": "security test",
				},
				"steps": []map[string]any{
					{
						"id": "s1", "title": "step", "type": "test",
						"method": "POST", "path": "/pets",
						"assertions": []map[string]any{
							{"target": "status_code", "operator": "eq", "expected": 400},
						},
					},
				},
				"generated_at": time.Now().Format(time.RFC3339),
			},
		},
	}
	data, _ := json.MarshalIndent(index, "", "  ")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "index.json"), data, 0644))
}

func TestScoreCommand_IsRegistered(t *testing.T) {
	found := false
	for _, sub := range rootCmd.Commands() {
		if sub.Name() == "score" {
			found = true
			break
		}
	}
	assert.True(t, found, "score command must be registered on rootCmd")
}

func TestScoreCommand_HasFlags(t *testing.T) {
	assert.NotNil(t, scoreCmd.Flags().Lookup("cases"))
	assert.NotNil(t, scoreCmd.Flags().Lookup("format"))
}

func TestScoreCommand_TerminalOutput(t *testing.T) {
	dir := t.TempDir()
	writeScoreTestIndex(t, dir)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"score", "--cases", dir})
	err := rootCmd.Execute()
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "Overall:")
	assert.Contains(t, out, "Coverage Breadth")
	assert.Contains(t, out, "Boundary Coverage")
	assert.Contains(t, out, "Security Coverage")
	assert.Contains(t, out, "Executability")
}

func TestScoreCommand_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	writeScoreTestIndex(t, dir)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"score", "--cases", dir, "--format", "json"})
	err := rootCmd.Execute()
	require.NoError(t, err)

	var report map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &report))
	assert.Contains(t, report, "overall")
	assert.Contains(t, report, "dimensions")
	assert.Contains(t, report, "total_cases")
}

func TestScoreCommand_InvalidFormat(t *testing.T) {
	dir := t.TempDir()
	writeScoreTestIndex(t, dir)

	// Reset format flag after this test so stale "xml" value doesn't bleed into
	// subsequent tests (cobra global command state persists between Execute calls).
	t.Cleanup(func() { _ = scoreCmd.Flags().Set("format", "terminal") })

	rootCmd.SetArgs([]string{"score", "--cases", dir, "--format", "xml"})
	err := rootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "xml")
}

func TestScoreCommand_MissingIndexJSON(t *testing.T) {
	dir := t.TempDir()
	rootCmd.SetArgs([]string{"score", "--cases", dir})
	err := rootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "index.json")
}

func TestScoreCommand_OutputContainsSuggestions(t *testing.T) {
	// Write index with only equivalence cases and no security cases → suggestions expected.
	dir := t.TempDir()
	index := map[string]any{
		"$schema":      "https://caseforge.dev/schema/v1/index.json",
		"version":      "1",
		"generated_at": time.Now().Format(time.RFC3339),
		"meta":         map[string]any{},
		"test_cases": []map[string]any{
			{
				"id": "tc-1", "title": "eq", "kind": "single", "priority": "P1",
				"tags": []string{},
				"source": map[string]any{
					"technique": "equivalence_partitioning",
					"spec_path": "POST /orders",
					"rationale": "eq",
				},
				"steps": []map[string]any{
					{
						"id": "s1", "title": "step", "type": "test",
						"method": "POST", "path": "/orders",
						"assertions": []map[string]any{
							{"target": "status_code", "operator": "eq", "expected": 201},
						},
					},
				},
				"generated_at": time.Now().Format(time.RFC3339),
			},
		},
	}
	data, _ := json.MarshalIndent(index, "", "  ")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "index.json"), data, 0644))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"score", "--cases", dir})
	err := rootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Suggestions")
	assert.Contains(t, buf.String(), "owasp")
}
