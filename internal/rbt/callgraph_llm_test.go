// internal/rbt/callgraph_llm_test.go
package rbt

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLLMCallGraphBuilder_ParsesValidResponse(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "service.go")
	require.NoError(t, os.WriteFile(srcPath, []byte("package svc\n"), 0644))

	resp := `{"definitions":["CreateUser","validateEmail"],"calls":[{"caller":"CreateUser","callees":["validateEmail","Save"]}]}`
	provider := &fakeLLMProvider{response: resp}
	builder := NewLLMCallGraphBuilder(NewLLMParser(provider, ""))

	defs, calls, err := builder.ExtractFuncs(srcPath)
	require.NoError(t, err)

	defNames := make(map[string]bool)
	for _, d := range defs {
		defNames[d.FuncName] = true
	}
	assert.True(t, defNames["CreateUser"])
	assert.True(t, defNames["validateEmail"])

	require.Len(t, calls, 2)
	callees := make(map[string]bool)
	for _, c := range calls {
		assert.Equal(t, "CreateUser", c.CallerFunc)
		callees[c.CalleeName] = true
	}
	assert.True(t, callees["validateEmail"])
	assert.True(t, callees["Save"])
}

func TestLLMCallGraphBuilder_MalformedJSON_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "service.go")
	require.NoError(t, os.WriteFile(srcPath, []byte("package svc\n"), 0644))

	provider := &fakeLLMProvider{response: "not json at all"}
	builder := NewLLMCallGraphBuilder(NewLLMParser(provider, ""))

	defs, calls, err := builder.ExtractFuncs(srcPath)
	require.NoError(t, err)
	assert.Empty(t, defs)
	assert.Empty(t, calls)
}

func TestLLMCallGraphBuilder_EmptyResponse_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "service.go")
	require.NoError(t, os.WriteFile(srcPath, []byte("package svc\n"), 0644))

	provider := &fakeLLMProvider{response: `{"definitions":[],"calls":[]}`}
	builder := NewLLMCallGraphBuilder(NewLLMParser(provider, ""))

	defs, calls, err := builder.ExtractFuncs(srcPath)
	require.NoError(t, err)
	assert.Empty(t, defs)
	assert.Empty(t, calls)
}

func TestLLMCallGraphBuilder_NilProvider_ReturnsEmpty(t *testing.T) {
	builder := NewLLMCallGraphBuilder(nil)
	defs, calls, err := builder.ExtractFuncs("/any/file.go")
	require.NoError(t, err)
	assert.Empty(t, defs)
	assert.Empty(t, calls)
}
