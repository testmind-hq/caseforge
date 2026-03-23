package render

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/output/schema"
)

func singleCase() schema.TestCase {
	return schema.TestCase{
		ID: "TC-0001", Title: "GET users list", Kind: "single", Priority: "P1",
		Steps: []schema.Step{{
			ID: "step-1", Method: "GET", Path: "/users",
			Assertions: []schema.Assertion{
				{Target: "status_code", Operator: "eq", Expected: 200},
			},
		}},
	}
}

func chainCase() schema.TestCase {
	return schema.TestCase{
		ID: "TC-0002", Title: "CRUD chain: /users", Kind: "chain", Priority: "P1",
		Steps: []schema.Step{
			{
				ID: "step-setup", Type: "setup", Method: "POST", Path: "/users",
				Headers: map[string]string{"Content-Type": "application/json"},
				Body: map[string]any{"name": "Alice"},
				Assertions: []schema.Assertion{{Target: "status_code", Operator: "eq", Expected: 201}},
				Captures: []schema.Capture{{Name: "userId", From: "jsonpath $.id"}},
			},
			{
				ID: "step-test", Type: "test", Method: "GET", Path: "/users/{{userId}}",
				Assertions: []schema.Assertion{{Target: "status_code", Operator: "eq", Expected: 200}},
			},
		},
	}
}

func TestPostmanRendererFormat(t *testing.T) {
	r := NewPostmanRenderer()
	assert.Equal(t, "postman", r.Format())
}

func TestPostmanRendererCreatesCollectionJSON(t *testing.T) {
	r := NewPostmanRenderer()
	dir := t.TempDir()
	err := r.Render([]schema.TestCase{singleCase()}, dir)
	require.NoError(t, err)
	_, statErr := os.Stat(filepath.Join(dir, "collection.json"))
	assert.NoError(t, statErr)
}

func TestPostmanCollectionTopLevel(t *testing.T) {
	r := NewPostmanRenderer()
	dir := t.TempDir()
	require.NoError(t, r.Render([]schema.TestCase{singleCase()}, dir))

	data, _ := os.ReadFile(filepath.Join(dir, "collection.json"))
	var coll map[string]any
	require.NoError(t, json.Unmarshal(data, &coll))

	info := coll["info"].(map[string]any)
	assert.Equal(t, "https://schema.getpostman.com/json/collection/v2.1.0/collection.json", info["schema"])
	assert.NotEmpty(t, info["_postman_id"])

	vars := coll["variable"].([]any)
	assert.Len(t, vars, 1)
	baseURL := vars[0].(map[string]any)
	assert.Equal(t, "base_url", baseURL["key"])
}

func TestPostmanSingleCaseItem(t *testing.T) {
	r := NewPostmanRenderer()
	dir := t.TempDir()
	require.NoError(t, r.Render([]schema.TestCase{singleCase()}, dir))

	data, _ := os.ReadFile(filepath.Join(dir, "collection.json"))
	var coll map[string]any
	require.NoError(t, json.Unmarshal(data, &coll))

	items := coll["item"].([]any)
	require.Len(t, items, 1)
	item := items[0].(map[string]any)
	req := item["request"].(map[string]any)
	assert.Equal(t, "GET", req["method"])

	urlObj := req["url"].(map[string]any)
	assert.Contains(t, urlObj["raw"].(string), "{{base_url}}")
}

func TestPostmanChainCaseIsFolder(t *testing.T) {
	r := NewPostmanRenderer()
	dir := t.TempDir()
	require.NoError(t, r.Render([]schema.TestCase{chainCase()}, dir))

	data, _ := os.ReadFile(filepath.Join(dir, "collection.json"))
	var coll map[string]any
	require.NoError(t, json.Unmarshal(data, &coll))

	items := coll["item"].([]any)
	require.Len(t, items, 1)
	folder := items[0].(map[string]any)
	// chain case becomes a folder with sub-items
	subItems := folder["item"].([]any)
	assert.Len(t, subItems, 2, "chain case has 2 steps")
}

func TestPostmanCaptureInSetupStep(t *testing.T) {
	r := NewPostmanRenderer()
	dir := t.TempDir()
	require.NoError(t, r.Render([]schema.TestCase{chainCase()}, dir))

	data, _ := os.ReadFile(filepath.Join(dir, "collection.json"))
	content := string(data)
	assert.Contains(t, content, "pm.environment.set")
	assert.Contains(t, content, "userId")
}

func TestPostmanStatusCodeTestScript(t *testing.T) {
	r := NewPostmanRenderer()
	dir := t.TempDir()
	require.NoError(t, r.Render([]schema.TestCase{singleCase()}, dir))

	data, _ := os.ReadFile(filepath.Join(dir, "collection.json"))
	content := string(data)
	assert.Contains(t, content, "pm.response.to.have.status(200)")
}

func TestPostmanHeaderEqTestScript(t *testing.T) {
	tc := schema.TestCase{
		ID: "TC-header", Title: "header eq check", Kind: "single",
		Steps: []schema.Step{{
			ID: "step-1", Method: "GET", Path: "/resource",
			Assertions: []schema.Assertion{
				{Target: "header Content-Type", Operator: "eq", Expected: "application/json"},
				{Target: "header X-Custom", Operator: "ne", Expected: "forbidden"},
			},
		}},
	}
	r := NewPostmanRenderer()
	dir := t.TempDir()
	require.NoError(t, r.Render([]schema.TestCase{tc}, dir))

	data, _ := os.ReadFile(filepath.Join(dir, "collection.json"))
	content := string(data)
	assert.Contains(t, content, `pm.response.headers.get(\"Content-Type\")).to.eql`)
	assert.Contains(t, content, `pm.response.headers.get(\"X-Custom\")).to.not.eql`)
}
