// internal/output/writer/writer.go
package writer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/testmind-hq/caseforge/internal/output/schema"
)

// IndexSchemaURL is the JSON Schema URL for the CaseForge index.json file.
const IndexSchemaURL = "https://caseforge.dev/schema/v1/index.json"

type SchemaWriter interface {
	Write(cases []schema.TestCase, outDir string) error
	Read(indexPath string) ([]schema.TestCase, error)
}

// IndexFile is the top-level structure of index.json.
type IndexFile struct {
	Schema      string            `json:"$schema"`
	Version     string            `json:"version"`
	GeneratedAt time.Time         `json:"generated_at"`
	TestCases   []schema.TestCase `json:"test_cases"`
}

type JSONSchemaWriter struct{}

func NewJSONSchemaWriter() *JSONSchemaWriter { return &JSONSchemaWriter{} }

// Write serializes cases to index.json in outDir. Any pre-existing index.json is overwritten.
func (w *JSONSchemaWriter) Write(cases []schema.TestCase, outDir string) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}
	index := IndexFile{
		Schema:      IndexSchemaURL,
		Version:     "1",
		GeneratedAt: time.Now(),
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
