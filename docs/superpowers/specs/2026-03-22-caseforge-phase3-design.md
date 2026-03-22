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
| `internal/spec/types.go` | 补充 `Operation.Security []string`、`Operation.Summary string` 字段 |
| `internal/spec/parser.go` | 解析 `security` 和 `summary` 字段 |

### 2.2 依赖关系

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

func (t *SecurityTechnique) Applies(op *spec.Operation) bool {
    // 至少一条规则触发时返回 true
}

func (t *SecurityTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
    // 遍历 8 条 per-op 规则，收集所有适用的用例
}
```

```go
// internal/methodology/owasp_spec.go

type SecuritySpecTechnique struct{}

func NewSecuritySpecTechnique() *SecuritySpecTechnique { return &SecuritySpecTechnique{} }
func (t *SecuritySpecTechnique) Name() string          { return "owasp_api_top10_spec" }

func (t *SecuritySpecTechnique) Generate(s *spec.ParsedSpec) ([]schema.TestCase, error) {
    // API5：功能级授权
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

#### API4 — 资源消耗无限制
- **触发**：操作有 `limit`/`size`/`per_page` query 参数
- **生成**：`limit=99999`，断言状态码 `200` 或 `400`（标注 rationale 为性能风险）

#### API5 — 功能级授权缺失（`SpecTechnique`）
- **触发**：Spec 同时含有低权限路径（`/users/me`）和高权限路径（含 `admin` 或 `DELETE /resource`）
- **生成**：用普通用户 token（`{{user_token}}`）访问高权限接口，期望 `403`

#### API6 — 批量赋值
- **触发**：`op.Method == "POST" || op.Method == "PUT"`，且有 request body
- **生成**：Body 额外添加只读字段（`id`、`createdAt`、`updatedAt`），断言响应中这些字段未被接受/修改

#### API7 — 注入（XSS/SQLi/路径遍历）
- **触发**：有 `string` 类型的 query/path/body 参数
- **生成**：三个子用例，分别注入 `"><script>alert(1)</script>`、`' OR 1=1--`、`../../../etc/passwd`，期望 `400` 或响应中无注入内容回显

#### API8 — 安全配置错误（CORS）
- **触发**：任意接口（`Applies` 每个操作都返回 true，但每个 path 只生成一次）
- **生成**：OPTIONS 请求，断言响应 `Access-Control-Allow-Origin` 不为 `*`

#### API9 — 资产管理不当（`SpecTechnique`）
- **触发**：Spec 路径中同时含有 `/v1/` 和 `/v2/`（或类似版本前缀）
- **生成**：访问旧版本路径，断言期望 `404` 或 `410`（标注旧 API 应已下线）

#### API10 — 不安全 API 消费（SSRF）
- **触发**：有 `url`/`callback`/`webhook`/`redirect` 名称的参数
- **生成**：注入 `http://127.0.0.1`，期望 `400` 或 `403`

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

在 `cmd/gen.go` 中：

```go
engine.AddSpecTechnique(methodology.NewSecuritySpecTechnique())
// SecurityTechnique 通过 NewEngine 的 techniques 参数传入
engine := methodology.NewEngine(provider,
    ...,
    methodology.NewSecurityTechnique(),
)
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
- **示例**：`userId` 和 `user_id` 同时出现 → `L008 warning: mixed naming styles (camelCase and snake_case)`

#### L009 — 分页参数不一致（InconsistentPagination）
- **严重性**：warning
- **检测**：有分页参数的操作中，检测 `page`/`size` 风格与 `offset`/`limit` 风格同时存在
- **示例**：`GET /users?page=1&size=10` 和 `GET /orders?offset=0&limit=10` 并存

#### L010 — 错误响应结构不统一（InconsistentErrorSchema）
- **严重性**：warning
- **检测**：收集所有 4xx/5xx 响应 Schema 的顶级字段名，若存在两种以上不同结构（如 `{error: ...}` vs `{message: ...}` vs `{code: ..., msg: ...}`）则报告
- **示例**：`{"error": "..."}` 和 `{"message": "..."}` 并存

#### L011 — 缺少安全声明（MissingSecurityScheme）
- **严重性**：error
- **检测**：非 GET 的接口（或根据 Spec 全局 security 配置判断）未声明 `security` 字段
- **排除**：路径含 `public`、`health`、`login`、`register` 的接口

#### L012 — 敏感字段暴露（SensitiveFieldExposed）
- **严重性**：error
- **检测**：响应 Schema（2xx）中含有字段名匹配 `password`、`secret`、`token`、`private_key`、`access_token`、`refresh_token`（大小写不敏感）
- **示例**：`GET /users/{id}` 响应含 `passwordHash` 字段

### 4.2 评分算法（`internal/lint/score.go`）

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

### 4.3 CLI 输出变更

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
| `Step.Path` | `request.url.raw`（`{{base_url}}/path`），`url.host`/`url.path` 分拆 |
| `Step.Headers` | `request.header[]`（`{key, value}` 数组） |
| `Step.Body`（JSON） | `request.body`（`mode: "raw"`, `raw: json`, `options.raw.language: "json"`） |
| `Step.Assertions` | `event[listen: "test"]` script |
| `Step.Captures` | `event[listen: "test"]` script（`pm.environment.set`） |
| `{{varName}}` | 直接保留（Postman 原生语法相同） |

### 5.3 Test Script 生成

每个 Step 生成 `event[listen: "test"]`：

```javascript
// 状态码断言（始终生成）
pm.test("status is 201", function () {
    pm.response.to.have.status(201);
});

// JSON body 存在性断言（有 body 类断言时）
pm.test("response is json", function () {
    pm.response.to.be.json;
});
```

Chain 用例 setup step 的 Captures（`userId: jsonpath $.id`）额外追加：

```javascript
// Capture: userId
var jsonData = pm.response.json();
pm.environment.set("userId", jsonData.id);
```

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

- `--cases` 可选，提供后输出受影响用例推断
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
| Path 重命名（旧路径消失，新路径出现，结构相似） | BREAKING |
| 新增 required 请求体字段 | POTENTIALLY_BREAKING |
| 响应字段类型放宽（`integer` → `number`） | POTENTIALLY_BREAKING |
| 新增 endpoint | NON_BREAKING |
| 新增可选参数 | NON_BREAKING |
| 新增响应字段 | NON_BREAKING |

### 6.3 受影响用例推断（`internal/diff/suggest.go`）

```go
// Suggest 读取 index.json，对每个 TestCase 检查：
// 1. Source.SpecPath 是否命中某条 Breaking/PotentiallyBreaking 变更的 path
// 2. Steps[].Path（去掉变量后）是否与被删除或修改的 endpoint 匹配
func Suggest(result DiffResult, cases []schema.TestCase) []AffectedCase

type AffectedCase struct {
    ID     string
    Title  string
    Reason string // 说明为什么受影响
}
```

### 6.4 CLI 输出示例

```
$ caseforge diff --old v1.yaml --new v2.yaml --cases ./cases/

BREAKING (2):
  ✗ DELETE /users/{id}              endpoint removed
  ✗ POST   /orders requestBody      field "customerId" type: integer → string

POTENTIALLY BREAKING (1):
  ⚠ GET    /users/{id} response.200  field "email" removed

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
