// internal/rbt/embedder.go
package rbt

import (
	"context"
	"math"
	"os"
	"sort"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// Embedder generates vector embeddings for text.
type Embedder interface {
	Embed(text string) ([]float32, error)
}

// NoopEmbedder returns a zero vector. Used when no embedding API key is set.
type NoopEmbedder struct{}

func (n *NoopEmbedder) Embed(_ string) ([]float32, error) {
	return []float32{0.0}, nil
}

// OpenAIEmbedder uses the OpenAI text-embedding-3-small model.
type OpenAIEmbedder struct {
	client *openai.Client
	model  string
}

// NewOpenAIEmbedder creates an embedder using OPENAI_API_KEY env var.
// Returns NoopEmbedder if the key is not set.
func NewOpenAIEmbedder() Embedder {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		return &NoopEmbedder{}
	}
	client := openai.NewClient(option.WithAPIKey(key))
	return &OpenAIEmbedder{client: &client, model: "text-embedding-3-small"}
}

// Embed generates an embedding vector for text.
func (e *OpenAIEmbedder) Embed(text string) ([]float32, error) {
	resp, err := e.client.Embeddings.New(context.Background(), openai.EmbeddingNewParams{
		Model: openai.EmbeddingModel(e.model),
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(text),
		},
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, nil
	}
	raw := resp.Data[0].Embedding
	out := make([]float32, len(raw))
	for i, v := range raw {
		out[i] = float32(v)
	}
	return out, nil
}

// cosineSimilarity computes the cosine similarity between two vectors.
func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

type chunkScore struct {
	IndexChunk
	Score float64
}

// TopKChunks returns the k most similar chunks to the query embedding.
func TopKChunks(query []float32, chunks []IndexChunk, k int) []IndexChunk {
	scored := make([]chunkScore, 0, len(chunks))
	for _, c := range chunks {
		if len(c.Embedding) == 0 {
			continue
		}
		scored = append(scored, chunkScore{c, cosineSimilarity(query, c.Embedding)})
	}
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})
	if k > len(scored) {
		k = len(scored)
	}
	out := make([]IndexChunk, k)
	for i := 0; i < k; i++ {
		out[i] = scored[i].IndexChunk
	}
	return out
}

// topKAboveThreshold returns up to k chunks whose cosine similarity to query
// is at least minScore. Chunks with empty embeddings are always skipped.
func topKAboveThreshold(query []float32, chunks []IndexChunk, k int, minScore float64) []IndexChunk {
	scored := make([]chunkScore, 0, len(chunks))
	for _, c := range chunks {
		if len(c.Embedding) == 0 {
			continue
		}
		if s := cosineSimilarity(query, c.Embedding); s >= minScore {
			scored = append(scored, chunkScore{c, s})
		}
	}
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})
	if k > len(scored) {
		k = len(scored)
	}
	out := make([]IndexChunk, k)
	for i := 0; i < k; i++ {
		out[i] = scored[i].IndexChunk
	}
	return out
}
