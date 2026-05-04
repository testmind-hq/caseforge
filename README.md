# CaseForge

> Spec-driven HTTP API test-case generator — turn OpenAPI into runnable Hurl, k6, or Postman cases

[![CI](https://github.com/testmind-hq/caseforge/actions/workflows/ci.yml/badge.svg)](https://github.com/testmind-hq/caseforge/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

[中文文档](README.zh-CN.md)

---

## What is CaseForge?

CaseForge reads your OpenAPI specification and generates structured, traceable test cases covering happy paths, edge cases, boundary values, and OWASP security scenarios. It outputs ready-to-run test files in multiple formats and can execute them against your API.

It works as a pure algorithmic generator out of the box (pairwise, boundary-value, combinatorial). Optionally, you can plug in an LLM (Anthropic / OpenAI / Gemini / any OpenAI-compatible API) to enrich edge-case discovery and mine response-body constraints.

## Features

- **Multiple LLM providers (optional)** — Anthropic, OpenAI, Gemini, or any OpenAI-compatible API (DeepSeek, Qwen, Moonshot, Azure). Disable entirely for pure-algorithm mode.
- **Multiple output formats** — Hurl, k6, Postman Collection v2.1, Markdown, CSV
- **OWASP security testing** — injection, auth bypass, and data exposure test cases
- **Spec linting** — validates OpenAPI specs with configurable severity thresholds and JSON output
- **Spec diff** — classifies breaking vs non-breaking changes; auto-generates cases for breaking ops
- **Risk-based testing** — detects which API operations are at risk from recent git changes via static analysis
- **Test case scoring** — multi-dimensional quality scoring with named coverage scenario tracking (breadth, boundary, security, execution, status coverage)
- **Natural language input** — `ask` generates cases from a plain-text description
- **Platform export** — exports to Allure, Xray, or TestRail
- **Webhook push** — fires `on_generate` / `on_run_complete` events to configured endpoints
- **Watch mode** — regenerates cases whenever the spec file changes
- **Checkpoint resume** — resumes interrupted `gen` runs from where they left off
- **Dynamic API exploration** — probes a live API to discover implicit validation rules (DEA)
- **Duplicate detection** — finds and removes structurally similar test cases
- **CI scaffolding** — generates GitHub Actions, GitLab CI, Jenkins, or shell workflow configs
- **MCP server** — exposes CaseForge as an MCP tool for AI agent pipelines
- **Onboarding wizard** — interactive `onboard` command walks through full setup in minutes
- **Auth bootstrap** — `--auth-bootstrap` prepends an auth setup step to all secured-endpoint cases so every technique works out of the box against authenticated APIs
- **Response body oracles** — `--with-oracles` uses two-step LLM prompting (Observation-Confirmation) to mine response body constraints and inject them as assertions
- **Coverage gap filling** — `score --fill-gaps` detects operations missing 2xx or 4xx coverage and auto-generates cases to close the gaps
- **Failure classification** — failed `run` cases are automatically tagged `server_error` / `missing_validation` / `auth_failure` / `security_regression`
- **Conformance checking** — `conformance` command mines oracle constraints for all operations and reports spec-vs-implementation mismatches against a live API
- **Pure algorithm mode** — works without an LLM key using pairwise, boundary-value, and combinatorial analysis

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
# Interactive setup (recommended for first use)
caseforge onboard

# Check your environment
caseforge doctor

# Generate test cases from an OpenAPI spec
caseforge gen --spec openapi.yaml --format hurl

# Run the generated tests
caseforge run --cases ./cases --target http://localhost:8080

# Lint the spec
caseforge lint --spec openapi.yaml
```

## Commands

### Core

| Command | Description |
|---------|-------------|
| `gen` | Generate test cases from an OpenAPI spec |
| `run` | Execute generated test cases (hurl or k6) |
| `ask` | Generate test cases from a natural language description |
| `lint` | Lint an OpenAPI spec for quality issues |
| `diff` | Compare two OpenAPI specs and classify breaking changes |
| `score` | Score the quality of generated test cases |
| `conformance` | Check spec-vs-implementation conformance using LLM-mined response body constraints |

### Analysis

| Command | Description |
|---------|-------------|
| `rbt` | Risk-based testing: assess which operations are at risk from recent git changes |
| `rbt index` | Auto-generate `caseforge-map.yaml` by analysing source code |
| `explore` | Dynamically probe a live API and infer implicit validation rules |
| `stats` | Show test case statistics for a cases directory |
| `dedupe` | Detect and optionally remove duplicate test cases |

### Workflow

| Command | Description |
|---------|-------------|
| `chain` | Generate multi-step chain cases via BFS over the dependency graph |
| `watch` | Watch a spec file and regenerate cases on change |
| `suite create` | Create a `suite.json` orchestration file |
| `suite validate` | Validate a `suite.json` against its `index.json` |
| `export` | Export `index.json` to Allure, Xray, or TestRail format |
| `ci init` | Generate a CI workflow config (GitHub Actions, GitLab CI, Jenkins, shell) |

### Utilities

| Command | Description |
|---------|-------------|
| `onboard` | Interactive setup wizard |
| `init` | Write a `.caseforge.yaml` config file |
| `config show` | Print the effective configuration |
| `doctor` | Check environment dependencies |
| `mcp` | Start CaseForge as an MCP server (stdio transport) |
| `pairwise` | Compute pairwise combinations for given parameters |
| `fake` | Generate fake data for a JSON schema |
| `completion` | Generate shell completion scripts |

---

## Command Reference

### `caseforge gen`

```
--spec string         OpenAPI spec file or URL (required)
--output string       Output directory (default: ./cases)
--format string       hurl | markdown | csv | postman | k6 (default: hurl)
--no-ai               Disable LLM; use pure algorithm mode
--technique string    Only run named techniques, comma-separated
                      e.g. equivalence_partitioning,boundary_value
--priority string     Filter output by minimum priority: P0|P1|P2|P3
--operations string   Comma-separated operationIds to process (default: all)
--concurrency int     Operations processed in parallel (default: 1)
--resume              Resume an interrupted run; skip completed operations
--tuple-level int     N-way coverage level for pairwise (2=pairwise, 3=3-way, default 2)
--seed int            Seed for deterministic generation (0 = random)
--max-cases-per-op int  Cap cases per operation by priority (0 = unlimited)
--include-path string   Regex to include operations by path (e.g. '^/users')
--exclude-path string   Regex to exclude operations by path (e.g. '^/admin')
--include-tag string    Comma-separated OpenAPI tags to include (e.g. 'users,auth')
--exclude-tag string    Comma-separated OpenAPI tags to exclude (e.g. 'deprecated')
--auth-bootstrap      Wrap all secured-endpoint cases with an auth setup step
--with-oracles        Mine response body constraints via LLM and inject as assertions (requires LLM)
```

### `caseforge run`

```
--cases string    Directory containing generated test files (required)
--format string   hurl | k6 (default: hurl)
--target string   API base URL, e.g. http://localhost:8080
--var key=value   Variables injected into test files (repeatable)
--output string   Directory to write run-report.json
```

### `caseforge lint`

```
--spec string           OpenAPI spec file or URL (required)
--min-score int         Fail if spec score is below threshold (0 = disabled)
--format string         terminal | json (default: terminal)
--output string         Write lint-report.json to this directory
--skip-rules string     Comma-separated rule IDs to skip, e.g. L001,L003
```

### `caseforge diff`

```
--old string        Old spec file (required)
--new string        New spec file (required)
--cases string      Cases directory; reads index.json to infer affected cases
--format string     text | json (default: text)
--gen-cases string  Generate test cases for breaking operations into this directory
```

### `caseforge score`

```
--cases string    Directory containing index.json (default: ./cases)
--format string   terminal | json (default: terminal)
--fill-gaps       Auto-generate cases for operations missing 2xx or 4xx coverage
--spec string     OpenAPI spec path (required for --fill-gaps)
--min-score int   Exit non-zero if overall score is below threshold (0 = disabled)
--save-history    Append score to .caseforge-conformance.json for trend tracking
```

### `caseforge conformance`

```
--spec string     OpenAPI spec file (required)
--target string   API base URL to test against (required)
--output string   Directory to write conformance-report.json (optional)
```

### `caseforge rbt`

```
--spec string       OpenAPI spec file (required)
--cases string      Directory containing test case JSON files (default: ./cases)
--src string        Source code root directory (default: ./)
--base string       Base git ref for diff (default: HEAD~1)
--head string       Head git ref for diff (default: HEAD)
--generate          Generate test cases for high-risk uncovered operations
--no-ai             Algorithm-only mode for both route inference and generation
--gen-format string Format for generated cases: hurl|postman|k6|markdown|csv
--output string     Directory to write rbt-report.json (default: ./reports)
--format string     terminal | json (default: terminal)
--fail-on string    Exit non-zero if risk >= level: none|low|medium|high (default: high)
--map string        Path to caseforge-map.yaml
--dry-run           Skip git diff; report all operations as risk=none
```

### `caseforge rbt index`

```
--spec string       OpenAPI spec file (required)
--src string        Source code root to analyse (default: ./)
--out string        Output map file (default: caseforge-map.yaml)
--strategy string   llm | embed | hybrid (default: llm)
--overwrite         Overwrite existing map file
--depth int         Call graph traversal depth (0 = dynamic)
--algo string       Go call graph algorithm: rta | vta (default: rta)
```

### `caseforge ask`

```
--output string   Output directory (default: ./cases)
--format string   hurl | markdown | csv | postman | k6 (default: hurl)
```

### `caseforge suite create`

```
--id string       Suite ID (required)
--title string    Suite title (required)
--kind string     sequential | chain (default: sequential)
--cases string    Comma-separated case IDs to include
--output string   Output file path (default: suite.json)
```

### `caseforge suite validate`

```
--suite string    Path to suite.json (required)
--cases string    Cases directory containing index.json
```

### `caseforge chain`

```
--spec string         OpenAPI spec file or URL (required)
--depth int           Maximum chain depth 1..4 (default: 2)
--output string       Output directory (default: ./chains)
--format string       hurl | markdown | csv | postman | k6 (default: hurl)
--data-pool string    JSON data pool file written by explore --export-pool
--seed-postman string Postman Collection v2.1 JSON to seed the data pool
```

Chain cases follow OpenAPI Links to wire producer `$response.body` fields into
consumer path/query parameters, and auto-append a DELETE teardown step for
depth-2 chains where the consumer is not a DELETE operation.

### `caseforge explore`

```
--spec string              OpenAPI spec file
--target string            Target API base URL (required without --dry-run)
--max-probes int           Maximum HTTP probes per run (default: 50)
--output string            Directory to write dea-report.json (default: ./reports)
--dry-run                  Seed hypotheses only; do not execute probes
--export-pool string       Write observed 2xx response field values to a JSON data pool file
--prioritize-uncovered     Two-pass scheduling: breadth-scan all ops in pass 1, then
                           focus remaining budget on operations that did not return 2xx
--max-failures int         Stop after discovering this many rules (0 = unlimited)
--include-path string      Regex to include operations by path (e.g. '^/users')
--exclude-path string      Regex to exclude operations by path (e.g. '^/admin')
--include-tag string       Comma-separated OpenAPI tags to include (e.g. 'users,auth')
--exclude-tag string       Comma-separated OpenAPI tags to exclude (e.g. 'deprecated')
```

The data pool written by `--export-pool` can be loaded into `caseforge chain`
via `--data-pool` to seed realistic field values into generated chain probes.

### `caseforge export`

```
--cases string    Directory containing index.json (required)
--format string   allure | xray | testrail (required)
--output string   Output directory (default: ./export)
```

### `caseforge ci init`

```
--platform string   github-actions | gitlab-ci | jenkins | shell (default: github-actions)
--spec string       OpenAPI spec path used in the generated workflow (default: openapi.yaml)
--output string     Output file path (default: platform-specific)
--force             Overwrite existing file without prompting
```

### `caseforge dedupe`

```
--cases string        Directory of test case JSON files (default: ./cases)
--threshold float     Jaccard similarity threshold 0.0–1.0 (default: 0.9)
--merge               Auto-delete lower-scoring duplicates
--dry-run             Report what would be deleted without deleting
--format string       terminal | json (default: terminal)
```

### `caseforge watch`

```
--spec string     OpenAPI spec file to watch (required, local file only)
--output string   Output directory (default: ./cases)
--format string   hurl | k6 | postman | markdown | csv
--no-ai           Disable LLM
```

### `caseforge stats`

```
--cases string    Directory containing index.json (default: ./cases)
--format string   terminal | json (default: terminal)
```

---

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
  dir: ./cases

lint:
  fail_on: error               # error | warning | info

# Webhook notifications (optional)
webhooks:
  - url: https://hooks.example.com/caseforge
    events: [on_generate, on_run_complete]
    secret: your-hmac-secret   # signs requests with X-CaseForge-Signature-256
    timeout_seconds: 10
    max_retries: 3
```

### LLM Providers

| Provider | `ai.provider` | Env var |
|----------|--------------|---------|
| Anthropic (default) | `anthropic` | `ANTHROPIC_API_KEY` |
| OpenAI | `openai` | `OPENAI_API_KEY` |
| DeepSeek / Qwen / Azure | `openai-compat` | `OPENAI_API_KEY` + `ai.base_url` |
| Google Gemini | `gemini` | `GEMINI_API_KEY` or `GOOGLE_API_KEY` |
| No AI | `noop` | — |

### Webhook Events

| Event | Fires when |
|-------|-----------|
| `on_generate` | Each operation completes generation (includes method, path, case count) |
| `on_run_complete` | The full `gen` run finishes (includes total cases, output directory) |

Requests are signed with HMAC-SHA256 when `secret` is set. Verify with:

```
X-CaseForge-Signature-256: sha256=<hex>
```

---

## Techniques

| Technique | Flag value |
|-----------|-----------|
| Equivalence Partitioning | `equivalence_partitioning` |
| Boundary Value Analysis | `boundary_value` |
| Decision Table | `decision_table` |
| State Transition | `state_transition` |
| Pairwise (IPOG) | `pairwise` |
| Idempotency | `idempotency` |
| OWASP API Top 10 (spec-based) | `owasp_api_top10` |
| OWASP API Top 10 (LLM-annotated) | `owasp_api_top10_spec` |
| Classification Tree (MBT) | `classification_tree` |
| Orthogonal Array | `orthogonal_array` |
| Example Extraction | `example_extraction` |
| Positive Examples | `positive_examples` |
| Isolated Negative | `isolated_negative` |
| Required Field Omission | `required_omission` |
| Field Boundary | `field_boundary` |
| Schema Violation | `schema_violation` |
| Constraint Mutation | `constraint_mutation` |
| Variable Irrelevance | `variable_irrelevance` |
| Mutation | `mutation` |
| Type Coercion | `type_coercion` |
| Unicode Fuzzing | `unicode_fuzzing` |
| Mass Assignment | `mass_assignment` |
| IDOR | `idor` |
| Semantic Annotation (nullable/readOnly/writeOnly) | `semantic_annotation` |
| Auth Chain | `auth_chain` |
| CRUD Chain | `chain_crud` |
| Chain Sequence (Jaccard similarity) | `chain_sequence` |
| Business Rule Violation | `business_rule_violation` |

---

## Requirements

- Go 1.26+ (build from source)
- [hurl](https://hurl.dev/docs/installation.html) — required for `caseforge run --format hurl`
- [k6](https://k6.io/docs/get-started/installation/) — required for `caseforge run --format k6`

Run `caseforge doctor` to verify your environment.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Acknowledgements

caseforge's design was informed by the following projects and academic research in the API testing space:

**Open-source projects (concept-level references; no source code derived):**

- [Schemathesis](https://github.com/schemathesis/schemathesis) — property-based testing patterns
- [CATS](https://github.com/Endava/cats) — fuzzing techniques
- [EvoMaster](https://github.com/EMResearch/EvoMaster) — coverage metric definition
- [Tcases](https://github.com/Cornutum/tcases) — isolated negative testing and N-way coverage principles
- [RESTler](https://github.com/microsoft/restler-fuzzer) — dependency-graph and N-step chain testing
- [Portman](https://github.com/apideck-libraries/portman) — semantic annotation and field-boundary patterns
- [Microcks](https://github.com/microcks/microcks) — HAR-based traffic import and conformance gating

**Academic research (concept-level references):**

- RBCTest — Observation-Confirmation prompting pattern for response-body oracle mining
- AutoRestTest — failure classification and coverage-gap filling
- RESTifAI — LLM-driven business-rule violation and chain-sequence inference

**Standards:**

- [OWASP API Security Top 10](https://owasp.org/API-Security/) — security category structure

None of these projects' source code is embedded in caseforge. See [NOTICE](NOTICE) for full attribution and the explicit "no source derived" declaration.

## License

Apache License 2.0 — see [LICENSE](LICENSE).
