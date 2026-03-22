# CaseForge Skill

Use this skill when asked to generate API test cases from an OpenAPI spec.

## Quick Start

```bash
# Generate test cases (AI-enhanced)
caseforge gen --spec openapi.yaml --output ./cases/

# Run generated cases
caseforge run --cases ./cases/ --var base_url=https://api.example.com

# Pure algorithm mode (no AI)
caseforge gen --spec openapi.yaml --output ./cases/ --no-ai

# Check spec quality
caseforge lint --spec openapi.yaml
```

## Commands

| Command | Purpose |
|---------|---------|
| `gen` | Generate test cases from OpenAPI spec → index.json + .hurl files |
| `lint` | Check spec for quality issues (missing operationId, no 2xx response, etc.) |
| `run` | Execute .hurl files against a real server via hurl |
| `doctor` | Check environment (hurl installed, API key configured) |
| `init` | Create `.caseforge.yaml` config in current directory |
| `pairwise` | Compute pairwise combinations without a spec |
| `fake` | Generate fake data for a JSON schema |

## Methodologies Applied

- **Equivalence Partitioning** — valid/invalid partitions for each field
- **Boundary Value Analysis** — min/max/min-1/max+1 for numeric and string fields
- **Decision Table** — one case per enum value combination
- **State Transition** — one case per state transition (requires AI annotation)
- **Pairwise (IPOG)** — covering array for 4+ independent parameters
- **Idempotency** — duplicate request case for write operations
