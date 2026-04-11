// internal/datagen/postman_seed.go
package datagen

import (
	"encoding/json"
	"fmt"
	"os"
)

type postmanCollection struct {
	Item []postmanItem `json:"item"`
}

type postmanItem struct {
	Name    string          `json:"name"`
	Item    []postmanItem   `json:"item"`    // folder — has nested items, no request
	Request *postmanRequest `json:"request"` // leaf request
}

type postmanRequest struct {
	Method string       `json:"method"`
	Body   *postmanBody `json:"body"`
}

type postmanBody struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}

// ParsePostmanCollection reads a Postman Collection v2.1 file and returns a DataPool
// populated with field values extracted from all JSON (mode=raw) request bodies.
// Non-JSON bodies and folders are handled gracefully.
func ParsePostmanCollection(path string) (*DataPool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read postman collection: %w", err)
	}
	var col postmanCollection
	if err := json.Unmarshal(data, &col); err != nil {
		return nil, fmt.Errorf("parse postman collection: %w", err)
	}
	pool := NewDataPool()
	extractPostmanItems(col.Item, pool)
	return pool, nil
}

func extractPostmanItems(items []postmanItem, pool *DataPool) {
	for _, item := range items {
		// Recurse into folders (items with sub-items and no request)
		if len(item.Item) > 0 {
			extractPostmanItems(item.Item, pool)
		}
		if item.Request == nil || item.Request.Body == nil {
			continue
		}
		if item.Request.Body.Mode != "raw" || item.Request.Body.Raw == "" {
			continue
		}
		var body map[string]any
		if err := json.Unmarshal([]byte(item.Request.Body.Raw), &body); err != nil {
			continue // non-JSON raw body
		}
		for k, v := range body {
			pool.Add(k, v)
		}
	}
}
