package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScore_NoIssues(t *testing.T) {
	assert.Equal(t, 100, Score(nil))
}

func TestScore_Errors(t *testing.T) {
	issues := []LintIssue{
		{Severity: "error"},
		{Severity: "error"},
	}
	assert.Equal(t, 80, Score(issues))
}

func TestScore_Warnings(t *testing.T) {
	issues := []LintIssue{
		{Severity: "warning"},
		{Severity: "warning"},
		{Severity: "warning"},
	}
	assert.Equal(t, 91, Score(issues))
}

func TestScore_Mixed(t *testing.T) {
	issues := []LintIssue{
		{Severity: "error"},
		{Severity: "warning"},
		{Severity: "warning"},
	}
	assert.Equal(t, 84, Score(issues))
}

func TestScore_ClampedAtZero(t *testing.T) {
	var issues []LintIssue
	for i := 0; i < 15; i++ {
		issues = append(issues, LintIssue{Severity: "error"})
	}
	assert.Equal(t, 0, Score(issues))
}
