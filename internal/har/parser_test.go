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
