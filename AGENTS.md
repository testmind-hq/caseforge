# CaseForge — Contribution Guidelines

These guidelines apply to anyone (human or AI assistant) working on caseforge.

## Repository Structure

- `cmd/` — CLI entry points and command handlers
- `internal/` — internal packages (not importable by external code)
- `docs/` — user-facing documentation
- `scripts/` — automation scripts (acceptance, release)
- `testdata/` — test fixtures
- `examples/` — usage examples

## Acceptance Tests

Every meaningful change must run the acceptance suite before commit:

```bash
./scripts/acceptance.sh
```

**Rules:**
- All scenarios must pass before pushing to `main`.
- When adding a new command or feature, add corresponding scenario(s) to `docs/acceptance/acceptance-tests.md` and a matching `contains` / `run` call in `scripts/acceptance.sh` **in the same commit**.
- If a scenario fails, determine the root cause:
  - **Code bug** → fix the code, re-run the suite.
  - **Environment issue** (no API key, no running server) → document in the Expected column as "expected failure" and ensure the scenario tests the error path, not the happy path.

## Local Development

```bash
go test ./...           # all packages must pass
go build -o caseforge . # binary builds clean
./scripts/acceptance.sh # behavior-level checks
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for the full contribution workflow.
