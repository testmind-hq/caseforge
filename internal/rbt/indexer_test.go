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

func TestRunGoCallGraphPhase_NoGoMod_ReturnsEmpty(t *testing.T) {
	// A directory without go.mod should return empty without panic.
	dir := t.TempDir()
	idx := &Indexer{SrcDir: dir, Algo: "rta"}
	unclaimed := []ChangedFile{{Path: filepath.Join(dir, "service.go")}}
	routeFiles := map[string][]RouteMapping{}

	mappings, claimed := idx.runGoCallGraphPhase(unclaimed, routeFiles)
	assert.Empty(t, mappings)
	assert.Empty(t, claimed)
}

// minimalSpec is a valid OpenAPI 3.0 YAML with a single GET /pets operation.
const minimalSpec = `openapi: "3.0.0"
info:
  title: Test API
  version: "1.0.0"
paths:
  /pets:
    get:
      operationId: listPets
      summary: List all pets
      responses:
        "200":
          description: OK
`

// callCountingEmbedder counts how many Embed calls are made.
type callCountingEmbedder struct {
	calls int
	dim   int
}

func (c *callCountingEmbedder) Embed(text string) ([]float32, error) {
	c.calls++
	v := make([]float32, c.dim)
	for i := range v {
		v[i] = float32(c.calls%5+1) / 5.0
	}
	return v, nil
}

func TestRunEmbedPhase_NoopEmbedder_FallsBackToRegex(t *testing.T) {
	dir := t.TempDir()
	srcFile := filepath.Join(dir, "handler.go")
	require.NoError(t, os.WriteFile(srcFile, []byte(`package handler

func Register(r *gin.Engine) {
    r.POST("/users", CreateUser)
}
`), 0644))

	specFile := filepath.Join(dir, "openapi.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte(minimalSpec), 0644))

	idx := &Indexer{
		SrcDir:   dir,
		SpecPath: specFile,
		Store:    NewIndexStore(filepath.Join(dir, ".caseforge-index")),
		Embedder: &NoopEmbedder{},
	}

	files := []ChangedFile{{Path: srcFile}}
	mappings, err := idx.runEmbedPhase(files)
	require.NoError(t, err)
	// NoopEmbedder triggers regex fallback; regex finds POST /users.
	require.NotEmpty(t, mappings)
	assert.Equal(t, "regex", mappings[0].Via)
}

func TestRunEmbedPhase_NoSpecPath_FallsBackToRegex(t *testing.T) {
	dir := t.TempDir()
	srcFile := filepath.Join(dir, "handler.go")
	require.NoError(t, os.WriteFile(srcFile, []byte(`package handler

func Register(r *gin.Engine) {
    r.GET("/items", listItems)
}
`), 0644))

	idx := &Indexer{
		SrcDir:   dir,
		SpecPath: "", // no spec
		Store:    NewIndexStore(filepath.Join(dir, ".caseforge-index")),
		Embedder: &fakeEmbedder{dim: 4},
	}

	files := []ChangedFile{{Path: srcFile}}
	mappings, err := idx.runEmbedPhase(files)
	require.NoError(t, err)
	// No spec ops → no embed mappings → falls back to regex.
	for _, m := range mappings {
		assert.NotEqual(t, "embed", m.Via, "expected regex fallback, not embed")
	}
}

func TestRunEmbedPhase_WithSpec_ReturnsEmbedMappings(t *testing.T) {
	dir := t.TempDir()
	srcFile := filepath.Join(dir, "handler.go")
	require.NoError(t, os.WriteFile(srcFile, []byte("package handler\n\nfunc listPets() {}\n"), 0644))

	specFile := filepath.Join(dir, "openapi.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte(minimalSpec), 0644))

	idx := &Indexer{
		SrcDir:   dir,
		SpecPath: specFile,
		Store:    NewIndexStore(filepath.Join(dir, ".caseforge-index")),
		Embedder: &fakeEmbedder{dim: 4},
	}

	files := []ChangedFile{{Path: srcFile}}
	mappings, err := idx.runEmbedPhase(files)
	require.NoError(t, err)
	require.NotEmpty(t, mappings)
	assert.Equal(t, "embed", mappings[0].Via)
	assert.Equal(t, "GET", mappings[0].Method)
	assert.Equal(t, "/pets", mappings[0].RoutePath)
	assert.Equal(t, srcFile, mappings[0].SourceFile)
}

func TestRunEmbedPhase_SavesSpecOpsAndChunksToIndex(t *testing.T) {
	dir := t.TempDir()
	srcFile := filepath.Join(dir, "svc.go")
	require.NoError(t, os.WriteFile(srcFile, []byte("package svc\n\nfunc List() {}\n"), 0644))

	specFile := filepath.Join(dir, "openapi.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte(minimalSpec), 0644))

	storeDir := filepath.Join(dir, ".caseforge-index")
	idx := &Indexer{
		SrcDir:   dir,
		SpecPath: specFile,
		Store:    NewIndexStore(storeDir),
		Embedder: &fakeEmbedder{dim: 4},
	}

	files := []ChangedFile{{Path: srcFile}}
	_, err := idx.runEmbedPhase(files)
	require.NoError(t, err)

	// Load index and verify chunks and spec ops were persisted.
	localIdx, err := idx.Store.Load()
	require.NoError(t, err)
	require.NotNil(t, localIdx)
	assert.NotEmpty(t, localIdx.Chunks, "file chunks should be saved")
	assert.NotEmpty(t, localIdx.SpecOps, "spec ops should be saved")
	assert.Equal(t, "GET /pets", localIdx.SpecOps[0].Operation)
}

func TestRunEmbedPhase_IncrementalSkipsUnchangedFile(t *testing.T) {
	dir := t.TempDir()
	srcFile := filepath.Join(dir, "handler.go")
	require.NoError(t, os.WriteFile(srcFile, []byte("package handler\n\nfunc Handle() {}\n"), 0644))

	specFile := filepath.Join(dir, "openapi.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte(minimalSpec), 0644))

	counter := &callCountingEmbedder{dim: 4}
	idx := &Indexer{
		SrcDir:   dir,
		SpecPath: specFile,
		Store:    NewIndexStore(filepath.Join(dir, ".caseforge-index")),
		Embedder: counter,
	}
	files := []ChangedFile{{Path: srcFile}}

	// First run: embeds file chunk + spec op.
	_, err := idx.runEmbedPhase(files)
	require.NoError(t, err)
	callsAfterFirst := counter.calls

	// Second run with same file: file unchanged, spec op cached → no new embed calls.
	_, err = idx.runEmbedPhase(files)
	require.NoError(t, err)
	assert.Equal(t, callsAfterFirst, counter.calls, "second run should not re-embed unchanged file or spec op")
}

func TestIsSpecOpStale_MissingFromIndex_ReturnsTrue(t *testing.T) {
	idx := &LocalIndex{}
	assert.True(t, isSpecOpStale(idx, "GET /pets"))
}

func TestIsSpecOpStale_PresentWithEmbedding_ReturnsFalse(t *testing.T) {
	idx := &LocalIndex{
		SpecOps: []IndexSpecOp{{Operation: "GET /pets", Embedding: []float32{0.1, 0.2}}},
	}
	assert.False(t, isSpecOpStale(idx, "GET /pets"))
}

func TestIsSpecOpStale_PresentWithEmptyEmbedding_ReturnsTrue(t *testing.T) {
	idx := &LocalIndex{
		SpecOps: []IndexSpecOp{{Operation: "GET /pets", Embedding: nil}},
	}
	assert.True(t, isSpecOpStale(idx, "GET /pets"))
}

func TestIsSpecOpStale_NilIndex_ReturnsTrue(t *testing.T) {
	assert.True(t, isSpecOpStale(nil, "GET /pets"))
}
