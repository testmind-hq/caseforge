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
- **Spec Lint** — 对 OpenAPI 规范进行质量检查，支持可配置的严重级别阈值
- **Spec Diff** — 对比两个规范版本，自动分类破坏性变更与非破坏性变更
- **MCP 服务器** — 将 CaseForge 作为 MCP 工具暴露，可接入 AI 智能体工作流
- **纯算法模式** — 无需 LLM Key，基于配对分析和边界值分析生成测试用例

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
# 检查环境依赖
caseforge doctor

# 从 OpenAPI 规范生成测试用例
caseforge gen --spec openapi.yaml --format hurl

# 执行生成的测试
caseforge run --cases ./cases --format hurl

# 检查规范质量
caseforge lint --spec openapi.yaml
```

## 命令说明

| 命令 | 描述 |
|------|------|
| `gen` | 从 OpenAPI 规范生成测试用例 |
| `run` | 执行生成的测试用例（hurl 或 k6） |
| `lint` | 检查 OpenAPI 规范的质量问题 |
| `diff` | 比较两个 OpenAPI 规范并分类破坏性变更 |
| `doctor` | 检查环境依赖 |
| `mcp` | 启动 CaseForge MCP 服务器（stdio 传输） |
| `init` | 在当前目录初始化 CaseForge |

### `caseforge gen`

```
--spec string     OpenAPI 规范文件路径或 URL（必填）
--output string   输出目录（默认：./cases）
--format string   输出格式：hurl|markdown|csv|postman|k6（默认：hurl）
--no-ai           禁用 LLM，使用纯算法模式
```

### `caseforge run`

```
--cases string    包含测试文件的目录（必填）
--format string   测试运行器格式：hurl|k6（默认：hurl）
```

### `caseforge lint`

```
--spec string       OpenAPI 规范文件路径或 URL（必填）
--min-score int     规范评分低于该阈值时失败（0 = 禁用）
```

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

lint:
  fail_on: error               # error | warning | info
```

### LLM Provider 配置

| Provider | `ai.provider` | 环境变量 |
|----------|--------------|---------|
| Anthropic（默认） | `anthropic` | `ANTHROPIC_API_KEY` |
| OpenAI | `openai` | `OPENAI_API_KEY` |
| DeepSeek / Qwen / Azure | `openai-compat` | `OPENAI_API_KEY` + `ai.base_url` |
| Google Gemini | `gemini` | `GEMINI_API_KEY` 或 `GOOGLE_API_KEY` |
| 无 AI | `noop` | — |

## 依赖要求

- Go 1.26+（源码构建）
- [hurl](https://hurl.dev/docs/installation.html) — `caseforge run --format hurl` 所需
- [k6](https://k6.io/docs/get-started/installation/) — `caseforge run --format k6` 所需

运行 `caseforge doctor` 检查当前环境。

## 贡献

请参阅 [CONTRIBUTING.md](CONTRIBUTING.md)。

## 开源协议

Apache License 2.0 — 详见 [LICENSE](LICENSE)。
