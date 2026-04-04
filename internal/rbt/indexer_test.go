// internal/rbt/indexer_test.go
package rbt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeEmbedder struct{ dim int }

func (f *fakeEmbedder) Embed(text string) ([]float32, error) {
	v := make([]float32, f.dim)
	for i := range v {
		v[i] = float32(len(text)%10) / 10.0
	}
	return v, nil
}

func TestIndexer_RunRegex_WritesMapFile(t *testing.T) {
	dir := t.TempDir()
	srcFile := filepath.Join(dir, "handler.go")
	require.NoError(t, os.WriteFile(srcFile, []byte(`package handler

func Register(r *gin.Engine) {
    r.POST("/users", CreateUser)
}
`), 0644))

	outPath := filepath.Join(dir, "caseforge-map.yaml")
	indexer := &Indexer{
		SrcDir:   dir,
		OutPath:  outPath,
		Store:    NewIndexStore(filepath.Join(dir, ".caseforge-index")),
		Embedder: &fakeEmbedder{dim: 3},
	}

	require.NoError(t, indexer.RunRegex())

	data, err := os.ReadFile(outPath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "mappings:")
}

func TestIndexer_DoesNotOverwrite_WhenFlagFalse(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "caseforge-map.yaml")
	require.NoError(t, os.WriteFile(outPath, []byte("existing: true\n"), 0644))

	indexer := &Indexer{
		SrcDir:    dir,
		OutPath:   outPath,
		Overwrite: false,
	}
	err := indexer.RunRegex()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestIndexer_Overwrites_WhenFlagTrue(t *testing.T) {
	dir := t.TempDir()
	srcFile := filepath.Join(dir, "handler.go")
	require.NoError(t, os.WriteFile(srcFile, []byte("package main\n"), 0644))
	outPath := filepath.Join(dir, "caseforge-map.yaml")
	require.NoError(t, os.WriteFile(outPath, []byte("existing: true\n"), 0644))

	indexer := &Indexer{
		SrcDir:    dir,
		OutPath:   outPath,
		Overwrite: true,
		Store:     NewIndexStore(filepath.Join(dir, ".caseforge-index")),
	}
	require.NoError(t, indexer.RunRegex())
	data, _ := os.ReadFile(outPath)
	assert.False(t, strings.Contains(string(data), "existing: true"))
}
