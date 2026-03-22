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
	_, _, err := r.Run(t.TempDir(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "hurl not found on PATH")
}

func TestParseHurlReport(t *testing.T) {
	reportJSON := `{
		"entries": [
			{"filename": "TC-abc.hurl", "success": true},
			{"filename": "TC-def.hurl", "success": false},
			{"filename": "TC-ghi.hurl", "success": true}
		]
	}`
	passed, failed := parseHurlReport([]byte(reportJSON))
	assert.Equal(t, 2, passed)
	assert.Equal(t, 1, failed)
}

func TestParseHurlReportInvalidJSON(t *testing.T) {
	passed, failed := parseHurlReport([]byte("not json"))
	assert.Equal(t, 0, passed)
	assert.Equal(t, 0, failed)
}
