# Contributing to CaseForge

Thank you for taking the time to contribute!

## Prerequisites

- Go 1.26+
- [hurl](https://hurl.dev/docs/installation.html) (required for integration tests)
- [k6](https://k6.io/docs/get-started/installation/) (optional, for k6 runner)

## Getting Started

```bash
git clone https://github.com/testmind-hq/caseforge.git
cd caseforge
go mod download
go test ./...
```

## Branching

- Branch off `main` for all new work.
- Use descriptive branch names: `feature/<topic>`, `fix/<topic>`, `chore/<topic>`.

## Making Changes

1. Write a failing test first (TDD).
2. Implement the minimal code to make the test pass.
3. Run `go test ./...` — all packages must pass.
4. Run `go vet ./...` — no warnings.
5. Commit with a [Conventional Commits](https://www.conventionalcommits.org/) message:
   - `feat:` new feature
   - `fix:` bug fix
   - `chore:` tooling / infra
   - `docs:` documentation only

## Pull Requests

- Keep PRs focused — one concern per PR.
- Fill in the PR template summary and test plan.
- All CI checks must pass before merge.
- At least one maintainer review is required.

## Code Style

- Follow standard Go formatting (`gofmt`).
- Exported identifiers must have doc comments.
- Avoid global state; prefer dependency injection.

## License

By contributing you agree that your contributions will be licensed under the [Apache License 2.0](LICENSE).
