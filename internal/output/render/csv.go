// internal/output/render/csv.go
package render

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/testmind-hq/caseforge/internal/output/schema"
)

type CSVRenderer struct{}

func NewCSVRenderer() *CSVRenderer { return &CSVRenderer{} }
func (r *CSVRenderer) Format() string { return "csv" }

// Render writes cases.csv to outDir with full field coverage.
func (r *CSVRenderer) Render(cases []schema.TestCase, outDir string) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}
	f, err := os.Create(filepath.Join(outDir, "cases.csv"))
	if err != nil {
		return fmt.Errorf("creating CSV: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	header := []string{
		"id", "title", "priority", "kind", "tags",
		"technique", "spec_path", "rationale",
		"method", "path", "steps_count", "generated_at",
	}
	if err := w.Write(header); err != nil {
		return err
	}
	for _, tc := range cases {
		method, path := "", ""
		if len(tc.Steps) > 0 {
			method = tc.Steps[0].Method
			path = tc.Steps[0].Path
		}
		row := []string{
			tc.ID,
			tc.Title,
			tc.Priority,
			tc.Kind,
			strings.Join(tc.Tags, ";"),
			tc.Source.Technique,
			tc.Source.SpecPath,
			tc.Source.Rationale,
			method,
			path,
			fmt.Sprintf("%d", len(tc.Steps)),
			tc.GeneratedAt.UTC().Format("2006-01-02T15:04:05Z"),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}
