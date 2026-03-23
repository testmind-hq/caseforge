package runner

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVars(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected map[string]string
	}{
		{
			name:     "empty input",
			input:    nil,
			expected: map[string]string{},
		},
		{
			name:     "single var",
			input:    []string{"base_url=http://localhost:8080"},
			expected: map[string]string{"base_url": "http://localhost:8080"},
		},
		{
			name:     "multiple vars",
			input:    []string{"base_url=http://localhost", "token=abc123"},
			expected: map[string]string{"base_url": "http://localhost", "token": "abc123"},
		},
		{
			name:     "value with equals sign",
			input:    []string{"key=val=with=equals"},
			expected: map[string]string{"key": "val=with=equals"},
		},
		{
			name:     "malformed entry without equals is skipped",
			input:    []string{"noequals"},
			expected: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseVars(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestHurlRunnerNoBinary(t *testing.T) {
	// Temporarily clear PATH to simulate hurl not being installed
	t.Setenv("PATH", "")

	r := NewHurlRunner()
	_, err := r.Run(t.TempDir(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "hurl not found on PATH")
}

func TestBuildRunResult(t *testing.T) {
	reportJSON := `{
		"entries": [
			{"filename": "TC-abc.hurl", "success": true},
			{"filename": "TC-def.hurl", "success": false},
			{"filename": "TC-ghi.hurl", "success": true}
		]
	}`
	result := buildRunResult([]byte(reportJSON))
	assert.Equal(t, 2, result.Passed)
	assert.Equal(t, 1, result.Failed)
	assert.Len(t, result.Cases, 3)
	assert.Equal(t, "TC-abc", result.Cases[0].ID)
	assert.True(t, result.Cases[0].Passed)
	assert.Equal(t, "TC-def", result.Cases[1].ID)
	assert.False(t, result.Cases[1].Passed)
}

func TestBuildRunResultInvalidJSON(t *testing.T) {
	result := buildRunResult([]byte("not json"))
	assert.Equal(t, 0, result.Passed)
	assert.Equal(t, 0, result.Failed)
	assert.Nil(t, result.Cases)
}
