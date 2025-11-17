# Repository Guidelines

## Project Structure & Module Organization
- Go service code lives under `cmd/` (entrypoints such as `cmd/app/main.go`) and `internal/` (domain layers: `app`, `application`, `core`, `infrastructure`, `pkg`).
- Configuration and assets: `configs/` for app config samples, `docs/` + `docs.go` for Swagger generation, `scripts/` for helper scripts, `Dockerfile` and `docker-compose.yml` for container workflows.
- Database schemas generated with Ent reside in `internal/infrastructure/database/ent`; new schemas go in `internal/infrastructure/database/schema`.

## Build, Test, and Development Commands
- `make build` — compile the service to `bin/fusion` with version metadata.
- `go test ./...` — run all tests; add `-race` for data race checks.
- `make lint` — run `golangci-lint` against the codebase.
- `make generate-ent` — regenerate Ent code; run after editing schema files.
- `make generate-swagger` — refresh Swagger docs into `./docs/api/` via `docs.go`.
- `make docker-up` / `make docker-down` — start/stop services with Docker Compose.

## Coding Style & Naming Conventions
- Go 1.25 with modules; follow standard Go formatting (`go fmt`) and lint rules from `golangci-lint`.
- Package names are lowercase, short, and avoid underscores; exported identifiers use Go-style CamelCase.
- Keep handlers thin; prefer business logic in `internal/application` and domain types in `internal/core`.
- Configuration structs live in `configs/`; keep environment variables uppercase with `_` separators.

## Testing Guidelines
- Use Go `testing` with `testify` helpers; table-driven tests preferred.
- Name files with `_test.go`; keep test packages close to the code under test.
- Aim for meaningful coverage on domain logic; exercise HTTP handlers via Fiber test utilities where possible.
- Use `go test ./... -run <Name>` to target suites; keep fixtures small and deterministic.

## Commit & Pull Request Guidelines
- Commit messages: imperative voice, concise summary (e.g., `fix login token refresh`); group related changes, use Conventional Commits.
- Pull requests should describe scope, include relevant `make`/`go test` outputs, and link issues if applicable.
- Add screenshots or curl examples for API changes when helpful; note schema or migration impacts explicitly.
- Before merging, ensure lint/tests pass and regenerated assets (Ent, Swagger, mocks) are committed.
