# CaseForge — Claude Code Instructions

## Acceptance Tests

Every development session must run the acceptance suite before committing and after completing any feature or fix:

```bash
./scripts/acceptance.sh
```

**Rules:**
- All 44 scenarios must pass before pushing to `main`.
- When adding a new command or feature, add the corresponding scenario(s) to `docs/acceptance/acceptance-tests.md` and a matching `contains` / `run` call in `scripts/acceptance.sh` **in the same commit**.
- If a scenario fails, determine root cause before proceeding:
  - **Code bug** → fix the code, re-run the suite.
  - **Environment issue** (no API key, no running server) → document in the Expected column as "expected failure" and ensure the scenario tests the error path, not the happy path.
