// internal/har/parser_test.go
package har

import (
	"testing"
)

var sampleHAR = []byte(`{
  "log": {
    "version": "1.2",
    "entries": [
      {
        "request": {
          "method": "POST",
          "url": "https://api.example.com/users",
          "headers": [
            {"name": "Content-Type", "value": "application/json"},
            {"name": "user-agent", "value": "Mozilla/5.0"},
            {"name": "sec-fetch-mode", "value": "cors"}
          ],
          "postData": {
            "mimeType": "application/json",
            "text": "{\"email\":\"a@b.com\"}"
          }
        },
        "response": {
          "status": 201,
          "content": {
            "mimeType": "application/json",
            "text": "{\"id\":42,\"email\":\"a@b.com\"}"
          }
        }
      }
    ]
  }
}`)

func TestHARParser_ParsesMethod(t *testing.T) {
	entries, err := Parse(sampleHAR)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one entry")
	}
	if entries[0].Request.Method != "POST" {
		t.Errorf("expected method POST, got %q", entries[0].Request.Method)
	}
}

func TestHARParser_ParsesURL_StripsBaseURL(t *testing.T) {
	entries, err := Parse(sampleHAR)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if entries[0].Request.URL != "/users" {
		t.Errorf("expected URL /users, got %q", entries[0].Request.URL)
	}
}

func TestHARParser_ParsesBody(t *testing.T) {
	entries, err := Parse(sampleHAR)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if entries[0].Request.Body != `{"email":"a@b.com"}` {
		t.Errorf("unexpected body: %q", entries[0].Request.Body)
	}
}

func TestHARParser_StripsNoiseHeaders(t *testing.T) {
	entries, err := Parse(sampleHAR)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	for _, h := range entries[0].Request.Headers {
		switch h.Name {
		case "user-agent", "sec-fetch-mode":
			t.Errorf("noise header %q should have been stripped", h.Name)
		}
	}
}

func TestHARParser_KeepsContentTypeHeader(t *testing.T) {
	entries, err := Parse(sampleHAR)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	for _, h := range entries[0].Request.Headers {
		if h.Name == "Content-Type" {
			return
		}
	}
	t.Error("Content-Type header should have been kept")
}

func TestHARParser_EmptyEntries(t *testing.T) {
	empty := []byte(`{"log":{"entries":[]}}`)
	entries, err := Parse(empty)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil slice for empty entries, got %v", entries)
	}
}

func TestHARParser_StripsCredentialHeaders(t *testing.T) {
	harData := []byte(`{
  "log": {"entries": [{
    "request": {
      "method": "GET",
      "url": "https://api.example.com/items",
      "headers": [
        {"name": "Authorization", "value": "Bearer secret-token"},
        {"name": "Cookie", "value": "session=abc123"},
        {"name": "X-Api-Key", "value": "key-value"},
        {"name": "Accept", "value": "application/json"}
      ]
    },
    "response": {"status": 200, "content": {"mimeType": "application/json", "text": "{}"}}
  }]}
}`)
	entries, err := Parse(harData)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	for _, h := range entries[0].Request.Headers {
		switch h.Name {
		case "Authorization", "Cookie", "X-Api-Key":
			t.Errorf("credential header %q should have been stripped", h.Name)
		}
	}
	// Accept should be kept
	kept := false
	for _, h := range entries[0].Request.Headers {
		if h.Name == "Accept" {
			kept = true
		}
	}
	if !kept {
		t.Error("Accept header should have been kept")
	}
}

func TestHARParser_StripsNoiseHeadersCaseInsensitive(t *testing.T) {
	harData := []byte(`{
  "log": {"entries": [{
    "request": {
      "method": "GET",
      "url": "https://api.example.com/items",
      "headers": [
        {"name": "User-Agent", "value": "Mozilla/5.0"},
        {"name": "SEC-FETCH-SITE", "value": "cross-site"},
        {"name": "Accept", "value": "application/json"}
      ]
    },
    "response": {"status": 200, "content": {"mimeType": "application/json", "text": "{}"}}
  }]}
}`)
	entries, err := Parse(harData)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	for _, h := range entries[0].Request.Headers {
		switch h.Name {
		case "User-Agent", "SEC-FETCH-SITE":
			t.Errorf("noise header %q should have been stripped (case-insensitive)", h.Name)
		}
	}
}

func TestHARParser_InvalidJSON(t *testing.T) {
	_, err := Parse([]byte(`not-json`))
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestHARParser_URLWithQueryPreservedInURL(t *testing.T) {
	harData := []byte(`{
  "log": {"entries": [{
    "request": {
      "method": "GET",
      "url": "https://api.example.com/users?page=2&limit=10",
      "headers": []
    },
    "response": {"status": 200, "content": {"mimeType": "application/json", "text": "{}"}}
  }]}
}`)
	entries, err := Parse(harData)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if entries[0].Request.URL != "/users?page=2&limit=10" {
		t.Errorf("expected URL with query preserved, got %q", entries[0].Request.URL)
	}
}
