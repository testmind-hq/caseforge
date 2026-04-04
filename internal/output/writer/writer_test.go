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

func makeCase(id, kind, priority, technique string) schema.TestCase {
	return schema.TestCase{
		Schema:      schema.SchemaBaseURL,
		Version:     "1",
		ID:          id,
		Title:       "Test " + id,
		Kind:        kind,
		Priority:    priority,
		GeneratedAt: time.Now(),
		Source:      schema.CaseSource{Technique: technique},
		Steps:       []schema.Step{{ID: "step-main", Type: "test", Method: "GET", Path: "/"}},
	}
}

func TestWriteCreatesIndexJSON(t *testing.T) {
	dir := t.TempDir()
	w := NewJSONSchemaWriter()
	cases := []schema.TestCase{makeCase("TC-0001", "single", "P1", "equivalence_partitioning")}
	require.NoError(t, w.Write(cases, dir, WriteOptions{}))

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
	original := []schema.TestCase{makeCase("TC-0001", "single", "P1", "")}
	require.NoError(t, w.Write(original, dir, WriteOptions{}))
	indexPath := filepath.Join(dir, "index.json")
	read, err := w.Read(indexPath)
	require.NoError(t, err)
	assert.Equal(t, 1, len(read))
	assert.Equal(t, "TC-0001", read[0].ID)
}

// --- metadata tests ---

func TestWriteMeta_SpecHashAndVersion(t *testing.T) {
	dir := t.TempDir()
	w := NewJSONSchemaWriter()
	cases := []schema.TestCase{makeCase("TC-1", "single", "P0", "equivalence_partitioning")}
	opts := WriteOptions{SpecHash: "abc123", CaseforgeVersion: "v1.2.3"}
	require.NoError(t, w.Write(cases, dir, opts))

	data, _ := os.ReadFile(filepath.Join(dir, "index.json"))
	var index IndexFile
	require.NoError(t, json.Unmarshal(data, &index))
	assert.Equal(t, "abc123", index.Meta.SpecHash)
	assert.Equal(t, "v1.2.3", index.Meta.CaseforgeVersion)
}

func TestWriteMeta_ByTechnique(t *testing.T) {
	dir := t.TempDir()
	w := NewJSONSchemaWriter()
	cases := []schema.TestCase{
		makeCase("TC-1", "single", "P0", "equivalence_partitioning"),
		makeCase("TC-2", "single", "P1", "equivalence_partitioning"),
		makeCase("TC-3", "single", "P1", "boundary_value"),
	}
	require.NoError(t, w.Write(cases, dir, WriteOptions{}))

	data, _ := os.ReadFile(filepath.Join(dir, "index.json"))
	var index IndexFile
	require.NoError(t, json.Unmarshal(data, &index))
	assert.Equal(t, 2, index.Meta.ByTechnique["equivalence_partitioning"])
	assert.Equal(t, 1, index.Meta.ByTechnique["boundary_value"])
}

func TestWriteMeta_ByPriority(t *testing.T) {
	dir := t.TempDir()
	w := NewJSONSchemaWriter()
	cases := []schema.TestCase{
		makeCase("TC-1", "single", "P0", "ep"),
		makeCase("TC-2", "single", "P0", "ep"),
		makeCase("TC-3", "single", "P2", "ep"),
	}
	require.NoError(t, w.Write(cases, dir, WriteOptions{}))

	data, _ := os.ReadFile(filepath.Join(dir, "index.json"))
	var index IndexFile
	require.NoError(t, json.Unmarshal(data, &index))
	assert.Equal(t, 2, index.Meta.ByPriority["P0"])
	assert.Equal(t, 1, index.Meta.ByPriority["P2"])
}

func TestWriteMeta_ByKind(t *testing.T) {
	dir := t.TempDir()
	w := NewJSONSchemaWriter()
	cases := []schema.TestCase{
		makeCase("TC-1", "single", "P0", "ep"),
		makeCase("TC-2", "chain", "P1", "chain_crud"),
	}
	require.NoError(t, w.Write(cases, dir, WriteOptions{}))

	data, _ := os.ReadFile(filepath.Join(dir, "index.json"))
	var index IndexFile
	require.NoError(t, json.Unmarshal(data, &index))
	assert.Equal(t, 1, index.Meta.ByKind["single"])
	assert.Equal(t, 1, index.Meta.ByKind["chain"])
}

func TestWriteMeta_EmptyOpts(t *testing.T) {
	dir := t.TempDir()
	w := NewJSONSchemaWriter()
	cases := []schema.TestCase{makeCase("TC-1", "single", "P1", "ep")}
	require.NoError(t, w.Write(cases, dir, WriteOptions{}))

	data, _ := os.ReadFile(filepath.Join(dir, "index.json"))
	var index IndexFile
	require.NoError(t, json.Unmarshal(data, &index))
	assert.Empty(t, index.Meta.SpecHash)
	assert.Empty(t, index.Meta.CaseforgeVersion)
	assert.NotEmpty(t, index.Meta.ByTechnique)
}

func TestHashBytes(t *testing.T) {
	h1 := HashBytes([]byte("hello"))
	h2 := HashBytes([]byte("hello"))
	assert.Equal(t, h1, h2)
	assert.Len(t, h1, 64) // sha256 hex = 32 bytes = 64 chars

	h3 := HashBytes([]byte("world"))
	assert.NotEqual(t, h1, h3)
}

func TestHashFile(t *testing.T) {
	f := filepath.Join(t.TempDir(), "spec.yaml")
	require.NoError(t, os.WriteFile(f, []byte("openapi: 3.0.0"), 0644))
	hash, err := HashFile(f)
	require.NoError(t, err)
	assert.Len(t, hash, 64)
	assert.Equal(t, HashBytes([]byte("openapi: 3.0.0")), hash)
}

func TestHashFile_Missing(t *testing.T) {
	_, err := HashFile("/nonexistent/file.yaml")
	assert.Error(t, err)
}
