// internal/mcp/server_test.go
package mcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServerReturnsNonNil(t *testing.T) {
	s, provider := NewServer()
	require.NotNil(t, s)
	require.NotNil(t, provider)
}

func TestSamplingProviderIsAvailable(t *testing.T) {
	_, provider := NewServer()
	assert.True(t, provider.IsAvailable())
	assert.Equal(t, "mcp-sampling", provider.Name())
}
