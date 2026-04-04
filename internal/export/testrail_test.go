package export_test

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/export"
)

func TestTestRailExporter_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	err := (&export.TestRailExporter{}).Export(sampleCases(), dir)
	require.NoError(t, err)
	_, err = os.Stat(filepath.Join(dir, "testrail-import.csv"))
	assert.NoError(t, err, "testrail-import.csv must exist")
}

func TestTestRailExporter_HeaderRow(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, (&export.TestRailExporter{}).Export(sampleCases(), dir))

	f, err := os.Open(filepath.Join(dir, "testrail-import.csv"))
	require.NoError(t, err)
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(records), 2)
	assert.Equal(t, []string{"Title", "Type", "Priority", "Section", "Steps", "Expected Results"}, records[0])
}

func TestTestRailExporter_DataRow(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, (&export.TestRailExporter{}).Export(sampleCases(), dir))

	f, _ := os.Open(filepath.Join(dir, "testrail-import.csv"))
	defer f.Close()
	records, _ := csv.NewReader(f).ReadAll()

	row := records[1]
	assert.Equal(t, "TC-0001 POST /users - valid email", row[0]) // Title
	assert.Equal(t, "Automated", row[1])                         // Type
	assert.Equal(t, "High", row[2])                              // Priority (P1 → High)
	assert.Equal(t, "POST /users", row[3])                       // Section (SpecPath)
	assert.Contains(t, row[4], "POST /users")                    // Steps
	assert.Contains(t, row[5], "status_code eq 201")             // Expected Results
}

func TestTestRailExporter_Format(t *testing.T) {
	assert.Equal(t, "testrail", (&export.TestRailExporter{}).Format())
}

func TestTestRailExporter_EmptyCases(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, (&export.TestRailExporter{}).Export(nil, dir))
	f, _ := os.Open(filepath.Join(dir, "testrail-import.csv"))
	defer f.Close()
	records, _ := csv.NewReader(f).ReadAll()
	assert.Len(t, records, 1, "only header row when no cases")
}
