// internal/dea/explorer_test.go
package dea

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// testServer returns a simple mock API:
// - POST /pets with name="" → 400
// - POST /pets with name=<anything non-empty> → 201
// - POST /pets without name → 400
func testServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/pets" {
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			name, hasName := body["name"]
			if !hasName || name == nil || name == "" {
				w.WriteHeader(400)
				w.Write([]byte(`{"error":"name is required and must be non-empty"}`))
				return
			}
			w.WriteHeader(201)
			w.Write([]byte(`{"id":1,"name":"test"}`))
			return
		}
		w.WriteHeader(404)
	}))
}

func testSpec() *spec.ParsedSpec {
	minL := int64(1)
	return &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{
				Method: "POST",
				Path:   "/pets",
				RequestBody: &spec.RequestBody{
					Required: true,
					Content: map[string]*spec.MediaType{
						"application/json": {
							Schema: &spec.Schema{
								Type:     "object",
								Required: []string{"name"},
								Properties: map[string]*spec.Schema{
									"name": {Type: "string", MinLength: &minL},
									"tag":  {Type: "string"},
								},
							},
						},
					},
				},
				Responses: map[string]*spec.Response{
					"201": {Description: "Created"},
					"400": {Description: "Invalid"},
				},
			},
		},
	}
}

func TestExplorer_ProducesReport(t *testing.T) {
	srv := testServer()
	defer srv.Close()

	explorer := NewExplorer(srv.URL, 50)
	report, err := explorer.Explore(context.Background(), testSpec())
	require.NoError(t, err)
	require.NotNil(t, report)
	assert.Equal(t, srv.URL, report.TargetURL)
	assert.Greater(t, report.TotalProbes, 0)
}

func TestExplorer_DiscoverImplicitMinRule(t *testing.T) {
	srv := testServer()
	defer srv.Close()

	// Use spec WITHOUT declaring minLength to trigger implicit-min probe
	op := &spec.Operation{
		Method: "POST",
		Path:   "/pets",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type:     "object",
						Required: []string{"name"},
						Properties: map[string]*spec.Schema{
							"name": {Type: "string"}, // no minLength
						},
					},
				},
			},
		},
		Responses: map[string]*spec.Response{"201": {Description: "ok"}},
	}

	explorer := NewExplorer(srv.URL, 50)
	report, err := explorer.Explore(context.Background(), &spec.ParsedSpec{Operations: []*spec.Operation{op}})
	require.NoError(t, err)

	// With rule deduplication enabled, an implicit-min rule for the same
	// (operation, category, fieldPath) may be suppressed when a non-implicit
	// required-field rule was already recorded first. The exploration must
	// still produce at least one rule capturing the server's constraint on
	// the 'name' field.
	var nameRules []DiscoveredRule
	for _, r := range report.Rules {
		if r.FieldPath == "requestBody.name" {
			nameRules = append(nameRules, r)
		}
	}
	assert.NotEmpty(t, nameRules, "server rejects empty name → must produce at least one rule for requestBody.name")
}

func TestExplorer_MaxProbesRespected(t *testing.T) {
	srv := testServer()
	defer srv.Close()

	explorer := NewExplorer(srv.URL, 2) // very low limit
	report, err := explorer.Explore(context.Background(), testSpec())
	require.NoError(t, err)
	assert.LessOrEqual(t, report.TotalProbes, 2)
}

func TestExplorer_DryRun_NoHTTPCalls(t *testing.T) {
	// Use an unreachable URL — dry run must not make HTTP calls
	explorer := NewExplorer("http://localhost:1", 50)
	explorer.DryRun = true
	report, err := explorer.Explore(context.Background(), testSpec())
	require.NoError(t, err)
	assert.Greater(t, report.TotalProbes, 0, "dry run counts planned probes")
	assert.NotEmpty(t, report.Rules, "dry run must still produce planned hypotheses as 'pending' rules")
}

func TestExplorer_PrioritizeUncovered_DryRun(t *testing.T) {
	// In dry-run, PrioritizeUncovered falls through to the standard path (priority
	// scheduling is skipped in dry-run since no real HTTP calls are made).
	// The test verifies that setting the field does not break explore.
	e := NewExplorer("", 100)
	e.DryRun = true
	e.PrioritizeUncovered = true

	report, err := e.Explore(context.Background(), testSpec())
	require.NoError(t, err)
	assert.Greater(t, report.TotalProbes, 0, "expected probes > 0 in dry-run with PrioritizeUncovered")
}

func TestExplorer_PrioritizeUncovered_LiveProbes(t *testing.T) {
	// This test exercises the actual exploreWithPriority two-pass scheduling path
	// (DryRun=false, PrioritizeUncovered=true). Pass 1 runs one probe per operation
	// for breadth coverage; Pass 2 allocates remaining budget to ops that didn't get 2xx.
	srv := testServer()
	defer srv.Close()

	e := NewExplorer(srv.URL, 50)
	e.PrioritizeUncovered = true

	report, err := e.Explore(context.Background(), testSpec())
	require.NoError(t, err)
	require.NotNil(t, report)
	assert.Greater(t, report.TotalProbes, 0, "expected live probes > 0")
	assert.Equal(t, srv.URL, report.TargetURL)
	// Pass 1 fires at least one probe per operation; rules should be inferred.
	assert.NotEmpty(t, report.Rules, "expected rules from live two-pass exploration")
}

func TestExplorer_MaxFailures_StopsEarly(t *testing.T) {
	srv := testServer()
	defer srv.Close()

	// testSpec() has fields with constraints; multiple rules will be inferred.
	// Set MaxFailures=1 — should stop after finding the first rule.
	explorer := NewExplorer(srv.URL, 50)
	explorer.MaxFailures = 1
	report, err := explorer.Explore(context.Background(), testSpec())
	require.NoError(t, err)
	// With MaxFailures=1, at most 1 rule should be in the report.
	assert.LessOrEqual(t, len(report.Rules), 1)
}

func TestExplorer_MaxFailures_Zero_NoLimit(t *testing.T) {
	srv := testServer()
	defer srv.Close()

	// MaxFailures=0 (default) → no limit; normal exploration
	explorer := NewExplorer(srv.URL, 50)
	explorer.MaxFailures = 0
	report, err := explorer.Explore(context.Background(), testSpec())
	require.NoError(t, err)
	assert.Greater(t, report.TotalProbes, 0)
}

func TestExplorer_RuleDeduplication(t *testing.T) {
	srv := testServer()
	defer srv.Close()

	op := &spec.Operation{
		Method: "POST",
		Path:   "/pets",
		RequestBody: &spec.RequestBody{
			Required: true,
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type:     "object",
						Required: []string{"name"},
						Properties: map[string]*spec.Schema{
							"name": {Type: "string"},
						},
					},
				},
			},
		},
		Responses: map[string]*spec.Response{
			"201": {Description: "Created"},
			"400": {Description: "Bad Request"},
		},
	}

	explorer := NewExplorer(srv.URL, 100)
	report, err := explorer.Explore(context.Background(), &spec.ParsedSpec{Operations: []*spec.Operation{op}})
	require.NoError(t, err)

	type key struct {
		op  string
		cat RuleCategory
		fp  string
	}
	seen := make(map[key]int)
	for _, r := range report.Rules {
		k := key{op: r.Operation, cat: r.Category, fp: r.FieldPath}
		seen[k]++
	}
	for k, count := range seen {
		assert.Equal(t, 1, count,
			"duplicate rule for (%s, %s, %s): got %d copies", k.op, k.cat, k.fp, count)
	}
}
