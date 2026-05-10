// internal/mutation/operator_status.go
package mutation

import (
	"net/http"
	"sort"
	"strings"
)

// --- StatusSwap2xxOperator ---

type statusSwap2xxOperator struct{}

func NewStatusSwap2xxOperator() Operator { return &statusSwap2xxOperator{} }
func (o *statusSwap2xxOperator) Name() string { return "status_swap_2xx" }

func (o *statusSwap2xxOperator) Apply(resp *http.Response, body []byte) ([]byte, error) {
	switch resp.StatusCode {
	case 200:
		resp.StatusCode = 201
	case 201:
		resp.StatusCode = 200
	case 204:
		resp.StatusCode = 200
	}
	return body, nil
}

// --- ErrorInflationOperator ---

type errorInflationOperator struct{}

func NewErrorInflationOperator() Operator { return &errorInflationOperator{} }
func (o *errorInflationOperator) Name() string { return "error_inflation" }

func (o *errorInflationOperator) Apply(resp *http.Response, body []byte) ([]byte, error) {
	if resp.StatusCode >= 400 {
		resp.StatusCode = 200
	}
	return body, nil
}

// --- ContentTypeSwapOperator ---

type contentTypeSwapOperator struct{}

func NewContentTypeSwapOperator() Operator { return &contentTypeSwapOperator{} }
func (o *contentTypeSwapOperator) Name() string { return "content_type_swap" }

func (o *contentTypeSwapOperator) Apply(resp *http.Response, body []byte) ([]byte, error) {
	ct := strings.ToLower(resp.Header.Get("Content-Type"))
	if strings.HasPrefix(ct, "application/json") {
		resp.Header.Set("Content-Type", "text/plain; charset=utf-8")
	}
	return body, nil
}

// --- HeaderDropOperator ---

// skipHeaders are not dropped because removing them would break hurl's response parsing.
var skipHeaders = map[string]bool{
	"Content-Length":    true,
	"Transfer-Encoding": true,
	"Connection":        true,
}

type headerDropOperator struct{}

func NewHeaderDropOperator() Operator { return &headerDropOperator{} }
func (o *headerDropOperator) Name() string { return "header_drop" }

func (o *headerDropOperator) Apply(resp *http.Response, body []byte) ([]byte, error) {
	keys := make([]string, 0, len(resp.Header))
	for k := range resp.Header {
		if !skipHeaders[k] {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	if len(keys) > 0 {
		resp.Header.Del(keys[0])
	}
	return body, nil
}
