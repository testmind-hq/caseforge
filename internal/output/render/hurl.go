// internal/output/render/hurl.go
package render

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/testmind-hq/caseforge/internal/output/schema"
)

type HurlRenderer struct {
	baseURL string
}

func NewHurlRenderer(baseURL string) *HurlRenderer {
	if baseURL == "" {
		baseURL = "{{base_url}}"
	}
	return &HurlRenderer{baseURL: baseURL}
}

func (r *HurlRenderer) Format() string { return "hurl" }

func (r *HurlRenderer) Render(cases []schema.TestCase, outDir string) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}
	for _, tc := range cases {
		content := r.renderCase(tc)
		filename := sanitizeFilename(tc.ID) + ".hurl"
		if err := os.WriteFile(filepath.Join(outDir, filename), []byte(content), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", filename, err)
		}
	}
	return nil
}

func (r *HurlRenderer) renderCase(tc schema.TestCase) string {
	var b strings.Builder

	// Machine-readable annotations (CaseForge whitelist keys)
	b.WriteString(fmt.Sprintf("# id=%s\n", tc.ID))
	b.WriteString(fmt.Sprintf("# priority=%s\n", tc.Priority))
	b.WriteString(fmt.Sprintf("# technique=%s\n", tc.Source.Technique))
	if tc.Source.SpecPath != "" {
		b.WriteString(fmt.Sprintf("# spec_path=%s\n", tc.Source.SpecPath))
	}
	// Human-readable separator and description
	b.WriteString("# ─────────────────────────────────────────────\n")
	b.WriteString(fmt.Sprintf("# %s\n", tc.Title))
	if tc.Source.Rationale != "" {
		b.WriteString(fmt.Sprintf("# %s\n", tc.Source.Rationale))
	}
	b.WriteString("\n")

	for _, step := range tc.Steps {
		b.WriteString(r.renderStep(step))
		b.WriteString("\n")
	}
	return b.String()
}

func (r *HurlRenderer) renderStep(step schema.Step) string {
	var b strings.Builder

	// Request line
	b.WriteString(fmt.Sprintf("%s %s%s\n", step.Method, r.baseURL, step.Path))

	// Headers — sort keys for deterministic output
	headerKeys := make([]string, 0, len(step.Headers))
	for k := range step.Headers {
		headerKeys = append(headerKeys, k)
	}
	sort.Strings(headerKeys)
	for _, k := range headerKeys {
		b.WriteString(fmt.Sprintf("%s: %s\n", k, step.Headers[k]))
	}

	// Body
	if step.Body != nil {
		data, err := json.MarshalIndent(step.Body, "", "  ")
		if err != nil {
			b.WriteString(fmt.Sprintf("# ERROR: body serialization failed: %v\n", err))
		} else {
			b.WriteString("```json\n")
			b.WriteString(string(data))
			b.WriteString("\n```\n")
		}
	}

	b.WriteString("\n")

	// Assertions (Hurl HTTP + asserts block)
	statusCode := 200
	var asserts []schema.Assertion
	for _, a := range step.Assertions {
		if a.Target == "status_code" {
			if code, ok := a.Expected.(float64); ok {
				statusCode = int(code)
			} else if code, ok := a.Expected.(int); ok {
				statusCode = code
			}
		} else {
			asserts = append(asserts, a)
		}
	}
	b.WriteString(fmt.Sprintf("HTTP %d\n", statusCode))

	// Captures block must come BEFORE [Asserts] (Hurl spec ordering requirement)
	if len(step.Captures) > 0 {
		b.WriteString("\n[Captures]\n")
		for _, cap := range step.Captures {
			b.WriteString(renderCapture(cap))
		}
	}

	if len(asserts) > 0 {
		b.WriteString("\n[Asserts]\n")
		for _, a := range asserts {
			b.WriteString(renderAssertion(a))
		}
	}

	return b.String()
}

func renderCapture(c schema.Capture) string {
	// From format: "jsonpath $.field" → `varName: jsonpath "$.field"`
	// From format: "header X-Name"   → `varName: header "X-Name"`
	parts := strings.SplitN(c.From, " ", 2)
	if len(parts) == 2 {
		return fmt.Sprintf("%s: %s %q\n", c.Name, parts[0], parts[1])
	}
	return fmt.Sprintf("%s: %s\n", c.Name, c.From)
}

func renderAssertion(a schema.Assertion) string {
	switch a.Target {
	case "duration_ms":
		if a.Operator == "lt" {
			return fmt.Sprintf("duration < %v\n", a.Expected)
		}
	default:
		// body.<field> or custom target
		target := a.Target
		if strings.HasPrefix(target, "body.") {
			field := strings.TrimPrefix(target, "body.")
			switch a.Operator {
			case "eq":
				return fmt.Sprintf("jsonpath \"$.%s\" == %s\n", field, formatHurlValue(a.Expected))
			case "exists":
				return fmt.Sprintf("jsonpath \"$.%s\" exists\n", field)
			case "contains":
				return fmt.Sprintf("jsonpath \"$.%s\" contains %s\n", field, formatHurlValue(a.Expected))
			}
		}
	}
	return fmt.Sprintf("# unrendered assertion: %s %s %v\n", a.Target, a.Operator, a.Expected)
}

func formatHurlValue(v any) string {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("%q", val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func sanitizeFilename(id string) string {
	return strings.ReplaceAll(id, "/", "_")
}
