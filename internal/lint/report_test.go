// internal/lint/report_test.go
package lint

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewReport_Counts(t *testing.T) {
	issues := []LintIssue{
		{RuleID: "L001", Severity: "warning", Message: "w1", Path: "GET /a"},
		{RuleID: "L004", Severity: "error", Message: "e1", Path: "POST /b"},
		{RuleID: "L002", Severity: "warning", Message: "w2", Path: "GET /c"},
	}
	r := NewReport(issues)
	assert.Equal(t, 1, r.ErrorCount)
	assert.Equal(t, 2, r.WarningCount)
	assert.Equal(t, 3, len(r.Issues))
	assert.Equal(t, Score(issues), r.Score)
}

func TestNewReport_EmptyIssues(t *testing.T) {
	r := NewReport(nil)
	assert.Equal(t, 0, r.ErrorCount)
	assert.Equal(t, 0, r.WarningCount)
	assert.Equal(t, 100, r.Score)
	assert.Empty(t, r.Issues)
}

func TestLintReport_ToJSON(t *testing.T) {
	r := NewReport([]LintIssue{
		{RuleID: "L001", Severity: "warning", Message: "missing operationId", Path: "GET /users"},
	})
	data, err := r.ToJSON()
	assert.NoError(t, err)
	var out map[string]any
	assert.NoError(t, json.Unmarshal(data, &out))
	assert.Equal(t, float64(r.Score), out["score"])
	assert.Equal(t, float64(1), out["warning_count"])
	assert.Equal(t, float64(0), out["error_count"])
	issues := out["issues"].([]any)
	assert.Len(t, issues, 1)
}
