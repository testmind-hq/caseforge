// internal/sandbox/state_test.go
package sandbox

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStateStore_WriteReadDelete(t *testing.T) {
	s := newMemStateStore()

	s.Write("pets", "1", map[string]any{"id": "1", "name": "Fido"})

	obj, ok := s.Read("pets", "1")
	require.True(t, ok)
	assert.Equal(t, "Fido", obj["name"])

	s.Delete("pets", "1")
	_, ok = s.Read("pets", "1")
	assert.False(t, ok)
}

func TestStateStore_List(t *testing.T) {
	s := newMemStateStore()
	s.Write("pets", "1", map[string]any{"id": "1"})
	s.Write("pets", "2", map[string]any{"id": "2"})

	list := s.List("pets")
	assert.Len(t, list, 2)
}

func TestStateStore_List_Empty(t *testing.T) {
	s := newMemStateStore()
	list := s.List("pets")
	assert.NotNil(t, list)
	assert.Len(t, list, 0)
}

func TestStateStore_ConcurrentReadWrite(t *testing.T) {
	s := newMemStateStore()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		id := fmt.Sprintf("%d", i)
		go func(id string) {
			defer wg.Done()
			s.Write("pets", id, map[string]any{"id": id})
		}(id)
		go func(id string) {
			defer wg.Done()
			_, _ = s.Read("pets", id)
		}(id)
	}
	wg.Wait()
}
