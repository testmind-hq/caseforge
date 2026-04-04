// internal/output/render/k6_test.go
package render

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/output/schema"
)

func TestK6RendererFormat(t *testing.T) {
	r := NewK6Renderer()
	assert.Equal(t, "k6", r.Format())
}

func TestK6RendererCreatesFile(t *testing.T) {
	r := NewK6Renderer()
	dir := t.TempDir()
	err := r.Render([]schema.TestCase{singleCase()}, dir)
	require.NoError(t, err)
	_, err = os.Stat(filepath.Join(dir, "k6_tests.js"))
	assert.NoError(t, err)
}

func TestK6RendererImports(t *testing.T) {
	r := NewK6Renderer()
	dir := t.TempDir()
	require.NoError(t, r.Render([]schema.TestCase{singleCase()}, dir))
	content := readFile(t, filepath.Join(dir, "k6_tests.js"))
	assert.Contains(t, content, `import http from 'k6/http'`)
	assert.Contains(t, content, `import { check, group } from 'k6'`)
	assert.Contains(t, content, `BASE_URL`)
}

func TestK6RendererSingleCase(t *testing.T) {
	r := NewK6Renderer()
	dir := t.TempDir()
	require.NoError(t, r.Render([]schema.TestCase{singleCase()}, dir))
	content := readFile(t, filepath.Join(dir, "k6_tests.js"))
	assert.Contains(t, content, `group('TC-0001`)
	assert.Contains(t, content, `http.get(`)
	assert.Contains(t, content, `r.status === 200`)
}

func TestK6RendererStatusCodeAssertion(t *testing.T) {
	tc := schema.TestCase{
		ID: "TC-status", Title: "status check", Kind: "single",
		Steps: []schema.Step{{
			ID: "step-1", Method: "POST", Path: "/items",
			Assertions: []schema.Assertion{
				{Target: "status_code", Operator: "eq", Expected: 201},
			},
		}},
	}
	r := NewK6Renderer()
	dir := t.TempDir()
	require.NoError(t, r.Render([]schema.TestCase{tc}, dir))
	content := readFile(t, filepath.Join(dir, "k6_tests.js"))
	assert.Contains(t, content, `r.status === 201`)
}

func TestK6RendererJSONPathAssertion(t *testing.T) {
	tc := schema.TestCase{
		ID: "TC-json", Title: "jsonpath check", Kind: "single",
		Steps: []schema.Step{{
			ID: "step-1", Method: "GET", Path: "/users/1",
			Assertions: []schema.Assertion{
				{Target: "jsonpath $.id", Operator: "exists"},
				{Target: "jsonpath $.email", Operator: "eq", Expected: "a@b.com"},
			},
		}},
	}
	r := NewK6Renderer()
	dir := t.TempDir()
	require.NoError(t, r.Render([]schema.TestCase{tc}, dir))
	content := readFile(t, filepath.Join(dir, "k6_tests.js"))
	assert.Contains(t, content, `r.json('id') !== undefined`)
	assert.Contains(t, content, `r.json('email') === "a@b.com"`)
}

func TestK6RendererHeaderAssertion(t *testing.T) {
	tc := schema.TestCase{
		ID: "TC-hdr", Title: "header check", Kind: "single",
		Steps: []schema.Step{{
			ID: "step-1", Method: "GET", Path: "/resource",
			Assertions: []schema.Assertion{
				{Target: "header Content-Type", Operator: "eq", Expected: "application/json"},
			},
		}},
	}
	r := NewK6Renderer()
	dir := t.TempDir()
	require.NoError(t, r.Render([]schema.TestCase{tc}, dir))
	content := readFile(t, filepath.Join(dir, "k6_tests.js"))
	assert.Contains(t, content, `r.headers['Content-Type']`)
}

func TestK6RendererChainCase(t *testing.T) {
	r := NewK6Renderer()
	dir := t.TempDir()
	require.NoError(t, r.Render([]schema.TestCase{chainCase()}, dir))
	content := readFile(t, filepath.Join(dir, "k6_tests.js"))
	assert.Contains(t, content, `group('TC-0002`)
	assert.Contains(t, content, `http.post(`)
	assert.Contains(t, content, `http.get(`)
}

func TestK6RendererCapture(t *testing.T) {
	tc := schema.TestCase{
		ID: "TC-cap", Title: "capture", Kind: "chain",
		Steps: []schema.Step{
			{
				ID: "step-1", Method: "POST", Path: "/users",
				Assertions: []schema.Assertion{{Target: "status_code", Operator: "eq", Expected: 201}},
				Captures:   []schema.Capture{{Name: "userId", From: "jsonpath $.id"}},
			},
			{
				ID: "step-2", Method: "GET", Path: "/users/{{userId}}",
				Assertions: []schema.Assertion{{Target: "status_code", Operator: "eq", Expected: 200}},
			},
		},
	}
	r := NewK6Renderer()
	dir := t.TempDir()
	require.NoError(t, r.Render([]schema.TestCase{tc}, dir))
	content := readFile(t, filepath.Join(dir, "k6_tests.js"))
	assert.Contains(t, content, `const userId = res`)
	assert.Contains(t, content, `${userId}`)
}

func TestK6RendererRequestBody(t *testing.T) {
	tc := schema.TestCase{
		ID: "TC-body", Title: "post body", Kind: "single",
		Steps: []schema.Step{{
			ID: "step-1", Method: "POST", Path: "/users",
			Headers: map[string]string{"Content-Type": "application/json"},
			Body:    map[string]any{"email": "x@y.com", "age": 30},
			Assertions: []schema.Assertion{
				{Target: "status_code", Operator: "eq", Expected: 201},
			},
		}},
	}
	r := NewK6Renderer()
	dir := t.TempDir()
	require.NoError(t, r.Render([]schema.TestCase{tc}, dir))
	content := readFile(t, filepath.Join(dir, "k6_tests.js"))
	assert.Contains(t, content, `JSON.stringify(`)
	assert.Contains(t, content, `email`)
}

func TestK6RendererHEADMethod(t *testing.T) {
	tc := schema.TestCase{
		ID: "TC-head", Title: "head resource", Kind: "single",
		Steps: []schema.Step{{
			ID: "step-1", Method: "HEAD", Path: "/ping",
			Assertions: []schema.Assertion{{Target: "status_code", Operator: "eq", Expected: 200}},
		}},
	}
	r := NewK6Renderer()
	dir := t.TempDir()
	require.NoError(t, r.Render([]schema.TestCase{tc}, dir))
	content := readFile(t, filepath.Join(dir, "k6_tests.js"))
	assert.Contains(t, content, `http.request("HEAD",`)
}

func TestK6RendererDELETEMethod(t *testing.T) {
	tc := schema.TestCase{
		ID: "TC-del", Title: "delete resource", Kind: "single",
		Steps: []schema.Step{{
			ID: "step-1", Method: "DELETE", Path: "/users/1",
			Assertions: []schema.Assertion{{Target: "status_code", Operator: "eq", Expected: 204}},
		}},
	}
	r := NewK6Renderer()
	dir := t.TempDir()
	require.NoError(t, r.Render([]schema.TestCase{tc}, dir))
	content := readFile(t, filepath.Join(dir, "k6_tests.js"))
	assert.Contains(t, content, `http.del(`)
}

func TestK6AssertionLine_NewOperators(t *testing.T) {
	cases := []struct {
		name     string
		a        schema.Assertion
		contains string
	}{
		{
			name:     "jsonpath lt",
			a:        schema.Assertion{Target: "jsonpath $.age", Operator: "lt", Expected: 100},
			contains: `r.json('age') < 100`,
		},
		{
			name:     "jsonpath gt",
			a:        schema.Assertion{Target: "jsonpath $.score", Operator: "gt", Expected: 0},
			contains: `r.json('score') > 0`,
		},
		{
			name:     "jsonpath matches",
			a:        schema.Assertion{Target: "jsonpath $.email", Operator: "matches", Expected: `^.+@.+$`},
			contains: `new RegExp`,
		},
		{
			name:     "jsonpath is_iso8601",
			a:        schema.Assertion{Target: "jsonpath $.created_at", Operator: "is_iso8601", Expected: nil},
			contains: `Date.parse`,
		},
		{
			name:     "jsonpath is_uuid",
			a:        schema.Assertion{Target: "jsonpath $.id", Operator: "is_uuid", Expected: nil},
			contains: `/^[0-9a-f]`,
		},
		{
			name:     "duration gt",
			a:        schema.Assertion{Target: "duration_ms", Operator: "gt", Expected: 50},
			contains: `r.timings.duration > 50`,
		},
		{
			name:     "header contains",
			a:        schema.Assertion{Target: "header Content-Type", Operator: "contains", Expected: "json"},
			contains: `includes`,
		},
		{
			name:     "header matches",
			a:        schema.Assertion{Target: "header Content-Type", Operator: "matches", Expected: `application/.*`},
			contains: `new RegExp`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := k6AssertionLine(tc.a)
			assert.Contains(t, got, tc.contains, "k6AssertionLine output")
			assert.NotContains(t, got, "// unrendered:", "must not fall through to unrendered")
		})
	}
}
