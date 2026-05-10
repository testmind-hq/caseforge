// internal/mutation/operator_body_test.go
package mutation_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/testmind-hq/caseforge/internal/mutation"
)

func parseBody(t *testing.T, b []byte) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("invalid JSON output %q: %v", b, err)
	}
	return m
}

func TestFieldDrop(t *testing.T) {
	op := mutation.NewFieldDropOperator()
	body := []byte(`{"id":1,"name":"Alice","email":"a@b.com"}`)
	out, err := op.Apply(&http.Response{}, body)
	if err != nil {
		t.Fatal(err)
	}
	m := parseBody(t, out)
	if len(m) != 2 {
		t.Fatalf("expected 2 fields after drop, got %d: %v", len(m), m)
	}
}

func TestFieldDrop_NonJSON(t *testing.T) {
	op := mutation.NewFieldDropOperator()
	body := []byte(`not json`)
	out, err := op.Apply(&http.Response{}, body)
	if err != nil {
		t.Fatal("Apply must not error on non-JSON body")
	}
	if string(out) != "not json" {
		t.Fatal("non-JSON body must pass through unchanged")
	}
}

func TestFieldTypeSwap(t *testing.T) {
	op := mutation.NewFieldTypeSwapOperator()
	body := []byte(`{"name":"Alice","count":5}`)
	out, err := op.Apply(&http.Response{}, body)
	if err != nil {
		t.Fatal(err)
	}
	m := parseBody(t, out)
	_, nameIsNum := m["name"].(float64)
	_, countIsStr := m["count"].(string)
	if !nameIsNum && !countIsStr {
		t.Fatalf("expected one field type swap, got: %v", m)
	}
}

func TestArrayToNull(t *testing.T) {
	op := mutation.NewArrayToNullOperator()
	body := []byte(`{"items":[1,2,3],"name":"x"}`)
	out, err := op.Apply(&http.Response{}, body)
	if err != nil {
		t.Fatal(err)
	}
	m := parseBody(t, out)
	if m["items"] != nil {
		t.Fatalf("expected items=null, got %v", m["items"])
	}
}

func TestNullToArray(t *testing.T) {
	op := mutation.NewNullToArrayOperator()
	body := []byte(`{"result":null}`)
	out, err := op.Apply(&http.Response{}, body)
	if err != nil {
		t.Fatal(err)
	}
	m := parseBody(t, out)
	arr, ok := m["result"].([]any)
	if !ok {
		t.Fatalf("expected result=[], got %v (%T)", m["result"], m["result"])
	}
	if len(arr) != 0 {
		t.Fatalf("expected empty array, got %v", arr)
	}
}

func TestPaginationOffByOne(t *testing.T) {
	op := mutation.NewPaginationOffByOneOperator()
	body := []byte(`{"total":10,"items":[]}`)
	out, err := op.Apply(&http.Response{}, body)
	if err != nil {
		t.Fatal(err)
	}
	m := parseBody(t, out)
	total, _ := m["total"].(float64)
	if total != 9 && total != 11 {
		t.Fatalf("expected total 9 or 11, got %v", m["total"])
	}
}

func TestEmptyResultInjection_Array(t *testing.T) {
	op := mutation.NewEmptyResultInjectionOperator()
	body := []byte(`[{"id":1},{"id":2}]`)
	out, err := op.Apply(&http.Response{}, body)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "[]" {
		t.Fatalf("expected [], got %s", out)
	}
}

func TestEmptyResultInjection_Object(t *testing.T) {
	op := mutation.NewEmptyResultInjectionOperator()
	body := []byte(`{"data":[1,2,3],"meta":{}}`)
	out, err := op.Apply(&http.Response{}, body)
	if err != nil {
		t.Fatal(err)
	}
	m := parseBody(t, out)
	arr, ok := m["data"].([]any)
	if !ok || len(arr) != 0 {
		t.Fatalf("expected data=[], got %v", m["data"])
	}
}

func TestDateFormatSwap(t *testing.T) {
	op := mutation.NewDateFormatSwapOperator()
	body := []byte(`{"created_at":"2024-01-15T10:30:00Z"}`)
	out, err := op.Apply(&http.Response{}, body)
	if err != nil {
		t.Fatal(err)
	}
	m := parseBody(t, out)
	if m["created_at"] == "2024-01-15T10:30:00Z" {
		t.Fatal("created_at should have been transformed to unix timestamp")
	}
}

func TestNumericPrecisionLoss(t *testing.T) {
	op := mutation.NewNumericPrecisionLossOperator()
	body := []byte(`{"price":9.99,"count":3}`)
	out, err := op.Apply(&http.Response{}, body)
	if err != nil {
		t.Fatal(err)
	}
	m := parseBody(t, out)
	price, _ := m["price"].(float64)
	if price != 9 {
		t.Fatalf("expected price truncated to 9, got %v", price)
	}
	count, _ := m["count"].(float64)
	if count != 3 {
		t.Fatal("integer field must not be changed by precision loss")
	}
}
