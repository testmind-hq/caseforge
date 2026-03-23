// internal/lint/score.go
package lint

// Score computes a spec quality score from 0 to 100.
// Each error subtracts 10 points; each warning subtracts 3 points.
func Score(issues []LintIssue) int {
	base := 100
	for _, issue := range issues {
		switch issue.Severity {
		case "error":
			base -= 10
		case "warning":
			base -= 3
		}
	}
	if base < 0 {
		return 0
	}
	return base
}
