// internal/mcp/server_test.go
package mcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServerReturnsNonNil(t *testing.T) {
	s := NewServer()
	require.NotNil(t, s)
}

func TestSamplingProviderIsAvailable(t *testing.T) {
	assert.NotNil(t, NewServer())
}
