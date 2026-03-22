package render

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/output/schema"
)

func TestCSVRendererFullExport(t *testing.T) {
	r := NewCSVRenderer()
	outDir := t.TempDir()

	cases := []schema.TestCase{
		{
			ID:       "TC-abc",
			Title:    "test title",
			Priority: "P1",
			Kind:     "single",
			Tags:     []string{"users", "auth"},
			Source: schema.CaseSource{
				Technique: "equivalence_partitioning",
				SpecPath:  "/users POST",
				Rationale: "valid email",
			},
			Steps: []schema.Step{
				{Method: "POST", Path: "/users"},
			},
			GeneratedAt: time.Date(2026, 3, 22, 0, 0, 0, 0, time.UTC),
		},
	}

	err := r.Render(cases, outDir)
	require.NoError(t, err)

	f, err := os.Open(filepath.Join(outDir, "cases.csv"))
	require.NoError(t, err)
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	require.NoError(t, err)
	require.Len(t, records, 2, "header + 1 data row")

	header := records[0]
	assert.Equal(t, []string{"id", "title", "priority", "kind", "tags", "technique",
		"spec_path", "rationale", "method", "path", "steps_count", "generated_at"}, header)

	row := records[1]
	assert.Equal(t, "TC-abc", row[0])
	assert.Equal(t, "test title", row[1])
	assert.Equal(t, "P1", row[2])
	assert.Equal(t, "single", row[3])
	assert.Contains(t, row[4], "users")   // tags joined
	assert.Equal(t, "equivalence_partitioning", row[5])
	assert.Equal(t, "/users POST", row[6])
	assert.Equal(t, "valid email", row[7])
	assert.Equal(t, "POST", row[8])        // first step method
	assert.Equal(t, "/users", row[9])      // first step path
	assert.Equal(t, "1", row[10])          // steps_count
	assert.Contains(t, row[11], "2026")    // generated_at
}

func TestCSVRendererHandlesMultipleSteps(t *testing.T) {
	r := NewCSVRenderer()
	outDir := t.TempDir()
	cases := []schema.TestCase{
		{
			ID: "TC-chain",
			Steps: []schema.Step{
				{Method: "POST", Path: "/users"},
				{Method: "GET", Path: "/users/{{userId}}"},
			},
		},
	}
	err := r.Render(cases, outDir)
	require.NoError(t, err)

	f, _ := os.Open(filepath.Join(outDir, "cases.csv"))
	defer f.Close()
	records, _ := csv.NewReader(f).ReadAll()
	require.Len(t, records, 2)
	assert.Equal(t, "2", records[1][10])  // steps_count
	assert.Equal(t, "POST", records[1][8])  // first step method
}

func TestCSVWriteError(t *testing.T) {
	r := NewCSVRenderer()
	err := r.Render(nil, "/nonexistent/path/that/does/not/exist")
	assert.Error(t, err)
}
