// internal/mutation/operator_status_test.go
package mutation_test

import (
	"net/http"
	"testing"

	"github.com/testmind-hq/caseforge/internal/mutation"
)

func makeResp(statusCode int, headers map[string]string) *http.Response {
	h := http.Header{}
	for k, v := range headers {
		h.Set(k, v)
	}
	return &http.Response{StatusCode: statusCode, Header: h}
}

func TestStatusSwap2xx_200to201(t *testing.T) {
	op := mutation.NewStatusSwap2xxOperator()
	r := makeResp(200, nil)
	_, err := op.Apply(r, []byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	if r.StatusCode != 201 {
		t.Fatalf("expected 201, got %d", r.StatusCode)
	}
}

func TestStatusSwap2xx_201to200(t *testing.T) {
	op := mutation.NewStatusSwap2xxOperator()
	r := makeResp(201, nil)
	_, err := op.Apply(r, []byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	if r.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", r.StatusCode)
	}
}

func TestStatusSwap2xx_NoOp(t *testing.T) {
	op := mutation.NewStatusSwap2xxOperator()
	r := makeResp(404, nil)
	_, err := op.Apply(r, []byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	if r.StatusCode != 404 {
		t.Fatal("non-2xx status must not be changed")
	}
}

func TestErrorInflation(t *testing.T) {
	op := mutation.NewErrorInflationOperator()
	r := makeResp(400, nil)
	body := []byte(`{"error":"bad request"}`)
	out, err := op.Apply(r, body)
	if err != nil {
		t.Fatal(err)
	}
	if r.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", r.StatusCode)
	}
	if string(out) != string(body) {
		t.Fatal("body must pass through unchanged")
	}
}

func TestErrorInflation_NoOp(t *testing.T) {
	op := mutation.NewErrorInflationOperator()
	r := makeResp(200, nil)
	_, err := op.Apply(r, []byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	if r.StatusCode != 200 {
		t.Fatal("2xx status must not be inflated")
	}
}

func TestContentTypeSwap(t *testing.T) {
	op := mutation.NewContentTypeSwapOperator()
	r := makeResp(200, map[string]string{"Content-Type": "application/json"})
	_, err := op.Apply(r, []byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		t.Fatalf("expected text/plain, got %s", r.Header.Get("Content-Type"))
	}
}

func TestHeaderDrop(t *testing.T) {
	op := mutation.NewHeaderDropOperator()
	r := makeResp(200, map[string]string{
		"X-Request-Id": "abc",
		"X-Custom":     "val",
	})
	_, err := op.Apply(r, []byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	// Exactly one header should have been dropped
	remaining := 0
	for k := range r.Header {
		if !skipHeaders[k] {
			remaining++
		}
	}
	if remaining != 1 {
		t.Fatalf("expected 1 non-skip header remaining, got %d: %v", remaining, r.Header)
	}
}

var skipHeaders = map[string]bool{
	"Content-Length":    true,
	"Transfer-Encoding": true,
	"Connection":        true,
}
