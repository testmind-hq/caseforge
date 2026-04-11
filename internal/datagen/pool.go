// internal/datagen/pool.go
package datagen

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// DataPool holds field values observed from live API responses or seeded from
// external sources (Postman collections, explore reports). Keys are lowercase
// field names; values are ordered slices of observed values.
type DataPool struct {
	entries map[string][]any
}

// NewDataPool creates an empty DataPool.
func NewDataPool() *DataPool {
	return &DataPool{entries: make(map[string][]any)}
}

// Add records a value for the given field name. Keys are case-insensitive.
func (p *DataPool) Add(field string, value any) {
	key := strings.ToLower(field)
	p.entries[key] = append(p.entries[key], value)
}

// ValueFor returns the first observed value for a field name, or (nil, false).
// Lookup is case-insensitive.
func (p *DataPool) ValueFor(field string) (any, bool) {
	vals := p.entries[strings.ToLower(field)]
	if len(vals) == 0 {
		return nil, false
	}
	return vals[0], true
}

// Len returns the number of distinct field names in the pool.
func (p *DataPool) Len() int { return len(p.entries) }

// Merge copies all entries from other into p (in-place).
func (p *DataPool) Merge(other *DataPool) {
	for k, vs := range other.entries {
		p.entries[k] = append(p.entries[k], vs...)
	}
}

// Save writes the pool to a JSON file as {"fieldName": [values...]}.
func (p *DataPool) Save(path string) error {
	data, err := json.MarshalIndent(p.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal pool: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadDataPool reads a pool from a JSON file previously written by Save.
func LoadDataPool(path string) (*DataPool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read pool file: %w", err)
	}
	var raw map[string][]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse pool file: %w", err)
	}
	p := NewDataPool()
	for k, vs := range raw {
		for _, v := range vs {
			p.Add(k, v)
		}
	}
	return p, nil
}
