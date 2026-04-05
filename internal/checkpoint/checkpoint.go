// internal/checkpoint/checkpoint.go
// Package checkpoint persists generation progress so interrupted runs can resume.
package checkpoint

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const stateFileName = ".state.json"

// State tracks which operations have been successfully generated.
type State struct {
	SpecHash  string            `json:"spec_hash"`
	StartedAt time.Time         `json:"started_at"`
	Completed map[string]bool   `json:"completed"` // key: "METHOD /path"
}

// Manager reads and writes checkpoint state for a given output directory.
type Manager struct {
	path string // absolute path to .state.json
}

// NewManager returns a Manager for the given output directory.
func NewManager(outputDir string) *Manager {
	return &Manager{path: filepath.Join(outputDir, stateFileName)}
}

// Load reads an existing checkpoint file. Returns nil, nil if none exists.
func (m *Manager) Load() (*State, error) {
	data, err := os.ReadFile(m.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	if s.Completed == nil {
		s.Completed = make(map[string]bool)
	}
	return &s, nil
}

// Save writes the checkpoint to disk. Creates the output directory if needed.
func (m *Manager) Save(s *State) error {
	if err := os.MkdirAll(filepath.Dir(m.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.path, data, 0o644)
}

// Delete removes the checkpoint file. Safe to call when none exists.
func (m *Manager) Delete() error {
	err := os.Remove(m.path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// NewState creates a fresh checkpoint for the given spec hash.
func NewState(specHash string) *State {
	return &State{
		SpecHash:  specHash,
		StartedAt: time.Now(),
		Completed: make(map[string]bool),
	}
}

// OperationKey returns the canonical key for an operation.
func OperationKey(method, path string) string {
	return method + " " + path
}
