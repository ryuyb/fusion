# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Fusion is a Go-based streaming platform built with Fiber v3, using Clean Architecture with Dependency Injection (uber-go/fx). The project uses EntGO for database ORM and PostgreSQL as the database.

## Build & Run Commands

### Development
```bash
# Install development tools (Air, Swag, golangci-lint)
make install-tools

# Run development server with hot reload
air

# Run server directly
go run cmd/server/main.go serve

# Run with custom config
go run cmd/server/main.go serve --config=./configs/config.yaml

# Run with custom port
go run cmd/server/main.go serve --port=3000
```

### Build
```bash
# Build production binary
make build
# Output: bin/fusion

# Run built binary
./bin/fusion serve
```

### Database Migrations
```bash
# Run migrations
go run cmd/server/main.go migrate up

# Generate new Ent schema
make ent-new name=ModelName

# Generate Ent code after schema changes
make generate-ent
# Or directly: go generate ./internal/infrastructure/database/ent
```

### Code Quality
```bash
# Run linter
make lint

# Generate Swagger documentation
make generate-swagger

# Format Swagger annotations
make format-swagger
```

## Docker

### Build Docker Image
```bash
# Build with default settings
docker build -t fusion:latest .

# Build with version info
docker build \
  --build-arg VERSION=$(git describe --tags --always --dirty) \
  --build-arg BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S') \
  --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) \
  --build-arg GO_VERSION=$(go version | awk '{print $3}') \
  -t fusion:latest .
```

### Run with Docker Compose
```bash
# Start all services (app + PostgreSQL)
docker-compose up

# Start in detached mode
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

### Run Docker Container Manually
```bash
# Run with environment variables
docker run -d \
  --name fusion-app \
  -p 8080:8080 \
  -e FUSION_DATABASE_DSN="postgres://user:pass@host:5432/fusion" \
  -e FUSION_JWT_SECRET="your-secret-key" \
  fusion:latest

# Run with env file
docker run -d \
  --name fusion-app \
  -p 8080:8080 \
  --env-file .env \
  fusion:latest
```

## Architecture

### Layer Structure (Clean Architecture)

```
internal/
├── domain/           # Business logic core (entities, domain services, repository interfaces)
├── application/      # Application services, DTOs, use cases
├── infrastructure/   # External implementations (database, config, logger)
├── interface/        # Delivery mechanisms (HTTP handlers, routes, middleware)
├── pkg/             # Shared utilities (auth, validator, errors)
├── app/             # Application bootstrap and dependency wiring
└── cli/             # CLI commands (serve, migrate)
```

### Dependency Flow

- **Domain** (innermost) - Pure business logic, no external dependencies
- **Application** - Orchestrates domain logic, depends on domain layer
- **Infrastructure** - Implements domain interfaces (repositories), depends on domain
- **Interface** - HTTP layer, depends on application services
- **App** - Wires everything together using uber-go/fx modules

### Key Architectural Patterns

**Dependency Injection with uber-go/fx:**
- Each layer has a `module.go` file that provides dependencies via `fx.Module()`
- Modules are composed in `internal/app/module.go`
- Lifecycle hooks handle startup/shutdown (database connections, server)

**Repository Pattern:**
- Domain defines interfaces in `internal/domain/repository/`
- Infrastructure implements them in `internal/infrastructure/repository/`
- Application services depend on domain repository interfaces

**EntGO Schema Location:**
- Schemas are in `internal/infrastructure/database/schema/`
- Generated code is in `internal/infrastructure/database/ent/`
- Use soft delete mixin from `internal/pkg/entgo/mixin/soft_delete.go`

**Error Handling:**
- Custom errors in `internal/pkg/errors/`
- Centralized error handler middleware at `internal/interface/http/middleware/error_handler.go`
- Use structured errors: `errors.BadRequest()`, `errors.Unauthorized()`, etc.

**Authentication:**
- JWT implementation in `internal/pkg/auth/jwt.go`
- Auth middleware in `internal/interface/http/middleware/auth.go`
- Two modes: `Handler()` (required) and `Optional()` (optional auth)
- User context accessible via `auth.UserContextKey` and `auth.UserIdContextKey`

## Configuration

**Config Loading:**
- Config file: `configs/config.yaml`
- Environment-based: `configs/config.{env}.yaml` (use `--env=dev`)
- Environment variables: prefix with `FUSION_`, use underscores (e.g., `FUSION_SERVER_PORT=8080`)
- Viper automatically merges config file and env vars

**Config Structure:**
- Server: host, port, timeouts
- Database: DSN, connection pool settings
- Logger: level, console/file output, rotation settings
- JWT: secret, expiration

## Development Workflow

**Adding a New Feature:**

1. **Domain Layer:** Define entity in `internal/domain/entity/`, repository interface in `internal/domain/repository/`
2. **Infrastructure:** Create Ent schema in `internal/infrastructure/database/schema/`, run `make generate-ent`, implement repository in `internal/infrastructure/repository/`
3. **Application:** Create DTOs in `internal/application/dto/`, implement service in `internal/application/service/`
4. **Interface:** Create handler in `internal/interface/http/handler/`, register routes in `internal/interface/http/route/`
5. **Wiring:** Add providers to appropriate `module.go` files if needed

**Adding Swagger Documentation:**
- Use Swag annotations in handler files
- Run `make generate-swagger` to regenerate docs
- Swagger UI available at `/swagger/*` route (configured in `swagger_route.go`)

**Middleware Order (important):**
1. RequestID
2. CORS
3. Recovery
4. Compress
5. Logger
6. Auth (per-route, applied in route definitions)

## Testing

Currently no test files exist. When adding tests:
- Unit tests for domain/application layers
- Integration tests with `enttest` package for repositories
- HTTP tests for handlers using Fiber's test utilities

## Important Notes

- Database migrations run automatically on server startup via fx lifecycle hook
- All modules use fx for dependency injection - avoid manual instantiation
- Ent generation requires `--feature intercept` flag (see `generate.go`)
- Routes are modular - each domain has its own route file in `internal/interface/http/route/`
- Password hashing utilities available in `internal/pkg/utils/password.go`
- Validation uses go-playground/validator with custom validator wrapper in `internal/pkg/validator/`