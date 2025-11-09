# Fusion Project Overview

## Project Purpose
Fusion is a Go-based streaming platform built with Clean Architecture and Dependency Injection. It uses modern technologies including Fiber v3 web framework, EntGO ORM, and PostgreSQL database.

## Tech Stack
- **Language**: Go 1.25.3
- **Web Framework**: Fiber v3 (github.com/gofiber/fiber/v3 v3.0.0-rc.2)
- **ORM**: EntGO (entgo.io/ent v0.14.5)
- **Database**: PostgreSQL
- **Dependency Injection**: Uber Fx (go.uber.org/fx v1.24.0)
- **Configuration**: Viper (github.com/spf13/viper v1.21.0)
- **Authentication**: JWT (golang-jwt/jwt/v5 v5.3.0)
- **Validation**: go-playground/validator v10
- **Logging**: Uber Zap (go.uber.org/zap v1.27.0)
- **Documentation**: Swagger (github.com/swaggo/swag v1.16.6)
- **HTTP Client**: Resty v3 (resty.dev/v3)

## Architecture Pattern
**Clean Architecture** with layered structure:
- `internal/domain/` - Business logic core (entities, services, repository interfaces)
- `internal/application/` - Application services, DTOs, use cases
- `internal/infrastructure/` - External implementations (database, config, logger)
- `internal/interface/` - HTTP handlers, routes, middleware
- `internal/pkg/` - Shared utilities (auth, validator, errors)
- `internal/app/` - Application bootstrap and dependency wiring
- `cmd/` - CLI commands (serve, migrate)

## Key Features
- Clean Architecture with dependency injection using uber-go/fx
- Repository pattern with domain interfaces
- EntGO for database operations with soft delete mixin
- JWT-based authentication with middleware
- Structured error handling with custom error types
- Swagger API documentation
- Docker and Docker Compose support
- Multi-architecture Docker builds (AMD64 + ARM64)
- Comprehensive testing (unit, integration, e2e)

## Environment
- **Platform**: Darwin (macOS)
- **Server Port**: 8080 (default)
- **Database**: PostgreSQL (via Docker Compose)
