// internal/output/writer/writer.go
package writer

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/testmind-hq/caseforge/internal/output/schema"
)

// IndexSchemaURL is the JSON Schema URL for the CaseForge index.json file.
const IndexSchemaURL = "https://caseforge.dev/schema/v1/index.json"

// WriteOptions carries optional metadata to embed in index.json.
// Zero value is safe — all fields are omitted when empty.
type WriteOptions struct {
	// SpecHash is the hex-encoded SHA-256 of the source spec file.
	// Compute with HashFile or HashBytes.
	SpecHash string
	// CaseforgeVersion is the binary version string (from ldflags).
	CaseforgeVersion string
}

// IndexMeta holds statistics and provenance data written to index.json.
type IndexMeta struct {
	SpecHash         string         `json:"spec_hash,omitempty"`
	CaseforgeVersion string         `json:"caseforge_version,omitempty"`
	ByTechnique      map[string]int `json:"by_technique,omitempty"`
	ByPriority       map[string]int `json:"by_priority,omitempty"`
	ByKind           map[string]int `json:"by_kind,omitempty"`
}

type SchemaWriter interface {
	Write(cases []schema.TestCase, outDir string, opts WriteOptions) error
	Read(indexPath string) ([]schema.TestCase, error)
}

// IndexFile is the top-level structure of index.json.
type IndexFile struct {
	Schema      string            `json:"$schema"`
	Version     string            `json:"version"`
	GeneratedAt time.Time         `json:"generated_at"`
	Meta        IndexMeta         `json:"meta"`
	TestCases   []schema.TestCase `json:"test_cases"`
}

type JSONSchemaWriter struct{}

func NewJSONSchemaWriter() *JSONSchemaWriter { return &JSONSchemaWriter{} }

// Write serializes cases to index.json in outDir with optional metadata.
// Any pre-existing index.json is overwritten.
func (w *JSONSchemaWriter) Write(cases []schema.TestCase, outDir string, opts WriteOptions) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}
	index := IndexFile{
		Schema:      IndexSchemaURL,
		Version:     "1",
		GeneratedAt: time.Now(),
		Meta:        buildMeta(cases, opts),
		TestCases:   cases,
	}
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling index: %w", err)
	}
	if err := os.WriteFile(filepath.Join(outDir, "index.json"), data, 0644); err != nil {
		return fmt.Errorf("writing index.json: %w", err)
	}
	return nil
}

func (w *JSONSchemaWriter) Read(indexPath string) ([]schema.TestCase, error) {
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("reading index: %w", err)
	}
	var index IndexFile
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("parsing index: %w", err)
	}
	return index.TestCases, nil
}

// buildMeta computes IndexMeta from cases and caller-supplied options.
// The stat maps are nil (and thus omitted from JSON via omitempty) when no
// cases contribute to that dimension.
func buildMeta(cases []schema.TestCase, opts WriteOptions) IndexMeta {
	var byTechnique, byPriority, byKind map[string]int
	for _, tc := range cases {
		if tc.Source.Technique != "" {
			if byTechnique == nil {
				byTechnique = make(map[string]int)
			}
			byTechnique[tc.Source.Technique]++
		}
		if tc.Priority != "" {
			if byPriority == nil {
				byPriority = make(map[string]int)
			}
			byPriority[tc.Priority]++
		}
		if tc.Kind != "" {
			if byKind == nil {
				byKind = make(map[string]int)
			}
			byKind[tc.Kind]++
		}
	}
	return IndexMeta{
		SpecHash:         opts.SpecHash,
		CaseforgeVersion: opts.CaseforgeVersion,
		ByTechnique:      byTechnique,
		ByPriority:       byPriority,
		ByKind:           byKind,
	}
}

// HashFile returns the hex-encoded SHA-256 hash of the file at path.
func HashFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return HashBytes(data), nil
}

// HashBytes returns the hex-encoded SHA-256 hash of b.
func HashBytes(b []byte) string {
	sum := sha256.Sum256(b)
	return fmt.Sprintf("%x", sum)
}
