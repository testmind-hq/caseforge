package dea

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunProbe_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write([]byte(`{"id":1}`))
	}))
	defer srv.Close()

	probe := Probe{
		Method:  "POST",
		Path:    "/pets",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    map[string]any{"name": "Fido"},
	}

	ev, err := RunProbe(context.Background(), srv.URL, probe)
	require.NoError(t, err)
	assert.Equal(t, 201, ev.ActualStatus)
	assert.Contains(t, ev.ActualBody, `"id"`)
	assert.Greater(t, ev.Duration.Nanoseconds(), int64(0))
}

func TestRunProbe_400Response(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte(`{"error":"name is required"}`))
	}))
	defer srv.Close()

	probe := Probe{Method: "POST", Path: "/pets", Body: map[string]any{}}
	ev, err := RunProbe(context.Background(), srv.URL, probe)
	require.NoError(t, err)
	assert.Equal(t, 400, ev.ActualStatus)
	assert.Contains(t, ev.ActualBody, "name is required")
}

func TestRunProbe_SendsBodyAsJSON(t *testing.T) {
	var received map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(200)
	}))
	defer srv.Close()

	probe := Probe{
		Method:  "POST",
		Path:    "/pets",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    map[string]any{"name": "Fluffy", "age": float64(3)},
	}
	_, err := RunProbe(context.Background(), srv.URL, probe)
	require.NoError(t, err)
	assert.Equal(t, "Fluffy", received["name"])
}

func TestRunProbe_SendsNilBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer srv.Close()

	probe := Probe{Method: "POST", Path: "/pets", Body: nil}
	ev, err := RunProbe(context.Background(), srv.URL, probe)
	require.NoError(t, err)
	assert.Equal(t, 400, ev.ActualStatus)
}

func TestRunProbe_ConnectError(t *testing.T) {
	probe := Probe{Method: "POST", Path: "/pets"}
	_, err := RunProbe(context.Background(), "http://localhost:1", probe)
	assert.Error(t, err, "must return error when server is unreachable")
}

func TestRunProbe_CapturesResponseHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-ID", "abc-123")
		w.WriteHeader(200)
	}))
	defer srv.Close()

	probe := Probe{Method: "GET", Path: "/pets"}
	ev, err := RunProbe(context.Background(), srv.URL, probe)
	require.NoError(t, err)
	assert.Equal(t, "abc-123", ev.ActualHeaders["X-Request-Id"])
}

func TestRunProbe_SendsQueryParams(t *testing.T) {
	var receivedLimit string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedLimit = r.URL.Query().Get("limit")
		w.WriteHeader(200)
	}))
	defer srv.Close()

	probe := Probe{
		Method:      "GET",
		Path:        "/pets",
		QueryParams: map[string]string{"limit": "5"},
	}
	_, err := RunProbe(context.Background(), srv.URL, probe)
	require.NoError(t, err)
	assert.Equal(t, "5", receivedLimit)
}
