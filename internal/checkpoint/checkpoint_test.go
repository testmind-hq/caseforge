// internal/checkpoint/checkpoint_test.go
package checkpoint

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewState(t *testing.T) {
	s := NewState("abc123")
	assert.Equal(t, "abc123", s.SpecHash)
	assert.NotNil(t, s.Completed)
	assert.WithinDuration(t, time.Now(), s.StartedAt, 2*time.Second)
}

func TestOperationKey(t *testing.T) {
	assert.Equal(t, "POST /pets", OperationKey("POST", "/pets"))
	assert.Equal(t, "GET /users", OperationKey("GET", "/users"))
}

func TestManager_LoadNonExistent(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(dir)
	s, err := m.Load()
	require.NoError(t, err)
	assert.Nil(t, s)
}

func TestManager_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(dir)

	original := NewState("hash42")
	original.Completed["GET /pets"] = true
	original.Completed["POST /pets"] = true

	require.NoError(t, m.Save(original))

	loaded, err := m.Load()
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, "hash42", loaded.SpecHash)
	assert.True(t, loaded.Completed["GET /pets"])
	assert.True(t, loaded.Completed["POST /pets"])
	assert.False(t, loaded.Completed["DELETE /pets"])
}

func TestManager_SaveCreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "output")
	m := NewManager(dir)
	s := NewState("h1")
	require.NoError(t, m.Save(s))
	_, err := os.Stat(filepath.Join(dir, stateFileName))
	require.NoError(t, err)
}

func TestManager_Delete(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(dir)

	require.NoError(t, m.Save(NewState("h1")))
	require.NoError(t, m.Delete())

	loaded, err := m.Load()
	require.NoError(t, err)
	assert.Nil(t, loaded)
}

func TestManager_DeleteNonExistent(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(dir)
	// Should not error when file doesn't exist.
	assert.NoError(t, m.Delete())
}

func TestManager_LoadCorruptFile(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, stateFileName), []byte("not json"), 0o644))
	_, err := m.Load()
	assert.Error(t, err)
}
