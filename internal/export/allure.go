// internal/export/allure.go
package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/testmind-hq/caseforge/internal/output/schema"
)

// AllureExporter writes one Allure result JSON file per TestCase.
type AllureExporter struct{}

func (e *AllureExporter) Format() string { return "allure" }

func (e *AllureExporter) Export(cases []schema.TestCase, outDir string) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("allure export: mkdir %s: %w", outDir, err)
	}
	for _, tc := range cases {
		if err := writeAllureResult(tc, outDir); err != nil {
			return err
		}
	}
	return nil
}

func writeAllureResult(tc schema.TestCase, outDir string) error {
	uuid := toUUID(tc.ID)
	result := map[string]any{
		"uuid":      uuid,
		"historyId": tc.ID,
		"name":      tc.Title,
		"fullName":  tc.ID + " " + tc.Title,
		"status":    "unknown",
		"labels":    allureLabels(tc),
		"steps":     allureSteps(tc.Steps),
		"start":     tc.GeneratedAt.UnixMilli(),
		"stop":      tc.GeneratedAt.UnixMilli(),
	}
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("allure marshal %s: %w", tc.ID, err)
	}
	return os.WriteFile(filepath.Join(outDir, uuid+"-result.json"), data, 0o644)
}

func allureLabels(tc schema.TestCase) []map[string]string {
	labels := []map[string]string{
		{"name": "severity", "value": PriorityAllure(tc.Priority)},
	}
	for _, tag := range tc.Tags {
		labels = append(labels, map[string]string{"name": "tag", "value": tag})
	}
	if tc.Source.SpecPath != "" {
		labels = append(labels, map[string]string{"name": "suite", "value": tc.Source.SpecPath})
	}
	return labels
}

func allureSteps(steps []schema.Step) []map[string]any {
	out := make([]map[string]any, len(steps))
	for i, s := range steps {
		out[i] = map[string]any{
			"name":       s.Method + " " + s.Path,
			"status":     "unknown",
			"parameters": []any{},
		}
	}
	return out
}
