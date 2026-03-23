# CaseForge

> AI-powered HTTP API test-case generator from OpenAPI specs

[![CI](https://github.com/testmind-hq/caseforge/actions/workflows/ci.yml/badge.svg)](https://github.com/testmind-hq/caseforge/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

[中文文档](README.zh-CN.md)

---

## What is CaseForge?

CaseForge reads your OpenAPI specification and generates structured, traceable test cases covering happy paths, edge cases, boundary values, and OWASP security scenarios. It outputs ready-to-run test files in multiple formats and can execute them against your API.

## Features

- **AI-powered generation** — uses Anthropic, OpenAI, Gemini, or any OpenAI-compatible API (DeepSeek, Qwen, Moonshot, Azure)
- **Multiple output formats** — Hurl, k6, Postman Collection v2.1, Markdown, CSV
- **OWASP security testing** — generates injection, auth bypass, and data exposure test cases
- **Spec linting** — validates OpenAPI specs for quality issues with configurable severity thresholds
- **Spec diff** — classifies breaking vs non-breaking changes between two spec versions
- **MCP server** — exposes CaseForge as an MCP tool for use in AI agent pipelines
- **Pure algorithm mode** — works without an LLM key using pairwise and boundary-value analysis

## Installation

### Homebrew (macOS / Linux)

```bash
brew tap testmind-hq/tap
brew install caseforge
```

### Go

```bash
go install github.com/testmind-hq/caseforge@latest
```

### From source

```bash
git clone https://github.com/testmind-hq/caseforge.git
cd caseforge
go build -o caseforge .
```

## Quick Start

```bash
# Check your environment
caseforge doctor

# Generate test cases from an OpenAPI spec
caseforge gen --spec openapi.yaml --format hurl

# Run the generated tests
caseforge run --cases ./cases --format hurl

# Lint the spec
caseforge lint --spec openapi.yaml
```

## Commands

| Command | Description |
|---------|-------------|
| `gen` | Generate test cases from an OpenAPI spec |
| `run` | Execute generated test cases (hurl or k6) |
| `lint` | Lint an OpenAPI spec for quality issues |
| `diff` | Compare two OpenAPI specs and classify breaking changes |
| `doctor` | Check environment dependencies |
| `mcp` | Start CaseForge as an MCP server (stdio transport) |
| `init` | Initialize CaseForge in the current directory |

### `caseforge gen`

```
--spec string     OpenAPI spec file or URL (required)
--output string   Output directory (default: ./cases)
--format string   Output format: hurl|markdown|csv|postman|k6 (default: hurl)
--no-ai           Disable LLM, use pure algorithm mode
```

### `caseforge run`

```
--cases string    Directory containing generated test files (required)
--format string   Test runner format: hurl|k6 (default: hurl)
```

### `caseforge lint`

```
--spec string       OpenAPI spec file or URL (required)
--min-score int     Fail if spec score is below this threshold (0 = disabled)
```

## Configuration

Create `.caseforge.yaml` in your project root:

```yaml
ai:
  provider: anthropic          # anthropic | openai | openai-compat | gemini | noop
  model: claude-sonnet-4-6     # model name for the chosen provider
  # api_key: ...               # or set via env var (see below)
  # base_url: ...              # openai-compat only (e.g. https://api.deepseek.com/v1)

output:
  default_format: hurl         # hurl | markdown | csv | postman | k6

lint:
  fail_on: error               # error | warning | info
```

### LLM Providers

| Provider | `ai.provider` | Env var |
|----------|--------------|---------|
| Anthropic (default) | `anthropic` | `ANTHROPIC_API_KEY` |
| OpenAI | `openai` | `OPENAI_API_KEY` |
| DeepSeek / Qwen / Azure | `openai-compat` | `OPENAI_API_KEY` + `ai.base_url` |
| Google Gemini | `gemini` | `GEMINI_API_KEY` or `GOOGLE_API_KEY` |
| No AI | `noop` | — |

## Requirements

- Go 1.26+ (build from source)
- [hurl](https://hurl.dev/docs/installation.html) — required for `caseforge run --format hurl`
- [k6](https://k6.io/docs/get-started/installation/) — required for `caseforge run --format k6`

Run `caseforge doctor` to verify your environment.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Apache License 2.0 — see [LICENSE](LICENSE).
