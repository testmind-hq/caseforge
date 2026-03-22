// internal/output/writer/writer_test.go
package writer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/output/schema"
)

func TestWriteCreatesIndexJSON(t *testing.T) {
	dir := t.TempDir()
	w := NewJSONSchemaWriter()
	cases := []schema.TestCase{
		{
			Schema:      schema.SchemaBaseURL,
			Version:     "1",
			ID:          "TC-0001",
			Title:       "Test case 1",
			Kind:        "single",
			Priority:    "P1",
			GeneratedAt: time.Now(),
			Steps: []schema.Step{
				{ID: "step-main", Type: "test", Method: "GET", Path: "/users"},
			},
		},
	}
	require.NoError(t, w.Write(cases, dir))

	indexPath := filepath.Join(dir, "index.json")
	assert.FileExists(t, indexPath)

	data, err := os.ReadFile(indexPath)
	require.NoError(t, err)

	var index IndexFile
	require.NoError(t, json.Unmarshal(data, &index))
	assert.Equal(t, 1, len(index.TestCases))
	assert.Equal(t, "TC-0001", index.TestCases[0].ID)
}

func TestReadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	w := NewJSONSchemaWriter()
	original := []schema.TestCase{
		{Schema: schema.SchemaBaseURL, Version: "1", ID: "TC-0001",
			Title: "Test", Kind: "single", Priority: "P1",
			GeneratedAt: time.Now(),
			Steps: []schema.Step{{ID: "step-main", Type: "test", Method: "GET", Path: "/"}},
		},
	}
	require.NoError(t, w.Write(original, dir))
	indexPath := filepath.Join(dir, "index.json")
	read, err := w.Read(indexPath)
	require.NoError(t, err)
	assert.Equal(t, 1, len(read))
	assert.Equal(t, "TC-0001", read[0].ID)
}
