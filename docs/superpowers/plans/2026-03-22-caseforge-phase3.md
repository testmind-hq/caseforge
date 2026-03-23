# CaseForge Phase 3 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement OWASP API Top 10 security test generation, Lint L007-L012 rules with scoring, Postman Collection v2.1 renderer, and Spec Diff breaking-change detection.

**Architecture:** Four independent features share a common `internal/security/helpers.go` foundation (OWASP + lint). SecurityTechnique (per-op) and SecuritySpecTechnique (cross-op) follow the existing Technique/SpecTechnique pattern. Postman renderer implements the existing Renderer interface. Spec Diff is a standalone `internal/diff` package with a new `cmd/diff.go` command.

**Tech Stack:** Go 1.26.1, `github.com/google/uuid`, `github.com/spf13/cobra`, `github.com/getkin/kin-openapi/openapi3`, `github.com/stretchr/testify`

---

## Task 1: Foundation — `Operation.Security` field + `internal/security/helpers.go`

**Files:**
- Modify: `internal/spec/loader.go` — add `Security []string` field to `Operation` struct
- Modify: `internal/spec/parser.go` — parse `op.Security` in `convertOperation`
- Create: `internal/security/helpers.go`
- Create: `internal/security/helpers_test.go`

**Context:** `Operation` struct is defined in `loader.go`. `Summary` already exists there — add `Security` alongside it. In kin-openapi, `openapi3.Operation.Security` is `*openapi3.SecurityRequirements` which is `*[]openapi3.SecurityRequirement`; each `SecurityRequirement` is `map[string][]string` (scheme name → scopes). Extract just the scheme names.

- [ ] **Step 1: Add `Security []string` field to `Operation` in `loader.go`**

Open `internal/spec/loader.go`. In the `Operation` struct (lines 23-34), add after `Tags []string`:
```go
Security []string // names of security schemes declared on this operation
```

- [ ] **Step 2: Parse `Security` in `convertOperation` in `parser.go`**

In `internal/spec/parser.go`, in `convertOperation` (after the `Tags: op.Tags` assignment at line 62), add:
```go
if op.Security != nil {
    for _, req := range *op.Security {
        for schemeName := range req {
            o.Security = append(o.Security, schemeName)
        }
    }
}
```

- [ ] **Step 3: Write failing tests for helpers**

Create `internal/security/helpers_test.go`:
```go
package security_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/testmind-hq/caseforge/internal/security"
    "github.com/testmind-hq/caseforge/internal/spec"
)

func TestHasIDPathParam(t *testing.T) {
    assert.True(t, security.HasIDPathParam(&spec.Operation{Path: "/users/{userId}"}))
    assert.True(t, security.HasIDPathParam(&spec.Operation{Path: "/items/{id}"}))
    assert.False(t, security.HasIDPathParam(&spec.Operation{Path: "/users"}))
    assert.False(t, security.HasIDPathParam(&spec.Operation{Path: "/health"}))
}

func TestFindSensitiveFields(t *testing.T) {
    schema := &spec.Schema{Properties: map[string]*spec.Schema{
        "email":        {Type: "string"},
        "passwordHash": {Type: "string"},
        "accessToken":  {Type: "string"},
        "name":         {Type: "string"},
    }}
    fields := security.FindSensitiveFields(schema)
    assert.Contains(t, fields, "passwordHash")
    assert.Contains(t, fields, "accessToken")
    assert.NotContains(t, fields, "email")
    assert.NotContains(t, fields, "name")
}

func TestIsAuthRequired(t *testing.T) {
    assert.True(t, security.IsAuthRequired(&spec.Operation{Security: []string{"bearerAuth"}}))
    assert.False(t, security.IsAuthRequired(&spec.Operation{}))
}

func TestFindVersionedPaths(t *testing.T) {
    ops := []*spec.Operation{
        {Path: "/v1/users"},
        {Path: "/v2/users"},
        {Path: "/health"},
    }
    v1, v2 := security.FindVersionedPaths(ops)
    assert.NotEmpty(t, v1)
    assert.NotEmpty(t, v2)
}

func TestFindVersionedPaths_NoPairs(t *testing.T) {
    ops := []*spec.Operation{{Path: "/users"}, {Path: "/orders"}}
    v1, v2 := security.FindVersionedPaths(ops)
    assert.Empty(t, v1)
    assert.Empty(t, v2)
}
```

- [ ] **Step 4: Run tests to verify they fail**

```bash
cd /Users/yuchou/Github/yuchou87/caseforge
go test ./internal/security/... 2>&1 | head -20
```
Expected: build error ("no Go files in internal/security")

- [ ] **Step 5: Implement `internal/security/helpers.go`**

Create `internal/security/helpers.go`:
```go
// internal/security/helpers.go
package security

import (
    "strings"

    "github.com/testmind-hq/caseforge/internal/spec"
)

// sensitiveKeywords are substrings that flag a field as sensitive (case-insensitive).
var sensitiveKeywords = []string{
    "password", "secret", "token", "private_key", "access_token", "refresh_token",
}

// HasIDPathParam returns true if the operation path contains a path parameter
// whose name ends with "id" (case-insensitive), e.g. {id}, {userId}, {orderId}.
func HasIDPathParam(op *spec.Operation) bool {
    for _, seg := range strings.Split(op.Path, "/") {
        if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
            name := strings.ToLower(seg[1 : len(seg)-1])
            if name == "id" || strings.HasSuffix(name, "id") {
                return true
            }
        }
    }
    return false
}

// FindSensitiveFields returns field names in schema.Properties that contain
// a sensitive keyword (case-insensitive substring match).
func FindSensitiveFields(s *spec.Schema) []string {
    if s == nil {
        return nil
    }
    var found []string
    for name := range s.Properties {
        lower := strings.ToLower(name)
        for _, kw := range sensitiveKeywords {
            if strings.Contains(lower, kw) {
                found = append(found, name)
                break
            }
        }
    }
    return found
}

// IsAuthRequired returns true if the operation declares at least one security scheme.
func IsAuthRequired(op *spec.Operation) bool {
    return len(op.Security) > 0
}

// FindVersionedPaths returns two slices: paths containing "/v1/" (or starting with "/v1/")
// and paths containing "/v2/" (or similar next-version prefix).
// Returns nil slices when no versioned pair exists.
func FindVersionedPaths(ops []*spec.Operation) (v1Paths, v2Paths []string) {
    seen := map[string][]string{} // version prefix → paths
    for _, op := range ops {
        p := op.Path
        for _, ver := range []string{"/v1/", "/v2/", "/v3/"} {
            if strings.Contains(p, ver) || strings.HasPrefix(p, ver[1:]) {
                seen[ver] = append(seen[ver], p)
                break
            }
        }
    }
    if len(seen["/v1/"]) > 0 && len(seen["/v2/"]) > 0 {
        return seen["/v1/"], seen["/v2/"]
    }
    return nil, nil
}
```

- [ ] **Step 6: Run tests to verify they pass**

```bash
go test ./internal/security/... -v
```
Expected: all 5 tests PASS

- [ ] **Step 7: Verify spec parsing compiles**

```bash
go build ./...
```
Expected: clean compile

- [ ] **Step 8: Commit**

```bash
git add internal/spec/loader.go internal/spec/parser.go internal/security/
git commit -m "feat: add Operation.Security field and security helpers package"
```

---

## Task 2: `SecurityTechnique` — per-op OWASP rules (API1/2/3/4/6/7/10)

**Files:**
- Create: `internal/methodology/owasp.go`
- Create: `internal/methodology/owasp_test.go`

**Context:** Follows the same pattern as `equivalence.go` / `boundary.go`. Uses `buildValidBody` from `helpers.go`. All test cases get `priority: "P0"` and `tags: ["security", "owasp", "api<N>-<name>"]`. For assertion operators, use the existing `"ne"` and `"eq"` — **not** "notEqual". For `Expected: nil`, the Hurl renderer will emit `not exists`.

- [ ] **Step 1: Write failing tests**

Create `internal/methodology/owasp_test.go`:
```go
package methodology

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/testmind-hq/caseforge/internal/spec"
)

func securitySpec() *spec.ParsedSpec {
    return &spec.ParsedSpec{
        Operations: []*spec.Operation{
            {
                Method: "GET", Path: "/users/{userId}",
                Security:   []string{"bearerAuth"},
                Parameters: []*spec.Parameter{{Name: "userId", In: "path", Required: true, Schema: &spec.Schema{Type: "integer"}}},
                Responses:  map[string]*spec.Response{"200": {Description: "OK"}},
            },
            {
                Method: "PATCH", Path: "/users/{userId}",
                Security: []string{"bearerAuth"},
                RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
                    "application/json": {Schema: &spec.Schema{Type: "object",
                        Properties: map[string]*spec.Schema{"name": {Type: "string"}}}},
                }},
                Responses: map[string]*spec.Response{"200": {Description: "OK"}},
            },
            {
                Method: "GET", Path: "/items",
                Parameters: []*spec.Parameter{
                    {Name: "limit", In: "query", Schema: &spec.Schema{Type: "integer"}},
                    {Name: "q", In: "query", Schema: &spec.Schema{Type: "string"}},
                },
                Responses: map[string]*spec.Response{"200": {Description: "OK"}},
            },
        },
    }
}

func TestSecurityTechniqueAPI1_BOLA(t *testing.T) {
    st := NewSecurityTechnique()
    op := securitySpec().Operations[0] // GET /users/{userId}
    require.True(t, st.Applies(op))
    cases, err := st.Generate(op)
    require.NoError(t, err)
    var found bool
    for _, tc := range cases {
        for _, tag := range tc.Tags {
            if tag == "api1-bola" {
                found = true
                assert.Equal(t, "P0", tc.Priority)
                assert.Contains(t, tc.Steps[0].Path, "{{other_resource_id}}")
                assert.Equal(t, 403, tc.Steps[0].Assertions[0].Expected)
            }
        }
    }
    assert.True(t, found, "should generate API1 BOLA case")
}

func TestSecurityTechniqueAPI2_Auth(t *testing.T) {
    st := NewSecurityTechnique()
    op := securitySpec().Operations[0] // GET /users/{userId} has Security
    cases, err := st.Generate(op)
    require.NoError(t, err)
    var found bool
    for _, tc := range cases {
        for _, tag := range tc.Tags {
            if tag == "api2-broken-auth" {
                found = true
                _, hasAuth := tc.Steps[0].Headers["Authorization"]
                assert.False(t, hasAuth, "API2 case must not have Authorization header")
                assert.Equal(t, 401, tc.Steps[0].Assertions[0].Expected)
            }
        }
    }
    assert.True(t, found)
}

func TestSecurityTechniqueAPI3_BOPLA(t *testing.T) {
    st := NewSecurityTechnique()
    op := securitySpec().Operations[1] // PATCH /users/{userId}
    require.True(t, st.Applies(op))
    cases, err := st.Generate(op)
    require.NoError(t, err)
    var found bool
    for _, tc := range cases {
        for _, tag := range tc.Tags {
            if tag == "api3-bopla" {
                found = true
                body, ok := tc.Steps[0].Body.(map[string]any)
                require.True(t, ok)
                assert.Equal(t, true, body["is_admin"])
                assert.Equal(t, "admin", body["role"])
            }
        }
    }
    assert.True(t, found)
}

func TestSecurityTechniqueAPI4_Pagination(t *testing.T) {
    st := NewSecurityTechnique()
    op := securitySpec().Operations[2] // GET /items has limit param
    require.True(t, st.Applies(op))
    cases, err := st.Generate(op)
    require.NoError(t, err)
    var found bool
    for _, tc := range cases {
        for _, tag := range tc.Tags {
            if tag == "api4-resource-consumption" {
                found = true
                assert.Contains(t, tc.Steps[0].Path, "limit=99999")
            }
        }
    }
    assert.True(t, found)
}

func TestSecurityTechniqueAPI7_Injection(t *testing.T) {
    st := NewSecurityTechnique()
    op := securitySpec().Operations[2] // GET /items has string param "q"
    cases, err := st.Generate(op)
    require.NoError(t, err)
    var count int
    for _, tc := range cases {
        for _, tag := range tc.Tags {
            if tag == "api7-injection" {
                count++
            }
        }
    }
    assert.Equal(t, 3, count, "should generate 3 injection sub-cases (XSS, SQLi, path traversal)")
}

func TestSecurityTechniqueAPI10_SSRF(t *testing.T) {
    st := NewSecurityTechnique()
    op := &spec.Operation{
        Method: "POST", Path: "/webhooks",
        RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
            "application/json": {Schema: &spec.Schema{Type: "object",
                Properties: map[string]*spec.Schema{"url": {Type: "string"}}}},
        }},
        Responses: map[string]*spec.Response{"200": {Description: "OK"}},
    }
    require.True(t, st.Applies(op))
    cases, err := st.Generate(op)
    require.NoError(t, err)
    var found bool
    for _, tc := range cases {
        for _, tag := range tc.Tags {
            if tag == "api10-ssrf" {
                found = true
            }
        }
    }
    assert.True(t, found)
}

func TestSecurityTechniqueNotApplies(t *testing.T) {
    st := NewSecurityTechnique()
    op := &spec.Operation{
        Method: "GET", Path: "/health",
        Responses: map[string]*spec.Response{"200": {}},
    }
    // No ID param, no security, no body, no pagination, no string params
    assert.False(t, st.Applies(op))
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/methodology/... -run TestSecurity 2>&1 | head -10
```
Expected: compile error (SecurityTechnique undefined)

- [ ] **Step 3: Implement `internal/methodology/owasp.go`**

Create `internal/methodology/owasp.go`:
```go
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
func (t *SecurityTechnique) Name() string      { return "owasp_api_top10" }

func (t *SecurityTechnique) Applies(op *spec.Operation) bool {
    return security.HasIDPathParam(op) || // API1
        security.IsAuthRequired(op) || // API2
        op.Method == "PATCH" || op.Method == "PUT" || // API3
        hasPaginationParam(op) || // API4
        (( op.Method == "POST" || op.Method == "PUT") && op.RequestBody != nil) || // API6
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
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/methodology/... -run TestSecurity -v
```
Expected: all tests PASS

- [ ] **Step 5: Verify full test suite still passes**

```bash
go test ./... 2>&1 | tail -20
```
Expected: all PASS

- [ ] **Step 6: Commit**

```bash
git add internal/methodology/owasp.go internal/methodology/owasp_test.go
git commit -m "feat: implement SecurityTechnique for OWASP API1/2/3/4/6/7/10"
```

---

## Task 3: `SecuritySpecTechnique` — cross-op OWASP rules (API5/8/9)

**Files:**
- Create: `internal/methodology/owasp_spec.go`
- Create: `internal/methodology/owasp_spec_test.go`

**Context:** Implements `SpecTechnique` interface (same as `ChainTechnique`). API8 uses a local `seenPaths` map to generate one OPTIONS case per unique path. API5 looks for low-privilege paths (`/me`, `/profile`) and high-privilege paths (`admin` or `DELETE`). API9 uses `security.FindVersionedPaths`.

- [ ] **Step 1: Write failing tests**

Create `internal/methodology/owasp_spec_test.go`:
```go
package methodology

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/testmind-hq/caseforge/internal/spec"
)

func apiSpecWithRBAC() *spec.ParsedSpec {
    return &spec.ParsedSpec{
        Operations: []*spec.Operation{
            {Method: "GET", Path: "/users/me", Responses: map[string]*spec.Response{"200": {}}},
            {Method: "GET", Path: "/admin/users", Responses: map[string]*spec.Response{"200": {}}},
            {Method: "DELETE", Path: "/users/{id}", Responses: map[string]*spec.Response{"204": {}}},
        },
    }
}

func apiSpecWithVersions() *spec.ParsedSpec {
    return &spec.ParsedSpec{
        Operations: []*spec.Operation{
            {Method: "GET", Path: "/v1/users", Responses: map[string]*spec.Response{"200": {}}},
            {Method: "GET", Path: "/v2/users", Responses: map[string]*spec.Response{"200": {}}},
            {Method: "POST", Path: "/v2/orders", Responses: map[string]*spec.Response{"201": {}}},
        },
    }
}

func TestSecuritySpecAPI5_FunctionLevel(t *testing.T) {
    sst := NewSecuritySpecTechnique()
    cases, err := sst.Generate(apiSpecWithRBAC())
    require.NoError(t, err)
    var found bool
    for _, tc := range cases {
        for _, tag := range tc.Tags {
            if tag == "api5-function-level-auth" {
                found = true
                assert.Equal(t, "P0", tc.Priority)
                assert.Contains(t, tc.Steps[0].Headers["Authorization"], "{{user_token}}")
                assert.Equal(t, 403, tc.Steps[0].Assertions[0].Expected)
            }
        }
    }
    assert.True(t, found, "should generate API5 cases when low+high privilege paths exist")
}

func TestSecuritySpecAPI8_CORS(t *testing.T) {
    sst := NewSecuritySpecTechnique()
    ps := &spec.ParsedSpec{
        Operations: []*spec.Operation{
            {Method: "GET", Path: "/users", Responses: map[string]*spec.Response{"200": {}}},
            {Method: "POST", Path: "/users", Responses: map[string]*spec.Response{"201": {}}},
            {Method: "GET", Path: "/orders", Responses: map[string]*spec.Response{"200": {}}},
        },
    }
    cases, err := sst.Generate(ps)
    require.NoError(t, err)
    var corsCases []string
    for _, tc := range cases {
        for _, tag := range tc.Tags {
            if tag == "api8-cors" {
                corsCases = append(corsCases, tc.Steps[0].Path)
            }
        }
    }
    // /users appears twice (GET+POST) but should produce only 1 CORS case
    assert.Equal(t, 2, len(corsCases), "one CORS case per unique path")
    assert.NotEqual(t, corsCases[0], corsCases[1])
}

func TestSecuritySpecAPI9_AssetManagement(t *testing.T) {
    sst := NewSecuritySpecTechnique()
    cases, err := sst.Generate(apiSpecWithVersions())
    require.NoError(t, err)
    var found bool
    for _, tc := range cases {
        for _, tag := range tc.Tags {
            if tag == "api9-asset-management" {
                found = true
                assert.Equal(t, 404, tc.Steps[0].Assertions[0].Expected)
                assert.Contains(t, tc.Steps[0].Path, "/v1/")
            }
        }
    }
    assert.True(t, found)
}

func TestSecuritySpecNoAPI5_WithoutLowPriv(t *testing.T) {
    sst := NewSecuritySpecTechnique()
    ps := &spec.ParsedSpec{
        Operations: []*spec.Operation{
            {Method: "GET", Path: "/users", Responses: map[string]*spec.Response{"200": {}}},
        },
    }
    cases, err := sst.Generate(ps)
    require.NoError(t, err)
    for _, tc := range cases {
        for _, tag := range tc.Tags {
            assert.NotEqual(t, "api5-function-level-auth", tag)
        }
    }
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/methodology/... -run TestSecuritySpec 2>&1 | head -10
```
Expected: compile error

- [ ] **Step 3: Implement `internal/methodology/owasp_spec.go`**

Create `internal/methodology/owasp_spec.go`:
```go
// internal/methodology/owasp_spec.go
package methodology

import (
    "fmt"
    "strings"
    "time"

    "github.com/google/uuid"
    "github.com/testmind-hq/caseforge/internal/output/schema"
    "github.com/testmind-hq/caseforge/internal/security"
    "github.com/testmind-hq/caseforge/internal/spec"
)

// SecuritySpecTechnique generates cross-operation OWASP cases.
// Covers: API5 (function-level auth), API8 (CORS), API9 (asset management).
type SecuritySpecTechnique struct{}

func NewSecuritySpecTechnique() *SecuritySpecTechnique { return &SecuritySpecTechnique{} }
func (t *SecuritySpecTechnique) Name() string          { return "owasp_api_top10_spec" }

func (t *SecuritySpecTechnique) Generate(s *spec.ParsedSpec) ([]schema.TestCase, error) {
    var cases []schema.TestCase
    cases = append(cases, buildAPI5Cases(s)...)
    cases = append(cases, buildAPI8Cases(s)...)
    cases = append(cases, buildAPI9Cases(s)...)
    return cases, nil
}

// API5: Function-Level Authorization
func buildAPI5Cases(s *spec.ParsedSpec) []schema.TestCase {
    var lowPrivPaths, highPrivPaths []*spec.Operation
    for _, op := range s.Operations {
        p := op.Path
        if strings.Contains(p, "/me") || strings.Contains(p, "/profile") {
            lowPrivPaths = append(lowPrivPaths, op)
        } else if strings.Contains(p, "admin") || op.Method == "DELETE" {
            highPrivPaths = append(highPrivPaths, op)
        }
    }
    if len(lowPrivPaths) == 0 || len(highPrivPaths) == 0 {
        return nil
    }

    var cases []schema.TestCase
    for _, op := range highPrivPaths {
        step := schema.Step{
            ID: "step-1", Title: "access privileged endpoint with regular user token",
            Type: "test", Method: op.Method, Path: op.Path,
            Headers: map[string]string{"Authorization": "Bearer {{user_token}}"},
            Assertions: []schema.Assertion{
                {Target: "status_code", Operator: "eq", Expected: 403},
            },
        }
        id := fmt.Sprintf("TC-%s", uuid.New().String()[:8])
        cases = append(cases, schema.TestCase{
            Schema: schema.SchemaBaseURL, Version: "1", ID: id,
            Title:    fmt.Sprintf("[OWASP-API5] %s %s — 功能级授权缺失", op.Method, op.Path),
            Kind:     "single", Priority: "P0",
            Tags:     []string{"security", "owasp", "api5-function-level-auth"},
            Source:   schema.CaseSource{Technique: "owasp_api_top10_spec", SpecPath: fmt.Sprintf("%s %s", op.Method, op.Path), Rationale: "用普通用户 token 访问高权限接口应返回 403"},
            Steps:       []schema.Step{step},
            GeneratedAt: time.Now(),
        })
    }
    return cases
}

// API8: CORS Misconfiguration (one OPTIONS case per unique path)
func buildAPI8Cases(s *spec.ParsedSpec) []schema.TestCase {
    seenPaths := map[string]bool{}
    var cases []schema.TestCase
    for _, op := range s.Operations {
        if seenPaths[op.Path] {
            continue
        }
        seenPaths[op.Path] = true
        step := schema.Step{
            ID: "step-1", Title: "OPTIONS preflight request",
            Type: "test", Method: "OPTIONS", Path: op.Path,
            Headers: map[string]string{"Origin": "https://evil.example.com"},
            Assertions: []schema.Assertion{
                {Target: "header Access-Control-Allow-Origin", Operator: "ne", Expected: "*"},
            },
        }
        id := fmt.Sprintf("TC-%s", uuid.New().String()[:8])
        cases = append(cases, schema.TestCase{
            Schema: schema.SchemaBaseURL, Version: "1", ID: id,
            Title:    fmt.Sprintf("[OWASP-API8] OPTIONS %s — CORS 安全配置", op.Path),
            Kind:     "single", Priority: "P0",
            Tags:     []string{"security", "owasp", "api8-cors"},
            Source:   schema.CaseSource{Technique: "owasp_api_top10_spec", SpecPath: fmt.Sprintf("OPTIONS %s", op.Path), Rationale: "CORS 响应头 Access-Control-Allow-Origin 不应为 *"},
            Steps:       []schema.Step{step},
            GeneratedAt: time.Now(),
        })
    }
    return cases
}

// API9: Improper Asset Management (old versioned paths should return 404)
func buildAPI9Cases(s *spec.ParsedSpec) []schema.TestCase {
    v1Paths, _ := security.FindVersionedPaths(s.Operations)
    if len(v1Paths) == 0 {
        return nil
    }

    var cases []schema.TestCase
    seen := map[string]bool{}
    for _, op := range s.Operations {
        if !seen[op.Path] && containsPath(v1Paths, op.Path) {
            seen[op.Path] = true
            step := schema.Step{
                ID: "step-1", Title: "access deprecated v1 endpoint",
                Type: "test", Method: op.Method, Path: op.Path,
                Headers: map[string]string{},
                Assertions: []schema.Assertion{
                    {Target: "status_code", Operator: "eq", Expected: 404},
                },
            }
            id := fmt.Sprintf("TC-%s", uuid.New().String()[:8])
            cases = append(cases, schema.TestCase{
                Schema: schema.SchemaBaseURL, Version: "1", ID: id,
                Title:    fmt.Sprintf("[OWASP-API9] %s %s — 旧版本 API 应已下线", op.Method, op.Path),
                Kind:     "single", Priority: "P0",
                Tags:     []string{"security", "owasp", "api9-asset-management"},
                Source:   schema.CaseSource{Technique: "owasp_api_top10_spec", SpecPath: fmt.Sprintf("%s %s", op.Method, op.Path), Rationale: "Spec 同时含有 v1/v2 路径，旧版本路径应已下线返回 404"},
                Steps:       []schema.Step{step},
                GeneratedAt: time.Now(),
            })
        }
    }
    return cases
}

func containsPath(paths []string, p string) bool {
    for _, v := range paths {
        if v == p {
            return true
        }
    }
    return false
}
```

- [ ] **Step 4: Run tests**

```bash
go test ./internal/methodology/... -run TestSecuritySpec -v
```
Expected: all PASS

- [ ] **Step 5: Run full suite**

```bash
go test ./...
```
Expected: all PASS

- [ ] **Step 6: Commit**

```bash
git add internal/methodology/owasp_spec.go internal/methodology/owasp_spec_test.go
git commit -m "feat: implement SecuritySpecTechnique for OWASP API5/8/9"
```

---

## Task 4: Register OWASP in Engine + Hurl renderer extension

**Files:**
- Modify: `cmd/gen.go` — add `SecurityTechnique` + `SecuritySpecTechnique`
- Modify: `internal/output/render/hurl.go` — extend `renderAssertion` for `jsonpath` and `header` targets
- Modify: `internal/output/schema/model.go` — update `Assertion.Target` comment

**Context:** In `cmd/gen.go`, add `methodology.NewSecurityTechnique()` to `NewEngine(...)` variadic args (line 91-98), and call `engine.AddSpecTechnique(methodology.NewSecuritySpecTechnique())` after `engine.AddSpecTechnique(methodology.NewChainTechnique())`. In `renderAssertion`, add two new cases before the `default:` branch.

- [ ] **Step 1: Write a test for Hurl rendering of new assertion targets**

Add to `internal/output/render/hurl_test.go` (find the file and append):
```go
func TestRenderAssertionJSONPath(t *testing.T) {
    tc := schema.TestCase{
        Steps: []schema.Step{{
            Method: "GET", Path: "/users/1",
            Assertions: []schema.Assertion{
                {Target: "status_code", Operator: "eq", Expected: 200},
                {Target: "jsonpath $.is_admin", Operator: "ne", Expected: true},
                {Target: "jsonpath $.role", Operator: "ne", Expected: nil},
            },
        }},
    }
    r := NewHurlRenderer("")
    dir := t.TempDir()
    require.NoError(t, r.Render([]schema.TestCase{tc}, dir))
    files, _ := filepath.Glob(filepath.Join(dir, "*.hurl"))
    require.Len(t, files, 1)
    content, _ := os.ReadFile(files[0])
    assert.Contains(t, string(content), `jsonpath "$.is_admin" != true`)
    assert.Contains(t, string(content), `jsonpath "$.role" not exists`)
}

func TestRenderAssertionHeaderTarget(t *testing.T) {
    tc := schema.TestCase{
        Steps: []schema.Step{{
            Method: "OPTIONS", Path: "/users",
            Assertions: []schema.Assertion{
                {Target: "status_code", Operator: "eq", Expected: 200},
                {Target: "header Access-Control-Allow-Origin", Operator: "ne", Expected: "*"},
            },
        }},
    }
    r := NewHurlRenderer("")
    dir := t.TempDir()
    require.NoError(t, r.Render([]schema.TestCase{tc}, dir))
    files, _ := filepath.Glob(filepath.Join(dir, "*.hurl"))
    require.Len(t, files, 1)
    content, _ := os.ReadFile(files[0])
    assert.Contains(t, string(content), `header "Access-Control-Allow-Origin" != "*"`)
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/output/render/... -run TestRenderAssertion -v 2>&1 | tail -20
```
Expected: FAIL (assertions containing `# unrendered assertion`)

- [ ] **Step 3: Extend `renderAssertion` in `hurl.go`**

In `internal/output/render/hurl.go`, replace the `renderAssertion` function (lines 141-163):
```go
func renderAssertion(a schema.Assertion) string {
    switch {
    case a.Target == "duration_ms":
        if a.Operator == "lt" {
            return fmt.Sprintf("duration < %v\n", a.Expected)
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
        case "exists":
            return fmt.Sprintf("jsonpath %q exists\n", expr)
        case "contains":
            return fmt.Sprintf("jsonpath %q contains %s\n", expr, formatHurlValue(a.Expected))
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
        }
    case strings.HasPrefix(a.Target, "body."):
        // Legacy target format — delegate to jsonpath
        field := strings.TrimPrefix(a.Target, "body.")
        switch a.Operator {
        case "eq":
            return fmt.Sprintf("jsonpath \"$.%s\" == %s\n", field, formatHurlValue(a.Expected))
        case "exists":
            return fmt.Sprintf("jsonpath \"$.%s\" exists\n", field)
        case "contains":
            return fmt.Sprintf("jsonpath \"$.%s\" contains %s\n", field, formatHurlValue(a.Expected))
        }
    }
    return fmt.Sprintf("# unrendered assertion: %s %s %v\n", a.Target, a.Operator, a.Expected)
}
```

- [ ] **Step 4: Update `Assertion.Target` comment in `model.go`**

In `internal/output/schema/model.go`, find the `Assertion` struct (line ~38) and update the comment:
```go
type Assertion struct {
    Target   string `json:"target"`   // "status_code"|"jsonpath $.<field>"|"header <HeaderName>"|"duration_ms"
    Operator string `json:"operator"` // "eq"|"ne"|"lt"|"gt"|"contains"|"matches"|"exists"
    Expected any    `json:"expected"`
}
```

- [ ] **Step 5: Run tests**

```bash
go test ./internal/output/render/... -v
```
Expected: all PASS including new tests

- [ ] **Step 6: Register OWASP techniques in `cmd/gen.go`**

In `cmd/gen.go`, replace the engine construction block (lines 91-99):
```go
engine := methodology.NewEngine(provider,
    methodology.NewEquivalenceTechnique(),
    methodology.NewBoundaryTechnique(),
    methodology.NewDecisionTechnique(),
    methodology.NewStateTechnique(),
    methodology.NewIdempotentTechnique(),
    methodology.NewPairwiseTechnique(),
    methodology.NewSecurityTechnique(),
)
engine.AddSpecTechnique(methodology.NewChainTechnique())
engine.AddSpecTechnique(methodology.NewSecuritySpecTechnique())
```

- [ ] **Step 7: Build**

```bash
go build ./...
```
Expected: clean compile

- [ ] **Step 8: Run full test suite**

```bash
go test ./...
```
Expected: all PASS

- [ ] **Step 9: Commit**

```bash
git add cmd/gen.go internal/output/render/hurl.go internal/output/schema/model.go
git commit -m "feat: register OWASP techniques, extend hurl renderer for jsonpath/header targets"
```

---

## Task 5: Lint L007-L010 (design consistency rules)

**Files:**
- Create: `internal/lint/consistency.go`
- Create: `internal/lint/consistency_test.go`

**Context:** Each rule is a struct implementing `LintRule` interface (ID(), Severity(), Check()). All 4 rules are registered via `init()` using the existing `register()` function in `rules.go`. L008 checks only top-level request body properties. L010 checks all 4xx/5xx response schemas.

- [ ] **Step 1: Write failing tests**

Create `internal/lint/consistency_test.go`:
```go
package lint

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/testmind-hq/caseforge/internal/spec"
)

func TestL007_VerbInPath(t *testing.T) {
    ps := &spec.ParsedSpec{Operations: []*spec.Operation{
        {Method: "POST", Path: "/createUser", Responses: map[string]*spec.Response{"200": {}}},
        {Method: "GET", Path: "/users", Responses: map[string]*spec.Response{"200": {}}},
    }}
    rule := &ruleL007{}
    issues := rule.Check(ps)
    assert.Len(t, issues, 1)
    assert.Equal(t, "L007", issues[0].RuleID)
    assert.Equal(t, "warning", issues[0].Severity)
    assert.Contains(t, issues[0].Message, "create")
}

func TestL008_InconsistentNaming(t *testing.T) {
    ps := &spec.ParsedSpec{Operations: []*spec.Operation{
        {Method: "POST", Path: "/users",
            RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
                "application/json": {Schema: &spec.Schema{Properties: map[string]*spec.Schema{
                    "userId": {Type: "string"},
                }}},
            }},
            Parameters: []*spec.Parameter{{Name: "user_id", In: "query", Schema: &spec.Schema{Type: "string"}}},
            Responses:  map[string]*spec.Response{"200": {}},
        },
    }}
    rule := &ruleL008{}
    issues := rule.Check(ps)
    assert.Len(t, issues, 1)
    assert.Equal(t, "L008", issues[0].RuleID)
}

func TestL008_ConsistentNaming_NoIssue(t *testing.T) {
    ps := &spec.ParsedSpec{Operations: []*spec.Operation{
        {Method: "GET", Path: "/users",
            Parameters: []*spec.Parameter{
                {Name: "userId", In: "query", Schema: &spec.Schema{Type: "string"}},
                {Name: "orderId", In: "query", Schema: &spec.Schema{Type: "string"}},
            },
            Responses: map[string]*spec.Response{"200": {}},
        },
    }}
    rule := &ruleL008{}
    assert.Empty(t, rule.Check(ps))
}

func TestL009_InconsistentPagination(t *testing.T) {
    ps := &spec.ParsedSpec{Operations: []*spec.Operation{
        {Method: "GET", Path: "/users",
            Parameters: []*spec.Parameter{{Name: "page", In: "query", Schema: &spec.Schema{Type: "integer"}}},
            Responses:  map[string]*spec.Response{"200": {}},
        },
        {Method: "GET", Path: "/orders",
            Parameters: []*spec.Parameter{{Name: "offset", In: "query", Schema: &spec.Schema{Type: "integer"}}},
            Responses:  map[string]*spec.Response{"200": {}},
        },
    }}
    rule := &ruleL009{}
    issues := rule.Check(ps)
    assert.Len(t, issues, 1)
    assert.Equal(t, "L009", issues[0].RuleID)
}

func TestL010_InconsistentErrorSchema(t *testing.T) {
    ps := &spec.ParsedSpec{Operations: []*spec.Operation{
        {Method: "POST", Path: "/users", Responses: map[string]*spec.Response{
            "400": {Content: map[string]*spec.MediaType{
                "application/json": {Schema: &spec.Schema{Properties: map[string]*spec.Schema{"error": {Type: "string"}}}},
            }},
        }},
        {Method: "POST", Path: "/orders", Responses: map[string]*spec.Response{
            "400": {Content: map[string]*spec.MediaType{
                "application/json": {Schema: &spec.Schema{Properties: map[string]*spec.Schema{"message": {Type: "string"}}}},
            }},
        }},
    }}
    rule := &ruleL010{}
    issues := rule.Check(ps)
    assert.Len(t, issues, 1)
    assert.Equal(t, "L010", issues[0].RuleID)
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/lint/... -run 'TestL00[7-9]|TestL010' 2>&1 | head -10
```
Expected: compile error (ruleL007 undefined)

- [ ] **Step 3: Implement `internal/lint/consistency.go`**

Create `internal/lint/consistency.go`:
```go
// internal/lint/consistency.go
package lint

import (
    "fmt"
    "sort"
    "strconv"
    "strings"

    "github.com/testmind-hq/caseforge/internal/spec"
)

func init() {
    register(&ruleL007{})
    register(&ruleL008{})
    register(&ruleL009{})
    register(&ruleL010{})
}

var verbSegments = []string{"get", "create", "update", "delete", "list", "fetch", "add", "remove"}

// L007: verb in path segment
type ruleL007 struct{}

func (r *ruleL007) ID() string       { return "L007" }
func (r *ruleL007) Severity() string { return "warning" }
func (r *ruleL007) Check(ps *spec.ParsedSpec) []LintIssue {
    var issues []LintIssue
    for _, op := range ps.Operations {
        for _, seg := range strings.Split(strings.Trim(op.Path, "/"), "/") {
            lower := strings.ToLower(seg)
            for _, verb := range verbSegments {
                if strings.HasPrefix(lower, verb) {
                    issues = append(issues, LintIssue{
                        RuleID:   "L007",
                        Severity: "warning",
                        Message:  fmt.Sprintf("verb %q found in path segment", verb),
                        Path:     fmt.Sprintf("%s %s", op.Method, op.Path),
                    })
                    break
                }
            }
        }
    }
    return issues
}

// L008: inconsistent naming style (camelCase vs snake_case)
type ruleL008 struct{}

func (r *ruleL008) ID() string       { return "L008" }
func (r *ruleL008) Severity() string { return "warning" }
func (r *ruleL008) Check(ps *spec.ParsedSpec) []LintIssue {
    var allNames []string
    for _, op := range ps.Operations {
        for _, p := range op.Parameters {
            allNames = append(allNames, p.Name)
        }
        if op.RequestBody != nil {
            if mt, ok := op.RequestBody.Content["application/json"]; ok && mt.Schema != nil {
                for name := range mt.Schema.Properties {
                    allNames = append(allNames, name)
                }
            }
        }
    }
    hasCamel, hasSnake := false, false
    for _, name := range allNames {
        if strings.Contains(name, "_") {
            hasSnake = true
        } else if name != strings.ToLower(name) {
            hasCamel = true
        }
    }
    if hasCamel && hasSnake {
        return []LintIssue{{
            RuleID:   "L008",
            Severity: "warning",
            Message:  "mixed naming styles: camelCase and snake_case both present",
            Path:     "spec",
        }}
    }
    return nil
}

// L009: inconsistent pagination style
type ruleL009 struct{}

func (r *ruleL009) ID() string       { return "L009" }
func (r *ruleL009) Severity() string { return "warning" }
func (r *ruleL009) Check(ps *spec.ParsedSpec) []LintIssue {
    pageStyle, offsetStyle := false, false
    pageStyleOps, offsetStyleOps := []string{}, []string{}
    for _, op := range ps.Operations {
        for _, p := range op.Parameters {
            if p.In != "query" {
                continue
            }
            n := strings.ToLower(p.Name)
            if n == "page" || n == "size" {
                pageStyle = true
                pageStyleOps = append(pageStyleOps, fmt.Sprintf("%s %s", op.Method, op.Path))
            }
            if n == "offset" || n == "limit" {
                offsetStyle = true
                offsetStyleOps = append(offsetStyleOps, fmt.Sprintf("%s %s", op.Method, op.Path))
            }
        }
    }
    if pageStyle && offsetStyle {
        return []LintIssue{{
            RuleID:   "L009",
            Severity: "warning",
            Message:  fmt.Sprintf("mixed pagination styles: page/size in [%s] and offset/limit in [%s]", strings.Join(pageStyleOps, ", "), strings.Join(offsetStyleOps, ", ")),
            Path:     "spec",
        }}
    }
    return nil
}

// L010: inconsistent error response schema
type ruleL010 struct{}

func (r *ruleL010) ID() string       { return "L010" }
func (r *ruleL010) Severity() string { return "warning" }
func (r *ruleL010) Check(ps *spec.ParsedSpec) []LintIssue {
    seen := map[string]bool{} // sorted field key → true
    for _, op := range ps.Operations {
        for code, resp := range op.Responses {
            n, err := strconv.Atoi(code)
            if err != nil || n < 400 {
                continue
            }
            if mt, ok := resp.Content["application/json"]; ok && mt.Schema != nil {
                fields := make([]string, 0, len(mt.Schema.Properties))
                for name := range mt.Schema.Properties {
                    fields = append(fields, name)
                }
                sort.Strings(fields)
                seen[strings.Join(fields, ",")] = true
            }
        }
    }
    if len(seen) >= 2 {
        var structures []string
        for k := range seen {
            structures = append(structures, "{"+k+"}")
        }
        sort.Strings(structures)
        return []LintIssue{{
            RuleID:   "L010",
            Severity: "warning",
            Message:  fmt.Sprintf("inconsistent error response schemas: %s", strings.Join(structures, " vs ")),
            Path:     "spec",
        }}
    }
    return nil
}
```

- [ ] **Step 4: Run tests**

```bash
go test ./internal/lint/... -v
```
Expected: all PASS

- [ ] **Step 5: Commit**

```bash
git add internal/lint/consistency.go internal/lint/consistency_test.go
git commit -m "feat: add lint rules L007-L010 (design consistency)"
```

---

## Task 6: Lint L011-L012 (security rules)

**Files:**
- Create: `internal/lint/security_rules.go`
- Create: `internal/lint/security_rules_test.go`

- [ ] **Step 1: Write failing tests**

Create `internal/lint/security_rules_test.go`:
```go
package lint

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/testmind-hq/caseforge/internal/spec"
)

func TestL011_MissingSecurityScheme(t *testing.T) {
    ps := &spec.ParsedSpec{Operations: []*spec.Operation{
        {Method: "POST", Path: "/users", Security: nil, Responses: map[string]*spec.Response{"201": {}}},
        {Method: "GET", Path: "/health", Security: nil, Responses: map[string]*spec.Response{"200": {}}},      // GET excluded
        {Method: "POST", Path: "/public/login", Security: nil, Responses: map[string]*spec.Response{"200": {}}}, // /public excluded
        {Method: "DELETE", Path: "/users/{id}", Security: []string{"bearerAuth"}, Responses: map[string]*spec.Response{"204": {}}}, // has security
    }}
    rule := &ruleL011{}
    issues := rule.Check(ps)
    assert.Len(t, issues, 1)
    assert.Equal(t, "L011", issues[0].RuleID)
    assert.Equal(t, "error", issues[0].Severity)
    assert.Contains(t, issues[0].Path, "POST /users")
}

func TestL011_ExcludedPaths(t *testing.T) {
    ps := &spec.ParsedSpec{Operations: []*spec.Operation{
        {Method: "POST", Path: "/health/check", Responses: map[string]*spec.Response{"200": {}}},
        {Method: "POST", Path: "/auth/login", Responses: map[string]*spec.Response{"200": {}}},
        {Method: "POST", Path: "/users/register", Responses: map[string]*spec.Response{"201": {}}},
        {Method: "POST", Path: "/public/verify", Responses: map[string]*spec.Response{"200": {}}},
    }}
    rule := &ruleL011{}
    assert.Empty(t, rule.Check(ps))
}

func TestL012_SensitiveFieldExposed(t *testing.T) {
    ps := &spec.ParsedSpec{Operations: []*spec.Operation{
        {Method: "GET", Path: "/users/{id}",
            Responses: map[string]*spec.Response{
                "200": {Content: map[string]*spec.MediaType{
                    "application/json": {Schema: &spec.Schema{Properties: map[string]*spec.Schema{
                        "id":           {Type: "integer"},
                        "passwordHash": {Type: "string"},
                        "name":         {Type: "string"},
                    }}},
                }},
            },
        },
    }}
    rule := &ruleL012{}
    issues := rule.Check(ps)
    assert.Len(t, issues, 1)
    assert.Equal(t, "L012", issues[0].RuleID)
    assert.Equal(t, "error", issues[0].Severity)
    assert.Contains(t, issues[0].Message, "passwordHash")
}

func TestL012_NoSensitiveFields(t *testing.T) {
    ps := &spec.ParsedSpec{Operations: []*spec.Operation{
        {Method: "GET", Path: "/users",
            Responses: map[string]*spec.Response{
                "200": {Content: map[string]*spec.MediaType{
                    "application/json": {Schema: &spec.Schema{Properties: map[string]*spec.Schema{
                        "id": {Type: "integer"}, "name": {Type: "string"},
                    }}},
                }},
            },
        },
    }}
    rule := &ruleL012{}
    assert.Empty(t, rule.Check(ps))
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/lint/... -run 'TestL011|TestL012' 2>&1 | head -10
```
Expected: compile error

- [ ] **Step 3: Implement `internal/lint/security_rules.go`**

Create `internal/lint/security_rules.go`:
```go
// internal/lint/security_rules.go
package lint

import (
    "fmt"
    "strconv"
    "strings"

    "github.com/testmind-hq/caseforge/internal/security"
    "github.com/testmind-hq/caseforge/internal/spec"
)

func init() {
    register(&ruleL011{})
    register(&ruleL012{})
}

var excludedPathSubstrings = []string{"/public", "/health", "/login", "/register"}

// L011: non-GET operation missing security scheme declaration
type ruleL011 struct{}

func (r *ruleL011) ID() string       { return "L011" }
func (r *ruleL011) Severity() string { return "error" }
func (r *ruleL011) Check(ps *spec.ParsedSpec) []LintIssue {
    var issues []LintIssue
    for _, op := range ps.Operations {
        if op.Method == "GET" {
            continue
        }
        excluded := false
        for _, sub := range excludedPathSubstrings {
            if strings.Contains(op.Path, sub) {
                excluded = true
                break
            }
        }
        if excluded {
            continue
        }
        if len(op.Security) == 0 {
            issues = append(issues, LintIssue{
                RuleID:   "L011",
                Severity: "error",
                Message:  "non-GET operation has no security scheme declared",
                Path:     fmt.Sprintf("%s %s", op.Method, op.Path),
            })
        }
    }
    return issues
}

// L012: sensitive field exposed in 2xx response schema
type ruleL012 struct{}

func (r *ruleL012) ID() string       { return "L012" }
func (r *ruleL012) Severity() string { return "error" }
func (r *ruleL012) Check(ps *spec.ParsedSpec) []LintIssue {
    var issues []LintIssue
    for _, op := range ps.Operations {
        for code, resp := range op.Responses {
            n, err := strconv.Atoi(code)
            if err != nil || n < 200 || n >= 300 {
                continue
            }
            if mt, ok := resp.Content["application/json"]; ok && mt.Schema != nil {
                for _, fieldName := range security.FindSensitiveFields(mt.Schema) {
                    issues = append(issues, LintIssue{
                        RuleID:   "L012",
                        Severity: "error",
                        Message:  fmt.Sprintf("sensitive field %q exposed in %s response", fieldName, code),
                        Path:     fmt.Sprintf("%s %s", op.Method, op.Path),
                    })
                }
            }
        }
    }
    return issues
}
```

- [ ] **Step 4: Run tests**

```bash
go test ./internal/lint/... -v
```
Expected: all PASS

- [ ] **Step 5: Commit**

```bash
git add internal/lint/security_rules.go internal/lint/security_rules_test.go
git commit -m "feat: add lint rules L011-L012 (security: missing auth, sensitive fields)"
```

---

## Task 7: Lint scoring + `--min-score` CLI flag

**Files:**
- Create: `internal/lint/score.go`
- Create: `internal/lint/score_test.go`
- Modify: `cmd/lint.go` — add `--min-score` flag + score output

- [ ] **Step 1: Write failing tests**

Create `internal/lint/score_test.go`:
```go
package lint

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestScore_NoIssues(t *testing.T) {
    assert.Equal(t, 100, Score(nil))
}

func TestScore_Errors(t *testing.T) {
    issues := []LintIssue{
        {Severity: "error"},
        {Severity: "error"},
    }
    assert.Equal(t, 80, Score(issues))
}

func TestScore_Warnings(t *testing.T) {
    issues := []LintIssue{
        {Severity: "warning"},
        {Severity: "warning"},
        {Severity: "warning"},
    }
    assert.Equal(t, 91, Score(issues))
}

func TestScore_Mixed(t *testing.T) {
    issues := []LintIssue{
        {Severity: "error"},
        {Severity: "warning"},
        {Severity: "warning"},
    }
    assert.Equal(t, 84, Score(issues))
}

func TestScore_ClampedAtZero(t *testing.T) {
    var issues []LintIssue
    for i := 0; i < 15; i++ {
        issues = append(issues, LintIssue{Severity: "error"})
    }
    assert.Equal(t, 0, Score(issues))
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/lint/... -run TestScore 2>&1 | head -5
```
Expected: compile error (Score undefined)

- [ ] **Step 3: Implement `internal/lint/score.go`**

Create `internal/lint/score.go`:
```go
// internal/lint/score.go
package lint

// Score computes a spec quality score from 0 to 100.
// Each error subtracts 10 points; each warning subtracts 3 points.
func Score(issues []LintIssue) int {
    base := 100
    for _, issue := range issues {
        switch issue.Severity {
        case "error":
            base -= 10
        case "warning":
            base -= 3
        }
    }
    if base < 0 {
        return 0
    }
    return base
}
```

- [ ] **Step 4: Run score tests**

```bash
go test ./internal/lint/... -run TestScore -v
```
Expected: all 5 PASS

- [ ] **Step 5: Update `cmd/lint.go`**

Replace `cmd/lint.go` content with:
```go
// cmd/lint.go
package cmd

import (
    "fmt"
    "os"

    "github.com/fatih/color"
    "github.com/spf13/cobra"
    "github.com/testmind-hq/caseforge/internal/config"
    "github.com/testmind-hq/caseforge/internal/lint"
    "github.com/testmind-hq/caseforge/internal/spec"
)

var lintCmd = &cobra.Command{
    Use:   "lint",
    Short: "Lint an OpenAPI spec for quality issues",
    RunE:  runLint,
}

var (
    lintSpec     string
    lintMinScore int
)

func init() {
    rootCmd.AddCommand(lintCmd)
    lintCmd.Flags().StringVar(&lintSpec, "spec", "", "OpenAPI spec file or URL (required)")
    lintCmd.Flags().IntVar(&lintMinScore, "min-score", 0, "Fail if spec score is below this threshold (0 = disabled)")
    _ = lintCmd.MarkFlagRequired("spec")
}

func runLint(cmd *cobra.Command, args []string) error {
    cfg, err := config.Load()
    if err != nil {
        return fmt.Errorf("loading config: %w", err)
    }

    loader := spec.NewLoader()
    parsedSpec, err := loader.Load(lintSpec)
    if err != nil {
        fmt.Fprintf(os.Stderr, "✗ Failed to parse spec: %v\n", err)
        os.Exit(2)
    }

    issues := lint.RunAll(parsedSpec)
    score := lint.Score(issues)

    if len(issues) == 0 {
        color.Green("✓ No lint issues found")
    } else {
        hasError := false
        for _, iss := range issues {
            switch iss.Severity {
            case "error":
                color.Red("  ✗ [%s] %s: %s", iss.RuleID, iss.Path, iss.Message)
                hasError = true
            case "warning":
                color.Yellow("  ⚠ [%s] %s: %s", iss.RuleID, iss.Path, iss.Message)
            }
        }
        errCount := 0
        warnCount := 0
        for _, iss := range issues {
            if iss.Severity == "error" {
                errCount++
            } else {
                warnCount++
            }
        }
        fmt.Fprintf(os.Stderr, "\nSpec Score: %d/100  (%d errors, %d warnings)\n", score, errCount, warnCount)

        shouldFail := hasError
        if cfg.Lint.FailOn == "warning" {
            shouldFail = len(issues) > 0
        }
        if !shouldFail && lintMinScore > 0 && score < lintMinScore {
            fmt.Fprintf(os.Stderr, "exit code 1 (score %d < min-score %d)\n", score, lintMinScore)
            shouldFail = true
        }
        if shouldFail {
            os.Exit(3)
        }
    }
    return nil
}
```

- [ ] **Step 6: Build and run tests**

```bash
go build ./... && go test ./...
```
Expected: clean compile, all PASS

- [ ] **Step 7: Commit**

```bash
git add internal/lint/score.go internal/lint/score_test.go cmd/lint.go
git commit -m "feat: add lint scoring and --min-score flag"
```

---

## Task 8: Postman Collection v2.1 Renderer

**Files:**
- Create: `internal/output/render/postman.go`
- Create: `internal/output/render/postman_test.go`
- Modify: `cmd/gen.go` — uncomment/add `postman` case in format switch

- [ ] **Step 1: Write failing tests**

Create `internal/output/render/postman_test.go`:
```go
package render

import (
    "encoding/json"
    "os"
    "path/filepath"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/testmind-hq/caseforge/internal/output/schema"
)

func singleCase() schema.TestCase {
    return schema.TestCase{
        ID: "TC-0001", Title: "GET users list", Kind: "single", Priority: "P1",
        Steps: []schema.Step{{
            ID: "step-1", Method: "GET", Path: "/users",
            Assertions: []schema.Assertion{
                {Target: "status_code", Operator: "eq", Expected: 200},
            },
        }},
    }
}

func chainCase() schema.TestCase {
    return schema.TestCase{
        ID: "TC-0002", Title: "CRUD chain: /users", Kind: "chain", Priority: "P1",
        Steps: []schema.Step{
            {
                ID: "step-setup", Type: "setup", Method: "POST", Path: "/users",
                Headers: map[string]string{"Content-Type": "application/json"},
                Body: map[string]any{"name": "Alice"},
                Assertions: []schema.Assertion{{Target: "status_code", Operator: "eq", Expected: 201}},
                Captures: []schema.Capture{{Name: "userId", From: "jsonpath $.id"}},
            },
            {
                ID: "step-test", Type: "test", Method: "GET", Path: "/users/{{userId}}",
                Assertions: []schema.Assertion{{Target: "status_code", Operator: "eq", Expected: 200}},
            },
        },
    }
}

func TestPostmanRendererFormat(t *testing.T) {
    r := NewPostmanRenderer()
    assert.Equal(t, "postman", r.Format())
}

func TestPostmanRendererCreatesCollectionJSON(t *testing.T) {
    r := NewPostmanRenderer()
    dir := t.TempDir()
    err := r.Render([]schema.TestCase{singleCase()}, dir)
    require.NoError(t, err)
    _, statErr := os.Stat(filepath.Join(dir, "collection.json"))
    assert.NoError(t, statErr)
}

func TestPostmanCollectionTopLevel(t *testing.T) {
    r := NewPostmanRenderer()
    dir := t.TempDir()
    require.NoError(t, r.Render([]schema.TestCase{singleCase()}, dir))

    data, _ := os.ReadFile(filepath.Join(dir, "collection.json"))
    var coll map[string]any
    require.NoError(t, json.Unmarshal(data, &coll))

    info := coll["info"].(map[string]any)
    assert.Equal(t, "https://schema.getpostman.com/json/collection/v2.1.0/collection.json", info["schema"])
    assert.NotEmpty(t, info["_postman_id"])

    vars := coll["variable"].([]any)
    assert.Len(t, vars, 1)
    baseURL := vars[0].(map[string]any)
    assert.Equal(t, "base_url", baseURL["key"])
}

func TestPostmanSingleCaseItem(t *testing.T) {
    r := NewPostmanRenderer()
    dir := t.TempDir()
    require.NoError(t, r.Render([]schema.TestCase{singleCase()}, dir))

    data, _ := os.ReadFile(filepath.Join(dir, "collection.json"))
    var coll map[string]any
    require.NoError(t, json.Unmarshal(data, &coll))

    items := coll["item"].([]any)
    require.Len(t, items, 1)
    item := items[0].(map[string]any)
    req := item["request"].(map[string]any)
    assert.Equal(t, "GET", req["method"])

    urlObj := req["url"].(map[string]any)
    assert.Contains(t, urlObj["raw"].(string), "{{base_url}}")
}

func TestPostmanChainCaseIsFolder(t *testing.T) {
    r := NewPostmanRenderer()
    dir := t.TempDir()
    require.NoError(t, r.Render([]schema.TestCase{chainCase()}, dir))

    data, _ := os.ReadFile(filepath.Join(dir, "collection.json"))
    var coll map[string]any
    require.NoError(t, json.Unmarshal(data, &coll))

    items := coll["item"].([]any)
    require.Len(t, items, 1)
    folder := items[0].(map[string]any)
    // chain case becomes a folder with sub-items
    subItems := folder["item"].([]any)
    assert.Len(t, subItems, 2, "chain case has 2 steps")
}

func TestPostmanCaptureInSetupStep(t *testing.T) {
    r := NewPostmanRenderer()
    dir := t.TempDir()
    require.NoError(t, r.Render([]schema.TestCase{chainCase()}, dir))

    data, _ := os.ReadFile(filepath.Join(dir, "collection.json"))
    content := string(data)
    assert.Contains(t, content, "pm.environment.set")
    assert.Contains(t, content, "userId")
}

func TestPostmanStatusCodeTestScript(t *testing.T) {
    r := NewPostmanRenderer()
    dir := t.TempDir()
    require.NoError(t, r.Render([]schema.TestCase{singleCase()}, dir))

    data, _ := os.ReadFile(filepath.Join(dir, "collection.json"))
    content := string(data)
    assert.Contains(t, content, "pm.response.to.have.status(200)")
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/output/render/... -run TestPostman 2>&1 | head -10
```
Expected: compile error (PostmanRenderer undefined)

- [ ] **Step 3: Implement `internal/output/render/postman.go`**

Create `internal/output/render/postman.go`:
```go
// internal/output/render/postman.go
package render

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/google/uuid"
    "github.com/testmind-hq/caseforge/internal/output/schema"
)

// PostmanRenderer renders test cases as a Postman Collection v2.1 JSON file.
type PostmanRenderer struct{}

func NewPostmanRenderer() *PostmanRenderer { return &PostmanRenderer{} }
func (r *PostmanRenderer) Format() string  { return "postman" }

func (r *PostmanRenderer) Render(cases []schema.TestCase, outDir string) error {
    if err := os.MkdirAll(outDir, 0755); err != nil {
        return fmt.Errorf("creating output dir: %w", err)
    }

    collection := map[string]any{
        "info": map[string]any{
            "name":        "CaseForge Generated",
            "_postman_id": uuid.New().String(),
            "schema":      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
        },
        "variable": []any{
            map[string]any{"key": "base_url", "value": "http://localhost", "type": "string"},
        },
        "item": buildPostmanItems(cases),
    }

    data, err := json.MarshalIndent(collection, "", "  ")
    if err != nil {
        return fmt.Errorf("marshaling collection: %w", err)
    }
    return os.WriteFile(filepath.Join(outDir, "collection.json"), data, 0644)
}

func buildPostmanItems(cases []schema.TestCase) []any {
    var items []any
    for _, tc := range cases {
        if tc.Kind == "chain" {
            items = append(items, buildPostmanFolder(tc))
        } else {
            if len(tc.Steps) > 0 {
                items = append(items, buildPostmanRequest(tc.Title, tc.Steps[0]))
            }
        }
    }
    return items
}

func buildPostmanFolder(tc schema.TestCase) map[string]any {
    var subItems []any
    for _, step := range tc.Steps {
        subItems = append(subItems, buildPostmanRequest(step.Title, step))
    }
    return map[string]any{
        "name": tc.Title,
        "item": subItems,
    }
}

func buildPostmanRequest(name string, step schema.Step) map[string]any {
    req := map[string]any{
        "method": step.Method,
        "url":    buildPostmanURL(step.Path),
        "header": buildPostmanHeaders(step.Headers),
    }
    if step.Body != nil {
        data, _ := json.Marshal(step.Body)
        req["body"] = map[string]any{
            "mode": "raw",
            "raw":  string(data),
            "options": map[string]any{
                "raw": map[string]any{"language": "json"},
            },
        }
    }

    item := map[string]any{
        "name":    name,
        "request": req,
    }

    script := buildTestScript(step)
    if script != "" {
        item["event"] = []any{
            map[string]any{
                "listen": "test",
                "script": map[string]any{
                    "type": "text/javascript",
                    "exec": strings.Split(script, "\n"),
                },
            },
        }
    }
    return item
}

func buildPostmanURL(path string) map[string]any {
    rawURL := "{{base_url}}" + path
    segments := strings.Split(strings.Trim(path, "/"), "/")
    pathParts := make([]any, 0, len(segments))
    for _, seg := range segments {
        if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
            // OpenAPI {param} → Postman :param
            pathParts = append(pathParts, ":"+seg[1:len(seg)-1])
        } else {
            pathParts = append(pathParts, seg)
        }
    }
    return map[string]any{
        "raw":  rawURL,
        "host": []any{"{{base_url}}"},
        "path": pathParts,
    }
}

func buildPostmanHeaders(headers map[string]string) []any {
    var result []any
    for k, v := range headers {
        result = append(result, map[string]any{"key": k, "value": v})
    }
    return result
}

func buildTestScript(step schema.Step) string {
    var lines []string

    for _, a := range step.Assertions {
        switch {
        case a.Target == "status_code":
            code := 200
            if f, ok := a.Expected.(float64); ok {
                code = int(f)
            } else if i, ok := a.Expected.(int); ok {
                code = i
            }
            lines = append(lines,
                fmt.Sprintf("pm.test(\"status is %d\", function () {", code),
                fmt.Sprintf("    pm.response.to.have.status(%d);", code),
                "});",
            )
        case strings.HasPrefix(a.Target, "jsonpath "):
            expr := strings.TrimPrefix(a.Target, "jsonpath ")
            fieldPath := strings.TrimPrefix(expr, "$.")
            fieldPath = strings.ReplaceAll(fieldPath, ".", ".")
            valJS := postmanJSValue(a.Expected)
            switch a.Operator {
            case "eq":
                lines = append(lines,
                    fmt.Sprintf("pm.test(%q, function () {", a.Target+" eq "+fmt.Sprint(a.Expected)),
                    "    var jsonData = pm.response.json();",
                    fmt.Sprintf("    pm.expect(jsonData.%s).to.eql(%s);", fieldPath, valJS),
                    "});",
                )
            case "ne":
                if a.Expected == nil {
                    lines = append(lines,
                        fmt.Sprintf("pm.test(%q, function () {", a.Target+" not exists"),
                        "    var jsonData = pm.response.json();",
                        fmt.Sprintf("    pm.expect(jsonData.%s).to.not.exist;", fieldPath),
                        "});",
                    )
                } else {
                    lines = append(lines,
                        fmt.Sprintf("pm.test(%q, function () {", a.Target+" ne "+fmt.Sprint(a.Expected)),
                        "    var jsonData = pm.response.json();",
                        fmt.Sprintf("    pm.expect(jsonData.%s).to.not.eql(%s);", fieldPath, valJS),
                        "});",
                    )
                }
            case "contains":
                lines = append(lines,
                    fmt.Sprintf("pm.test(%q, function () {", a.Target+" contains "+fmt.Sprint(a.Expected)),
                    "    var jsonData = pm.response.json();",
                    fmt.Sprintf("    pm.expect(String(jsonData.%s)).to.include(%s);", fieldPath, valJS),
                    "});",
                )
            }
        case strings.HasPrefix(a.Target, "header "):
            headerName := strings.TrimPrefix(a.Target, "header ")
            valJS := postmanJSValue(a.Expected)
            if a.Operator == "ne" {
                lines = append(lines,
                    fmt.Sprintf("pm.test(%q, function () {", a.Target+" ne "+fmt.Sprint(a.Expected)),
                    fmt.Sprintf("    pm.expect(pm.response.headers.get(%q)).to.not.eql(%s);", headerName, valJS),
                    "});",
                )
            }
        }
    }

    // Add capture scripts
    if len(step.Captures) > 0 {
        lines = append(lines, "var jsonData = pm.response.json();")
        for _, cap := range step.Captures {
            if strings.HasPrefix(cap.From, "jsonpath $.") {
                fieldPath := strings.TrimPrefix(cap.From, "jsonpath $.")
                lines = append(lines,
                    fmt.Sprintf("pm.environment.set(%q, jsonData.%s);", cap.Name, fieldPath),
                )
            }
        }
    }

    return strings.Join(lines, "\n")
}

func postmanJSValue(v any) string {
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
```

- [ ] **Step 4: Enable `postman` case in `cmd/gen.go`**

In `cmd/gen.go`, add to the format switch (before `default:`):
```go
case "postman":
    renderer = render.NewPostmanRenderer()
```

Also update the format flag description in `init()` to include `postman`:
```go
genCmd.Flags().StringVar(&genFormat, "format", "hurl", "Output format: hurl|markdown|csv|postman")
```

- [ ] **Step 5: Run tests**

```bash
go test ./internal/output/render/... -run TestPostman -v
```
Expected: all 7 tests PASS

- [ ] **Step 6: Run full suite**

```bash
go build ./... && go test ./...
```
Expected: all PASS

- [ ] **Step 7: Commit**

```bash
git add internal/output/render/postman.go internal/output/render/postman_test.go cmd/gen.go
git commit -m "feat: add Postman Collection v2.1 renderer"
```

---

## Task 9: `internal/diff/diff.go` — Spec diff core logic

**Files:**
- Create: `internal/diff/diff.go`
- Create: `internal/diff/diff_test.go`
- Create: `testdata/petstore_v1.yaml`
- Create: `testdata/petstore_v2.yaml`

- [ ] **Step 1: Write failing tests**

Create `internal/diff/diff_test.go`:
```go
package diff

import (
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/testmind-hq/caseforge/internal/spec"
)

func v1Spec() *spec.ParsedSpec {
    return &spec.ParsedSpec{
        Operations: []*spec.Operation{
            {Method: "GET", Path: "/users", Parameters: []*spec.Parameter{{Name: "limit", In: "query", Required: false, Schema: &spec.Schema{Type: "integer"}}}, Responses: map[string]*spec.Response{"200": {Content: map[string]*spec.MediaType{"application/json": {Schema: &spec.Schema{Properties: map[string]*spec.Schema{"id": {Type: "integer"}, "email": {Type: "string"}}}}}}}},
            {Method: "DELETE", Path: "/users/{id}", Responses: map[string]*spec.Response{"204": {}}},
            {Method: "POST", Path: "/orders", RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{"application/json": {Schema: &spec.Schema{Properties: map[string]*spec.Schema{"customerId": {Type: "integer"}}}}}}, Responses: map[string]*spec.Response{"201": {}}},
        },
    }
}

func v2Spec() *spec.ParsedSpec {
    return &spec.ParsedSpec{
        Operations: []*spec.Operation{
            // GET /users: response field "email" removed, new optional param added
            {Method: "GET", Path: "/users", Parameters: []*spec.Parameter{{Name: "limit", In: "query", Required: false, Schema: &spec.Schema{Type: "integer"}}, {Name: "includeDeleted", In: "query", Required: false, Schema: &spec.Schema{Type: "boolean"}}}, Responses: map[string]*spec.Response{"200": {Content: map[string]*spec.MediaType{"application/json": {Schema: &spec.Schema{Properties: map[string]*spec.Schema{"id": {Type: "integer"}}}}}}}},
            // DELETE /users/{id} removed
            // POST /orders: customerId type changed
            {Method: "POST", Path: "/orders", RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{"application/json": {Schema: &spec.Schema{Properties: map[string]*spec.Schema{"customerId": {Type: "string"}}}}}}, Responses: map[string]*spec.Response{"201": {}}},
            // New endpoint
            {Method: "POST", Path: "/payments", Responses: map[string]*spec.Response{"201": {}}},
        },
    }
}

func TestDiffEndpointRemoved(t *testing.T) {
    result := Diff(v1Spec(), v2Spec())
    var found bool
    for _, c := range result.Changes {
        if c.Kind == Breaking && c.Method == "DELETE" && c.Path == "/users/{id}" {
            found = true
            assert.Contains(t, c.Description, "endpoint removed")
        }
    }
    assert.True(t, found, "should detect endpoint removal as BREAKING")
}

func TestDiffResponseFieldRemoved(t *testing.T) {
    result := Diff(v1Spec(), v2Spec())
    var found bool
    for _, c := range result.Changes {
        if c.Kind == Breaking && c.Path == "/users" && strings.Contains(c.Description, "email") {
            found = true
        }
    }
    assert.True(t, found, "response field deletion should be BREAKING")
}

func TestDiffParamTypeChanged(t *testing.T) {
    result := Diff(v1Spec(), v2Spec())
    var found bool
    for _, c := range result.Changes {
        if c.Kind == Breaking && c.Path == "/orders" && strings.Contains(c.Description, "customerId") {
            found = true
        }
    }
    assert.True(t, found, "param type change should be BREAKING")
}

func TestDiffNewEndpoint(t *testing.T) {
    result := Diff(v1Spec(), v2Spec())
    var found bool
    for _, c := range result.Changes {
        if c.Kind == NonBreaking && c.Path == "/payments" {
            found = true
        }
    }
    assert.True(t, found, "new endpoint should be NON_BREAKING")
}

func TestDiffNewOptionalParam(t *testing.T) {
    result := Diff(v1Spec(), v2Spec())
    var found bool
    for _, c := range result.Changes {
        if c.Kind == NonBreaking && c.Path == "/users" && strings.Contains(c.Description, "includeDeleted") {
            found = true
        }
    }
    assert.True(t, found, "new optional param should be NON_BREAKING")
}
```

Add `"strings"` to the import in the test file.

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/diff/... 2>&1 | head -10
```
Expected: compile error (package not found)

- [ ] **Step 3: Implement `internal/diff/diff.go`**

Create `internal/diff/diff.go`:
```go
// internal/diff/diff.go
package diff

import (
    "fmt"
    "sort"
    "strings"

    "github.com/testmind-hq/caseforge/internal/spec"
)

// ChangeKind classifies the impact of a spec change on API consumers.
type ChangeKind string

const (
    Breaking            ChangeKind = "BREAKING"
    PotentiallyBreaking ChangeKind = "POTENTIALLY_BREAKING"
    NonBreaking         ChangeKind = "NON_BREAKING"
)

// Change describes a single detected difference between two spec versions.
type Change struct {
    Kind        ChangeKind
    Method      string
    Path        string
    Location    string // "requestBody", "response.200", "param.limit", etc.
    Description string
}

// DiffResult holds all detected changes.
type DiffResult struct {
    Changes []Change
}

// HasBreaking returns true if any change is Breaking or PotentiallyBreaking.
func (r DiffResult) HasBreaking() bool {
    for _, c := range r.Changes {
        if c.Kind == Breaking || c.Kind == PotentiallyBreaking {
            return true
        }
    }
    return false
}

// Diff compares oldSpec against newSpec and returns all detected changes.
func Diff(oldSpec, newSpec *spec.ParsedSpec) DiffResult {
    var changes []Change

    oldOps := opMap(oldSpec)
    newOps := opMap(newSpec)

    // Detect removed endpoints, param changes, response field changes
    for key, oldOp := range oldOps {
        newOp, exists := newOps[key]
        if !exists {
            changes = append(changes, Change{
                Kind: Breaking, Method: oldOp.Method, Path: oldOp.Path,
                Description: "endpoint removed",
            })
            continue
        }
        changes = append(changes, diffParams(oldOp, newOp)...)
        changes = append(changes, diffResponseFields(oldOp, newOp)...)
        changes = append(changes, diffRequestBodyFields(oldOp, newOp)...)
    }

    // Detect path renames (BREAKING) and new endpoints (NON_BREAKING)
    uniqueOldPaths := sortedUniqueNewPaths(oldSpec, newSpec)
    uniqueNewPaths := uniqueNewPathsSet(oldSpec, newSpec)
    renamedNewPaths := map[string]bool{}

    for _, oldPath := range uniqueOldPaths {
        best, diffCount := findBestRename(oldPath, uniqueNewPaths)
        if best != "" && diffCount > 0 {
            changes = append(changes, Change{
                Kind: Breaking, Path: oldPath,
                Description: fmt.Sprintf("path renamed: %s → %s", oldPath, best),
            })
            renamedNewPaths[best] = true
        }
    }

    // New endpoints (not renames)
    for _, newPath := range sortedSlice(uniqueNewPaths) {
        if renamedNewPaths[newPath] {
            continue
        }
        // find all ops on this new path
        for _, newOp := range newSpec.Operations {
            if newOp.Path == newPath {
                changes = append(changes, Change{
                    Kind: NonBreaking, Method: newOp.Method, Path: newPath,
                    Description: "new endpoint",
                })
            }
        }
    }

    return DiffResult{Changes: changes}
}

// opMap builds a "METHOD /path" → *Operation map.
func opMap(ps *spec.ParsedSpec) map[string]*spec.Operation {
    m := map[string]*spec.Operation{}
    for _, op := range ps.Operations {
        m[op.Method+" "+op.Path] = op
    }
    return m
}

func diffParams(oldOp, newOp *spec.Operation) []Change {
    var changes []Change
    oldParams := paramMap(oldOp.Parameters)
    newParams := paramMap(newOp.Parameters)

    for name, oldP := range oldParams {
        newP, exists := newParams[name]
        if !exists {
            // Parameter removed — could break clients sending it; treat as POTENTIALLY_BREAKING
            changes = append(changes, Change{
                Kind: PotentiallyBreaking, Method: oldOp.Method, Path: oldOp.Path,
                Location:    "param." + name,
                Description: fmt.Sprintf("parameter %q removed", name),
            })
            continue
        }
        // Required flag: optional → required is BREAKING
        if !oldP.Required && newP.Required {
            changes = append(changes, Change{
                Kind: Breaking, Method: oldOp.Method, Path: oldOp.Path,
                Location:    "param." + name,
                Description: fmt.Sprintf("parameter %q changed from optional to required", name),
            })
        }
        // Type change is BREAKING
        if oldP.Schema != nil && newP.Schema != nil && oldP.Schema.Type != newP.Schema.Type {
            changes = append(changes, Change{
                Kind: Breaking, Method: oldOp.Method, Path: oldOp.Path,
                Location:    "param." + name,
                Description: fmt.Sprintf("parameter %q type changed: %s → %s", name, oldP.Schema.Type, newP.Schema.Type),
            })
        }
    }

    // New optional params are NON_BREAKING; new required params are POTENTIALLY_BREAKING
    for name, newP := range newParams {
        if _, exists := oldParams[name]; !exists {
            kind := NonBreaking
            desc := fmt.Sprintf("new optional parameter %q added", name)
            if newP.Required {
                kind = PotentiallyBreaking
                desc = fmt.Sprintf("new required parameter %q added", name)
            }
            changes = append(changes, Change{
                Kind: kind, Method: oldOp.Method, Path: oldOp.Path,
                Location: "param." + name, Description: desc,
            })
        }
    }
    return changes
}

func diffResponseFields(oldOp, newOp *spec.Operation) []Change {
    var changes []Change
    for code, oldResp := range oldOp.Responses {
        newResp, exists := newOp.Responses[code]
        if !exists {
            continue
        }
        oldSchema := responseJSONSchema(oldResp)
        newSchema := responseJSONSchema(newResp)
        if oldSchema == nil || newSchema == nil {
            continue
        }
        for fieldName, oldField := range oldSchema.Properties {
            newField, exists := newSchema.Properties[fieldName]
            if !exists {
                changes = append(changes, Change{
                    Kind: Breaking, Method: oldOp.Method, Path: oldOp.Path,
                    Location:    fmt.Sprintf("response.%s", code),
                    Description: fmt.Sprintf("response field %q removed", fieldName),
                })
                continue
            }
            if oldField.Type != newField.Type && newField.Type != "" {
                // Type widening (integer→number) is POTENTIALLY_BREAKING; other changes are BREAKING
                kind := Breaking
                if oldField.Type == "integer" && newField.Type == "number" {
                    kind = PotentiallyBreaking
                }
                changes = append(changes, Change{
                    Kind: kind, Method: oldOp.Method, Path: oldOp.Path,
                    Location:    fmt.Sprintf("response.%s", code),
                    Description: fmt.Sprintf("response field %q type changed: %s → %s", fieldName, oldField.Type, newField.Type),
                })
            }
        }
        // New response fields are NON_BREAKING
        for fieldName := range newSchema.Properties {
            if _, exists := oldSchema.Properties[fieldName]; !exists {
                changes = append(changes, Change{
                    Kind: NonBreaking, Method: oldOp.Method, Path: oldOp.Path,
                    Location:    fmt.Sprintf("response.%s", code),
                    Description: fmt.Sprintf("new response field %q added", fieldName),
                })
            }
        }
    }
    return changes
}

func diffRequestBodyFields(oldOp, newOp *spec.Operation) []Change {
    var changes []Change
    oldSchema := requestJSONSchema(oldOp)
    newSchema := requestJSONSchema(newOp)
    if oldSchema == nil || newSchema == nil {
        return nil
    }
    // New required fields in request body are POTENTIALLY_BREAKING
    oldRequired := stringSet(oldSchema.Required)
    for _, req := range newSchema.Required {
        if !oldRequired[req] {
            changes = append(changes, Change{
                Kind: PotentiallyBreaking, Method: oldOp.Method, Path: oldOp.Path,
                Location:    "requestBody",
                Description: fmt.Sprintf("new required request body field %q added", req),
            })
        }
    }
    // Field type changes in request body are BREAKING
    for fieldName, oldField := range oldSchema.Properties {
        if newField, exists := newSchema.Properties[fieldName]; exists {
            if oldField.Type != newField.Type && newField.Type != "" {
                changes = append(changes, Change{
                    Kind: Breaking, Method: oldOp.Method, Path: oldOp.Path,
                    Location:    "requestBody",
                    Description: fmt.Sprintf("request body field %q type changed: %s → %s", fieldName, oldField.Type, newField.Type),
                })
            }
        }
    }
    return changes
}

// --- path rename helpers ---

func sortedUniqueNewPaths(oldSpec, newSpec *spec.ParsedSpec) []string {
    oldPaths := pathSet(oldSpec)
    newPaths := pathSet(newSpec)
    var result []string
    for p := range oldPaths {
        if !newPaths[p] {
            result = append(result, p)
        }
    }
    sort.Strings(result)
    return result
}

func uniqueNewPathsSet(oldSpec, newSpec *spec.ParsedSpec) map[string]bool {
    oldPaths := pathSet(oldSpec)
    newPaths := pathSet(newSpec)
    result := map[string]bool{}
    for p := range newPaths {
        if !oldPaths[p] {
            result[p] = true
        }
    }
    return result
}

func findBestRename(oldPath string, candidates map[string]bool) (string, int) {
    oldSegs := splitPath(oldPath)
    best := ""
    bestDiff := -1
    for cand := range candidates {
        newSegs := splitPath(cand)
        if len(oldSegs) != len(newSegs) {
            continue
        }
        paramSame := true
        diffCount := 0
        for i := range oldSegs {
            oldIsParam := isParamSeg(oldSegs[i])
            newIsParam := isParamSeg(newSegs[i])
            if oldIsParam != newIsParam {
                paramSame = false
                break
            }
            if !oldIsParam && !newIsParam && oldSegs[i] != newSegs[i] {
                diffCount++
            }
        }
        if !paramSame || diffCount == 0 {
            continue
        }
        if bestDiff < 0 || diffCount < bestDiff || (diffCount == bestDiff && cand < best) {
            bestDiff = diffCount
            best = cand
        }
    }
    return best, bestDiff
}

// --- small helpers ---

func paramMap(params []*spec.Parameter) map[string]*spec.Parameter {
    m := map[string]*spec.Parameter{}
    for _, p := range params {
        m[p.Name] = p
    }
    return m
}

func pathSet(ps *spec.ParsedSpec) map[string]bool {
    m := map[string]bool{}
    for _, op := range ps.Operations {
        m[op.Path] = true
    }
    return m
}

func responseJSONSchema(resp *spec.Response) *spec.Schema {
    if resp == nil {
        return nil
    }
    if mt, ok := resp.Content["application/json"]; ok {
        return mt.Schema
    }
    return nil
}

func requestJSONSchema(op *spec.Operation) *spec.Schema {
    if op.RequestBody == nil {
        return nil
    }
    if mt, ok := op.RequestBody.Content["application/json"]; ok {
        return mt.Schema
    }
    return nil
}

func stringSet(ss []string) map[string]bool {
    m := map[string]bool{}
    for _, s := range ss {
        m[s] = true
    }
    return m
}

func splitPath(p string) []string {
    return strings.Split(strings.Trim(p, "/"), "/")
}

func isParamSeg(seg string) bool {
    return strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}")
}

func sortedSlice(m map[string]bool) []string {
    var result []string
    for k := range m {
        result = append(result, k)
    }
    sort.Strings(result)
    return result
}
```

- [ ] **Step 4: Run tests**

```bash
go test ./internal/diff/... -run TestDiff -v
```
Expected: all 5 tests PASS

- [ ] **Step 5: Create testdata for integration test**

Create `testdata/petstore_v1.yaml`:
```yaml
openapi: "3.0.3"
info:
  title: Petstore v1
  version: "1.0"
paths:
  /pets:
    get:
      operationId: listPets
      parameters:
        - name: limit
          in: query
          schema:
            type: integer
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: integer
                  name:
                    type: string
                  email:
                    type: string
  /pets/{petId}:
    delete:
      operationId: deletePet
      parameters:
        - name: petId
          in: path
          required: true
          schema:
            type: integer
      responses:
        "204":
          description: No Content
```

Create `testdata/petstore_v2.yaml`:
```yaml
openapi: "3.0.3"
info:
  title: Petstore v2
  version: "2.0"
paths:
  /pets:
    get:
      operationId: listPets
      parameters:
        - name: limit
          in: query
          schema:
            type: integer
        - name: includeArchived
          in: query
          schema:
            type: boolean
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: integer
                  name:
                    type: string
  /pets/{petId}/photos:
    get:
      operationId: listPetPhotos
      parameters:
        - name: petId
          in: path
          required: true
          schema:
            type: integer
      responses:
        "200":
          description: OK
```

- [ ] **Step 6: Commit**

```bash
git add internal/diff/diff.go internal/diff/diff_test.go testdata/petstore_v1.yaml testdata/petstore_v2.yaml
git commit -m "feat: implement Spec Diff core logic with BREAKING/POTENTIALLY_BREAKING/NON_BREAKING classification"
```

---

## Task 10: `internal/diff/suggest.go` — affected test case inference

**Files:**
- Create: `internal/diff/suggest.go`
- Create: `internal/diff/suggest_test.go`

**Context:** `Suggest` receives a `DiffResult` and `[]schema.TestCase`. It normalizes step paths (replacing `{{x}}` with `{x}`) and matches against `Change.Path`. Use `writer.NewJSONSchemaWriter().Read(indexPath)` to read index.json in `cmd/diff.go` (not here; the function receives already-loaded cases).

- [ ] **Step 1: Write failing tests**

Create `internal/diff/suggest_test.go`:
```go
package diff

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/testmind-hq/caseforge/internal/output/schema"
)

func TestSuggest_SpecPathMatch(t *testing.T) {
    result := DiffResult{Changes: []Change{
        {Kind: Breaking, Method: "DELETE", Path: "/users/{id}", Description: "endpoint removed"},
    }}
    cases := []schema.TestCase{
        {ID: "TC-0001", Title: "delete user", Source: schema.CaseSource{SpecPath: "DELETE /users/{id}"}},
        {ID: "TC-0002", Title: "create user", Source: schema.CaseSource{SpecPath: "POST /users"}},
    }
    affected := Suggest(result, cases)
    assert.Len(t, affected, 1)
    assert.Equal(t, "TC-0001", affected[0].ID)
    assert.NotEmpty(t, affected[0].Reason)
}

func TestSuggest_StepPathMatch(t *testing.T) {
    result := DiffResult{Changes: []Change{
        {Kind: Breaking, Method: "GET", Path: "/users/{id}", Description: "endpoint removed"},
    }}
    cases := []schema.TestCase{
        {ID: "TC-0003", Title: "chain case", Source: schema.CaseSource{SpecPath: "POST /users"},
            Steps: []schema.Step{
                {Path: "/users/{{userId}}"},
            },
        },
    }
    affected := Suggest(result, cases)
    assert.Len(t, affected, 1)
    assert.Equal(t, "TC-0003", affected[0].ID)
}

func TestSuggest_NonBreakingNotIncluded(t *testing.T) {
    result := DiffResult{Changes: []Change{
        {Kind: NonBreaking, Method: "GET", Path: "/users", Description: "new optional param"},
    }}
    cases := []schema.TestCase{
        {ID: "TC-0001", Source: schema.CaseSource{SpecPath: "GET /users"}},
    }
    affected := Suggest(result, cases)
    assert.Empty(t, affected)
}

func TestSuggest_Empty(t *testing.T) {
    result := DiffResult{}
    affected := Suggest(result, nil)
    assert.Empty(t, affected)
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/diff/... -run TestSuggest 2>&1 | head -5
```
Expected: compile error

- [ ] **Step 3: Implement `internal/diff/suggest.go`**

Create `internal/diff/suggest.go`:
```go
// internal/diff/suggest.go
package diff

import (
    "fmt"
    "strings"

    "github.com/testmind-hq/caseforge/internal/output/schema"
)

// AffectedCase describes a test case that may be broken by a spec change.
type AffectedCase struct {
    ID     string
    Title  string
    Reason string
}

// Suggest returns the test cases likely affected by Breaking or PotentiallyBreaking changes.
// cases should be loaded via writer.NewJSONSchemaWriter().Read(indexPath) before calling.
func Suggest(result DiffResult, cases []schema.TestCase) []AffectedCase {
    // Build a set of breaking change paths for fast lookup
    breakingPaths := map[string]Change{}
    for _, c := range result.Changes {
        if c.Kind == Breaking || c.Kind == PotentiallyBreaking {
            breakingPaths[c.Path] = c
        }
    }
    if len(breakingPaths) == 0 {
        return nil
    }

    var affected []AffectedCase
    for _, tc := range cases {
        reason := ""
        // Match via Source.SpecPath ("METHOD /path" → extract path)
        specPath := extractPath(tc.Source.SpecPath)
        if change, ok := breakingPaths[specPath]; ok {
            reason = fmt.Sprintf("%s change on %s: %s", change.Kind, change.Path, change.Description)
        }
        // Match via Steps[].Path (normalize {{x}} → {x})
        if reason == "" {
            for _, step := range tc.Steps {
                normalizedPath := normalizeTemplatePath(step.Path)
                if change, ok := breakingPaths[normalizedPath]; ok {
                    reason = fmt.Sprintf("step path %s matches %s change: %s", step.Path, change.Kind, change.Description)
                    break
                }
            }
        }
        if reason != "" {
            affected = append(affected, AffectedCase{ID: tc.ID, Title: tc.Title, Reason: reason})
        }
    }
    return affected
}

// extractPath extracts the path from a SpecPath string like "GET /users/{id}" → "/users/{id}".
func extractPath(specPath string) string {
    parts := strings.SplitN(specPath, " ", 2)
    if len(parts) == 2 {
        return parts[1]
    }
    return specPath
}

// normalizeTemplatePath converts Hurl/Postman {{varName}} to OpenAPI {varName}.
func normalizeTemplatePath(path string) string {
    // Replace {{x}} with {x}
    result := strings.Builder{}
    i := 0
    for i < len(path) {
        if i+1 < len(path) && path[i] == '{' && path[i+1] == '{' {
            // find closing }}
            end := strings.Index(path[i+2:], "}}")
            if end >= 0 {
                varName := path[i+2 : i+2+end]
                result.WriteString("{")
                result.WriteString(varName)
                result.WriteString("}")
                i = i + 2 + end + 2
                continue
            }
        }
        result.WriteByte(path[i])
        i++
    }
    return result.String()
}
```

- [ ] **Step 4: Run tests**

```bash
go test ./internal/diff/... -v
```
Expected: all PASS

- [ ] **Step 5: Commit**

```bash
git add internal/diff/suggest.go internal/diff/suggest_test.go
git commit -m "feat: implement Suggest for affected test case inference"
```

---

## Task 11: `cmd/diff.go` — CLI command + integration test

**Files:**
- Create: `cmd/diff.go`
- Modify: `cmd/root.go` — `diffCmd` is registered via `init()` in `diff.go`, no change needed to `root.go`

**Context:** The `cmd/diff.go` `init()` registers `diffCmd` with `rootCmd.AddCommand(diffCmd)`. Output is text by default; `--format json` outputs the JSON schema from the spec.

- [ ] **Step 1: Write a CLI integration test**

Add to `cmd/diff_test.go` (create new file) — follows the same `//go:build integration` pattern as `gen_integration_test.go`:
```go
// cmd/diff_test.go
//go:build integration

package cmd

import (
    "bytes"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestDiffCommand_BasicText(t *testing.T) {
    diffOld = "../testdata/petstore_v1.yaml"
    diffNew = "../testdata/petstore_v2.yaml"
    diffFormat = "text"
    diffCases = ""
    t.Cleanup(func() {
        diffOld = ""
        diffNew = ""
        diffFormat = "text"
        diffCases = ""
    })
    buf := &bytes.Buffer{}
    diffCmd.SetOut(buf)
    diffCmd.SetErr(buf)
    // breaking changes → runDiff returns errBreakingChanges; that's expected
    err := runDiff(diffCmd, nil)
    require.ErrorIs(t, err, errBreakingChanges)
    output := buf.String()
    assert.Contains(t, output, "BREAKING")
    assert.Contains(t, output, "/pets/{petId}")
}

func TestDiffCommand_JSONFormat(t *testing.T) {
    diffOld = "../testdata/petstore_v1.yaml"
    diffNew = "../testdata/petstore_v2.yaml"
    diffFormat = "json"
    diffCases = ""
    t.Cleanup(func() {
        diffOld = ""
        diffNew = ""
        diffFormat = "text"
        diffCases = ""
    })
    buf := &bytes.Buffer{}
    diffCmd.SetOut(buf)
    diffCmd.SetErr(buf)
    err := runDiff(diffCmd, nil)
    require.ErrorIs(t, err, errBreakingChanges)
    output := buf.String()
    assert.Contains(t, output, `"kind"`)
    assert.Contains(t, output, "BREAKING")
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test -tags integration ./cmd/... -run TestDiffCommand 2>&1 | head -10
```
Expected: compile error (diffCmd, diffOld, diffNew, etc. undefined — diff.go doesn't exist yet)

- [ ] **Step 3: Implement `cmd/diff.go`**

Create `cmd/diff.go`:
```go
// cmd/diff.go
package cmd

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
    "github.com/testmind-hq/caseforge/internal/diff"
    "github.com/testmind-hq/caseforge/internal/output/writer"
    "github.com/testmind-hq/caseforge/internal/spec"
)

// errBreakingChanges is returned by runDiff when breaking changes are detected.
// root.Execute() converts any non-nil error to os.Exit(1).
var errBreakingChanges = errors.New("breaking changes detected")

var diffCmd = &cobra.Command{
    Use:          "diff",
    Short:        "Compare two OpenAPI specs and classify breaking changes",
    RunE:         runDiff,
    SilenceErrors: true, // suppress cobra's "Error: breaking changes detected" message
}

var (
    diffOld    string
    diffNew    string
    diffCases  string
    diffFormat string
)

func init() {
    rootCmd.AddCommand(diffCmd)
    diffCmd.Flags().StringVar(&diffOld, "old", "", "Old spec file (required)")
    diffCmd.Flags().StringVar(&diffNew, "new", "", "New spec file (required)")
    diffCmd.Flags().StringVar(&diffCases, "cases", "", "Cases output dir (optional; reads index.json to infer affected cases)")
    diffCmd.Flags().StringVar(&diffFormat, "format", "text", "Output format: text|json")
    _ = diffCmd.MarkFlagRequired("old")
    _ = diffCmd.MarkFlagRequired("new")
}

func runDiff(cmd *cobra.Command, args []string) error {
    loader := spec.NewLoader()

    oldSpec, err := loader.Load(diffOld)
    if err != nil {
        return fmt.Errorf("loading old spec: %w", err)
    }
    newSpec, err := loader.Load(diffNew)
    if err != nil {
        return fmt.Errorf("loading new spec: %w", err)
    }

    result := diff.Diff(oldSpec, newSpec)

    // Load affected cases if --cases provided
    var affected []diff.AffectedCase
    if diffCases != "" {
        indexPath := filepath.Join(diffCases, "index.json")
        w := writer.NewJSONSchemaWriter()
        cases, readErr := w.Read(indexPath)
        if readErr != nil {
            fmt.Fprintf(os.Stderr, "warning: could not read %s: %v\n", indexPath, readErr)
        } else {
            affected = diff.Suggest(result, cases)
        }
    }

    switch diffFormat {
    case "json":
        if err := printDiffJSON(cmd, result, affected); err != nil {
            return err
        }
    default:
        printDiffText(cmd, result, affected)
    }

    if result.HasBreaking() {
        return errBreakingChanges
    }
    return nil
}

func printDiffText(cmd *cobra.Command, result diff.DiffResult, affected []diff.AffectedCase) {
    // Group by kind
    var breaking, potBreaking, nonBreaking []diff.Change
    for _, c := range result.Changes {
        switch c.Kind {
        case diff.Breaking:
            breaking = append(breaking, c)
        case diff.PotentiallyBreaking:
            potBreaking = append(potBreaking, c)
        case diff.NonBreaking:
            nonBreaking = append(nonBreaking, c)
        }
    }

    out := cmd.OutOrStdout()

    if len(breaking) > 0 {
        fmt.Fprintf(out, "\nBREAKING (%d):\n", len(breaking))
        for _, c := range breaking {
            loc := ""
            if c.Location != "" {
                loc = " " + c.Location
            }
            fmt.Fprintf(out, "  ✗ %-8s %-30s %s\n", c.Method, c.Path+loc, c.Description)
        }
    }
    if len(potBreaking) > 0 {
        fmt.Fprintf(out, "\nPOTENTIALLY BREAKING (%d):\n", len(potBreaking))
        for _, c := range potBreaking {
            loc := ""
            if c.Location != "" {
                loc = " " + c.Location
            }
            fmt.Fprintf(out, "  ⚠ %-8s %-30s %s\n", c.Method, c.Path+loc, c.Description)
        }
    }
    if len(nonBreaking) > 0 {
        fmt.Fprintf(out, "\nNON-BREAKING (%d):\n", len(nonBreaking))
        for _, c := range nonBreaking {
            fmt.Fprintf(out, "  + %-8s %-30s %s\n", c.Method, c.Path, c.Description)
        }
    }

    if len(affected) > 0 {
        fmt.Fprintf(out, "\nAffected test cases:\n")
        for _, a := range affected {
            fmt.Fprintf(out, "  %s  %s — %s\n", a.ID, a.Title, a.Reason)
        }
    }

    if len(result.Changes) == 0 {
        fmt.Fprintln(out, "No changes detected.")
    }
}

type jsonDiffOutput struct {
    Summary struct {
        Breaking            int `json:"breaking"`
        PotentiallyBreaking int `json:"potentially_breaking"`
        NonBreaking         int `json:"non_breaking"`
    } `json:"summary"`
    Changes       []diff.Change       `json:"changes"`
    AffectedCases []diff.AffectedCase `json:"affected_cases,omitempty"`
}

func printDiffJSON(cmd *cobra.Command, result diff.DiffResult, affected []diff.AffectedCase) error {
    out := jsonDiffOutput{}
    out.Changes = result.Changes
    out.AffectedCases = affected
    for _, c := range result.Changes {
        switch c.Kind {
        case diff.Breaking:
            out.Summary.Breaking++
        case diff.PotentiallyBreaking:
            out.Summary.PotentiallyBreaking++
        case diff.NonBreaking:
            out.Summary.NonBreaking++
        }
    }
    data, err := json.MarshalIndent(out, "", "  ")
    if err != nil {
        return err
    }
    fmt.Fprintln(cmd.OutOrStdout(), string(data))
    return nil
}
```

- [ ] **Step 4: Run tests**

```bash
go test -tags integration ./cmd/... -run TestDiffCommand -v
```
Expected: PASS

- [ ] **Step 5: Run full test suite**

```bash
go build ./... && go test ./...
```
Expected: all PASS

- [ ] **Step 6: Smoke test manually**

```bash
go run ./main.go diff --old testdata/petstore_v1.yaml --new testdata/petstore_v2.yaml
```
Expected: output shows BREAKING changes for /pets/{petId} removal

- [ ] **Step 7: Commit**

```bash
git add cmd/diff.go
git commit -m "feat: add caseforge diff command for breaking change detection"
```

---

## Final: Full verification

- [ ] **Step 1: Run race detector**

```bash
go test -race ./...
```
Expected: no data races

- [ ] **Step 2: Verify build**

```bash
go build -o /dev/null ./...
```
Expected: clean

- [ ] **Step 3: Run full test suite one more time**

```bash
go test ./... -count=1
```
Expected: all PASS

- [ ] **Step 4: Final commit if any cleanup needed**

```bash
git add -A && git status
# commit only if there are uncommitted changes
```
