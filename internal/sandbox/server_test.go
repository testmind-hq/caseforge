// internal/sandbox/server_test.go
package sandbox

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// petstoreSpec returns a minimal ParsedSpec modelling GET/POST /pets and GET/DELETE /pets/{petId}.
func petstoreSpec() *spec.ParsedSpec {
	idSchema := &spec.Schema{Type: "string", Format: "uuid"}
	nameSchema := &spec.Schema{Type: "string"}
	petSchema := &spec.Schema{
		Type:       "object",
		Properties: map[string]*spec.Schema{"id": idSchema, "name": nameSchema},
	}
	jsonContent := func(s *spec.Schema) map[string]*spec.MediaType {
		return map[string]*spec.MediaType{"application/json": {Schema: s}}
	}
	return &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{
				OperationID: "listPets",
				Method:      "GET",
				Path:        "/pets",
				Responses: map[string]*spec.Response{
					"200": {Content: jsonContent(&spec.Schema{Type: "array", Items: petSchema})},
				},
			},
			{
				OperationID: "createPet",
				Method:      "POST",
				Path:        "/pets",
				Responses: map[string]*spec.Response{
					"201": {Content: jsonContent(petSchema)},
				},
			},
			{
				OperationID: "showPetById",
				Method:      "GET",
				Path:        "/pets/{petId}",
				Parameters: []*spec.Parameter{
					{Name: "petId", In: "path", Required: true, Schema: &spec.Schema{Type: "string"}},
				},
				Responses: map[string]*spec.Response{
					"200": {Content: jsonContent(petSchema)},
					"404": {Description: "Not found"},
				},
			},
			{
				OperationID: "deletePet",
				Method:      "DELETE",
				Path:        "/pets/{petId}",
				Parameters: []*spec.Parameter{
					{Name: "petId", In: "path", Required: true, Schema: &spec.Schema{Type: "string"}},
				},
				Responses: map[string]*spec.Response{
					"204": {Description: "Deleted"},
					"404": {Description: "Not found"},
				},
			},
		},
	}
}

// newTestMux builds an http.ServeMux from a ParsedSpec exactly as NewSandboxServer does,
// but returns it directly so httptest.NewServer can wrap it.
func newTestMux(ps *spec.ParsedSpec) http.Handler {
	logger, _ := buildLogger(Options{LogLevel: "silent"})
	gen := newResponseGenerator("auto")
	store := newMemStateStore()
	mux := http.NewServeMux()
	for _, op := range ps.Operations {
		pattern := strings.ToUpper(op.Method) + " " + op.Path
		mux.HandleFunc(pattern, makeHandler(op, store, gen, logger))
	}
	return mux
}

func TestSandbox_ListPets_Empty(t *testing.T) {
	ts := httptest.NewServer(newTestMux(petstoreSpec()))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/pets")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var list []any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&list))
	assert.NotNil(t, list)
}

func TestSandbox_CRUD_Flow(t *testing.T) {
	ts := httptest.NewServer(newTestMux(petstoreSpec()))
	defer ts.Close()

	// POST → 201
	body := bytes.NewBufferString(`{"name":"Fido"}`)
	resp, err := http.Post(ts.URL+"/pets", "application/json", body)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	data, _ := io.ReadAll(resp.Body)
	var created map[string]any
	require.NoError(t, json.Unmarshal(data, &created))
	id, ok := created["id"].(string)
	require.True(t, ok, "POST response must include id field")
	require.NotEmpty(t, id)

	// GET /{id} → 200
	resp2, err := http.Get(ts.URL + "/pets/" + id)
	require.NoError(t, err)
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	// DELETE /{id} → 204
	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/pets/"+id, nil)
	resp3, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp3.Body.Close()
	assert.Equal(t, http.StatusNoContent, resp3.StatusCode)

	// GET /{id} after delete → 404
	resp4, err := http.Get(ts.URL + "/pets/" + id)
	require.NoError(t, err)
	defer resp4.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp4.StatusCode)
}

func TestSandbox_WrongContentType_Returns415(t *testing.T) {
	ts := httptest.NewServer(newTestMux(petstoreSpec()))
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/pets", "text/plain", strings.NewReader("name=Fido"))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
}

func TestSandbox_Start_Addr(t *testing.T) {
	ps := petstoreSpec()
	srv := NewSandboxServer(ps, Options{LogLevel: "silent"})
	err := srv.Start("127.0.0.1", 0)
	require.NoError(t, err)
	assert.NotEmpty(t, srv.Addr())

	// Shutdown to release the port
	_ = srv.Shutdown(t.Context())
}
