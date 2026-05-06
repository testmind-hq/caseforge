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
		filename := FilenameFor(tc) + ".hurl"
		if err := os.WriteFile(filepath.Join(outDir, filename), []byte(content), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", filename, err)
		}
	}
	return nil
}

// singleLineSep returns "# ── {title} " padded to ~50 chars with "─" characters.
func singleLineSep(title string) string {
	prefix := fmt.Sprintf("# ── %s ", title)
	const totalWidth = 50
	padding := totalWidth - len(prefix)
	if padding < 2 {
		padding = 2
	}
	return prefix + strings.Repeat("─", padding)
}

const chainSep = "# ══════════════════════════════════════════════════"

func (r *HurlRenderer) renderCase(tc schema.TestCase) string {
	if tc.Kind == "chain" {
		return r.renderChainCase(tc)
	}
	return r.renderSingleCase(tc)
}

func (r *HurlRenderer) renderSingleCase(tc schema.TestCase) string {
	var b strings.Builder

	// Use title from case, or from first step if case title empty.
	caseTitle := tc.Title
	if caseTitle == "" && len(tc.Steps) > 0 {
		caseTitle = tc.Steps[0].Title
	}

	for _, step := range tc.Steps {
		// Appendix B single case header: "# ── {title} ────..."
		b.WriteString(singleLineSep(caseTitle))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("# case_id=%s\n", tc.ID))
		b.WriteString(fmt.Sprintf("# case_name=%s\n", caseTitle))
		b.WriteString(fmt.Sprintf("# step_id=%s\n", step.ID))
		b.WriteString(fmt.Sprintf("# step_type=%s\n", step.Type))
		if tc.Source.Technique != "" {
			b.WriteString(fmt.Sprintf("# technique=%s\n", tc.Source.Technique))
		}
		b.WriteString(fmt.Sprintf("# priority=%s\n", tc.Priority))
		b.WriteString("\n")

		b.WriteString(r.renderStep(step))
		b.WriteString("\n")
	}
	return b.String()
}

func (r *HurlRenderer) renderChainCase(tc schema.TestCase) string {
	var b strings.Builder

	// Chain file header with double-line separator
	b.WriteString(chainSep)
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("# %s\n", tc.Title))
	b.WriteString(fmt.Sprintf("# case_id=%s\n", tc.ID))
	b.WriteString(fmt.Sprintf("# case_name=%s\n", tc.Title))
	b.WriteString("# case_kind=chain\n")
	b.WriteString(fmt.Sprintf("# priority=%s\n", tc.Priority))
	b.WriteString(chainSep)
	b.WriteString("\n\n")

	for _, step := range tc.Steps {
		// Per-step separator: "# ── {step.Title} [{step.Type}] ──..."
		stepTitle := fmt.Sprintf("%s [%s]", step.Title, step.Type)
		b.WriteString(singleLineSep(stepTitle))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("# step_id=%s\n", step.ID))
		b.WriteString(fmt.Sprintf("# step_type=%s\n", step.Type))
		b.WriteString(fmt.Sprintf("# title=%s\n", step.Title))
		for _, dep := range step.DependsOn {
			b.WriteString(fmt.Sprintf("# depends_on=%s\n", dep))
		}
		b.WriteString("\n")

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
	// Exact status (eq) goes into the HTTP line; range operators (gte/lte/gt/lt/ne)
	// go into the [Asserts] block as `status <op> N` so Hurl can evaluate them.
	exactStatus := -1
	var asserts []schema.Assertion
	for _, a := range step.Assertions {
		if a.Target == "status_code" {
			if a.Operator == "eq" {
				if code, ok := a.Expected.(float64); ok {
					exactStatus = int(code)
				} else if code, ok := a.Expected.(int); ok {
					exactStatus = code
				}
			} else {
				// Range assertion — render as `status <op> N` in [Asserts].
				asserts = append(asserts, a)
			}
		} else {
			asserts = append(asserts, a)
		}
	}
	if exactStatus >= 0 {
		b.WriteString(fmt.Sprintf("HTTP %d\n", exactStatus))
	} else {
		// No exact status constraint; use wildcard so range asserts in [Asserts] govern.
		b.WriteString("HTTP *\n")
	}

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

// formatHurlRegex wraps a pattern in Hurl's /regex/ literal syntax.
// Hurl's matches predicate requires /pattern/, not "pattern".
func formatHurlRegex(v any) string {
	if v == nil {
		return "/./"
	}
	return fmt.Sprintf("/%v/", v)
}

func renderAssertion(a schema.Assertion) string {
	// Normalise legacy "body.field" → "jsonpath $.field" to avoid duplicating
	// the full operator set in a separate case block.
	if strings.HasPrefix(a.Target, "body.") {
		a.Target = "jsonpath $." + strings.TrimPrefix(a.Target, "body.")
	}
	switch {
	case a.Target == "status_code":
		// Non-eq status operators (gte/gt/lte/lt/ne) routed here from step rendering.
		switch a.Operator {
		case "gte":
			return fmt.Sprintf("status >= %v\n", a.Expected)
		case "lte":
			return fmt.Sprintf("status <= %v\n", a.Expected)
		case "gt":
			return fmt.Sprintf("status > %v\n", a.Expected)
		case "lt":
			return fmt.Sprintf("status < %v\n", a.Expected)
		case "ne":
			return fmt.Sprintf("status != %v\n", a.Expected)
		}
	case a.Target == "duration_ms":
		switch a.Operator {
		case "lt":
			return fmt.Sprintf("duration < %v\n", a.Expected)
		case "lte":
			return fmt.Sprintf("duration <= %v\n", a.Expected)
		case "gt":
			return fmt.Sprintf("duration > %v\n", a.Expected)
		case "gte":
			return fmt.Sprintf("duration >= %v\n", a.Expected)
		}
	case strings.HasPrefix(a.Target, "jsonpath "):
		expr := strings.TrimPrefix(a.Target, "jsonpath ")
		switch a.Operator {
		case "eq":
			return fmt.Sprintf("jsonpath %q == %s\n", expr, formatHurlValue(a.Expected))
		case "ne":
			if a.Expected == nil {
				return fmt.Sprintf("jsonpath %q not exists\n", expr)
			}
			return fmt.Sprintf("jsonpath %q != %s\n", expr, formatHurlValue(a.Expected))
		case "lt":
			return fmt.Sprintf("jsonpath %q < %v\n", expr, a.Expected)
		case "lte":
			return fmt.Sprintf("jsonpath %q <= %v\n", expr, a.Expected)
		case "gt":
			return fmt.Sprintf("jsonpath %q > %v\n", expr, a.Expected)
		case "gte":
			return fmt.Sprintf("jsonpath %q >= %v\n", expr, a.Expected)
		case "exists":
			return fmt.Sprintf("jsonpath %q exists\n", expr)
		case "not_exists":
			return fmt.Sprintf("jsonpath %q not exists\n", expr)
		case "contains":
			return fmt.Sprintf("jsonpath %q contains %s\n", expr, formatHurlValue(a.Expected))
		case "matches":
			return fmt.Sprintf("jsonpath %q matches %s\n", expr, formatHurlRegex(a.Expected))
		case "is_iso8601":
			return fmt.Sprintf("jsonpath %q isDate\n", expr)
		case "is_uuid":
			return fmt.Sprintf("jsonpath %q matches /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i\n", expr)
		}
	case strings.HasPrefix(a.Target, "header "):
		headerName := strings.TrimPrefix(a.Target, "header ")
		switch a.Operator {
		case "eq":
			return fmt.Sprintf("header %q == %s\n", headerName, formatHurlValue(a.Expected))
		case "ne":
			if a.Expected == nil {
				return fmt.Sprintf("header %q not exists\n", headerName)
			}
			return fmt.Sprintf("header %q != %s\n", headerName, formatHurlValue(a.Expected))
		case "exists":
			return fmt.Sprintf("header %q exists\n", headerName)
		case "contains":
			return fmt.Sprintf("header %q contains %s\n", headerName, formatHurlValue(a.Expected))
		case "matches":
			return fmt.Sprintf("header %q matches %s\n", headerName, formatHurlRegex(a.Expected))
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
