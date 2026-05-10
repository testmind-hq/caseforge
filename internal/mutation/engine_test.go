// internal/mutation/engine_test.go
package mutation_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/testmind-hq/caseforge/internal/mutation"
)

func buildTestHurlFile(t *testing.T, dir, caseID string, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, caseID+".hurl"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func buildTestIndexJSON(t *testing.T, dir, caseID, title string) {
	t.Helper()
	data, _ := json.Marshal(map[string]any{
		"test_cases": []map[string]any{{"id": caseID, "title": title}},
	})
	if err := os.WriteFile(filepath.Join(dir, "index.json"), data, 0644); err != nil {
		t.Fatal(err)
	}
}

func TestEngineRun_KilledMutation(t *testing.T) {
	if _, err := exec.LookPath("hurl"); err != nil {
		t.Skip("hurl not installed")
	}

	// Backend always returns 200
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"id":1}`))
	}))
	defer backend.Close()

	dir := t.TempDir()
	// Hurl file asserts HTTP 200 — status_swap_2xx changes it to 201 → hurl FAILS → mutation KILLED
	buildTestHurlFile(t, dir, "TC-0001", "GET {{base_url}}/\nHTTP 200\n")
	buildTestIndexJSON(t, dir, "TC-0001", "root GET")

	opts := mutation.RunOptions{
		Target:      backend.URL,
		CasesDir:    dir,
		Operators:   []mutation.Operator{mutation.NewStatusSwap2xxOperator()},
		Concurrency: 1,
	}
	run, err := mutation.Run(opts)
	if err != nil {
		t.Fatal(err)
	}
	if run.Survivors != 0 {
		t.Fatalf("status_swap_2xx must be killed by HTTP 200 assertion, got %d survivors", run.Survivors)
	}
	if run.Killed != 1 {
		t.Fatalf("expected 1 killed, got %d", run.Killed)
	}
}

func TestEngineRun_SurvivorMutation(t *testing.T) {
	if _, err := exec.LookPath("hurl"); err != nil {
		t.Skip("hurl not installed")
	}

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":1,"name":"test"}`))
	}))
	defer backend.Close()

	dir := t.TempDir()
	// Hurl only checks status == 200, not body fields → field_drop survives
	buildTestHurlFile(t, dir, "TC-0002", "GET {{base_url}}/\nHTTP 200\n")
	buildTestIndexJSON(t, dir, "TC-0002", "root GET body")

	opts := mutation.RunOptions{
		Target:      backend.URL,
		CasesDir:    dir,
		Operators:   []mutation.Operator{mutation.NewFieldDropOperator()},
		Concurrency: 1,
	}
	run, err := mutation.Run(opts)
	if err != nil {
		t.Fatal(err)
	}
	if run.Survivors != 1 {
		t.Fatalf("field_drop must survive when no body assertion exists, got %d survivors", run.Survivors)
	}
}

func TestEngineRun_NoHurlFiles(t *testing.T) {
	dir := t.TempDir()
	// Only index.json, no .hurl files — the case should be counted as killed (not error)
	buildTestIndexJSON(t, dir, "TC-MISSING", "no hurl file")

	opts := mutation.RunOptions{
		Target:      "http://127.0.0.1:1", // unreachable, but won't be called
		CasesDir:    dir,
		Operators:   []mutation.Operator{mutation.NewFieldDropOperator()},
		Concurrency: 1,
	}
	run, err := mutation.Run(opts)
	if err != nil {
		t.Fatal(err)
	}
	// Missing hurl file → runOnce returns false (not survived) → killed
	if run.TotalRuns != 1 {
		t.Fatalf("expected 1 total run, got %d", run.TotalRuns)
	}
}
