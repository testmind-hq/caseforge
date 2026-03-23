package export_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/export"
	"github.com/testmind-hq/caseforge/internal/output/schema"
)

// sampleCases is the shared fixture used by allure, xray, and testrail tests.
func sampleCases() []schema.TestCase {
	return []schema.TestCase{
		{
			ID:       "TC-0001",
			Title:    "POST /users - valid email",
			Priority: "P1",
			Tags:     []string{"users", "auth"},
			Source:   schema.CaseSource{Technique: "equivalence_partitioning", SpecPath: "POST /users"},
			Steps: []schema.Step{
				{
					ID:     "step-main",
					Title:  "send request",
					Type:   "test",
					Method: "POST",
					Path:   "/users",
					Assertions: []schema.Assertion{
						{Target: "status_code", Operator: "eq", Expected: 201},
					},
				},
			},
			GeneratedAt: time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC),
		},
	}
}

func TestAllureExporter_CreatesResultFile(t *testing.T) {
	dir := t.TempDir()
	err := (&export.AllureExporter{}).Export(sampleCases(), dir)
	require.NoError(t, err)

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Contains(t, entries[0].Name(), "-result.json")
}

func TestAllureExporter_ResultFileContent(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, (&export.AllureExporter{}).Export(sampleCases(), dir))

	entries, _ := os.ReadDir(dir)
	data, err := os.ReadFile(filepath.Join(dir, entries[0].Name()))
	require.NoError(t, err)

	var result map[string]any
	require.NoError(t, json.Unmarshal(data, &result))

	assert.Equal(t, "POST /users - valid email", result["name"])
	assert.Equal(t, "unknown", result["status"])

	labels := result["labels"].([]any)
	severityFound := false
	for _, l := range labels {
		lm := l.(map[string]any)
		if lm["name"] == "severity" {
			assert.Equal(t, "critical", lm["value"])
			severityFound = true
		}
	}
	assert.True(t, severityFound, "severity label must be present")

	steps := result["steps"].([]any)
	require.Len(t, steps, 1)
	step := steps[0].(map[string]any)
	assert.Equal(t, "POST /users", step["name"])
	assert.Equal(t, "unknown", step["status"])
	assert.Contains(t, step, "parameters", "parameters key must be present for Allure")
}

func TestAllureExporter_Format(t *testing.T) {
	assert.Equal(t, "allure", (&export.AllureExporter{}).Format())
}

func TestAllureExporter_EmptyCases(t *testing.T) {
	dir := t.TempDir()
	err := (&export.AllureExporter{}).Export(nil, dir)
	require.NoError(t, err)
	entries, _ := os.ReadDir(dir)
	assert.Len(t, entries, 0)
}
