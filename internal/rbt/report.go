// internal/rbt/report.go
package rbt

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

// PrintTerminal writes a formatted risk report table to w.
func PrintTerminal(w io.Writer, report RiskReport) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "Operation\tRisk\tAffected\tCases\tSource")
	fmt.Fprintln(tw, strings.Repeat("-", 20)+"\t"+strings.Repeat("-", 8)+"\t"+strings.Repeat("-", 8)+"\t"+strings.Repeat("-", 5)+"\t"+strings.Repeat("-", 30))

	for _, op := range report.Operations {
		affected := "-"
		if op.Affected {
			affected = "✓"
		}
		source := "-"
		if len(op.SourceRefs) > 0 {
			ref := op.SourceRefs[0]
			source = ref.SourceFile
			if ref.Line > 0 {
				source = fmt.Sprintf("%s:%d", ref.SourceFile, ref.Line)
			}
			if ref.Via != "" {
				source = fmt.Sprintf("%s (%s)", source, ref.Via)
			}
		}
		fmt.Fprintf(tw, "%s %s\t%s\t%s\t%d\t%s\n",
			op.Method, op.Path,
			strings.ToUpper(string(op.Risk)),
			affected,
			len(op.TestCases),
			source,
		)
	}
	tw.Flush()

	fmt.Fprintf(w, "\nRisk Score: %.2f  (%d uncovered / %d affected)\n",
		report.RiskScore, report.TotalUncovered, report.TotalAffected)
}

// ShouldFail returns true if the report has any operation with risk >= failOn level.
// failOn values: "none" (never fail), "low", "medium", "high".
func ShouldFail(report RiskReport, failOn string) bool {
	if failOn == "none" {
		return false
	}
	threshold := riskOrder(RiskLevel(failOn))
	for _, op := range report.Operations {
		if riskOrder(op.Risk) >= threshold && op.Risk != RiskNone && op.Risk != RiskUncertain {
			return true
		}
	}
	return false
}

func riskOrder(r RiskLevel) int {
	switch r {
	case RiskLow:
		return 1
	case RiskMedium:
		return 2
	case RiskHigh:
		return 3
	}
	return 0
}

type reportJSON struct {
	DiffBase       string              `json:"diff_base"`
	DiffHead       string              `json:"diff_head"`
	ChangedFiles   []ChangedFile       `json:"changed_files"`
	Operations     []OperationCoverage `json:"operations"`
	TotalAffected  int                 `json:"total_affected"`
	TotalCovered   int                 `json:"total_covered"`
	TotalUncovered int                 `json:"total_uncovered"`
	RiskScore      float64             `json:"risk_score"`
	GeneratedAt    string              `json:"generated_at"`
}

// MarshalReportJSON serializes a RiskReport to indented JSON using a consistent
// snake_case schema. Used by both --format json stdout and WriteReportJSON.
func MarshalReportJSON(report RiskReport) ([]byte, error) {
	rj := reportJSON{
		DiffBase:       report.DiffBase,
		DiffHead:       report.DiffHead,
		ChangedFiles:   report.ChangedFiles,
		Operations:     report.Operations,
		TotalAffected:  report.TotalAffected,
		TotalCovered:   report.TotalCovered,
		TotalUncovered: report.TotalUncovered,
		RiskScore:      report.RiskScore,
		GeneratedAt:    report.GeneratedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
	return json.MarshalIndent(rj, "", "  ")
}

// WriteReportJSON writes the report as rbt-report.json to outputDir.
func WriteReportJSON(outputDir string, report RiskReport) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("create output dir: %w", err)
	}

	data, err := MarshalReportJSON(report)
	if err != nil {
		return "", err
	}

	outPath := filepath.Join(outputDir, "rbt-report.json")
	if err := os.WriteFile(outPath, data, 0644); err != nil {
		return "", err
	}
	return outPath, nil
}
