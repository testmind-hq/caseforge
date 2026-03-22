# CaseForge Phase 3 Design Spec

**Date:** 2026-03-22
**Branch:** phase3
**Scope:** OWASP API Top 10 用例生成、Lint 完整版（设计一致性 + 安全声明 + 评分）、Postman Collection v2.1 渲染器、Spec Diff 破坏性变更检测

---

## 1. 背景与目标

Phase 1 完成了核心方法论引擎（等价类、边界值、决策表、状态转换、Pairwise、幂等性）和基础 CLI。Phase 2 完成了链式用例、事件总线、TUI 进度显示、完整 CSV 导出、URL 编码修复和 HurlRunner report-json 解析。

Phase 3 补齐 Phase 2 规划中尚未实现的四个高价值特性：

1. **OWASP API Top 10 安全测试用例生成** — 针对 API 安全威胁自动生成专项用例
2. **Lint 完整版** — 设计一致性规则 + 安全声明规则 + Spec 质量评分
3. **Postman Collection v2.1 渲染器** — 支持 `gen --format postman` 输出可导入的 Collection
4. **Spec Diff** — 两版本对比，分类破坏性/非破坏性变更，推断受影响用例

---

## 2. 整体架构

### 2.1 文件地图

**新增文件：**

| 文件 | 职责 |
|------|------|
| `internal/security/helpers.go` | 共享 helpers：`HasIDPathParam`、`FindSensitiveFields`、`IsAuthRequired`、`FindVersionedPaths` |
| `internal/security/helpers_test.go` | helpers 单元测试 |
| `internal/methodology/owasp.go` | `SecurityTechnique`：per-op 实现 `Technique`，覆盖 API1/2/3/4/6/7/8/10 |
| `internal/methodology/owasp_test.go` | SecurityTechnique 单元测试 |
| `internal/methodology/owasp_spec.go` | `SecuritySpecTechnique`：跨操作规则，实现 `SpecTechnique`，覆盖 API5/9 |
| `internal/methodology/owasp_spec_test.go` | SecuritySpecTechnique 单元测试 |
| `internal/lint/consistency.go` | L007-L010：设计一致性规则 |
| `internal/lint/security_rules.go` | L011-L012：安全声明规则 |
| `internal/lint/score.go` | Spec 质量评分（0-100） |
| `internal/lint/score_test.go` | 评分逻辑单元测试 |
| `internal/output/render/postman.go` | Postman Collection v2.1 渲染器 |
| `internal/output/render/postman_test.go` | 渲染器单元测试 |
| `internal/diff/diff.go` | 两版本 Spec 对比，输出 `DiffResult` |
| `internal/diff/diff_test.go` | Diff 逻辑单元测试 |
| `internal/diff/suggest.go` | 根据 DiffResult 推断受影响的 TestCase ID |
| `internal/diff/suggest_test.go` | 推断逻辑单元测试 |
| `cmd/diff.go` | `caseforge diff --old --new [--cases]` 命令 |
| `testdata/petstore_v1.yaml` | Diff 集成测试用的旧版 Spec |
| `testdata/petstore_v2.yaml` | Diff 集成测试用的新版 Spec（含破坏性变更） |

**修改文件：**

| 文件 | 变更 |
|------|------|
| `cmd/gen.go` | 注册 `SecurityTechnique` + `SecuritySpecTechnique`，加入 `postman` format case |
| `cmd/lint.go` | 输出评分，新增 `--min-score` flag |
| `cmd/root.go` | 注册 `diffCmd` |
| `internal/output/render/hurl.go` | 扩展 `renderAssertion`，新增 target 和 operator 映射（见下表） |
| `internal/output/schema/model.go` | 更新 `Assertion.Target` 字段注释，列出规范值：`"status_code"`、`"jsonpath $.<field>"`、`"header <HeaderName>"`、`"duration_ms"` |
| `internal/spec/loader.go` | 在 `Operation` 结构体（定义于此文件）中新增 `Security []string` 字段；`Summary string` 字段已存在，无需改动 |
| `internal/spec/parser.go` | 在 `convertOperation` 中解析 `op.Security`（kin-openapi `openapi3.SecurityRequirements`），提取所有 scheme name 填充到 `Security []string` |

### 2.2 Hurl Renderer 扩展（`renderAssertion`）

`internal/output/render/hurl.go` 的 `renderAssertion` 函数新增以下 target 支持：

| Target | Operator | Expected | Hurl 输出 |
|--------|----------|----------|-----------|
| `jsonpath $.<field>` | `eq` | any | `jsonpath "$.<field>" == <value>` |
| `jsonpath $.<field>` | `ne` | non-nil | `jsonpath "$.<field>" != <value>` |
| `jsonpath $.<field>` | `ne` | nil | `jsonpath "$.<field>" not exists` |
| `jsonpath $.<field>` | `contains` | string | `jsonpath "$.<field>" contains "<value>"` |
| `header <HeaderName>` | `eq` | any | `header "<HeaderName>" == <value>` |
| `header <HeaderName>` | `ne` | non-nil | `header "<HeaderName>" != <value>` |

其余 operator（`lt`、`gt`、`matches`）对这两种 target 暂不支持，跳过（输出注释 `# unrendered assertion`）。

### 2.3 依赖关系

```
internal/security/helpers.go
    ← internal/methodology/owasp.go
    ← internal/lint/security_rules.go

internal/diff/diff.go
    ← internal/diff/suggest.go
    ← cmd/diff.go
```

OWASP technique 和 lint 安全规则共享 `internal/security` helpers，避免重复逻辑。其余三个特性完全独立。

---

## 3. OWASP API Top 10 用例生成

### 3.1 接口设计

```go
// internal/methodology/owasp.go

type SecurityTechnique struct{}

func NewSecurityTechnique() *SecurityTechnique { return &SecurityTechnique{} }
func (t *SecurityTechnique) Name() string      { return "owasp_api_top10" }

// Applies 返回 true 的条件（满足任意一条即可）：
//   - HasIDPathParam(op) — API1 BOLA
//   - op.Security 非空 — API2 认证
//   - op.Method 为 PATCH/PUT — API3 BOPLA
//   - 操作有 limit/size/per_page 参数 — API4
//   - op.Method 为 POST/PUT 且有 body — API6
//   - 有 string 类型参数 — API7
//   - 始终 true（用于 API8，但 API8 在 SecuritySpecTechnique 中处理，见下）
//   - 有 url/callback/webhook/redirect 参数 — API10
// 实际实现：返回 HasIDPathParam(op) || len(op.Security)>0 || ...（检查以上条件）
func (t *SecurityTechnique) Applies(op *spec.Operation) bool { ... }

func (t *SecurityTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
    // 遍历适用的 per-op 规则（API1/2/3/4/6/7/10），收集所有用例
    // API8 CORS 测试移至 SecuritySpecTechnique，避免 per-op 去重问题
}
```

```go
// internal/methodology/owasp_spec.go

type SecuritySpecTechnique struct{}

func NewSecuritySpecTechnique() *SecuritySpecTechnique { return &SecuritySpecTechnique{} }
func (t *SecuritySpecTechnique) Name() string          { return "owasp_api_top10_spec" }

func (t *SecuritySpecTechnique) Generate(s *spec.ParsedSpec) ([]schema.TestCase, error) {
    // API5：功能级授权
    // API8：CORS OPTIONS（按 path 去重，收集所有唯一 path，每个只生成一次 OPTIONS 用例）
    // API9：资产管理（旧版路径仍可访问）
}
```

### 3.2 十条规则详细设计

#### API1 — BOLA（对象级授权缺失）
- **触发**：`security.HasIDPathParam(op)` 为 true（路径含 `{id}`、`{userId}` 等）
- **生成**：将 ID 参数替换为 `{{other_resource_id}}`，断言期望 `403`
- **变量说明**：用例 title 含 setup note 提示需要通过 `--var other_resource_id=<值>` 注入

#### API2 — 认证机制失效
- **触发**：`op.Security` 非空（接口声明了 security scheme）
- **生成**：去除 Authorization header，断言期望 `401`

#### API3 — BOPLA（对象属性级授权缺失）
- **触发**：`op.Method == "PATCH" || op.Method == "PUT"`
- **生成**：Body 额外注入 `{"is_admin": true, "role": "admin"}`，断言响应中这些字段值未被修改
- **断言格式**：使用 `schema.Assertion{Target: "jsonpath $.is_admin", Operator: "ne", Expected: true}` 和 `{Target: "jsonpath $.role", Operator: "ne", Expected: "admin"}`（`"ne"` 是现有 schema 包定义的操作符）；若注入字段不应出现在响应中（`Expected: nil` 语义），改用 `{Target: "jsonpath $.is_admin", Operator: "ne", Expected: nil}`，Hurl 渲染器对 `Expected == nil` + `Operator == "ne"` 的 jsonpath 断言输出 `jsonpath "$.is_admin" not exists`

#### API4 — 资源消耗无限制
- **触发**：操作有 `limit`/`size`/`per_page` query 参数
- **生成**：`limit=99999`，断言状态码为 `200`（rationale 注明"服务器接受了极大分页值，存在性能风险"；正确实现应返回 `400`，但生成的用例为探测性测试，记录实际行为供审查）
- **断言**：`{Target: "status_code", Operator: "eq", Expected: 200}`（即检查服务器**是否允许**此值，而非期望拒绝）

#### API5 — 功能级授权缺失（`SpecTechnique`）
- **触发**：Spec 中同时存在低权限路径和高权限路径（定义如下）
  - **低权限路径**：`op.Path` 包含 `/me` 或 `/profile`（如 `/users/me`、`/profile`）
  - **高权限路径**：`op.Path` 包含 `admin`，或 `op.Method == "DELETE"`
- **生成条件**：低权限路径和高权限路径均至少存在一条时，对每条高权限路径生成一个用例：用普通用户 token（`{{user_token}}`）访问该接口，期望 `403`

#### API6 — 批量赋值
- **触发**：`op.Method == "POST" || op.Method == "PUT"`，且有 request body
- **生成**：Body 额外添加只读字段（`id`、`createdAt`、`updatedAt`），断言响应中这些字段未被接受/修改
- **断言格式**：同 API3。对每个注入的只读字段生成 `{Target: "jsonpath $.<field>", Operator: "ne", Expected: <injected_value>}`；若 `Expected == nil`，Hurl 渲染输出 `jsonpath "$.<field>" not exists`

#### API7 — 注入（XSS/SQLi/路径遍历）
- **触发**：有 `string` 类型的 query/path/body 参数
- **生成**：三个子用例，分别注入 `"><script>alert(1)</script>`、`' OR 1=1--`、`../../../etc/passwd`
- **断言**：每个子用例生成单一状态码断言 `{Target: "status_code", Operator: "eq", Expected: 400}`（期望服务器拒绝注入内容）

#### API8 — 安全配置错误（CORS）
- **实现位置**：**`SecuritySpecTechnique`**（非 per-op），因为需要跨所有操作去重
- **触发**：总是（`SecuritySpecTechnique.Generate` 始终执行 API8 检查）
- **去重**：在 `Generate` 内部遍历 `s.Operations`，维护局部 `seenPaths := map[string]bool{}`，以 `op.Path` 为键；同一 path 只生成一次 OPTIONS 用例
- **生成**：OPTIONS 请求到该 path，断言响应头 `Access-Control-Allow-Origin` 不为 `*`
- **断言格式**：`schema.Assertion{Target: "header Access-Control-Allow-Origin", Operator: "ne", Expected: "*"}`
- **新 Target 类型**：`"header <HeaderName>"` 是 Phase 3 新增的 Assertion target 格式（现有代码只有 `"status_code"` 和 `"jsonpath $.<field>"`）；Hurl 渲染器需在 Phase 3 中扩展以支持 `header` target；Postman 渲染器映射为 `pm.expect(pm.response.headers.get("<HeaderName>")).to.not.eql("*")`

#### API9 — 资产管理不当（`SpecTechnique`）
- **触发**：Spec 路径中同时含有 `/v1/` 和 `/v2/`（或类似版本前缀）
- **生成**：访问旧版本路径，断言期望 `404`（旧 API 应已下线；`410 Gone` 也是正确实现，但测试用例使用 `404` 作为单一期望值）
- **断言**：`{Target: "status_code", Operator: "eq", Expected: 404}`

#### API10 — 不安全 API 消费（SSRF）
- **触发**：有 `url`/`callback`/`webhook`/`redirect` 名称的参数
- **生成**：注入 `http://127.0.0.1`，断言期望 `400`（期望服务器验证并拒绝内网 URL）
- **断言**：`{Target: "status_code", Operator: "eq", Expected: 400}`

### 3.3 生成用例结构

```json
{
  "id": "TC-xxxx",
  "title": "[OWASP-API1] GET /users/{id} — BOLA 越权访问",
  "kind": "single",
  "priority": "P0",
  "tags": ["security", "owasp", "api1-bola"],
  "source": {
    "technique": "owasp_api_top10",
    "spec_path": "GET /users/{id}",
    "rationale": "路径含 ID 参数，需验证对象级授权：用合法 token 访问他人资源应返回 403"
  },
  "steps": [...]
}
```

### 3.4 注册到 Engine

在 `cmd/gen.go` 中（注意：先构造 engine，再调用方法）：

```go
engine := methodology.NewEngine(provider,
    ...,
    methodology.NewSecurityTechnique(),  // per-op：API1/2/3/4/6/7/10
)
engine.AddSpecTechnique(methodology.NewSecuritySpecTechnique())  // cross-op：API5/8/9
```

---

## 4. Lint 完整版

### 4.1 新增规则（L007-L012）

#### L007 — 路径含动词（VerbInPath）
- **严重性**：warning
- **检测**：路径各 segment 匹配常见动词前缀（`get`、`create`、`update`、`delete`、`list`、`fetch`、`add`、`remove`）
- **示例**：`POST /createUser` → `L007 warning: verb "create" in path segment`

#### L008 — 命名风格不一致（InconsistentNaming）
- **严重性**：warning
- **检测**：跨操作收集所有参数名和 body 字段名，检测同时出现 camelCase 和 snake_case
- **字段遍历深度**：只检查 request body schema 的**顶层**属性（`schema.Properties` 一层，不递归嵌套对象）
- **示例**：`userId` 和 `user_id` 同时出现 → `L008 warning: mixed naming styles (camelCase and snake_case)`

#### L009 — 分页参数不一致（InconsistentPagination）
- **严重性**：warning
- **检测**：有分页参数的操作中，检测 `page`/`size` 风格与 `offset`/`limit` 风格同时存在
- **示例**：`GET /users?page=1&size=10` 和 `GET /orders?offset=0&limit=10` 并存

#### L010 — 错误响应结构不统一（InconsistentErrorSchema）
- **严重性**：warning
- **检测**：收集所有 4xx/5xx 响应 Schema 的顶级字段名集合（`sorted(keys(schema.Properties))`），若存在两种及以上**不同的顶级字段名集合**则报告（如 `{"error"}` vs `{"message"}` vs `{"code","msg"}`）
- **"不同结构"定义**：字段名集合的字符串表示不同（例如 `"code,error"` vs `"message"`），顺序无关（排序后比较）
- **示例**：`{"error": "..."}` 和 `{"message": "..."}` 并存

#### L011 — 缺少安全声明（MissingSecurityScheme）
- **严重性**：error
- **检测**：`op.Method != "GET"` 且 `len(op.Security) == 0`（`op.Security` 为空切片或 nil）
- **排除**：路径（`op.Path`）包含以下任一子串：`/public`、`/health`、`/login`、`/register`（大小写敏感，前缀斜杠确保 segment 级别匹配）
- **不考虑**全局 security 配置（简化实现，只检查 operation 级别字段）

#### L012 — 敏感字段暴露（SensitiveFieldExposed）
- **严重性**：error
- **检测**：响应 Schema（2xx）中含有字段名**包含**以下任一关键词（`strings.Contains(strings.ToLower(fieldName), keyword)`）：`password`、`secret`、`token`、`private_key`、`access_token`、`refresh_token`
- **示例**：`GET /users/{id}` 响应含 `passwordHash` 字段，因为 `strings.Contains("passwordhash", "password")` 为 true

### 4.2 规则注册模式

新的 lint 规则文件（`consistency.go`、`security_rules.go`）通过 `init()` 函数自动注册到 `internal/lint` 包的全局规则列表：

```go
// 在 consistency.go / security_rules.go 文件末尾
func init() {
    register(newVerbInPathRule())
    register(newInconsistentNamingRule())
    // ...
}
```

`register` 是 `internal/lint/rules.go` 中现有的**未导出**函数（小写 `r`），维护一个包级 `[]LintRule` 切片。同包内的 `init()` 调用它即可，不需要导出。

### 4.3 评分算法（`internal/lint/score.go`）

```go
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

### 4.4 CLI 输出变更

```
caseforge lint --spec openapi.yaml [--min-score 80]

✗ L004  error    GET /users — no 2xx response defined
✗ L011  error    DELETE /users/{id} — no security scheme declared
⚠ L007  warning  POST /getUsers — verb in path segment "getUsers"
⚠ L008  warning  mixed naming styles: userId (camel) vs user_id (snake)

Spec Score: 64/100  (2 errors, 2 warnings)
exit code 1 (score 64 < min-score 80)
```

- 有 error 时 exit code 1（现有行为不变）
- 新增：`--min-score <n>`，分数低于阈值也 exit code 1

---

## 5. Postman Collection v2.1 渲染器

### 5.1 输出格式

**命令：** `caseforge gen --spec api.yaml --format postman --output ./cases/`
**输出文件：** `<outDir>/collection.json`

**顶层结构：**

```json
{
  "info": {
    "name": "CaseForge Generated",
    "_postman_id": "<uuid>",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "variable": [
    { "key": "base_url", "value": "http://localhost", "type": "string" }
  ],
  "item": [ ...items... ]
}
```

### 5.2 Schema 映射

| CaseForge | Postman v2.1 |
|-----------|-------------|
| `TestCase`（`kind: "single"`） | `item`（Request） |
| `TestCase`（`kind: "chain"`） | `item`（Folder），`item.item[]` 为子 requests |
| `Step.Method` | `request.method` |
| `Step.Path` | `request.url.raw`（`"{{base_url}}" + step.Path`），`url.host: ["{{base_url}}"]`，`url.path`：对 `strings.Split(strings.Trim(step.Path, "/"), "/")` 结果中每个 segment，若为 `{param}` 形式则转换为 `:param`，否则保留原值；模板变量 `{{param}}` 在路径中直接保留（Postman 原生支持） |
| `Step.Headers` | `request.header[]`（`{key, value}` 数组） |
| `Step.Body`（JSON） | `request.body`（`mode: "raw"`, `raw: json`, `options.raw.language: "json"`） |
| `Step.Assertions` | `event[listen: "test"]` script |
| `Step.Captures` | `event[listen: "test"]` script（`pm.environment.set`） |
| `{{varName}}` | 直接保留（Postman 原生语法相同） |

### 5.3 Test Script 生成

每个 Step 生成 `event[listen: "test"]`，脚本由以下各部分拼接：

**1. 状态码断言（来自 `Assertion{Target: "status_code", Operator: "eq", Expected: 201}`，由 `BasicAssertions` 生成）：**
```javascript
pm.test("status is 201", function () {
    pm.response.to.have.status(201);
});
```

**2. JSONPath 断言（来自 `Assertion{Target: "jsonpath $.<field>", Operator: "eq"|"ne"|..., Expected: ...}`）：**
```javascript
pm.test("jsonpath $.id eq ...", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.id).to.eql(<expected>);  // "eq"
    // pm.expect(jsonData.id).to.not.eql(<expected>);  // "ne"
    // pm.expect(String(jsonData.id)).to.include(<expected>);  // "contains"
});
```

**3. Captures（来自 `Capture{Name: "userId", From: "jsonpath $.id"}`）：**
```javascript
var jsonData = pm.response.json();
pm.environment.set("userId", jsonData.id);
```

**规则摘要：**
- `Target: "status_code"` → `pm.response.to.have.status(n)`（注意：现有代码全部使用 `"status_code"`，不是 `"status"`）
- `Target: "jsonpath $.<path>"` + `Operator: "eq"` → `pm.expect(val).to.eql(expected)`
- `Target: "jsonpath $.<path>"` + `Operator: "ne"` + `Expected != nil` → `pm.expect(val).to.not.eql(expected)`
- `Target: "jsonpath $.<path>"` + `Operator: "ne"` + `Expected == nil` → `pm.expect(val).to.not.exist` （Postman 的 `not.exist` 断言字段不存在或为 null）
- `Target: "jsonpath $.<path>"` + `Operator: "contains"` → `pm.expect(String(val)).to.include(expected)`
- `Target: "header <HeaderName>"` + `Operator: "ne"` → `pm.expect(pm.response.headers.get("<HeaderName>")).to.not.eql(expected)`
- `Capture.From: "jsonpath $.<path>"` → `pm.environment.set(name, jsonData.<path>)`（只支持单层 path；多层路径 `$.a.b` 生成 `jsonData.a.b`）

### 5.4 接口注册

`internal/output/render/postman.go` 实现现有 `Renderer` 接口：

```go
type PostmanRenderer struct{}

func NewPostmanRenderer() *PostmanRenderer { return &PostmanRenderer{} }
func (r *PostmanRenderer) Format() string  { return "postman" }
func (r *PostmanRenderer) Render(cases []schema.TestCase, outDir string) error { ... }
```

在 `cmd/gen.go` switch 中加入：

```go
case "postman":
    renderer = render.NewPostmanRenderer()
```

---

## 6. Spec Diff

### 6.1 命令

```bash
caseforge diff --old v1.yaml --new v2.yaml [--cases ./cases/] [--format text|json]
```

- `--cases` 可选，值为 `caseforge gen` 的输出目录；命令自动读取 `<dir>/index.json` 推断受影响用例
- `--format text`（默认）：人类可读；`--format json`：机器可读，供 CI 消费
- **Exit code**：有 BREAKING 或 POTENTIALLY_BREAKING 变更时返回 1，仅 NON_BREAKING 时返回 0

### 6.2 变更分类（`internal/diff/diff.go`）

```go
type ChangeKind string

const (
    Breaking           ChangeKind = "BREAKING"
    PotentiallyBreaking ChangeKind = "POTENTIALLY_BREAKING"
    NonBreaking        ChangeKind = "NON_BREAKING"
)

type Change struct {
    Kind        ChangeKind
    Method      string
    Path        string
    Location    string // "requestBody", "response.200", "param.limit", etc.
    Description string
}

type DiffResult struct {
    Changes []Change
}
```

**Breaking 变更检测规则：**

| 变更 | 分类 |
|------|------|
| Endpoint 删除 | BREAKING |
| 参数从 optional → required | BREAKING |
| 参数类型变更（`integer` → `string`） | BREAKING |
| 响应字段删除（2xx schema） | BREAKING |
| 响应字段类型变更（2xx schema） | BREAKING |
| Path 重命名（旧路径消失 + 新路径出现 + 非 segment 数量变化，且新旧路径 segment 结构相同） | BREAKING |
| 新增 required 请求体字段 | POTENTIALLY_BREAKING |
| 响应字段类型放宽（`integer` → `number`） | POTENTIALLY_BREAKING |
| 新增 endpoint | NON_BREAKING |
| 新增可选参数 | NON_BREAKING |
| 新增响应字段 | NON_BREAKING |

**Path 重命名判定算法（与 HTTP method 无关，只比较 path）：**

```
uniqueOldPaths = sorted(set(op.Path for op in v1.Operations) - set(op.Path for op in v2.Operations))
uniqueNewPaths = set(op.Path for op in v2.Operations) - set(op.Path for op in v1.Operations)

// 排序 uniqueOldPaths 保证确定性输出（Go: sort.Strings(oldPaths)）
for oldPath in uniqueOldPaths:
    bestCandidate = ""
    bestDiff = MaxInt
    for newPath in uniqueNewPaths:
        segs1 = strings.Split(strings.Trim(oldPath, "/"), "/")
        segs2 = strings.Split(strings.Trim(newPath, "/"), "/")
        if len(segs1) != len(segs2): continue
        // 参数位置必须完全相同（isParam 对两边都检查）
        paramSame = all(isParam(segs1[i]) == isParam(segs2[i]) for i in range(len(segs1)))
        if !paramSame: continue
        // 只统计两边都是静态 segment 但值不同的位置
        diffCount = count(i where !isParam(segs1[i]) && !isParam(segs2[i]) && segs1[i] != segs2[i])
        if diffCount == 0: continue  // 路径完全相同，不是重命名
        if diffCount < bestDiff || (diffCount == bestDiff && newPath < bestCandidate):
            bestDiff = diffCount
            bestCandidate = newPath
    if bestCandidate != "":
        emit BREAKING change: "path renamed: {oldPath} → {bestCandidate}"
        remove bestCandidate from uniqueNewPaths  // 避免同一 newPath 被多次匹配
```

其中 `isParam(seg) = strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}")`.

多候选时选 `diffCount` 最小的；`diffCount` 相同时按 `newPath` 字典序选最小的（确定性）。未被匹配的 `newPath` 保留为 NON_BREAKING 新增 endpoint。

### 6.3 `--format json` 输出 Schema

`--format json` 输出单个 JSON 对象，结构如下：

```json
{
  "summary": {
    "breaking": 2,
    "potentially_breaking": 1,
    "non_breaking": 2
  },
  "changes": [
    {
      "kind": "BREAKING",
      "method": "DELETE",
      "path": "/users/{id}",
      "location": "",
      "description": "endpoint removed"
    }
  ],
  "affected_cases": [
    {
      "id": "TC-0007",
      "title": "POST /orders — boundary values",
      "reason": "requestBody field type change may invalidate generated body"
    }
  ]
}
```

- `affected_cases` 仅在提供了 `--cases` 时出现
- `--cases <dir>` 读取 `<dir>/index.json`（由 `caseforge gen` 输出的 case index 文件）

### 6.4 `index.json` 格式（`--cases` 输入）

`caseforge gen` 通过 `internal/output/writer` 包写入 `index.json`，其顶层结构为 `writer.IndexFile` 封装类型（非裸数组）：

```json
{
  "$schema": "https://caseforge.dev/schema/v1/index.json",
  "version": "1",
  "generated_at": "...",
  "test_cases": [
    {
      "id": "TC-0007",
      "title": "POST /orders — boundary values",
      "kind": "single",
      "priority": "P1",
      "source": { "technique": "boundary", "spec_path": "POST /orders" },
      "steps": [...]
    }
  ]
}
```

`Suggest` 函数通过调用 `writer.NewJSONSchemaWriter().Read(indexPath)` 读取用例列表（`Read` 是 `*JSONSchemaWriter` 上的方法，非包级函数），通过 `Source.SpecPath` 和 `Steps[].Path` 与 `DiffResult.Changes` 中的 `Change.Path` 匹配。

`cmd/diff.go` 中传入路径为：`filepath.Join(casesDir, "index.json")`。

### 6.5 受影响用例推断（`internal/diff/suggest.go`）

**路径规范化：** 在与 `Change.Path`（OpenAPI 风格：`/users/{id}`）比较之前，需将 Step 路径中的 Hurl/Postman 模板变量规范化：
- `{{varName}}` → `{varName}`（双花括号 → 单花括号）
- 然后对规范化后的路径与 `Change.Path` 做字符串相等比较

```go
// Suggest 接收 DiffResult 和用例列表，返回受影响的用例。
// 调用方通过 writer.NewJSONSchemaWriter().Read(indexPath) 读取用例列表后传入。
// 1. Source.SpecPath（格式：`"METHOD /path"`）中的 path 部分命中某条 BREAKING/POTENTIALLY_BREAKING 变更的 Change.Path
// 2. 或 Steps[].Path 规范化后（{{x}}→{x}）与某条变更的 Change.Path 相等
func Suggest(result DiffResult, cases []schema.TestCase) []AffectedCase

type AffectedCase struct {
    ID     string
    Title  string
    Reason string // 说明为什么受影响
}
```

### 6.6 CLI 输出示例

```
$ caseforge diff --old v1.yaml --new v2.yaml --cases ./cases/

BREAKING (3):
  ✗ DELETE /users/{id}              endpoint removed
  ✗ POST   /orders requestBody      field "customerId" type: integer → string
  ✗ GET    /users/{id} response.200  field "email" removed

POTENTIALLY BREAKING (1):
  ⚠ POST   /users requestBody       new required field "phone" added

NON-BREAKING (2):
  + POST   /payments                new endpoint
  + GET    /users                   optional param "includeDeleted" added

Affected test cases:
  TC-0007  POST /orders — request body type change may invalidate generated body
  TC-0012  GET  /users/{id} — response assertion on "email" will fail

exit code 1
```

---

## 7. 测试策略

| 特性 | 单元测试 | 集成测试 |
|------|----------|----------|
| OWASP | 每条规则的触发/不触发 + 生成内容断言 | `gen --spec petstore.yaml --no-ai` 含 OWASP 用例 |
| Lint 新规则 | 每条规则独立测试 | `lint --spec inconsistent.yaml` 全规则命中 |
| Lint 评分 | 纯函数，各扣分场景 | 集成在 lint 命令测试中 |
| Postman 渲染 | single/chain TestCase → JSON 结构断言 | `gen --format postman` 输出可被 `json.Unmarshal` 解析 |
| Spec Diff | 各变更类型的独立检测 | `diff --old petstore_v1.yaml --new petstore_v2.yaml` 全输出验证 |

---

## 8. Definition of Done

- [ ] `go test ./...` 全部通过
- [ ] `go test -race ./...` 无数据竞争
- [ ] `go build ./...` 干净编译
- [ ] OWASP 10 条规则全部覆盖（每条至少一个测试验证触发条件和生成内容）
- [ ] Lint L007-L012 全部实现并有对应测试
- [ ] `caseforge lint --min-score 80` 在低质量 spec 上 exit code 1
- [ ] `gen --format postman` 输出合法 Postman Collection v2.1（可通过 JSON Schema 验证）
- [ ] Chain 用例的 Captures 在 Postman collection 中以 `pm.environment.set` 正确渲染
- [ ] `caseforge diff` 正确分类 BREAKING/POTENTIALLY_BREAKING/NON_BREAKING
- [ ] `caseforge diff --cases` 推断出受影响用例
