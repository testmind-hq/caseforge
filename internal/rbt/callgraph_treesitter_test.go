// internal/rbt/callgraph_treesitter_test.go
package rbt

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func callgraphFixtureDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "testdata", "callgraph")
}

func TestTreeSitterCallGraphBuilder_ExtractFuncs_Go(t *testing.T) {
	builder := NewTreeSitterCallGraphBuilder()
	handlerPath := filepath.Join(callgraphFixtureDir(), "handler.go")

	defs, calls, err := builder.ExtractFuncs(handlerPath)
	require.NoError(t, err)

	defNames := make(map[string]bool)
	for _, d := range defs {
		defNames[d.FuncName] = true
	}
	assert.True(t, defNames["Register"], "expected Register in defs")
	assert.True(t, defNames["CreateUser"], "expected CreateUser in defs")

	callNames := make(map[string]bool)
	for _, c := range calls {
		callNames[c.CalleeName] = true
	}
	assert.True(t, len(callNames) > 0, "expected some calls extracted")
}

func TestTreeSitterCallGraphBuilder_UnsupportedExt_ReturnsEmpty(t *testing.T) {
	builder := NewTreeSitterCallGraphBuilder()
	defs, calls, err := builder.ExtractFuncs("/some/file.php")
	require.NoError(t, err)
	assert.Empty(t, defs)
	assert.Empty(t, calls)
}

func TestTreeSitterCallGraphBuilder_NonexistentFile_ReturnsEmpty(t *testing.T) {
	builder := NewTreeSitterCallGraphBuilder()
	defs, calls, err := builder.ExtractFuncs("/nonexistent/file.go")
	require.NoError(t, err)
	assert.Empty(t, defs)
	assert.Empty(t, calls)
}
