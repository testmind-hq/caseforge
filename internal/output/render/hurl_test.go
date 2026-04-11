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

// TestHurlRendererStatusCodeRangeAssertions verifies that gte/lt status_code operators
// render as Hurl [Asserts] block entries with `HTTP *` wildcard rather than a mis-matched
// exact status line.
func TestHurlRendererStatusCodeRangeAssertions(t *testing.T) {
	r := NewHurlRenderer("{{base_url}}")
	tc := schema.TestCase{
		ID: "TC-idor", Title: "IDOR range check", Kind: "single",
		Steps: []schema.Step{{
			ID: "step-1", Method: "GET", Path: "/users/99999",
			Assertions: []schema.Assertion{
				{Target: "status_code", Operator: "gte", Expected: 400},
				{Target: "status_code", Operator: "lt", Expected: 500},
			},
		}},
	}
	content := r.renderCase(tc)
	assert.Contains(t, content, "HTTP *")
	assert.Contains(t, content, "status >= 400")
	assert.Contains(t, content, "status < 500")
	// Must not render a concrete status that conflicts with the range
	assert.NotContains(t, content, "HTTP 200")
	assert.NotContains(t, content, "HTTP 500")
}

// TestHurlSingleCaseAppendixBFormat verifies the Appendix B comment format for single cases.
func TestHurlSingleCaseAppendixBFormat(t *testing.T) {
	r := NewHurlRenderer("{{base_url}}")
	tc := sampleCase()
	content := r.renderCase(tc)

	// New required fields
	assert.Contains(t, content, "# case_id=TC-0001")
	assert.Contains(t, content, "# step_id=step-main")
	assert.Contains(t, content, "# step_type=test")
	assert.Contains(t, content, "# priority=P0")
	assert.Contains(t, content, "# title=")
	assert.Contains(t, content, "# technique=equivalence_partitioning")

	// Old format must be absent
	assert.NotContains(t, content, "# id=TC-0001")
	assert.NotContains(t, content, "# spec_path=")
	assert.NotContains(t, content, "# ─────")

	// Single-line separator style
	assert.Contains(t, content, "# ──")
}

// TestHurlSingleCaseNoChainHeaders ensures single cases don't get chain-style headers.
func TestHurlSingleCaseNoChainHeaders(t *testing.T) {
	r := NewHurlRenderer("{{base_url}}")
	tc := sampleCase()
	content := r.renderCase(tc)

	assert.NotContains(t, content, "# ══════")
	assert.NotContains(t, content, "# case_kind=chain")
}

// TestHurlSingleCaseTechniqueOmittedWhenEmpty verifies technique line is omitted when empty.
func TestHurlSingleCaseTechniqueOmittedWhenEmpty(t *testing.T) {
	r := NewHurlRenderer("{{base_url}}")
	tc := sampleCase()
	tc.Source.Technique = ""
	content := r.renderCase(tc)

	assert.NotContains(t, content, "# technique=")
}

// TestHurlChainCaseAppendixBFormat verifies the Appendix B comment format for chain cases.
func TestHurlChainCaseAppendixBFormat(t *testing.T) {
	r := NewHurlRenderer("")
	tc := schema.TestCase{
		ID:       "TC-chain01",
		Title:    "CRUD user lifecycle",
		Priority: "P1",
		Kind:     "chain",
		Source:   schema.CaseSource{Technique: "chain_crud"},
		Steps: []schema.Step{
			{
				ID:     "step-setup",
				Title:  "create user",
				Type:   "setup",
				Method: "POST",
				Path:   "/users",
				Headers: map[string]string{"Content-Type": "application/json"},
				Body:    map[string]any{"name": "Alice"},
				Assertions: []schema.Assertion{
					{Target: "status_code", Operator: "eq", Expected: 201},
					{Target: "jsonpath $.id", Operator: "exists", Expected: nil},
				},
				Captures: []schema.Capture{
					{Name: "userId", From: "jsonpath $.id"},
				},
			},
			{
				ID:        "step-test",
				Title:     "get user",
				Type:      "test",
				Method:    "GET",
				Path:      "/users/{{userId}}",
				DependsOn: []string{"step-setup"},
				Assertions: []schema.Assertion{
					{Target: "status_code", Operator: "eq", Expected: 200},
				},
			},
		},
	}

	content := r.renderCase(tc)

	// Chain header
	assert.Contains(t, content, "# ══════")
	assert.Contains(t, content, "# case_id=TC-chain01")
	assert.Contains(t, content, "# case_kind=chain")
	assert.Contains(t, content, "# priority=P1")

	// Per-step headers (with [type] in the title separator)
	assert.Contains(t, content, "# ── create user [setup]")
	assert.Contains(t, content, "# step_id=step-setup")
	assert.Contains(t, content, "# step_type=setup")
	assert.Contains(t, content, "# title=create user")

	assert.Contains(t, content, "# ── get user [test]")
	assert.Contains(t, content, "# step_id=step-test")
	assert.Contains(t, content, "# step_type=test")
	assert.Contains(t, content, "# title=get user")

	// depends_on present for step-test
	assert.Contains(t, content, "# depends_on=step-setup")

	// technique NOT present in chain case
	assert.NotContains(t, content, "# technique=")

	// spec_path NOT present
	assert.NotContains(t, content, "# spec_path=")

	// Captures and HTTP still present
	assert.Contains(t, content, "[Captures]")
	assert.Contains(t, content, `userId: jsonpath "$.id"`)
	assert.Contains(t, content, "GET {{base_url}}/users/{{userId}}")

	// [Captures] must appear before [Asserts]
	capturesIdx := strings.Index(content, "[Captures]")
	assertsIdx := strings.Index(content, "[Asserts]")
	require.True(t, capturesIdx >= 0, "[Captures] block must be present in output")
	require.True(t, assertsIdx >= 0, "[Asserts] block must be present in output")
	assert.Less(t, capturesIdx, assertsIdx, "[Captures] must appear before [Asserts]")
}

func TestHurlRenderAssertion_NotExists(t *testing.T) {
	r := NewHurlRenderer("{{base_url}}")
	tc := schema.TestCase{
		ID: "TC-wo", Title: "writeonly check", Kind: "single",
		Steps: []schema.Step{{
			ID: "step-1", Method: "GET", Path: "/users/1",
			Assertions: []schema.Assertion{
				{Target: "status_code", Operator: "eq", Expected: 200},
				{Target: "jsonpath $.password", Operator: "not_exists"},
			},
		}},
	}
	content := r.renderCase(tc)
	assert.Contains(t, content, `jsonpath "$.password" not exists`)
}

// TestHurlChainCaseDependsOnOmittedWhenEmpty verifies depends_on not emitted when empty.
func TestHurlChainCaseDependsOnOmittedWhenEmpty(t *testing.T) {
	r := NewHurlRenderer("")
	tc := schema.TestCase{
		ID:       "TC-chain02",
		Title:    "simple chain",
		Priority: "P2",
		Kind:     "chain",
		Steps: []schema.Step{
			{
				ID:     "step-only",
				Title:  "only step",
				Type:   "test",
				Method: "GET",
				Path:   "/ping",
				Assertions: []schema.Assertion{
					{Target: "status_code", Operator: "eq", Expected: 200},
				},
			},
		},
	}

	content := r.renderCase(tc)
	assert.NotContains(t, content, "# depends_on=")
}

func TestHurlRendererRendersCaptureBlock(t *testing.T) {
	r := NewHurlRenderer("")
	tc := schema.TestCase{
		ID:       "TC-chain01",
		Priority: "P1",
		Kind:     "chain",
		Source:   schema.CaseSource{Technique: "chain_crud"},
		Steps: []schema.Step{
			{
				ID:     "step-setup",
				Title:  "create user",
				Type:   "setup",
				Method: "POST",
				Path:   "/users",
				Headers: map[string]string{"Content-Type": "application/json"},
				Body:    map[string]any{"name": "Alice"},
				Assertions: []schema.Assertion{
					{Target: "status_code", Operator: "eq", Expected: 201},
					{Target: "jsonpath $.id", Operator: "exists", Expected: nil},
				},
				Captures: []schema.Capture{
					{Name: "userId", From: "jsonpath $.id"},
				},
			},
			{
				ID:     "step-test",
				Title:  "get user",
				Type:   "test",
				Method: "GET",
				Path:   "/users/{{userId}}",
				Assertions: []schema.Assertion{
					{Target: "status_code", Operator: "eq", Expected: 200},
				},
			},
		},
	}

	content := r.renderCase(tc)
	assert.Contains(t, content, "[Captures]")
	assert.Contains(t, content, `userId: jsonpath "$.id"`)
	assert.Contains(t, content, "GET {{base_url}}/users/{{userId}}")
	// [Captures] must appear before [Asserts]
	capturesIdx := strings.Index(content, "[Captures]")
	assertsIdx := strings.Index(content, "[Asserts]")
	require.True(t, capturesIdx >= 0, "[Captures] block must be present in output")
	require.True(t, assertsIdx >= 0, "[Asserts] block must be present in output")
	assert.Less(t, capturesIdx, assertsIdx, "[Captures] must appear before [Asserts]")
}

func TestHurlRendererOutputContainsSourceAnnotation(t *testing.T) {
	dir := t.TempDir()
	r := NewHurlRenderer("{{base_url}}")
	require.NoError(t, r.Render([]schema.TestCase{sampleCase()}, dir))

	files, _ := filepath.Glob(filepath.Join(dir, "*.hurl"))
	content, _ := os.ReadFile(files[0])
	s := string(content)
	// Machine-readable annotation (Appendix B format)
	assert.Contains(t, s, "technique=equivalence_partitioning")
	assert.Contains(t, s, "priority=P0")
}

func TestRenderAssertionJSONPath(t *testing.T) {
	tc := schema.TestCase{
		Steps: []schema.Step{{
			Method: "GET", Path: "/users/1",
			Assertions: []schema.Assertion{
				{Target: "status_code", Operator: "eq", Expected: 200},
				{Target: "jsonpath $.is_admin", Operator: "ne", Expected: true},
				{Target: "jsonpath $.role", Operator: "ne", Expected: nil},
			},
		}},
	}
	r := NewHurlRenderer("")
	dir := t.TempDir()
	require.NoError(t, r.Render([]schema.TestCase{tc}, dir))
	files, _ := filepath.Glob(filepath.Join(dir, "*.hurl"))
	require.Len(t, files, 1)
	content, _ := os.ReadFile(files[0])
	assert.Contains(t, string(content), `jsonpath "$.is_admin" != true`)
	assert.Contains(t, string(content), `jsonpath "$.role" not exists`)
}

func TestRenderAssertionHeaderTarget(t *testing.T) {
	tc := schema.TestCase{
		Steps: []schema.Step{{
			Method: "OPTIONS", Path: "/users",
			Assertions: []schema.Assertion{
				{Target: "status_code", Operator: "eq", Expected: 200},
				{Target: "header Access-Control-Allow-Origin", Operator: "ne", Expected: "*"},
			},
		}},
	}
	r := NewHurlRenderer("")
	dir := t.TempDir()
	require.NoError(t, r.Render([]schema.TestCase{tc}, dir))
	files, _ := filepath.Glob(filepath.Join(dir, "*.hurl"))
	require.Len(t, files, 1)
	content, _ := os.ReadFile(files[0])
	assert.Contains(t, string(content), `header "Access-Control-Allow-Origin" != "*"`)
}

func TestRenderAssertion_NewOperators_Hurl(t *testing.T) {
	cases := []struct {
		name     string
		a        schema.Assertion
		contains string
	}{
		{
			name:     "jsonpath lt",
			a:        schema.Assertion{Target: "jsonpath $.age", Operator: "lt", Expected: 100},
			contains: `jsonpath "$.age" < 100`,
		},
		{
			name:     "jsonpath gt",
			a:        schema.Assertion{Target: "jsonpath $.score", Operator: "gt", Expected: 0},
			contains: `jsonpath "$.score" > 0`,
		},
		{
			name:     "jsonpath matches",
			a:        schema.Assertion{Target: "jsonpath $.email", Operator: "matches", Expected: `^.+@.+\..+$`},
			contains: `jsonpath "$.email" matches /^.+@.+\..+$/`,
		},
		{
			name:     "jsonpath is_iso8601",
			a:        schema.Assertion{Target: "jsonpath $.created_at", Operator: "is_iso8601", Expected: nil},
			contains: `jsonpath "$.created_at" isDate`,
		},
		{
			name:     "jsonpath is_uuid",
			a:        schema.Assertion{Target: "jsonpath $.id", Operator: "is_uuid", Expected: nil},
			contains: `jsonpath "$.id" matches /^[0-9a-f]`,
		},
		{
			name:     "duration gt",
			a:        schema.Assertion{Target: "duration_ms", Operator: "gt", Expected: 0},
			contains: `duration > 0`,
		},
		{
			name:     "header exists",
			a:        schema.Assertion{Target: "header Content-Type", Operator: "exists", Expected: nil},
			contains: `header "Content-Type" exists`,
		},
		{
			name:     "header contains",
			a:        schema.Assertion{Target: "header Content-Type", Operator: "contains", Expected: "json"},
			contains: `header "Content-Type" contains "json"`,
		},
		{
			name:     "header matches",
			a:        schema.Assertion{Target: "header Content-Type", Operator: "matches", Expected: `application/.*`},
			contains: `header "Content-Type" matches /application/.*/`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderAssertion(tc.a)
			assert.Contains(t, got, tc.contains, "renderAssertion output")
			assert.NotContains(t, got, "# unrendered assertion", "must not fall through to unrendered")
		})
	}
}
