// internal/export/xray.go
package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/testmind-hq/caseforge/internal/output/schema"
)

// XrayExporter writes a single xray-import.json for Jira Xray Cloud import.
type XrayExporter struct{}

func (e *XrayExporter) Format() string { return "xray" }

func (e *XrayExporter) Export(cases []schema.TestCase, outDir string) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("xray export: mkdir %s: %w", outDir, err)
	}
	tests := make([]map[string]any, 0, len(cases))
	for _, tc := range cases {
		tests = append(tests, xrayTest(tc))
	}
	out := map[string]any{"tests": tests}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return fmt.Errorf("xray marshal: %w", err)
	}
	return os.WriteFile(filepath.Join(outDir, "xray-import.json"), data, 0o644)
}

func xrayTest(tc schema.TestCase) map[string]any {
	desc := fmt.Sprintf("Technique: %s\nSpec: %s\nRationale: %s",
		tc.Source.Technique, tc.Source.SpecPath, tc.Source.Rationale)
	labels := tc.Tags
	if labels == nil {
		labels = []string{}
	}
	return map[string]any{
		"summary":     tc.ID + " " + tc.Title,
		"description": desc,
		"testType":    "Generic",
		"priority":    PriorityXray(tc.Priority),
		"labels":      labels,
		"steps":       xraySteps(tc.Steps),
	}
}

func xraySteps(steps []schema.Step) []map[string]string {
	out := make([]map[string]string, len(steps))
	for i, s := range steps {
		action := s.Method + " " + s.Path
		if s.Body != nil {
			if b, err := json.Marshal(s.Body); err == nil {
				action += " body: " + string(b)
			}
		}
		out[i] = map[string]string{
			"action": action,
			"result": AssertionsSummary(s.Assertions),
		}
	}
	return out
}
