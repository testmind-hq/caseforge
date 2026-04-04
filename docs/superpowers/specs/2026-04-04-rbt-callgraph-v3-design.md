# RBT V3 — Go 类型感知调用图设计文档

**日期**：2026-04-04
**作者**：caseforge team
**状态**：已批准，待实现

---

## 目标

V2 的 name-based matching 在 Go Repository Pattern（接口多态）处断链：`service.go` 通过 `UserRepo` 接口调用 `MySQLRepo.Save`，V2 无法解析接口派发，导致 service 层改动无法追踪到路由。

V3 对 Go 项目引入 `golang.org/x/tools/go/callgraph`，做完整类型感知的调用图分析，将 `Via:"go-callgraph"` 置信度提升至 0.9（RTA）/ 0.95（PTA）。

---

## 背景与限制

### V2 的盲点

```
service.go:
    func CreateUser(repo UserRepo) {
        repo.Save(...)      // V2 name-matching 断链：只有 "Save" 名字，不知道实现
    }
```

V3 能解决：RTA 分析实际被实例化的类型，确定 `repo.Save` 调用的是 `MySQLRepo.Save`，从而完整追踪调用链。

### V3 有效范围

- ✅ 接口派发（Go Repository Pattern）：RTA/PTA 完整支持
- ✅ 直接函数调用：与 V2 等效，但精度更高
- ✅ 方法调用（struct method）：完整支持
- ⚠️ `reflect` 调用：不支持，超出范围
- ⚠️ `go generate` 生成代码：需项目可编译
- ❌ 非 Go 语言：不在范围内，继续走 V2

### 前提条件

- 项目根目录存在 `go.mod`（否则跳过 V3，静默 fallback 到 V2）
- 项目可编译（`go/packages` 需要成功加载）
- 编译失败 → 静默 fallback 到 V2，不影响用户

---

## 语言支持

仅 Go（`.go` 文件）。其他语言继续走 V2 的 tree-sitter + LLM 路径。

---

## 架构

### 数据流

```
git diff → ChangedFile 列表
    ↓
① runTreeSitterPhase（现有，不改）
   → RouteMapping{Via:"treesitter", Confidence:1.0}
    ↓ unclaimed
② [V3 NEW] runGoCallGraphPhase
   → 检测 SrcDir/go.mod；不存在则跳过（返回空）
   → go/packages.Load("./...", dir=SrcDir) 加载整个 module
   → buildssa.BuildPackage → SSA 形式
   → rta.Analyze(roots) / pointer.Analyze(config) 构建调用图
   → 从 unclaimed 中 .go 文件的所有函数出发 BFS 向上
   → 命中 routeFileMappings 节点 → RouteMapping{Via:"go-callgraph"}
   → 任何错误 → 静默 fallback，返回 nil（V2 兜底）
    ↓ 仍 unclaimed（非 Go 文件 + Go 分析失败的文件）
③ runCallGraphPhase（V2，现有，不改）
    ↓ 仍 unclaimed
④ runEmbedPhase（现有，不改）
    ↓
writeMapFile
```

### RunHybrid 新增 Phase

```go
func (idx *Indexer) RunHybrid(llmParser *LLMParser) error {
    files, _ := findSourceFiles(idx.SrcDir)

    // Phase 1: tree-sitter 直接路由检测（现有）
    mappings, routeFileMappings := idx.runTreeSitterPhase(files)
    unclaimed := subtractFiles(files, mappings)

    // Phase 2: [V3 NEW] Go 类型感知调用图
    goMappings, goClaimed := idx.runGoCallGraphPhase(unclaimed, routeFileMappings)
    mappings = append(mappings, goMappings...)
    unclaimed = subtractChangedFiles(unclaimed, goClaimed)

    // Phase 3: V2 name-based 调用图（现有）
    cgMappings, cgClaimed := idx.runCallGraphPhase(files, unclaimed, routeFileMappings, llmParser)
    mappings = append(mappings, cgMappings...)
    unclaimed = subtractChangedFiles(unclaimed, cgClaimed)

    // Phase 4: embedding（现有）
    embedMappings, _ := idx.runEmbedPhase(unclaimed)
    mappings = append(mappings, embedMappings...)

    return idx.writeMapFile(mappings, "hybrid")
}
```

---

## 新增文件

### `internal/rbt/callgraph_go.go`

```go
// GoCallGraphBuilder 使用 go/packages + go/callgraph 做类型感知分析。
// 与 V2 的逐文件 CallGraphBuilder 独立：一次性加载整个 module。
type GoCallGraphBuilder struct {
    SrcDir string
    Algo   string // "rta"（默认）| "pta"
}

// BuildAndTrace 一步完成：加载 packages → SSA → 调用图 → BFS 追踪。
// unclaimed: 待解析的变更文件（只处理 .go 文件）
// routeFileMappings: Phase 1 产出的路由文件映射
// maxDepth: BFS 深度上限（0 = 动态截断）
// 返回：找到的路由映射、已解决的 unclaimed 文件
// 任何内部错误（编译失败、go/packages 错误）返回 (nil, nil, err)，调用方静默 fallback
func (b *GoCallGraphBuilder) BuildAndTrace(
    unclaimed []ChangedFile,
    routeFileMappings map[string][]RouteMapping,
    maxDepth int,
) ([]RouteMapping, []ChangedFile, error)
```

**实现要点：**

1. 用 `go/packages.Load` 加载 `./...`，设置 `NeedSyntax | NeedTypes | NeedTypesInfo | NeedDeps`
2. 用 `golang.org/x/tools/go/ssa/ssautil.AllPackages` + `prog.Build()` 构建 SSA
3. 根据 `Algo` 选择算法：
   - `"rta"`：`golang.org/x/tools/go/callgraph/rta.Analyze(mains, true)`
   - `"pta"`：`golang.org/x/tools/go/pointer.Analyze(config)`
4. 遍历调用图，构建 `file::funcName → []callerNode` 倒排图（与 V2 `CallGraph.Edges` 结构相同）
5. BFS 从 unclaimed `.go` 文件的函数出发，追踪至 `routeFileMappings` 中的节点
6. 返回 `RouteMapping{Via: via, Confidence: confidence}` 及 covered files

**Via / Confidence：**

| Algo | Via | Confidence |
|------|-----|------------|
| rta | `"go-callgraph"` | `0.9` |
| pta | `"go-callgraph-pta"` | `0.95` |

### `internal/rbt/callgraph_go_test.go`

使用 `testdata/callgraph_go/` fixture（真实可编译的 Go module）：

```
testdata/callgraph_go/
  go.mod          // module testapp
  main.go         // main() → router.Register()
  handler.go      // Register() → userService.CreateUser(repo)
  service.go      // CreateUser(repo UserRepo) → repo.Save()
  repo.go         // type UserRepo interface { Save(...) error }
  repo_impl.go    // type MySQLRepo struct{}; func (r *MySQLRepo) Save(...) error
```

fixture 专门覆盖接口派发：`service.go` 通过 `UserRepo` 接口调用，V2 在此断链，V3 能追踪。

| 测试 | 覆盖点 |
|------|--------|
| `TestGoCallGraphBuilder_RTA_TracesInterface` | service.go 改动 → 接口 → 找到路由 |
| `TestGoCallGraphBuilder_Fallback_WhenNoGoMod` | 无 go.mod → 返回空，不报错 |
| `TestGoCallGraphBuilder_DepthCap` | maxDepth=1，3 层链截断 |
| `TestRunGoCallGraphPhase_IntegratesWithIndexer` | phase 方法返回正确 claimed files |

---

## 修改现有文件

### `internal/rbt/indexer.go`

`Indexer` struct 新增字段：

```go
Algo string // "rta"（默认）| "pta"；空字符串等同于 "rta"
```

新增 phase 方法：

```go
func (idx *Indexer) runGoCallGraphPhase(
    unclaimed []ChangedFile,
    routeFileMappings map[string][]RouteMapping,
) ([]RouteMapping, []ChangedFile)
```

内部构造 `GoCallGraphBuilder{SrcDir: idx.SrcDir, Algo: idx.Algo}`，调用 `BuildAndTrace`；错误时返回空（静默 fallback）。

### `cmd/rbt_index.go`

新增 flag：

```
--algo string   Go call graph algorithm: rta or pta (default "rta")
```

`runRBTIndex` 中读取并传入 `Indexer.Algo`。

---

## 依赖

```
golang.org/x/tools  (go get golang.org/x/tools@latest)
```

标准 Go 生态工具链，无外部服务依赖，无 API key 要求。

---

## 测试策略

### 单元测试

见上文 `callgraph_go_test.go` 说明。

### 验收测试（新增 AT-064 ~ AT-066）

| ID | Scenario | Command | Expected |
|----|----------|---------|----------|
| AT-064 | `--algo` flag 注册 | `caseforge rbt index --help` | `--algo` 出现 |
| AT-065 | RTA 追踪接口调用 | `caseforge rbt index --strategy hybrid --algo rta`（fixture Go module）| report 含 `via: go-callgraph` 条目 |
| AT-066 | 非 Go 项目正常降级 | `caseforge rbt index --strategy hybrid`（无 go.mod 的目录）| 正常运行，无错误 |

---

## 不在范围内

- `reflect` / `unsafe` 动态调用
- 非 Go 语言（由 V2 负责）
- 调用图结果持久化缓存（可作 V3.1）
- `go/callgraph/vta`（Refinement-based，可作 V3.1 的第三个 algo 选项）
- 增量分析（仅分析变更 package）

---

## 文件变更汇总

| 文件 | 动作 |
|------|------|
| `internal/rbt/callgraph_go.go` | 新建 |
| `internal/rbt/callgraph_go_test.go` | 新建 |
| `internal/rbt/testdata/callgraph_go/` | 新建（go.mod + 5个 .go 文件） |
| `internal/rbt/indexer.go` | 修改：新增 `Algo` 字段 + `runGoCallGraphPhase` |
| `internal/rbt/indexer_test.go` | 修改：新增 `runGoCallGraphPhase` 测试 |
| `cmd/rbt_index.go` | 修改：新增 `--algo` flag |
| `docs/acceptance/acceptance-tests.md` | 修改：新增 AT-064~066 |
| `scripts/acceptance.sh` | 修改：新增对应检查 |
| `CLAUDE.md` | 修改：scenario 数更新到 66 |
| `go.mod` / `go.sum` | 修改：新增 `golang.org/x/tools` |
