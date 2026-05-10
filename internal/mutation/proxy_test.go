// internal/mutation/proxy_test.go
package mutation_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/testmind-hq/caseforge/internal/mutation"
)

func TestProxyPassthrough(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":1,"name":"test"}`))
	}))
	defer backend.Close()

	proxy, err := mutation.NewProxy(backend.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer proxy.Close()

	resp, err := http.Get("http://" + proxy.Addr())
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if string(body) != `{"id":1,"name":"test"}` {
		t.Fatalf("passthrough: expected original body, got %s", body)
	}
}

func TestProxyAppliesMutation(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":1,"name":"test","email":"a@b.com"}`))
	}))
	defer backend.Close()

	proxy, err := mutation.NewProxy(backend.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer proxy.Close()
	proxy.SetActive(mutation.NewFieldDropOperator())

	resp, err := http.Get("http://" + proxy.Addr())
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var m map[string]any
	if err := json.Unmarshal(body, &m); err != nil {
		t.Fatalf("invalid JSON: %s", body)
	}
	if len(m) != 2 {
		t.Fatalf("expected 2 fields after field_drop, got %d: %v", len(m), m)
	}
}

func TestProxyStatusMutation(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer backend.Close()

	proxy, err := mutation.NewProxy(backend.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer proxy.Close()
	proxy.SetActive(mutation.NewStatusSwap2xxOperator())

	resp, err := http.Get("http://" + proxy.Addr())
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 201 {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}

func TestRegistry(t *testing.T) {
	ops := mutation.Registry()
	if len(ops) != 12 {
		t.Fatalf("expected 12 operators, got %d", len(ops))
	}
	names := map[string]bool{}
	for _, op := range ops {
		if names[op.Name()] {
			t.Fatalf("duplicate operator name: %s", op.Name())
		}
		names[op.Name()] = true
	}
}
