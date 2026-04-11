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

	var implicitRules []DiscoveredRule
	for _, r := range report.Rules {
		if r.Implicit {
			implicitRules = append(implicitRules, r)
		}
	}
	assert.NotEmpty(t, implicitRules, "server rejects empty name → must produce implicit min rule")
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
