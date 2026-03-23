// internal/runner/k6_test.go
package runner

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestK6RunnerParseSummary(t *testing.T) {
	summaryJSON := `{
		"metrics": {
			"checks": {
				"passes": 8,
				"fails":  2
			}
		}
	}`
	result, err := parseK6Summary([]byte(summaryJSON))
	require.NoError(t, err)
	assert.Equal(t, 8, result.Passed)
	assert.Equal(t, 2, result.Failed)
	assert.Empty(t, result.Cases)
}

func TestK6RunnerParseSummaryAllPass(t *testing.T) {
	summaryJSON := `{"metrics":{"checks":{"passes":5,"fails":0}}}`
	result, err := parseK6Summary([]byte(summaryJSON))
	require.NoError(t, err)
	assert.Equal(t, 5, result.Passed)
	assert.Equal(t, 0, result.Failed)
}

func TestK6RunnerParseSummaryNoChecks(t *testing.T) {
	summaryJSON := `{"metrics":{}}`
	result, err := parseK6Summary([]byte(summaryJSON))
	require.NoError(t, err)
	assert.Equal(t, 0, result.Passed)
	assert.Equal(t, 0, result.Failed)
}

func TestK6RunnerNotFound(t *testing.T) {
	r := &K6Runner{k6Bin: "/nonexistent/k6"}
	dir := t.TempDir()
	_, err := r.Run(dir, nil)
	assert.Error(t, err)
}
