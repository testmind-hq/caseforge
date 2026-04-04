// internal/rbt/embedder_test.go
package rbt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCosine_IdenticalVectors_IsOne(t *testing.T) {
	a := []float32{1.0, 0.0, 0.0}
	b := []float32{1.0, 0.0, 0.0}
	assert.InDelta(t, 1.0, cosineSimilarity(a, b), 0.001)
}

func TestCosine_OrthogonalVectors_IsZero(t *testing.T) {
	a := []float32{1.0, 0.0}
	b := []float32{0.0, 1.0}
	assert.InDelta(t, 0.0, cosineSimilarity(a, b), 0.001)
}

func TestCosine_OppositeVectors_IsNegativeOne(t *testing.T) {
	a := []float32{1.0, 0.0}
	b := []float32{-1.0, 0.0}
	assert.InDelta(t, -1.0, cosineSimilarity(a, b), 0.001)
}

func TestCosine_ZeroVector_IsZero(t *testing.T) {
	a := []float32{0.0, 0.0}
	b := []float32{1.0, 0.0}
	assert.InDelta(t, 0.0, cosineSimilarity(a, b), 0.001)
}

func TestTopK_ReturnsTopCandidates(t *testing.T) {
	queryEmb := []float32{1.0, 0.0, 0.0}
	chunks := []IndexChunk{
		{File: "a.go", Embedding: []float32{0.9, 0.1, 0.0}},
		{File: "b.go", Embedding: []float32{0.0, 1.0, 0.0}},
		{File: "c.go", Embedding: []float32{0.95, 0.0, 0.1}},
	}
	top := TopKChunks(queryEmb, chunks, 2)
	assert.Len(t, top, 2)
	files := map[string]bool{top[0].File: true, top[1].File: true}
	assert.False(t, files["b.go"], "b.go should not be in top 2")
}

func TestNoopEmbedder_ReturnsZeroVector(t *testing.T) {
	e := &NoopEmbedder{}
	emb, err := e.Embed("hello world")
	assert.NoError(t, err)
	assert.Len(t, emb, 1)
	assert.InDelta(t, 0.0, emb[0], 0.001)
}
