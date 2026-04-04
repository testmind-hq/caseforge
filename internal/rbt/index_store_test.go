// internal/rbt/index_store_test.go
package rbt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexStore_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	store := NewIndexStore(dir)

	idx := &LocalIndex{
		Chunks: []IndexChunk{
			{File: "service.go", Fn: "CreateUser", Hash: "abc123", Embedding: []float32{0.1, 0.2, 0.3}},
		},
		SpecOps: []IndexSpecOp{
			{Operation: "POST /users", Description: "Create user", Embedding: []float32{0.15, 0.22, 0.31}},
		},
	}

	require.NoError(t, store.Save(idx))

	loaded, err := store.Load()
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Len(t, loaded.Chunks, 1)
	assert.Equal(t, "service.go", loaded.Chunks[0].File)
	assert.InDelta(t, 0.1, loaded.Chunks[0].Embedding[0], 0.001)
}

func TestIndexStore_LoadMissing_ReturnsNil(t *testing.T) {
	store := NewIndexStore("/nonexistent/path")
	idx, err := store.Load()
	require.NoError(t, err)
	assert.Nil(t, idx)
}

func TestIndexStore_IsStale_ChangedHash(t *testing.T) {
	dir := t.TempDir()
	store := NewIndexStore(dir)
	idx := &LocalIndex{
		Chunks: []IndexChunk{{File: "a.go", Hash: "old123"}},
	}
	require.NoError(t, store.Save(idx))

	assert.True(t, store.IsChunkStale("a.go", "new456"))
	assert.False(t, store.IsChunkStale("a.go", "old123"))
}
