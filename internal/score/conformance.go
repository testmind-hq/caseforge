// internal/score/conformance.go
package score

import (
	"encoding/json"
	"os"
	"time"
)

// ConformanceRecord is a single history entry.
type ConformanceRecord struct {
	Date    time.Time `json:"date"`
	Overall int       `json:"overall"`
}

// ConformanceHistory is the full history file format.
type ConformanceHistory struct {
	Records []ConformanceRecord `json:"records"`
}

// ConformanceTrend summarises how the score changed vs previous run.
type ConformanceTrend string

const (
	TrendImproving ConformanceTrend = "improving"
	TrendDegrading ConformanceTrend = "degrading"
	TrendStable    ConformanceTrend = "stable"
	TrendUnknown   ConformanceTrend = "unknown"
)

// ConformanceResult is the conformance block returned in score output.
type ConformanceResult struct {
	Score int              `json:"score"`
	Trend ConformanceTrend `json:"trend"`
	// ByOperation maps "METHOD /path" to its individual score (0-100).
	ByOperation map[string]int `json:"by_operation,omitempty"`
}

// ComputeConformance computes a per-operation conformance score and trend.
// Score: report.Overall (the existing weighted overall score).
// Trend: compare current Overall score to the last record in history.
func ComputeConformance(report Report, history ConformanceHistory) ConformanceResult {
	current := report.Overall

	var trend ConformanceTrend
	if len(history.Records) == 0 {
		trend = TrendUnknown
	} else {
		last := history.Records[len(history.Records)-1].Overall
		switch {
		case last < current:
			trend = TrendImproving
		case last > current:
			trend = TrendDegrading
		default:
			trend = TrendStable
		}
	}

	return ConformanceResult{
		Score: current,
		Trend: trend,
	}
}

// LoadHistory reads history from path. Returns empty history if path does not exist.
func LoadHistory(path string) (ConformanceHistory, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ConformanceHistory{}, nil
		}
		return ConformanceHistory{}, err
	}
	var h ConformanceHistory
	if err := json.Unmarshal(data, &h); err != nil {
		return ConformanceHistory{}, err
	}
	return h, nil
}

// SaveHistory appends the current Overall score to history at path.
func SaveHistory(path string, history ConformanceHistory, overall int) error {
	history.Records = append(history.Records, ConformanceRecord{
		Date:    time.Now(),
		Overall: overall,
	})
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
