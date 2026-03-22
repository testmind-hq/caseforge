// internal/methodology/owasp.go
package methodology

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/security"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// SecurityTechnique generates per-operation OWASP API Top 10 security test cases.
// Covers: API1, API2, API3, API4, API6, API7, API10.
// API5, API8, API9 are cross-operation and live in SecuritySpecTechnique (owasp_spec.go).
type SecurityTechnique struct {
	gen *datagen.Generator
}

func NewSecurityTechnique() *SecurityTechnique {
	return &SecurityTechnique{gen: datagen.NewGenerator(nil)}
}
func (t *SecurityTechnique) Name() string { return "owasp_api_top10" }

func (t *SecurityTechnique) Applies(op *spec.Operation) bool {
	return security.HasIDPathParam(op) || // API1
		security.IsAuthRequired(op) || // API2
		op.Method == "PATCH" || op.Method == "PUT" || // API3
		hasPaginationParam(op) || // API4
		((op.Method == "POST" || op.Method == "PUT") && op.RequestBody != nil) || // API6
		hasStringParam(op) || // API7
		hasSSRFParam(op) // API10
}

func (t *SecurityTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	var cases []schema.TestCase

	if security.HasIDPathParam(op) {
		cases = append(cases, buildAPI1Case(op))
	}
	if security.IsAuthRequired(op) {
		cases = append(cases, buildAPI2Case(op))
	}
	if op.Method == "PATCH" || op.Method == "PUT" {
		cases = append(cases, t.buildAPI3Case(op))
	}
	if hasPaginationParam(op) {
		cases = append(cases, buildAPI4Case(op))
	}
	if (op.Method == "POST" || op.Method == "PUT") && op.RequestBody != nil {
		cases = append(cases, t.buildAPI6Case(op))
	}
	if hasStringParam(op) {
		cases = append(cases, buildAPI7Cases(op)...)
	}
	if hasSSRFParam(op) {
		cases = append(cases, buildAPI10Case(op))
	}

	return cases, nil
}

// --- API1: BOLA ---

func buildAPI1Case(op *spec.Operation) schema.TestCase {
	// Replace ID path param with template variable for attacker-controlled ID
	path := op.Path
	for _, seg := range strings.Split(op.Path, "/") {
		if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
			name := strings.ToLower(seg[1 : len(seg)-1])
			if name == "id" || strings.HasSuffix(name, "id") {
				path = strings.ReplaceAll(op.Path, seg, "{{other_resource_id}}")
				break
			}
		}
	}
	step := schema.Step{
		ID: "step-1", Title: "access other user's resource",
		Type: "test", Method: op.Method, Path: path,
		Headers: map[string]string{},
		Assertions: []schema.Assertion{
			{Target: "status_code", Operator: "eq", Expected: 403},
		},
	}
	return owaspCase(op, "api1-bola",
		fmt.Sprintf("[OWASP-API1] %s %s — BOLA 越权访问", op.Method, op.Path),
		"路径含 ID 参数，需验证对象级授权：用合法 token 访问他人资源应返回 403",
		step)
}

// --- API2: Broken Authentication ---

func buildAPI2Case(op *spec.Operation) schema.TestCase {
	step := schema.Step{
		ID: "step-1", Title: "request without auth token",
		Type: "test", Method: op.Method, Path: op.Path,
		Headers: map[string]string{},
		Assertions: []schema.Assertion{
			{Target: "status_code", Operator: "eq", Expected: 401},
		},
	}
	return owaspCase(op, "api2-broken-auth",
		fmt.Sprintf("[OWASP-API2] %s %s — 认证机制失效", op.Method, op.Path),
		"接口声明了安全方案，去除 Authorization header 应返回 401",
		step)
}

// --- API3: BOPLA ---

func (t *SecurityTechnique) buildAPI3Case(op *spec.Operation) schema.TestCase {
	body := map[string]any{}
	if op.RequestBody != nil {
		if mt, ok := op.RequestBody.Content["application/json"]; ok && mt.Schema != nil {
			for k, v := range buildValidBody(t.gen, op) {
				body[k] = v
			}
		}
	}
	body["is_admin"] = true
	body["role"] = "admin"

	step := schema.Step{
		ID: "step-1", Title: "inject privileged fields in body",
		Type: "test", Method: op.Method, Path: op.Path,
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    body,
		Assertions: []schema.Assertion{
			{Target: "status_code", Operator: "eq", Expected: 200},
			{Target: "jsonpath $.is_admin", Operator: "ne", Expected: true},
			{Target: "jsonpath $.role", Operator: "ne", Expected: "admin"},
		},
	}
	return owaspCase(op, "api3-bopla",
		fmt.Sprintf("[OWASP-API3] %s %s — BOPLA 属性级越权", op.Method, op.Path),
		"PATCH/PUT 注入特权字段，响应中这些字段不应被修改或返回",
		step)
}

// --- API4: Resource Consumption ---

func buildAPI4Case(op *spec.Operation) schema.TestCase {
	// Build path with limit=99999 injected
	paramName := "limit"
	for _, p := range op.Parameters {
		if p.In == "query" {
			n := strings.ToLower(p.Name)
			if n == "limit" || n == "size" || n == "per_page" {
				paramName = p.Name
				break
			}
		}
	}
	path := op.Path + "?" + paramName + "=99999"

	step := schema.Step{
		ID: "step-1", Title: "request with extreme pagination value",
		Type: "test", Method: op.Method, Path: path,
		Headers: map[string]string{},
		Assertions: []schema.Assertion{
			{Target: "status_code", Operator: "eq", Expected: 200},
		},
	}
	return owaspCase(op, "api4-resource-consumption",
		fmt.Sprintf("[OWASP-API4] %s %s — 资源消耗无限制", op.Method, op.Path),
		"注入极大分页值 99999，服务器若返回 200 则存在性能风险（正确实现应返回 400）",
		step)
}

// --- API6: Mass Assignment ---

func (t *SecurityTechnique) buildAPI6Case(op *spec.Operation) schema.TestCase {
	body := map[string]any{}
	if b := buildValidBody(t.gen, op); b != nil {
		for k, v := range b {
			body[k] = v
		}
	}
	readonlyFields := map[string]any{"id": 99999, "createdAt": "2000-01-01T00:00:00Z", "updatedAt": "2000-01-01T00:00:00Z"}
	for k, v := range readonlyFields {
		body[k] = v
	}

	assertions := []schema.Assertion{{Target: "status_code", Operator: "eq", Expected: 201}}
	for field, val := range readonlyFields {
		assertions = append(assertions, schema.Assertion{
			Target: fmt.Sprintf("jsonpath $.%s", field), Operator: "ne", Expected: val,
		})
	}

	step := schema.Step{
		ID: "step-1", Title: "inject read-only fields in body",
		Type: "test", Method: op.Method, Path: op.Path,
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    body, Assertions: assertions,
	}
	return owaspCase(op, "api6-mass-assignment",
		fmt.Sprintf("[OWASP-API6] %s %s — 批量赋值", op.Method, op.Path),
		"注入只读字段 id/createdAt/updatedAt，响应中这些字段不应被接受并返回注入值",
		step)
}

// --- API7: Injection ---

var injectionPayloads = []struct {
	payload string
	label   string
}{
	{`"><script>alert(1)</script>`, "xss"},
	{`' OR 1=1--`, "sqli"},
	{`../../../etc/passwd`, "path-traversal"},
}

func buildAPI7Cases(op *spec.Operation) []schema.TestCase {
	// Find first string param (query, path, or body)
	paramName := ""
	paramIn := ""
	for _, p := range op.Parameters {
		if p.Schema != nil && p.Schema.Type == "string" {
			paramName = p.Name
			paramIn = p.In
			break
		}
	}
	if paramName == "" {
		// Try body
		if op.RequestBody != nil {
			if mt, ok := op.RequestBody.Content["application/json"]; ok && mt.Schema != nil {
				for name, s := range mt.Schema.Properties {
					if s.Type == "string" {
						paramName = name
						paramIn = "body"
						break
					}
				}
			}
		}
	}
	if paramName == "" {
		return nil
	}

	var cases []schema.TestCase
	for _, inj := range injectionPayloads {
		var path string
		var body map[string]any
		headers := map[string]string{}

		switch paramIn {
		case "query":
			path = op.Path + "?" + paramName + "=" + url.QueryEscape(inj.payload)
		case "path":
			path = strings.ReplaceAll(op.Path,
				fmt.Sprintf("{%s}", paramName),
				url.PathEscape(inj.payload))
		case "body":
			path = op.Path
			body = map[string]any{paramName: inj.payload}
			headers["Content-Type"] = "application/json"
		default:
			path = op.Path
		}
		if path == "" {
			path = op.Path
		}

		step := schema.Step{
			ID: "step-1", Title: fmt.Sprintf("inject %s payload", inj.label),
			Type: "test", Method: op.Method, Path: path,
			Headers: headers, Body: body,
			Assertions: []schema.Assertion{
				{Target: "status_code", Operator: "eq", Expected: 400},
			},
		}
		tc := owaspCase(op, "api7-injection",
			fmt.Sprintf("[OWASP-API7] %s %s — 注入 (%s)", op.Method, op.Path, inj.label),
			fmt.Sprintf("注入 %s payload，服务器应拒绝并返回 400", inj.label),
			step)
		cases = append(cases, tc)
	}
	return cases
}

// --- API10: SSRF ---

func buildAPI10Case(op *spec.Operation) schema.TestCase {
	// Find the URL-like param
	paramName := ""
	paramIn := ""
	for _, p := range op.Parameters {
		n := strings.ToLower(p.Name)
		if n == "url" || n == "callback" || n == "webhook" || n == "redirect" {
			paramName = p.Name
			paramIn = p.In
			break
		}
	}
	if paramName == "" && op.RequestBody != nil {
		if mt, ok := op.RequestBody.Content["application/json"]; ok && mt.Schema != nil {
			for name := range mt.Schema.Properties {
				n := strings.ToLower(name)
				if n == "url" || n == "callback" || n == "webhook" || n == "redirect" {
					paramName = name
					paramIn = "body"
					break
				}
			}
		}
	}

	var path string
	var body map[string]any
	headers := map[string]string{}
	const ssrfPayload = "http://127.0.0.1"

	switch paramIn {
	case "query":
		path = op.Path + "?" + paramName + "=" + url.QueryEscape(ssrfPayload)
	case "body":
		path = op.Path
		body = map[string]any{paramName: ssrfPayload}
		headers["Content-Type"] = "application/json"
	default:
		path = op.Path
	}

	step := schema.Step{
		ID: "step-1", Title: "inject internal URL for SSRF",
		Type: "test", Method: op.Method, Path: path,
		Headers: headers, Body: body,
		Assertions: []schema.Assertion{
			{Target: "status_code", Operator: "eq", Expected: 400},
		},
	}
	return owaspCase(op, "api10-ssrf",
		fmt.Sprintf("[OWASP-API10] %s %s — SSRF", op.Method, op.Path),
		"注入内网 URL http://127.0.0.1，服务器应验证并拒绝 (400)",
		step)
}

// --- helpers ---

func owaspCase(op *spec.Operation, tag, title, rationale string, step schema.Step) schema.TestCase {
	id := fmt.Sprintf("TC-%s", uuid.New().String()[:8])
	return schema.TestCase{
		Schema:  schema.SchemaBaseURL,
		Version: "1",
		ID:      id,
		Title:   title,
		Kind:    "single",
		Priority: "P0",
		Tags:    append([]string{"security", "owasp"}, tag),
		Source: schema.CaseSource{
			Technique: "owasp_api_top10",
			SpecPath:  fmt.Sprintf("%s %s", op.Method, op.Path),
			Rationale: rationale,
		},
		Steps:       []schema.Step{step},
		GeneratedAt: time.Now(),
	}
}

func hasPaginationParam(op *spec.Operation) bool {
	for _, p := range op.Parameters {
		n := strings.ToLower(p.Name)
		if p.In == "query" && (n == "limit" || n == "size" || n == "per_page") {
			return true
		}
	}
	return false
}

func hasStringParam(op *spec.Operation) bool {
	for _, p := range op.Parameters {
		if p.Schema != nil && p.Schema.Type == "string" {
			return true
		}
	}
	if op.RequestBody != nil {
		if mt, ok := op.RequestBody.Content["application/json"]; ok && mt.Schema != nil {
			for _, s := range mt.Schema.Properties {
				if s.Type == "string" {
					return true
				}
			}
		}
	}
	return false
}

func hasSSRFParam(op *spec.Operation) bool {
	for _, p := range op.Parameters {
		n := strings.ToLower(p.Name)
		if n == "url" || n == "callback" || n == "webhook" || n == "redirect" {
			return true
		}
	}
	if op.RequestBody != nil {
		if mt, ok := op.RequestBody.Content["application/json"]; ok && mt.Schema != nil {
			for name := range mt.Schema.Properties {
				n := strings.ToLower(name)
				if n == "url" || n == "callback" || n == "webhook" || n == "redirect" {
					return true
				}
			}
		}
	}
	return false
}
