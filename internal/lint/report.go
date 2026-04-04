// internal/lint/report.go
package lint

import "encoding/json"

// LintReport is the structured output of a lint run.
type LintReport struct {
	Score        int         `json:"score"`
	ErrorCount   int         `json:"error_count"`
	WarningCount int         `json:"warning_count"`
	Issues       []LintIssue `json:"issues"`
}

// NewReport builds a LintReport from a slice of issues.
func NewReport(issues []LintIssue) LintReport {
	r := LintReport{Issues: issues}
	if r.Issues == nil {
		r.Issues = []LintIssue{}
	}
	for _, iss := range issues {
		switch iss.Severity {
		case "error":
			r.ErrorCount++
		case "warning":
			r.WarningCount++
		}
	}
	r.Score = Score(issues)
	return r
}

// ToJSON serialises the report to indented JSON.
func (r LintReport) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}
