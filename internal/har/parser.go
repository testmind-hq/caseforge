// internal/har/parser.go
package har

import (
	"encoding/json"
	"net/url"
	"strings"
)

// Entry represents one request-response pair in a HAR log.
type Entry struct {
	Request  Request
	Response Response
}

// Request holds the relevant fields of a HAR request entry.
type Request struct {
	Method   string
	URL      string      // stripped to /path?query by Parse
	Headers  []NameValue
	MIMEType string
	Body     string // raw body text (from postData.text)
}

// Response holds the relevant fields of a HAR response entry.
type Response struct {
	Status   int
	MIMEType string
	Body     string // raw body text (from content.text)
}

// NameValue is a generic name-value pair used for headers.
type NameValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// noiseHeaders is the set of header names (lowercase) to strip.
// Includes browser fingerprinting headers and credential headers (to prevent
// accidental secret leakage into generated output files).
var noiseHeaders = map[string]bool{
	"user-agent":                true,
	"accept-encoding":           true,
	"accept-language":           true,
	"cache-control":             true,
	"pragma":                    true,
	"connection":                true,
	"upgrade-insecure-requests": true,
	// Credential headers — strip to prevent secret leakage in generated files
	"authorization": true,
	"cookie":        true,
	"set-cookie":    true,
	"x-api-key":     true,
}

// isNoise reports whether the header with the given lowercase name should be stripped.
func isNoise(lower string) bool {
	if noiseHeaders[lower] {
		return true
	}
	return strings.HasPrefix(lower, "sec-")
}

// harFile is the internal structure used to unmarshal HAR JSON.
type harFile struct {
	Log struct {
		Entries []struct {
			Request struct {
				Method   string      `json:"method"`
				URL      string      `json:"url"`
				Headers  []NameValue `json:"headers"`
				PostData *struct {
					MIMEType string `json:"mimeType"`
					Text     string `json:"text"`
				} `json:"postData"`
			} `json:"request"`
			Response struct {
				Status  int `json:"status"`
				Content struct {
					MIMEType string `json:"mimeType"`
					Text     string `json:"text"`
				} `json:"content"`
			} `json:"response"`
		} `json:"entries"`
	} `json:"log"`
}

// Parse reads HAR JSON bytes and returns the list of entries.
// Noise headers are stripped from each request.
func Parse(data []byte) ([]Entry, error) {
	var hf harFile
	if err := json.Unmarshal(data, &hf); err != nil {
		return nil, err
	}

	if len(hf.Log.Entries) == 0 {
		return nil, nil
	}

	entries := make([]Entry, 0, len(hf.Log.Entries))
	for _, raw := range hf.Log.Entries {
		// Strip noise headers
		var headers []NameValue
		for _, h := range raw.Request.Headers {
			if !isNoise(strings.ToLower(h.Name)) {
				headers = append(headers, h)
			}
		}

		req := Request{
			Method:  raw.Request.Method,
			URL:     StripBaseURL(raw.Request.URL),
			Headers: headers,
		}
		if raw.Request.PostData != nil {
			req.MIMEType = raw.Request.PostData.MIMEType
			req.Body = raw.Request.PostData.Text
		}

		resp := Response{
			Status:   raw.Response.Status,
			MIMEType: raw.Response.Content.MIMEType,
			Body:     raw.Response.Content.Text,
		}

		entries = append(entries, Entry{Request: req, Response: resp})
	}

	return entries, nil
}

// StripBaseURL removes the scheme and host from rawURL, returning just path+query.
// Returns rawURL unchanged if parsing fails.
func StripBaseURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	result := u.Path
	if u.RawQuery != "" {
		result += "?" + u.RawQuery
	}
	return result
}
