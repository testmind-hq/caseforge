# CaseForge Skill

Use this skill when asked to generate, score, analyse, or manage API test cases using CaseForge.

## Quick Start

```bash
# First time: interactive setup
caseforge onboard

# Generate test cases (AI-enhanced)
caseforge gen --spec openapi.yaml --output ./cases/

# Run generated cases against a live API
caseforge run --cases ./cases/ --target https://api.example.com

# Pure algorithm mode (no AI required)
caseforge gen --spec openapi.yaml --output ./cases/ --no-ai

# Check spec quality
caseforge lint --spec openapi.yaml
```

## Commands

### Core

| Command | Purpose |
|---------|---------|
| `gen` | Generate test cases from OpenAPI spec → `index.json` + format files |
| `run` | Execute `.hurl` or `.js` (k6) test files against a real server |
| `ask` | Generate test cases from a natural language description |
| `lint` | Check spec for quality issues (missing operationId, no 2xx, etc.) |
| `diff` | Compare two spec versions; classify breaking changes |
| `score` | Multi-dimensional quality scoring of generated test cases |

### Analysis

| Command | Purpose |
|---------|---------|
| `rbt` | Risk-based testing: detect which API ops are at risk from recent git changes |
| `rbt index` | Auto-generate `caseforge-map.yaml` by analysing source code |
| `explore` | Probe a live API to discover implicit validation rules (DEA) |
| `stats` | Show test case statistics (total, by technique, by priority, by format) |
| `dedupe` | Detect and optionally remove structurally duplicate test cases |

### Workflow

| Command | Purpose |
|---------|---------|
| `watch` | Watch spec file and regenerate on change |
| `suite create` | Create `suite.json` for cross-case DAG orchestration |
| `suite validate` | Validate `suite.json` against `index.json` |
| `export` | Export `index.json` to Allure / Xray / TestRail |
| `ci init` | Generate CI workflow config (GitHub Actions, GitLab CI, Jenkins, shell) |

### Utilities

| Command | Purpose |
|---------|---------|
| `onboard` | Interactive wizard: config, LLM provider, MCP server, skill install |
| `init` | Write `.caseforge.yaml` in current directory |
| `config show` | Print effective configuration |
| `doctor` | Check environment (hurl, k6, API keys) |
| `mcp` | Start CaseForge as an MCP server (stdio) |
| `pairwise` | Compute pairwise combinations without a spec |
| `fake` | Generate fake data for a JSON schema |

## Key Flags for `gen`

```
--technique string    Comma-separated techniques to run:
                        equivalence_partitioning, boundary_value,
                        decision_table, state_transition, pairwise,
                        idempotency, owasp, classification_tree,
                        orthogonal_array
--priority string     Minimum priority to include: P0 | P1 | P2 | P3
--operations string   Comma-separated operationIds (default: all)
--concurrency int     Parallel operations (default: 1)
--resume              Resume interrupted run from checkpoint
--no-ai               Algorithm-only mode
```

## Methodologies Applied

- **Equivalence Partitioning** — valid/invalid partitions for each field
- **Boundary Value Analysis** — min/max/min-1/max+1 for numeric and string fields
- **Decision Table** — one case per enum value combination
- **State Transition** — one case per state transition (requires AI annotation)
- **Pairwise (IPOG)** — covering array for 4+ independent parameters
- **Idempotency** — duplicate request case for write operations
- **OWASP Security** — injection (SQLi, XSS, path traversal), auth bypass, CORS
- **Classification Tree (MBT)** — structured classification of valid/invalid inputs
- **Orthogonal Array** — strength-2 orthogonal arrays for large parameter spaces

## Configuration (`.caseforge.yaml`)

```yaml
ai:
  provider: anthropic   # anthropic | openai | openai-compat | gemini | noop
  model: claude-sonnet-4-6

output:
  default_format: hurl  # hurl | markdown | csv | postman | k6
  dir: ./cases

lint:
  fail_on: error        # error | warning | info

# Optional: push events to a webhook
webhooks:
  - url: https://hooks.example.com/caseforge
    events: [on_generate, on_run_complete]
    secret: your-hmac-secret
```

## Webhook Events

| Event | Payload |
|-------|---------|
| `on_generate` | `{event, timestamp, operation: {id, method, path}, case_count}` |
| `on_run_complete` | `{event, timestamp, total_cases, output_dir}` |

Signed with `X-CaseForge-Signature-256: sha256=<hmac>` when `secret` is configured.

## Common Workflows

### Targeted generation

```bash
# Only generate for specific operations
caseforge gen --spec api.yaml --operations createUser,deleteUser

# Only security cases at P0/P1
caseforge gen --spec api.yaml --technique owasp --priority P1

# Resume after interruption
caseforge gen --spec api.yaml --resume
```

### CI integration

```bash
# Generate CI config
caseforge ci init --platform github-actions --spec openapi.yaml

# Full CI pipeline
caseforge gen --spec openapi.yaml --no-ai
caseforge lint --spec openapi.yaml --min-score 70
caseforge run --cases ./cases --target $API_URL
caseforge score --cases ./cases
```

### Risk-based testing

```bash
# Build source-to-route map (once per project)
caseforge rbt index --spec openapi.yaml --src ./src

# Assess risk for current branch changes
caseforge rbt --spec openapi.yaml --cases ./cases --generate
```

### Export to test management

```bash
caseforge export --cases ./cases --format allure --output ./allure-results
caseforge export --cases ./cases --format xray --output ./xray-export
```
