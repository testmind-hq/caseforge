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
		fmt.Sprintf("[OWASP-API1] %s %s — BOLA unauthorized access", op.Method, op.Path),
		"Path contains an ID parameter; verify object-level authorization: accessing another user's resource with a valid token should return 403",
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
		fmt.Sprintf("[OWASP-API2] %s %s — broken authentication", op.Method, op.Path),
		"Endpoint declares a security scheme; removing the Authorization header should return 401",
		step)
}

// --- API3: BOPLA ---

func (t *SecurityTechnique) buildAPI3Case(op *spec.Operation) schema.TestCase {
	body := map[string]any{}
	for k, v := range buildValidBody(t.gen, op) {
		body[k] = v
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
		fmt.Sprintf("[OWASP-API3] %s %s — BOPLA property-level access", op.Method, op.Path),
		"PATCH/PUT with injected privileged fields; those fields must not be modified or reflected in the response",
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
		fmt.Sprintf("[OWASP-API4] %s %s — unrestricted resource consumption", op.Method, op.Path),
		"Inject extreme pagination value 99999; a 200 response indicates a performance risk (correct implementation should return 400)",
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

	expectedStatus := 201
	if strings.EqualFold(op.Method, "PUT") || strings.EqualFold(op.Method, "PATCH") {
		expectedStatus = 200
	}
	assertions := []schema.Assertion{{Target: "status_code", Operator: "eq", Expected: expectedStatus}}
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
		fmt.Sprintf("[OWASP-API6] %s %s — mass assignment", op.Method, op.Path),
		"Inject read-only fields id/createdAt/updatedAt; the response must not accept or reflect the injected values",
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
			fmt.Sprintf("[OWASP-API7] %s %s — injection (%s)", op.Method, op.Path, inj.label),
			fmt.Sprintf("Inject %s payload; server must reject and return 400", inj.label),
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
		"Inject internal URL http://127.0.0.1; server must validate and reject (400)",
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
