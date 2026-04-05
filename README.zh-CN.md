# CaseForge

> 基于 OpenAPI 规范的 AI 驱动 HTTP API 测试用例生成工具

[![CI](https://github.com/testmind-hq/caseforge/actions/workflows/ci.yml/badge.svg)](https://github.com/testmind-hq/caseforge/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

[English](README.md)

---

## 简介

CaseForge 读取你的 OpenAPI 规范，自动生成结构化、可追溯的测试用例，覆盖正常路径、边界值、边缘场景和 OWASP 安全测试场景。支持多种输出格式，并可直接对 API 执行测试。

## 功能特性

- **AI 驱动生成** — 支持 Anthropic、OpenAI、Gemini，以及任何 OpenAI 兼容 API（DeepSeek、Qwen、Moonshot、Azure）
- **多种输出格式** — Hurl、k6、Postman Collection v2.1、Markdown、CSV
- **OWASP 安全测试** — 自动生成注入攻击、认证绕过、数据暴露等安全测试用例
- **Spec Lint** — 对 OpenAPI 规范进行质量检查，支持可配置的严重级别阈值和 JSON 输出
- **Spec Diff** — 对比两个规范版本，自动分类破坏性变更，并为破坏性操作生成测试用例
- **风险驱动测试（RBT）** — 通过静态分析检测最近 git 变更影响的 API 操作
- **测试用例评分** — 从覆盖率、方法论、优先级分布等维度对测试用例进行质量评分
- **自然语言输入** — `ask` 命令支持从纯文本描述生成测试用例
- **平台导出** — 导出为 Allure、Xray 或 TestRail 格式
- **Webhook 推送** — 生成完成时触发 `on_generate` / `on_run_complete` 事件通知
- **Watch 模式** — 监听规范文件变更，自动重新生成测试用例
- **断点续跑** — `gen` 中断后可从已完成的操作处继续
- **动态 API 探索（DEA）** — 主动探测在线 API，发现隐式验证规则
- **重复检测** — 检测并清理结构相似的冗余测试用例
- **CI 工作流生成** — 生成 GitHub Actions、GitLab CI、Jenkins 或 Shell 工作流配置
- **MCP 服务器** — 将 CaseForge 作为 MCP 工具暴露，可接入 AI 智能体工作流
- **交互式引导** — `onboard` 命令引导完成全部初始化配置
- **纯算法模式** — 无需 LLM Key，基于配对分析、边界值分析和组合分析生成测试用例

## 安装

### Homebrew（macOS / Linux）

```bash
brew tap testmind-hq/tap
brew install caseforge
```

### Go 安装

```bash
go install github.com/testmind-hq/caseforge@latest
```

### 从源码构建

```bash
git clone https://github.com/testmind-hq/caseforge.git
cd caseforge
go build -o caseforge .
```

## 快速开始

```bash
# 交互式初始化（首次使用推荐）
caseforge onboard

# 检查环境依赖
caseforge doctor

# 从 OpenAPI 规范生成测试用例
caseforge gen --spec openapi.yaml --format hurl

# 执行生成的测试
caseforge run --cases ./cases --target http://localhost:8080

# 检查规范质量
caseforge lint --spec openapi.yaml
```

## 命令说明

### 核心命令

| 命令 | 描述 |
|------|------|
| `gen` | 从 OpenAPI 规范生成测试用例 |
| `run` | 执行生成的测试用例（hurl 或 k6） |
| `ask` | 从自然语言描述生成测试用例 |
| `lint` | 检查 OpenAPI 规范质量 |
| `diff` | 对比两个规范版本，分类破坏性变更 |
| `score` | 对生成的测试用例进行质量评分 |

### 分析命令

| 命令 | 描述 |
|------|------|
| `rbt` | 风险驱动测试：检测最近 git 变更影响的 API 操作 |
| `rbt index` | 分析源码，自动生成 `caseforge-map.yaml` |
| `explore` | 主动探测在线 API，发现隐式验证规则 |
| `stats` | 显示 cases 目录的测试用例统计信息 |
| `dedupe` | 检测并清理重复的测试用例 |

### 工作流命令

| 命令 | 描述 |
|------|------|
| `watch` | 监听规范文件变更，自动重新生成 |
| `suite create` | 创建测试套件编排文件 `suite.json` |
| `suite validate` | 验证 `suite.json` 与 `index.json` 的一致性 |
| `export` | 导出 `index.json` 为 Allure / Xray / TestRail 格式 |
| `ci init` | 生成 CI 工作流配置文件 |

### 工具命令

| 命令 | 描述 |
|------|------|
| `onboard` | 交互式初始化向导 |
| `init` | 在当前目录写入 `.caseforge.yaml` 配置 |
| `config show` | 打印当前生效的配置 |
| `doctor` | 检查环境依赖 |
| `mcp` | 启动 CaseForge MCP 服务器（stdio 传输） |
| `pairwise` | 计算给定参数的配对组合 |
| `fake` | 为 JSON Schema 生成 Fake 数据 |
| `completion` | 生成 Shell 补全脚本 |

---

## 命令参数参考

### `caseforge gen`

```
--spec string         OpenAPI 规范文件路径或 URL（必填）
--output string       输出目录（默认：./cases）
--format string       hurl | markdown | csv | postman | k6（默认：hurl）
--no-ai               禁用 LLM，使用纯算法模式
--technique string    只运行指定的测试技术，逗号分隔
                      如 equivalence_partitioning,boundary_value
--priority string     按最低优先级过滤输出：P0|P1|P2|P3
--operations string   只处理指定的 operationId（逗号分隔，默认全部）
--concurrency int     并发处理的操作数量（默认：1）
--resume              从上次中断处继续，跳过已完成的操作
```

### `caseforge run`

```
--cases string    包含测试文件的目录（必填）
--format string   hurl | k6（默认：hurl）
--target string   API 基础 URL，如 http://localhost:8080
--var key=value   注入到测试文件的变量（可重复）
--output string   写入 run-report.json 的目录
```

### `caseforge lint`

```
--spec string           OpenAPI 规范文件路径或 URL（必填）
--min-score int         规范评分低于该阈值时失败（0 = 禁用）
--format string         terminal | json（默认：terminal）
--output string         将 lint-report.json 写入此目录
--skip-rules string     跳过的规则 ID，逗号分隔，如 L001,L003
```

### `caseforge diff`

```
--old string        旧规范文件（必填）
--new string        新规范文件（必填）
--cases string      Cases 目录；读取 index.json 以推断受影响的用例
--format string     text | json（默认：text）
--gen-cases string  为破坏性操作生成测试用例到此目录
```

### `caseforge score`

```
--cases string    包含 index.json 的目录（默认：./cases）
--format string   terminal | json（默认：terminal）
```

### `caseforge rbt`

```
--spec string       OpenAPI 规范文件（必填）
--cases string      包含测试用例 JSON 文件的目录（默认：./cases）
--src string        源码根目录（默认：./）
--base string       用于 diff 的基准 git 引用（默认：HEAD~1）
--head string       用于 diff 的目标 git 引用（默认：HEAD）
--generate          为高风险未覆盖操作生成测试用例
--no-ai             路由推断和生成均使用纯算法模式
--gen-format string 生成用例的格式：hurl|postman|k6|markdown|csv
--output string     写入 rbt-report.json 的目录（默认：./reports）
--format string     terminal | json（默认：terminal）
--fail-on string    风险 >= 指定级别时退出非零：none|low|medium|high（默认：high）
--map string        caseforge-map.yaml 文件路径
--dry-run           跳过 git diff，将所有操作报告为 risk=none
```

### `caseforge rbt index`

```
--spec string       OpenAPI 规范文件（必填）
--src string        要分析的源码根目录（默认：./）
--out string        输出 Map 文件路径（默认：caseforge-map.yaml）
--strategy string   llm | embed | hybrid（默认：llm）
--overwrite         覆盖已有 Map 文件
--depth int         调用图遍历深度（0 = 动态）
--algo string       Go 调用图算法：rta | vta（默认：rta）
```

### `caseforge ask`

```
--output string   输出目录（默认：./cases）
--format string   hurl | markdown | csv | postman | k6（默认：hurl）
```

### `caseforge suite create`

```
--id string       套件 ID（必填）
--title string    套件标题（必填）
--kind string     sequential | chain（默认：sequential）
--cases string    要包含的用例 ID，逗号分隔
--output string   输出文件路径（默认：suite.json）
```

### `caseforge suite validate`

```
--suite string    suite.json 文件路径（必填）
--cases string    包含 index.json 的 cases 目录
```

### `caseforge explore`

```
--spec string       OpenAPI 规范文件
--target string     目标 API 基础 URL（非 --dry-run 时必填）
--max-probes int    每次运行最大 HTTP 探测数（默认：50）
--output string     写入 dea-report.json 的目录（默认：./reports）
--dry-run           只生成假设，不执行探测
```

### `caseforge export`

```
--cases string    包含 index.json 的目录（必填）
--format string   allure | xray | testrail（必填）
--output string   输出目录（默认：./export）
```

### `caseforge ci init`

```
--platform string   github-actions | gitlab-ci | jenkins | shell（默认：github-actions）
--spec string       生成工作流中使用的规范路径（默认：openapi.yaml）
--output string     输出文件路径（默认由平台决定）
--force             不提示直接覆盖已有文件
```

---

## 配置

在项目根目录创建 `.caseforge.yaml`：

```yaml
ai:
  provider: anthropic          # anthropic | openai | openai-compat | gemini | noop
  model: claude-sonnet-4-6     # 对应 provider 的模型名称
  # api_key: ...               # 也可以通过环境变量设置（见下表）
  # base_url: ...              # 仅 openai-compat 使用（如 https://api.deepseek.com/v1）

output:
  default_format: hurl         # hurl | markdown | csv | postman | k6
  dir: ./cases

lint:
  fail_on: error               # error | warning | info

# Webhook 通知（可选）
webhooks:
  - url: https://hooks.example.com/caseforge
    events: [on_generate, on_run_complete]
    secret: your-hmac-secret   # 使用 X-CaseForge-Signature-256 签名请求
    timeout_seconds: 10
    max_retries: 3
```

### LLM Provider 配置

| Provider | `ai.provider` | 环境变量 |
|----------|--------------|---------|
| Anthropic（默认） | `anthropic` | `ANTHROPIC_API_KEY` |
| OpenAI | `openai` | `OPENAI_API_KEY` |
| DeepSeek / Qwen / Azure | `openai-compat` | `OPENAI_API_KEY` + `ai.base_url` |
| Google Gemini | `gemini` | `GEMINI_API_KEY` 或 `GOOGLE_API_KEY` |
| 无 AI | `noop` | — |

### Webhook 事件

| 事件 | 触发时机 |
|------|---------|
| `on_generate` | 每个操作生成完成时（包含方法、路径、用例数量） |
| `on_run_complete` | 整个 `gen` 运行完成后（包含总用例数、输出目录） |

配置了 `secret` 时，请求使用 HMAC-SHA256 签名，头部为 `X-CaseForge-Signature-256: sha256=<hex>`。

---

## 测试技术

| 技术 | `--technique` 参数值 |
|------|---------------------|
| 等价类划分 | `equivalence_partitioning` |
| 边界值分析 | `boundary_value` |
| 判定表 | `decision_table` |
| 状态转换 | `state_transition` |
| 配对测试（IPOG） | `pairwise` |
| 幂等性 | `idempotency` |
| OWASP API Top 10 | `owasp_api_top10` |
| 分类树（MBT） | `classification_tree` |
| 正交数组 | `orthogonal_array` |
| 示例提取 | `example_extraction` |

---

## 依赖要求

- Go 1.26+（源码构建）
- [hurl](https://hurl.dev/docs/installation.html) — `caseforge run --format hurl` 所需
- [k6](https://k6.io/docs/get-started/installation/) — `caseforge run --format k6` 所需

运行 `caseforge doctor` 检查当前环境。

## 贡献

请参阅 [CONTRIBUTING.md](CONTRIBUTING.md)。

## 开源协议

Apache License 2.0 — 详见 [LICENSE](LICENSE)。
