// internal/output/render/k6.go
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

type K6Renderer struct{}

func NewK6Renderer() *K6Renderer { return &K6Renderer{} }

func (r *K6Renderer) Format() string { return "k6" }

func (r *K6Renderer) Render(cases []schema.TestCase, outDir string) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}
	var sb strings.Builder
	sb.WriteString("import http from 'k6/http';\n")
	sb.WriteString("import { check, group } from 'k6';\n\n")
	sb.WriteString("const BASE_URL = __ENV.BASE_URL || 'http://localhost';\n\n")
	sb.WriteString("export const options = {\n")
	sb.WriteString("  scenarios: {\n")
	sb.WriteString("    default: { executor: 'shared-iterations', vus: 1, iterations: 1 },\n")
	sb.WriteString("  },\n")
	sb.WriteString("};\n\n")
	sb.WriteString("export default function () {\n")
	for _, tc := range cases {
		sb.WriteString(renderK6Group(tc))
	}
	sb.WriteString("}\n")
	out := filepath.Join(outDir, "k6_tests.js")
	return os.WriteFile(out, []byte(sb.String()), 0o644)
}

func renderK6Group(tc schema.TestCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("  group('%s', function () {\n", tc.Title))
	for i, step := range tc.Steps {
		resVar := fmt.Sprintf("res%d", i+1)
		if i == 0 {
			resVar = "res"
		}
		path := interpolateK6Path(step.Path)
		method := strings.ToUpper(step.Method)
		urlExpr := fmt.Sprintf("`${BASE_URL}%s`", path)
		var callLines []string
		switch method {
		case "GET", "DELETE", "HEAD":
			callLines = buildK6NoBodyCall(method, urlExpr, step.Headers)
		default:
			callLines = buildK6BodyCall(method, urlExpr, step.Headers, step.Body)
		}
		sb.WriteString(fmt.Sprintf("    const %s = %s\n", resVar, strings.Join(callLines, "\n    ")))
		checks := buildK6Checks(step.Assertions)
		if len(checks) > 0 {
			sb.WriteString(fmt.Sprintf("    check(%s, {\n", resVar))
			for _, c := range checks {
				sb.WriteString("      " + c + "\n")
			}
			sb.WriteString("    });\n")
		}
		for _, cap := range step.Captures {
			if strings.HasPrefix(cap.From, "jsonpath $.") {
				field := strings.TrimPrefix(cap.From, "jsonpath $.")
				sb.WriteString(fmt.Sprintf("    const %s = %s.json('%s');\n", cap.Name, resVar, field))
			}
		}
	}
	sb.WriteString("  });\n")
	return sb.String()
}

func interpolateK6Path(path string) string {
	var result strings.Builder
	i := 0
	for i < len(path) {
		if i+1 < len(path) && path[i] == '{' && path[i+1] == '{' {
			end := strings.Index(path[i+2:], "}}")
			if end >= 0 {
				varName := path[i+2 : i+2+end]
				result.WriteString("${" + varName + "}")
				i = i + 2 + end + 2
				continue
			}
		}
		result.WriteByte(path[i])
		i++
	}
	return result.String()
}

func buildK6NoBodyCall(method, urlExpr string, headers map[string]string) []string {
	headerJS := buildK6Headers(headers)
	switch strings.ToUpper(method) {
	case "HEAD":
		if headerJS != "" {
			return []string{fmt.Sprintf(`http.request("HEAD", %s, null, %s);`, urlExpr, headerJS)}
		}
		return []string{fmt.Sprintf(`http.request("HEAD", %s, null);`, urlExpr)}
	case "DELETE":
		if headerJS != "" {
			return []string{fmt.Sprintf("http.del(%s, null, %s);", urlExpr, headerJS)}
		}
		return []string{fmt.Sprintf("http.del(%s);", urlExpr)}
	default: // GET, OPTIONS, etc.
		if headerJS != "" {
			return []string{fmt.Sprintf("http.%s(%s, %s);", strings.ToLower(method), urlExpr, headerJS)}
		}
		return []string{fmt.Sprintf("http.%s(%s);", strings.ToLower(method), urlExpr)}
	}
}

func buildK6BodyCall(method, urlExpr string, headers map[string]string, body any) []string {
	bodyJS := "null"
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			bodyJS = "/* ERROR: body serialization failed: " + err.Error() + " */"
		} else {
			bodyJS = "JSON.stringify(" + string(b) + ")"
		}
	}
	headerJS := buildK6Headers(headers)
	if headerJS == "" {
		return []string{fmt.Sprintf("http.%s(%s, %s);", strings.ToLower(method), urlExpr, bodyJS)}
	}
	return []string{fmt.Sprintf("http.%s(%s, %s, %s);", strings.ToLower(method), urlExpr, bodyJS, headerJS)}
}

func buildK6Headers(headers map[string]string) string {
	if len(headers) == 0 {
		return ""
	}
	keys := make([]string, 0, len(headers))
	for k := range headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("'%s': '%s'", k, headers[k]))
	}
	return "{ headers: { " + strings.Join(parts, ", ") + " } }"
}

func buildK6Checks(assertions []schema.Assertion) []string {
	var checks []string
	for _, a := range assertions {
		line := k6AssertionLine(a)
		if line != "" {
			checks = append(checks, line)
		}
	}
	return checks
}

func k6AssertionLine(a schema.Assertion) string {
	// Normalise legacy "body.field" → "jsonpath $.field" so all operators are handled uniformly.
	if strings.HasPrefix(a.Target, "body.") {
		a.Target = "jsonpath $." + strings.TrimPrefix(a.Target, "body.")
	}
	switch {
	case a.Target == "status_code":
		switch a.Operator {
		case "eq":
			return fmt.Sprintf("'status is %v': (r) => r.status === %v,", a.Expected, a.Expected)
		case "ne":
			return fmt.Sprintf("'status is not %v': (r) => r.status !== %v,", a.Expected, a.Expected)
		case "gte":
			return fmt.Sprintf("'status >= %v': (r) => r.status >= %v,", a.Expected, a.Expected)
		case "lte":
			return fmt.Sprintf("'status <= %v': (r) => r.status <= %v,", a.Expected, a.Expected)
		case "gt":
			return fmt.Sprintf("'status > %v': (r) => r.status > %v,", a.Expected, a.Expected)
		case "lt":
			return fmt.Sprintf("'status < %v': (r) => r.status < %v,", a.Expected, a.Expected)
		}
	case strings.HasPrefix(a.Target, "jsonpath $."):
		field := strings.TrimPrefix(a.Target, "jsonpath $.")
		switch a.Operator {
		case "exists":
			return fmt.Sprintf("'jsonpath $.%s exists': (r) => r.json('%s') !== undefined,", field, field)
		case "not_exists":
			return fmt.Sprintf("'jsonpath $.%s not exists': (r) => r.json('%s') === undefined,", field, field)
		case "eq":
			return fmt.Sprintf("'jsonpath $.%s eq %v': (r) => r.json('%s') === %s,", field, a.Expected, field, k6JSValue(a.Expected))
		case "ne":
			if a.Expected == nil {
				return fmt.Sprintf("'jsonpath $.%s not null': (r) => r.json('%s') !== undefined,", field, field)
			}
			return fmt.Sprintf("'jsonpath $.%s ne %v': (r) => r.json('%s') !== %s,", field, a.Expected, field, k6JSValue(a.Expected))
		case "lt":
			return fmt.Sprintf("'jsonpath $.%s lt %v': (r) => r.json('%s') < %v,", field, a.Expected, field, a.Expected)
		case "gt":
			return fmt.Sprintf("'jsonpath $.%s gt %v': (r) => r.json('%s') > %v,", field, a.Expected, field, a.Expected)
		case "contains":
			return fmt.Sprintf("'jsonpath $.%s contains %v': (r) => String(r.json('%s')).includes(%s),", field, a.Expected, field, k6JSValue(a.Expected))
		case "matches":
			return fmt.Sprintf("'jsonpath $.%s matches %v': (r) => new RegExp(%s).test(String(r.json('%s'))),", field, a.Expected, k6JSValue(a.Expected), field)
		case "is_iso8601":
			return fmt.Sprintf("'jsonpath $.%s is_iso8601': (r) => !isNaN(Date.parse(r.json('%s'))),", field, field)
		case "is_uuid":
			return fmt.Sprintf("'jsonpath $.%s is_uuid': (r) => /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(r.json('%s')),", field, field)
		}
	case strings.HasPrefix(a.Target, "header "):
		name := strings.TrimPrefix(a.Target, "header ")
		switch a.Operator {
		case "eq":
			return fmt.Sprintf("'header %s eq %v': (r) => r.headers['%s'] === %s,", name, a.Expected, name, k6JSValue(a.Expected))
		case "ne":
			return fmt.Sprintf("'header %s ne %v': (r) => r.headers['%s'] !== %s,", name, a.Expected, name, k6JSValue(a.Expected))
		case "contains":
			return fmt.Sprintf("'header %s contains %v': (r) => String(r.headers['%s']).includes(%s),", name, a.Expected, name, k6JSValue(a.Expected))
		case "matches":
			return fmt.Sprintf("'header %s matches %v': (r) => new RegExp(%s).test(r.headers['%s']),", name, a.Expected, k6JSValue(a.Expected), name)
		}
	case a.Target == "duration_ms":
		switch a.Operator {
		case "lt":
			return fmt.Sprintf("'duration < %vms': (r) => r.timings.duration < %v,", a.Expected, a.Expected)
		case "gt":
			return fmt.Sprintf("'duration > %vms': (r) => r.timings.duration > %v,", a.Expected, a.Expected)
		}
	}
	return fmt.Sprintf("// unrendered: %s %s %v", a.Target, a.Operator, a.Expected)
}

func k6JSValue(v any) string {
	if v == nil {
		return "null"
	}
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
