// internal/dea/response_check_test.go
package dea

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestFindResponseSchema_ExactMatch(t *testing.T) {
	op := makeSchemaOp()
	s := findResponseSchema(op, 201)
	require.NotNil(t, s)
	assert.Equal(t, "object", s.Type)
}

func TestFindResponseSchema_FallbackTo2xx(t *testing.T) {
	op := &spec.Operation{
		Method: "GET", Path: "/items",
		Responses: map[string]*spec.Response{
			"200": {Content: map[string]*spec.MediaType{
				"application/json": {Schema: &spec.Schema{Type: "object"}},
			}},
		},
	}
	// 201 not declared — should fall back to 200
	s := findResponseSchema(op, 201)
	require.NotNil(t, s)
}

func TestFindResponseSchema_NoSchema_ReturnsNil(t *testing.T) {
	op := &spec.Operation{
		Method:    "DELETE", Path: "/items/1",
		Responses: map[string]*spec.Response{"204": {Description: "no content"}},
	}
	assert.Nil(t, findResponseSchema(op, 204))
}

func TestCheckResponseBody_MissingRequired(t *testing.T) {
	s := makeSchemaOp().Responses["201"].Content["application/json"].Schema
	violations := checkResponseBody(`{"id": 1}`, s) // missing "name"
	assert.Len(t, violations, 1)
	assert.Contains(t, violations[0], "name")
}

func TestCheckResponseBody_WrongType(t *testing.T) {
	s := &spec.Schema{
		Type:       "object",
		Properties: map[string]*spec.Schema{"count": {Type: "integer"}},
	}
	violations := checkResponseBody(`{"count": "not-a-number"}`, s)
	assert.Len(t, violations, 1)
	assert.Contains(t, violations[0], "count")
}

func TestCheckResponseBody_Valid_NoViolations(t *testing.T) {
	s := makeSchemaOp().Responses["201"].Content["application/json"].Schema
	violations := checkResponseBody(`{"id": 1, "name": "rex"}`, s)
	assert.Empty(t, violations)
}

func TestCheckResponseBody_NonJSON_ReturnsNil(t *testing.T) {
	violations := checkResponseBody("not json", &spec.Schema{Type: "object"})
	assert.Nil(t, violations)
}

func TestCheckResponseBody_EmptyBody_ReturnsNil(t *testing.T) {
	violations := checkResponseBody("", &spec.Schema{Type: "object"})
	assert.Nil(t, violations)
}

func TestValidateProbeResponse_ReturnsMismatchRule(t *testing.T) {
	op := makeSchemaOp()
	probe := Probe{Method: "POST", Path: "/pets", ExpectedStatus: 201}
	ev := &Evidence{ActualStatus: 201, ActualBody: `{"id": 1}`} // missing "name"
	rule := validateProbeResponse(op, probe, ev)
	require.NotNil(t, rule)
	assert.Equal(t, CategoryResponseSchemaMismatch, rule.Category)
	assert.Contains(t, rule.Description, "name")
}

func TestValidateProbeResponse_ValidBody_ReturnsNil(t *testing.T) {
	op := makeSchemaOp()
	probe := Probe{Method: "POST", Path: "/pets", ExpectedStatus: 201}
	ev := &Evidence{ActualStatus: 201, ActualBody: `{"id": 1, "name": "rex"}`}
	assert.Nil(t, validateProbeResponse(op, probe, ev))
}

func TestValidateProbeResponse_NoSchema_ReturnsNil(t *testing.T) {
	op := &spec.Operation{
		Method:    "DELETE",
		Path:      "/pets/1",
		Responses: map[string]*spec.Response{"204": {}},
	}
	probe := Probe{Method: "DELETE", Path: "/pets/1", ExpectedStatus: 204}
	ev := &Evidence{ActualStatus: 204, ActualBody: ""}
	assert.Nil(t, validateProbeResponse(op, probe, ev))
}

func TestExplorer_ResponseSchemaMismatch_ProducesRule(t *testing.T) {
	// Server returns 201 but omits the required "name" field from the response body
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/pets" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(201)
			// Missing "name" field — should trigger response schema mismatch rule
			_, _ = w.Write([]byte(`{"id": 1}`))
			return
		}
		w.WriteHeader(404)
	}))
	defer srv.Close()

	explorer := NewExplorer(srv.URL, 50)
	parsedSpec := &spec.ParsedSpec{Operations: []*spec.Operation{makeSchemaOp()}}
	report, err := explorer.Explore(context.Background(), parsedSpec)
	require.NoError(t, err)

	var mismatchRules []DiscoveredRule
	for _, r := range report.Rules {
		if r.Category == CategoryResponseSchemaMismatch {
			mismatchRules = append(mismatchRules, r)
		}
	}
	assert.NotEmpty(t, mismatchRules, "expected at least one response schema mismatch rule")
}

// makeSchemaOp returns a POST /pets operation with a declared 201 response schema.
func makeSchemaOp() *spec.Operation {
	minL := int64(1)
	return &spec.Operation{
		Method: "POST", Path: "/pets",
		RequestBody: &spec.RequestBody{
			Required: true,
			Content: map[string]*spec.MediaType{
				"application/json": {Schema: &spec.Schema{
					Type:     "object",
					Required: []string{"name"},
					Properties: map[string]*spec.Schema{
						"name": {Type: "string", MinLength: &minL},
					},
				}},
			},
		},
		Responses: map[string]*spec.Response{
			"201": {
				Description: "Created",
				Content: map[string]*spec.MediaType{
					"application/json": {Schema: &spec.Schema{
						Type:     "object",
						Required: []string{"id", "name"},
						Properties: map[string]*spec.Schema{
							"id":   {Type: "integer"},
							"name": {Type: "string", MinLength: &minL},
						},
					}},
				},
			},
		},
	}
}
