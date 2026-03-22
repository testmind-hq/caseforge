// internal/output/render/hurl_test.go
package render

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/output/schema"
)

func sampleCase() schema.TestCase {
	return schema.TestCase{
		Schema:   schema.SchemaBaseURL,
		Version:  "1",
		ID:       "TC-0001",
		Title:    "POST /users - valid email",
		Kind:     "single",
		Priority: "P0",
		Source: schema.CaseSource{
			Technique: "equivalence_partitioning",
			SpecPath:  "POST /users requestBody.email",
			Rationale: "valid email in valid partition",
		},
		Steps: []schema.Step{
			{
				ID:     "step-main",
				Title:  "send valid user creation request",
				Type:   "test",
				Method: "POST",
				Path:   "/users",
				Headers: map[string]string{"Content-Type": "application/json"},
				Body:   map[string]any{"email": "test@example.com", "age": 25},
				Assertions: []schema.Assertion{
					{Target: "status_code", Operator: "eq", Expected: 201},
					{Target: "duration_ms", Operator: "lt", Expected: 2000},
					{Target: "body.id", Operator: "exists", Expected: true},
				},
			},
		},
		GeneratedAt: time.Now(),
	}
}

func TestHurlRendererCreatesFiles(t *testing.T) {
	dir := t.TempDir()
	r := NewHurlRenderer("{{base_url}}")
	err := r.Render([]schema.TestCase{sampleCase()}, dir)
	require.NoError(t, err)

	files, _ := filepath.Glob(filepath.Join(dir, "*.hurl"))
	assert.NotEmpty(t, files)
}

func TestHurlRendererOutputContainsMethod(t *testing.T) {
	dir := t.TempDir()
	r := NewHurlRenderer("{{base_url}}")
	require.NoError(t, r.Render([]schema.TestCase{sampleCase()}, dir))

	files, _ := filepath.Glob(filepath.Join(dir, "*.hurl"))
	require.NotEmpty(t, files)
	content, _ := os.ReadFile(files[0])
	s := string(content)
	assert.Contains(t, s, "POST {{base_url}}/users")
}

func TestHurlRendererOutputContainsAssertions(t *testing.T) {
	dir := t.TempDir()
	r := NewHurlRenderer("{{base_url}}")
	require.NoError(t, r.Render([]schema.TestCase{sampleCase()}, dir))

	files, _ := filepath.Glob(filepath.Join(dir, "*.hurl"))
	content, _ := os.ReadFile(files[0])
	s := string(content)
	assert.True(t, strings.Contains(s, "HTTP 201") || strings.Contains(s, "status == 201"))
}

func TestHurlRendererOutputContainsSourceAnnotation(t *testing.T) {
	dir := t.TempDir()
	r := NewHurlRenderer("{{base_url}}")
	require.NoError(t, r.Render([]schema.TestCase{sampleCase()}, dir))

	files, _ := filepath.Glob(filepath.Join(dir, "*.hurl"))
	content, _ := os.ReadFile(files[0])
	s := string(content)
	// Machine-readable annotation
	assert.Contains(t, s, "technique=equivalence_partitioning")
	assert.Contains(t, s, "priority=P0")
}
