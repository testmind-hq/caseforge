// internal/llm/gemini_test.go
package llm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeminiProviderName(t *testing.T) {
	p := &GeminiProvider{model: "gemini-2.5-flash"}
	assert.Equal(t, "gemini:gemini-2.5-flash", p.Name())
	assert.True(t, p.IsAvailable())
}

func TestNewProvider_GeminiWithKey(t *testing.T) {
	p := NewProvider("fake-gemini-key", "gemini", "gemini-2.5-flash")
	require.NotNil(t, p)
	assert.Equal(t, "gemini:gemini-2.5-flash", p.Name())
}

func TestNewProvider_GeminiNoKey(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "")
	t.Setenv("GOOGLE_API_KEY", "")
	p := NewProvider("", "gemini", "gemini-2.5-flash")
	assert.Equal(t, "noop", p.Name())
}

func TestNewProvider_GeminiDefaultModel(t *testing.T) {
	p := NewProvider("fake-key", "gemini", "")
	assert.Equal(t, "gemini:gemini-2.5-flash", p.Name())
}
