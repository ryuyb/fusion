# Code Style and Conventions

## Naming Conventions

### Package and Module Names
- Use singular nouns for package names (e.g., `repository`, `service`, `handler`)
- Use lowerCamelCase or snake_case for package names
- Group related functionality into cohesive packages

### File Naming
- Use snake_case for file names: `user_service.go`
- End handler files with `_handler.go`: `user_handler.go`
- End service files with `_service.go`: `user_service.go`
- End repository files with `_repository.go`: `user_repository.go`
- End DTO files with `_dto.go`: `create_user_dto.go`

### Variable and Function Names
- Use CamelCase for exported functions and variables
- Use camelCase for internal functions and variables
- Use meaningful, descriptive names
- Avoid abbreviations unless widely known (e.g., `ID`, `API`, `HTTP`)

## Code Structure

### Layer Architecture
Follow Clean Architecture layers strictly:
1. **Domain Layer** (`internal/domain/`)
   - Pure business logic
   - Entities (structs)
   - Repository interfaces
   - Domain services
   - No external dependencies

2. **Application Layer** (`internal/application/`)
   - Use cases
   - DTOs (Data Transfer Objects)
   - Application services
   - Depends on domain layer only

3. **Infrastructure Layer** (`internal/infrastructure/`)
   - Implements repository interfaces from domain
   - Database implementations (EntGO)
   - External service clients
   - Config and logger implementations

4. **Interface Layer** (`internal/interface/`)
   - HTTP handlers
   - Routes
   - Middleware
   - Depends on application layer

### Dependency Injection Pattern
- Use uber-go/fx for dependency injection
- Each layer has a `module.go` file
- Modules are composed in `internal/app/module.go`
- Use lifecycle hooks for startup/shutdown
- Example:
```go
// In module.go
func Module() fx.Option {
    return fx.Options(
        fx.Provide(NewService),
        fx.Invoke(RegisterHooks),
    )
}
```

## Error Handling

### Custom Error Types
- Use custom error types from `internal/pkg/errors/`
- Structured errors: `errors.BadRequest()`, `errors.Unauthorized()`, etc.
- Centralized error handler middleware at `internal/interface/http/middleware/error_handler.go`
- Avoid returning raw errors; wrap with context

### Error Example
```go
if err != nil {
    return nil, errors.BadRequest("invalid user ID", err)
}
```

## Authentication

### JWT Implementation
- JWT implementation: `internal/pkg/auth/jwt.go`
- Auth middleware: `internal/interface/http/middleware/auth.go`
- Two modes:
  - `Handler()` - Required authentication
  - `Optional()` - Optional authentication

### Accessing User Context
```go
// From handlers
userID := auth.GetUserID(c.Context())
user := auth.GetUser(c.Context())
```

## Validation

### Validator Pattern
- Uses go-playground/validator
- Custom validator wrapper: `internal/pkg/validator/`
- Tag-based validation on DTOs
- Example:
```go
type CreateUserDTO struct {
    Email    string `validate:"required,email"`
    Password string `validate:"required,min=8"`
    Name     string `validate:"required,min=2"`
}
```

## Database (EntGO)

### Schema Location
- Schemas: `internal/infrastructure/database/schema/`
- Generated code: `internal/infrastructure/database/ent/`
- Use soft delete mixin: `internal/pkg/entgo/mixin/soft_delete.go`

### Schema Example
```go
// internal/infrastructure/database/schema/user.go
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("email"),
        field.String("password"),
    }
}
```

### Ent Generation
```bash
make ent-new name=ModelName
make generate-ent
```

## HTTP Routing

### Route Structure
- Routes are modular in `internal/interface/http/route/`
- Each domain has its own route file
- Use Fiber v3 syntax
- Group routes with versioning: `/api/v1/`

### Route Example
```go
// User routes
api.Post("/users", handler.CreateUser)
api.Get("/users/:id", handler.GetUser)
api.Put("/users/:id", handler.UpdateUser)
api.Delete("/users/:id", handler.DeleteUser)
```

## Middleware Order
Important: Apply middleware in this order:
1. RequestID
2. CORS
3. Recovery
4. Compress
5. Logger
6. Auth (per-route, applied in route definitions)

## Documentation

### Swagger Annotations
- Use Swag annotations in handler files
- Format: `// @Summary Description`
- Place directly above handler functions
- Generate: `make generate-swagger`
- Available at `/swagger/*`

### Swagger Example
```go
// @Summary Create a new user
// @Description Creates a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param CreateUserDTO body CreateUserDTO true "User data"
// @Success 201 {object} User
// @Failure 400 {object} Error
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    // Handler logic
}
```

## Testing
- No tests exist yet, but use these patterns when adding:
  - Unit tests for domain/application layers
  - Integration tests with `enttest` package
  - HTTP tests using Fiber's test utilities
  - Use table-driven tests where appropriate
  - Use `require` or `assert` from testify for assertions

## Utilities
- Password hashing: `internal/pkg/utils/password.go`
- Custom error types: `internal/pkg/errors/`
- Auth utilities: `internal/pkg/auth/`
- Validator wrapper: `internal/pkg/validator/`

## General Best Practices
- Keep functions small and focused (single responsibility)
- Use interfaces for abstraction
- Dependency flow: inner layers don't know about outer layers
- Return early to avoid nested conditionals
- Use meaningful error messages
- Log with context (use structured logging with Zap)
- Write comments for non-obvious logic
- Keep dependencies minimal
