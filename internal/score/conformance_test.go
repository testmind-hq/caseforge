// internal/score/conformance_test.go
package score

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComputeConformance_UnknownTrend_NoHistory(t *testing.T) {
	report := Report{Overall: 72}
	result := ComputeConformance(report, ConformanceHistory{})
	assert.Equal(t, 72, result.Score)
	assert.Equal(t, TrendUnknown, result.Trend)
}

func TestComputeConformance_ImprovingTrend(t *testing.T) {
	history := ConformanceHistory{
		Records: []ConformanceRecord{
			{Date: time.Now().Add(-24 * time.Hour), Overall: 60},
		},
	}
	report := Report{Overall: 75}
	result := ComputeConformance(report, history)
	assert.Equal(t, 75, result.Score)
	assert.Equal(t, TrendImproving, result.Trend)
}

func TestComputeConformance_DegradingTrend(t *testing.T) {
	history := ConformanceHistory{
		Records: []ConformanceRecord{
			{Date: time.Now().Add(-24 * time.Hour), Overall: 80},
		},
	}
	report := Report{Overall: 65}
	result := ComputeConformance(report, history)
	assert.Equal(t, 65, result.Score)
	assert.Equal(t, TrendDegrading, result.Trend)
}

func TestComputeConformance_StableTrend(t *testing.T) {
	history := ConformanceHistory{
		Records: []ConformanceRecord{
			{Date: time.Now().Add(-24 * time.Hour), Overall: 72},
		},
	}
	report := Report{Overall: 72}
	result := ComputeConformance(report, history)
	assert.Equal(t, 72, result.Score)
	assert.Equal(t, TrendStable, result.Trend)
}

func TestLoadHistory_NonExistentFile_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent-history.json")
	h, err := LoadHistory(path)
	require.NoError(t, err)
	assert.Empty(t, h.Records)
}

func TestSaveHistory_AppendsRecord(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "conformance.json")

	// Save first record.
	err := SaveHistory(path, ConformanceHistory{}, 70)
	require.NoError(t, err)

	// Load and verify.
	h, err := LoadHistory(path)
	require.NoError(t, err)
	require.Len(t, h.Records, 1)
	assert.Equal(t, 70, h.Records[0].Overall)

	// Save second record.
	err = SaveHistory(path, h, 80)
	require.NoError(t, err)

	// Load and verify both records.
	h2, err := LoadHistory(path)
	require.NoError(t, err)
	require.Len(t, h2.Records, 2)
	assert.Equal(t, 70, h2.Records[0].Overall)
	assert.Equal(t, 80, h2.Records[1].Overall)

	// Verify file is valid JSON.
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	var raw map[string]any
	require.NoError(t, json.Unmarshal(data, &raw))
}
