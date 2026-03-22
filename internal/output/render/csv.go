// internal/output/render/csv.go
// CSVRenderer is a Phase 2 stub. It writes a minimal CSV with ID and title only.
// Full implementation (all fields, proper escaping) is scheduled for Phase 2.
package render

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/testmind-hq/caseforge/internal/output/schema"
)

type CSVRenderer struct{}

func NewCSVRenderer() *CSVRenderer { return &CSVRenderer{} }
func (r *CSVRenderer) Format() string { return "csv" }

// Render writes a minimal CSV. Phase 2 will add full field coverage.
func (r *CSVRenderer) Render(cases []schema.TestCase, outDir string) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(outDir, "cases.csv"))
	if err != nil {
		return fmt.Errorf("creating CSV: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	_ = w.Write([]string{"id", "title", "priority", "technique", "spec_path", "rationale"})
	for _, tc := range cases {
		_ = w.Write([]string{tc.ID, tc.Title, tc.Priority,
			tc.Source.Technique, tc.Source.SpecPath, tc.Source.Rationale})
	}
	w.Flush()
	return w.Error()
}
