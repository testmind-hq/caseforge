// internal/rbt/report_test.go
package rbt

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sampleReport() RiskReport {
	return RiskReport{
		DiffBase: "HEAD~1",
		DiffHead: "HEAD",
		Operations: []OperationCoverage{
			{
				Method: "POST", Path: "/users",
				Affected: true, Risk: RiskHigh,
				TestCases:  nil,
				SourceRefs: []RouteMapping{{SourceFile: "handler/user.go", Line: 42, Via: "regex"}},
			},
			{
				Method: "GET", Path: "/users/{id}",
				Affected: true, Risk: RiskMedium,
				TestCases:  []TestCaseRef{{CaseID: "c1"}},
				SourceRefs: []RouteMapping{{SourceFile: "handler/user.go", Line: 18}},
			},
			{
				Method: "DELETE", Path: "/users/{id}",
				Affected: false, Risk: RiskNone,
				TestCases: []TestCaseRef{{CaseID: "c2"}, {CaseID: "c3"}},
			},
		},
		TotalAffected:  2,
		TotalCovered:   1,
		TotalUncovered: 1,
		RiskScore:      0.5,
		GeneratedAt:    time.Now(),
	}
}

func TestPrintTerminal_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	PrintTerminal(&buf, sampleReport())
	out := buf.String()
	assert.Contains(t, out, "Operation")
	assert.Contains(t, out, "Risk")
	assert.Contains(t, out, "Cases")
}

func TestPrintTerminal_ContainsOperations(t *testing.T) {
	var buf bytes.Buffer
	PrintTerminal(&buf, sampleReport())
	out := buf.String()
	assert.Contains(t, out, "POST /users")
	assert.Contains(t, out, "HIGH")
	assert.Contains(t, out, "GET /users/{id}")
	assert.Contains(t, out, "MEDIUM")
}

func TestPrintTerminal_ContainsRiskScore(t *testing.T) {
	var buf bytes.Buffer
	PrintTerminal(&buf, sampleReport())
	out := buf.String()
	assert.Contains(t, out, "0.50")
}

func TestShouldFailOn_High(t *testing.T) {
	report := sampleReport()
	assert.True(t, ShouldFail(report, "high"))
	assert.True(t, ShouldFail(report, "medium"))
	assert.False(t, ShouldFail(report, "none"))
}

func TestShouldFailOn_NoneLevel_NeverFails(t *testing.T) {
	report := sampleReport()
	assert.False(t, ShouldFail(report, "none"))
}

func TestWriteReportJSON_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	report := sampleReport()
	path, err := WriteReportJSON(dir, report)
	require.NoError(t, err)
	assert.True(t, strings.HasSuffix(path, "rbt-report.json"))
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"diff_base"`)
}
