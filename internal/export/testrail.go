// internal/export/testrail.go
package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/testmind-hq/caseforge/internal/output/schema"
)

// TestRailExporter writes a single testrail-import.csv for TestRail import.
type TestRailExporter struct{}

func (e *TestRailExporter) Format() string { return "testrail" }

func (e *TestRailExporter) Export(cases []schema.TestCase, outDir string) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("testrail export: mkdir %s: %w", outDir, err)
	}
	f, err := os.Create(filepath.Join(outDir, "testrail-import.csv"))
	if err != nil {
		return fmt.Errorf("testrail export: create file: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)

	if err := w.Write([]string{"Title", "Type", "Priority", "Section", "Steps", "Expected Results"}); err != nil {
		return err
	}
	for _, tc := range cases {
		if err := w.Write(testRailRow(tc)); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func testRailRow(tc schema.TestCase) []string {
	steps, expected := testRailStepsAndExpected(tc.Steps)
	return []string{
		tc.ID + " " + tc.Title,
		"Automated",
		PriorityTestRail(tc.Priority),
		tc.Source.SpecPath,
		steps,
		expected,
	}
}

func testRailStepsAndExpected(steps []schema.Step) (string, string) {
	var stepsText, expectedText string
	for i, s := range steps {
		prefix := ""
		if len(steps) > 1 {
			prefix = fmt.Sprintf("Step %d: ", i+1)
		}
		stepsText += prefix + s.Method + " " + s.Path
		if i < len(steps)-1 {
			stepsText += "\n"
		}
		expectedText += AssertionsSummary(s.Assertions)
		if i < len(steps)-1 {
			expectedText += "\n"
		}
	}
	return stepsText, expectedText
}
