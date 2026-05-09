// internal/sandbox/state.go
package sandbox

import "sync"

// StateStore tracks in-memory resource objects for stateful CRUD flows.
type StateStore interface {
	Write(resourceType, id string, body map[string]any)
	Read(resourceType, id string) (map[string]any, bool)
	Delete(resourceType, id string)
	List(resourceType string) []map[string]any
}

type memStateStore struct {
	mu    sync.RWMutex
	store map[string]map[string]map[string]any // resourceType → id → object
}

func newMemStateStore() *memStateStore {
	return &memStateStore{store: make(map[string]map[string]map[string]any)}
}

func (m *memStateStore) Write(resourceType, id string, body map[string]any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.store[resourceType] == nil {
		m.store[resourceType] = make(map[string]map[string]any)
	}
	m.store[resourceType][id] = body
}

func (m *memStateStore) Read(resourceType, id string) (map[string]any, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if rt, ok := m.store[resourceType]; ok {
		if obj, ok := rt[id]; ok {
			return obj, true
		}
	}
	return nil, false
}

func (m *memStateStore) Delete(resourceType, id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if rt, ok := m.store[resourceType]; ok {
		delete(rt, id)
	}
}

func (m *memStateStore) List(resourceType string) []map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()
	rt, ok := m.store[resourceType]
	if !ok {
		return []map[string]any{}
	}
	out := make([]map[string]any, 0, len(rt))
	for _, v := range rt {
		out = append(out, v)
	}
	return out
}
