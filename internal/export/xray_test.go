package export_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/export"
)

func TestXrayExporter_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	err := (&export.XrayExporter{}).Export(sampleCases(), dir)
	require.NoError(t, err)
	_, err = os.Stat(filepath.Join(dir, "xray-import.json"))
	assert.NoError(t, err, "xray-import.json must exist")
}

func TestXrayExporter_FileContent(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, (&export.XrayExporter{}).Export(sampleCases(), dir))

	data, err := os.ReadFile(filepath.Join(dir, "xray-import.json"))
	require.NoError(t, err)
	s := string(data)

	assert.Contains(t, s, `"TC-0001 POST /users - valid email"`)
	assert.Contains(t, s, `"Manual"`)
	assert.Contains(t, s, `"High"`)          // P1 → High
	assert.Contains(t, s, `"POST /users"`)   // step action
	assert.Contains(t, s, `"status_code eq 201"`) // step result
}

func TestXrayExporter_Format(t *testing.T) {
	assert.Equal(t, "xray", (&export.XrayExporter{}).Format())
}

func TestXrayExporter_EmptyCases_ProducesEmptyArray(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, (&export.XrayExporter{}).Export(nil, dir))
	data, err := os.ReadFile(filepath.Join(dir, "xray-import.json"))
	require.NoError(t, err)
	// Verify "tests":[] not "tests":null
	assert.True(t, strings.Contains(string(data), `"tests":[]`) ||
		strings.Contains(string(data), `"tests": []`),
		"empty cases must produce 'tests':[] not null; got: %s", string(data))
}
