// internal/mutation/operator_body.go
package mutation

import (
	"encoding/json"
	"net/http"
	"time"
)

// --- FieldDropOperator ---

type fieldDropOperator struct{}

func NewFieldDropOperator() Operator { return &fieldDropOperator{} }
func (o *fieldDropOperator) Name() string { return "field_drop" }

func (o *fieldDropOperator) Apply(_ *http.Response, body []byte) ([]byte, error) {
	m, keys, err := parseObjectBody(body)
	if err != nil || len(keys) == 0 {
		return body, nil
	}
	delete(m, keys[0])
	return json.Marshal(m)
}

// --- FieldTypeSwapOperator ---

type fieldTypeSwapOperator struct{}

func NewFieldTypeSwapOperator() Operator { return &fieldTypeSwapOperator{} }
func (o *fieldTypeSwapOperator) Name() string { return "field_type_swap" }

func (o *fieldTypeSwapOperator) Apply(_ *http.Response, body []byte) ([]byte, error) {
	m, keys, err := parseObjectBody(body)
	if err != nil {
		return body, nil
	}
	for _, k := range keys {
		switch v := m[k].(type) {
		case string:
			m[k] = float64(len(v))
			return json.Marshal(m)
		case float64:
			m[k] = "mutated"
			return json.Marshal(m)
		}
	}
	return body, nil
}

// --- ArrayToNullOperator ---

type arrayToNullOperator struct{}

func NewArrayToNullOperator() Operator { return &arrayToNullOperator{} }
func (o *arrayToNullOperator) Name() string { return "array_to_null" }

func (o *arrayToNullOperator) Apply(_ *http.Response, body []byte) ([]byte, error) {
	m, keys, err := parseObjectBody(body)
	if err != nil {
		return body, nil
	}
	for _, k := range keys {
		if _, ok := m[k].([]any); ok {
			m[k] = nil
			return json.Marshal(m)
		}
	}
	return body, nil
}

// --- NullToArrayOperator ---

type nullToArrayOperator struct{}

func NewNullToArrayOperator() Operator { return &nullToArrayOperator{} }
func (o *nullToArrayOperator) Name() string { return "null_to_array" }

func (o *nullToArrayOperator) Apply(_ *http.Response, body []byte) ([]byte, error) {
	m, keys, err := parseObjectBody(body)
	if err != nil {
		return body, nil
	}
	for _, k := range keys {
		if m[k] == nil {
			m[k] = []any{}
			return json.Marshal(m)
		}
	}
	return body, nil
}

// --- PaginationOffByOneOperator ---

var paginationFields = []string{"total", "count", "total_count", "totalCount", "size"}

type paginationOffByOneOperator struct{}

func NewPaginationOffByOneOperator() Operator { return &paginationOffByOneOperator{} }
func (o *paginationOffByOneOperator) Name() string { return "pagination_off_by_one" }

func (o *paginationOffByOneOperator) Apply(_ *http.Response, body []byte) ([]byte, error) {
	m, _, err := parseObjectBody(body)
	if err != nil {
		return body, nil
	}
	for _, field := range paginationFields {
		if v, ok := m[field].(float64); ok {
			m[field] = v - 1
			return json.Marshal(m)
		}
	}
	return body, nil
}

// --- EmptyResultInjectionOperator ---

type emptyResultInjectionOperator struct{}

func NewEmptyResultInjectionOperator() Operator { return &emptyResultInjectionOperator{} }
func (o *emptyResultInjectionOperator) Name() string { return "empty_result_injection" }

func (o *emptyResultInjectionOperator) Apply(_ *http.Response, body []byte) ([]byte, error) {
	var arr []any
	if json.Unmarshal(body, &arr) == nil {
		return []byte("[]"), nil
	}
	m, keys, err := parseObjectBody(body)
	if err != nil {
		return body, nil
	}
	for _, k := range keys {
		if _, ok := m[k].([]any); ok {
			m[k] = []any{}
			return json.Marshal(m)
		}
	}
	return body, nil
}

// --- DateFormatSwapOperator ---

type dateFormatSwapOperator struct{}

func NewDateFormatSwapOperator() Operator { return &dateFormatSwapOperator{} }
func (o *dateFormatSwapOperator) Name() string { return "date_format_swap" }

func (o *dateFormatSwapOperator) Apply(_ *http.Response, body []byte) ([]byte, error) {
	m, keys, err := parseObjectBody(body)
	if err != nil {
		return body, nil
	}
	for _, k := range keys {
		s, ok := m[k].(string)
		if !ok {
			continue
		}
		t, parseErr := time.Parse(time.RFC3339, s)
		if parseErr != nil {
			continue
		}
		m[k] = t.Unix()
		return json.Marshal(m)
	}
	return body, nil
}

// --- NumericPrecisionLossOperator ---

type numericPrecisionLossOperator struct{}

func NewNumericPrecisionLossOperator() Operator { return &numericPrecisionLossOperator{} }
func (o *numericPrecisionLossOperator) Name() string { return "numeric_precision_loss" }

func (o *numericPrecisionLossOperator) Apply(_ *http.Response, body []byte) ([]byte, error) {
	m, keys, err := parseObjectBody(body)
	if err != nil {
		return body, nil
	}
	mutated := false
	for _, k := range keys {
		v, ok := m[k].(float64)
		if !ok {
			continue
		}
		truncated := float64(int64(v))
		if truncated != v {
			m[k] = truncated
			mutated = true
		}
	}
	if !mutated {
		return body, nil
	}
	return json.Marshal(m)
}

// --- shared helper ---

// parseObjectBody decodes body as a JSON object and returns the map and a stable key slice.
func parseObjectBody(body []byte) (map[string]any, []string, error) {
	var m map[string]any
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, nil, err
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return m, keys, nil
}
