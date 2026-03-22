// internal/spec/loader_test.go
package spec

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalLoaderLoadsYAML(t *testing.T) {
	loader := NewLocalLoader()
	spec, err := loader.Load("../../testdata/petstore.yaml")
	require.NoError(t, err)
	assert.NotEmpty(t, spec.Title)
	assert.NotEmpty(t, spec.Operations)
}

func TestLocalLoaderMissingFile(t *testing.T) {
	loader := NewLocalLoader()
	_, err := loader.Load("/nonexistent/path.yaml")
	assert.Error(t, err)
}

func TestURLLoaderFetchesRemote(t *testing.T) {
	// Serve a local YAML file over HTTP for testing
	content, err := os.ReadFile("../../testdata/petstore.yaml")
	require.NoError(t, err)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		w.Write(content)
	}))
	defer srv.Close()

	loader := NewURLLoader()
	spec, err := loader.Load(srv.URL + "/openapi.yaml")
	require.NoError(t, err)
	assert.NotEmpty(t, spec.Operations)
}

func TestNewLoaderDispatchesByPrefix(t *testing.T) {
	l := NewLoader()
	assert.IsType(t, &URLLoader{}, l.loaderFor("https://example.com/api.yaml"))
	assert.IsType(t, &LocalLoader{}, l.loaderFor("/local/path.yaml"))
}
