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

func TestRunTreeSitterPhase_NoTreeSitter_ReturnsEmpty(t *testing.T) {
	t.Setenv("PATH", t.TempDir()) // make tree-sitter unavailable
	dir := t.TempDir()
	srcFile := filepath.Join(dir, "handler.go")
	require.NoError(t, os.WriteFile(srcFile, []byte("package handler\n"), 0644))

	idx := &Indexer{SrcDir: dir}
	files := []ChangedFile{{Path: srcFile}}
	mappings, routeFiles := idx.runTreeSitterPhase(files)
	assert.Empty(t, mappings)
	assert.Empty(t, routeFiles)
}

func TestRunCallGraphPhaseWithBuilder_MockBuilder_TracesChain(t *testing.T) {
	handlerPath := "/tmp/cg-test-idx/handler.go"
	servicePath := "/tmp/cg-test-idx/service.go"
	utilsPath := "/tmp/cg-test-idx/utils.go"

	routeFiles := map[string][]RouteMapping{
		handlerPath: {{SourceFile: handlerPath, Method: "POST", RoutePath: "/users", Via: "treesitter"}},
	}

	allFiles := []ChangedFile{{Path: handlerPath}, {Path: servicePath}, {Path: utilsPath}}
	unclaimed := []ChangedFile{{Path: utilsPath}}

	builder := &mockCallGraphBuilder{
		data: map[string]struct {
			defs  []CallNode
			calls []CallEdge
		}{
			handlerPath: {
				defs:  []CallNode{{File: handlerPath, FuncName: "Register"}},
				calls: []CallEdge{{CallerFile: handlerPath, CallerFunc: "Register", CalleeName: "CreateUser"}},
			},
			servicePath: {
				defs:  []CallNode{{File: servicePath, FuncName: "CreateUser"}},
				calls: []CallEdge{{CallerFile: servicePath, CallerFunc: "CreateUser", CalleeName: "validate"}},
			},
			utilsPath: {
				defs:  []CallNode{{File: utilsPath, FuncName: "validate"}},
				calls: nil,
			},
		},
	}

	idx := &Indexer{Depth: 0}
	mappings, claimed := idx.runCallGraphPhaseWithBuilder(allFiles, unclaimed, routeFiles, builder)
	require.Len(t, mappings, 1)
	assert.Equal(t, "POST", mappings[0].Method)
	assert.Equal(t, "/users", mappings[0].RoutePath)
	assert.Equal(t, "callgraph", mappings[0].Via)
	// utils.go was resolved through the call chain — it should appear in claimed
	require.Len(t, claimed, 1)
	assert.Equal(t, utilsPath, claimed[0].Path)
}
