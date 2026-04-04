// internal/rbt/index_store.go
package rbt

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// IndexChunk is one function-level code chunk with its embedding.
type IndexChunk struct {
	File      string    `json:"file"`
	Fn        string    `json:"fn"`
	Hash      string    `json:"hash"`
	Embedding []float32 `json:"embedding"`
}

// IndexSpecOp is one spec operation with its embedding.
type IndexSpecOp struct {
	Operation   string    `json:"operation"`
	Description string    `json:"description"`
	Embedding   []float32 `json:"embedding"`
}

// LocalIndex is the full .caseforge-index/index.json content.
type LocalIndex struct {
	Chunks  []IndexChunk  `json:"chunks"`
	SpecOps []IndexSpecOp `json:"spec_ops"`
}

// IndexStore reads and writes the .caseforge-index/ directory.
type IndexStore struct {
	dir string
}

// NewIndexStore creates an IndexStore rooted at dir.
func NewIndexStore(dir string) *IndexStore {
	return &IndexStore{dir: dir}
}

func (s *IndexStore) indexPath() string {
	return filepath.Join(s.dir, "index.json")
}

// Load reads index.json. Returns nil (no error) if it doesn't exist.
func (s *IndexStore) Load() (*LocalIndex, error) {
	data, err := os.ReadFile(s.indexPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var idx LocalIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, err
	}
	return &idx, nil
}

// Save writes idx to index.json, creating the directory if needed.
func (s *IndexStore) Save(idx *LocalIndex) error {
	if err := os.MkdirAll(s.dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.indexPath(), data, 0644)
}

// IsChunkStale returns true if the stored hash for file does not match newHash.
// NOTE: this method reads index.json from disk on each call. For bulk use over
// many files, load the index once with Load() and use the package-level
// isChunkStale(localIdx, file, hash) helper instead.
func (s *IndexStore) IsChunkStale(file, newHash string) bool {
	idx, err := s.Load()
	if err != nil || idx == nil {
		return true
	}
	for _, c := range idx.Chunks {
		if c.File == file {
			return c.Hash != newHash
		}
	}
	return true
}
