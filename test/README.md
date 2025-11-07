# Integration Tests

This directory contains integration tests for the Fusion application. Integration tests verify the complete interaction between application layers, including service logic, repository implementations, and actual database operations.

## Directory Structure

```
test/
├── integration/
│   ├── testutil/              # Test utilities and helpers
│   │   ├── container.go      # TestContainer - all dependencies wired (recommended)
│   │   ├── db.go             # Database setup and management
│   │   └── fixtures.go       # Test data factories and fixtures
│   ├── repository/           # Repository layer integration tests
│   │   └── user_repository_integration_test.go
│   └── service/              # Service layer end-to-end tests
│       └── user_service_integration_test.go
└── README.md                 # This file
```

## Test Types

### Repository Integration Tests
Located in `test/integration/repository/`

These tests verify:
- CRUD operations with actual database
- Data persistence and retrieval
- Database constraints (unique keys, foreign keys)
- Soft delete behavior
- Pagination and ordering
- Concurrent operations

### Service Integration Tests
Located in `test/integration/service/`

These tests verify:
- Complete service workflows (Create → Read → Update → Delete)
- Business logic with real database state
- Error handling across layers
- Password hashing and validation
- Conflict detection (duplicate usernames/emails)
- Transaction behavior

## Running Tests

### Using Makefile (Recommended)

```bash
# Run only unit tests (fast, no database)
make test-unit

# Run integration tests with SQLite (fast, no setup required)
make test-integration

# Run integration tests with PostgreSQL (requires PostgreSQL running)
make test-integration-postgres

# Run all tests (unit + integration)
make test-all
# or simply:
make test

# Run tests with coverage report
make test-coverage
```

### Using Go Command Directly

```bash
# Run unit tests
go test -v -race ./internal/...

# Run integration tests with SQLite (default)
go test -v -race ./test/integration/...

# Run integration tests with PostgreSQL
TEST_DB_DRIVER=postgres \
TEST_DB_DSN="postgres://postgres:postgres@localhost:5432/fusion_test?sslmode=disable" \
go test -v -race ./test/integration/...

# Run all tests
go test -v -race ./...
```

## Database Configuration

### SQLite (Default)

By default, integration tests use SQLite in-memory database via `modernc.org/sqlite` (pure Go, no CGO required).

**Advantages:**
- ✅ Fast execution
- ✅ No external dependencies
- ✅ Works in any environment (CI/CD friendly)
- ✅ Automatic cleanup (in-memory)
- ✅ Cross-platform (no CGO)

**Limitations:**
- ⚠️ Minor SQL dialect differences from PostgreSQL
- ⚠️ Different concurrency behavior

### PostgreSQL

For production-parity testing, you can run tests against PostgreSQL.

**Setup:**

1. Start PostgreSQL (via Docker Compose):
   ```bash
   docker-compose up -d
   ```

2. Create test database:
   ```bash
   docker exec -it fusion-postgres psql -U postgres -c "CREATE DATABASE fusion_test;"
   ```

3. Run tests:
   ```bash
   make test-integration-postgres
   ```

**Advantages:**
- ✅ Production-environment parity
- ✅ Tests PostgreSQL-specific features
- ✅ Identical SQL dialect
- ✅ Real concurrency behavior

**Disadvantages:**
- ⚠️ Requires PostgreSQL running
- ⚠️ Slower than SQLite
- ⚠️ Requires Docker (typically)

### Environment Variables

Configure test database using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `TEST_DB_DRIVER` | Database driver (`sqlite3` or `postgres`) | `sqlite3` |
| `TEST_DB_DSN` | Database connection string | `file:ent?mode=memory&cache=shared&_fk=1` |

**Examples:**

```bash
# SQLite in-memory (default)
go test ./test/integration/...

# SQLite file-based
export TEST_DB_DRIVER=sqlite3
export TEST_DB_DSN="file:test.db?cache=shared&_fk=1"
go test ./test/integration/...

# PostgreSQL
export TEST_DB_DRIVER=postgres
export TEST_DB_DSN="postgres://postgres:postgres@localhost:5432/fusion_test?sslmode=disable"
go test ./test/integration/...
```

## Test Utilities

### testutil.SetupTestDB(t)

Creates a test database client with automatic migration and cleanup.

```go
client := testutil.SetupTestDB(t)
// Client automatically closed when test completes
```

### testutil.SetupSQLiteDB(t) / testutil.SetupPostgresDB(t)

Explicitly choose database type:

```go
client := testutil.SetupSQLiteDB(t)  // Always use SQLite
client := testutil.SetupPostgresDB(t) // Always use PostgreSQL
```

### testutil.CleanupDB(t, client)

Delete all data from database (useful for test isolation):

```go
testutil.CleanupDB(t, client)
```

### testutil.NewTestContainer(t) (Recommended)

**The recommended way to set up integration tests.** Creates a complete test environment with all dependencies wired:

```go
func TestUserService_Create(t *testing.T) {
    container := testutil.NewTestContainer(t)

    // Access any component you need:
    svc := container.UserService      // For service tests
    repo := container.UserRepo        // For repository tests
    client := container.Client        // For direct database access
    logger := container.Logger        // For logging

    // Automatic cleanup on test completion
}
```

**What TestContainer provides:**
- ✅ Database client (automatically configured and migrated)
- ✅ Logger (test logger with output to testing.T)
- ✅ All repositories (UserRepo, etc.)
- ✅ All services (UserService, etc.)
- ✅ Automatic cleanup (no need for defer cleanup())
- ✅ Single source of truth for test dependencies

**Benefits:**
- No manual dependency wiring
- Consistent setup across all tests
- Easy to extend when adding new components
- Automatic cleanup registration

**Example for repository tests:**
```go
func TestUserRepository_Create(t *testing.T) {
    container := testutil.NewTestContainer(t)
    repo := container.UserRepo

    user, err := repo.Create(ctx, testUser)
    require.NoError(t, err)
}
```

**Example for service tests:**
```go
func TestUserService_Create(t *testing.T) {
    container := testutil.NewTestContainer(t)
    svc := container.UserService

    user, err := svc.Create(ctx, request)
    require.NoError(t, err)
}
```

### Test Data Fixtures

Create test users easily:

```go
// Basic test user
user := testutil.NewTestUser("username", "email@example.com")

// User with custom password
user, err := testutil.NewTestUserWithPassword("user", "email@example.com", "mypassword")

// User with custom status
user := testutil.NewTestUserWithStatus("user", "email@example.com", entity.UserStatusInactive)

// Random user (for bulk testing)
user := testutil.RandomUser()

// Create user in database
created := testutil.CreateTestUser(t, repo, "username", "email@example.com")

// Create multiple users
users := testutil.CreateTestUsers(t, repo, 10) // Creates 10 users
```

## Test Coverage

Current integration test coverage:

### UserRepository
- ✅ Create (success, duplicate username, duplicate email)
- ✅ FindByID (success, not found)
- ✅ FindByUsername (success, not found, case sensitivity)
- ✅ FindByEmail (success, not found)
- ✅ Update (success, partial updates, not found, constraints)
- ✅ Delete (success, soft delete verification, not found)
- ✅ List (pagination, offset/limit, ordering, empty database)
- ✅ Concurrent operations

### UserService
- ✅ Create (success, password hashing, conflict detection)
- ✅ Update (success, status changes, conflict detection, not found)
- ✅ Delete (success, soft delete, not found)
- ✅ GetByID (success, not found)
- ✅ List (pagination, validation, edge cases)
- ✅ Complete workflow (CRUD lifecycle)

**Total:** 45+ integration test cases

## Best Practices

### 1. Test Isolation

Each test should be independent. Use TestContainer for automatic cleanup:

```go
func TestExample(t *testing.T) {
    container := testutil.NewTestContainer(t)
    // Cleanup automatically handled by TestContainer

    // Test code...
}
```

### 2. Use Fixtures

Leverage test utilities instead of manual data creation:

```go
// Good
user := testutil.CreateTestUser(t, repo, "testuser", "test@example.com")

// Avoid
user := &entity.User{Username: "testuser", Email: "test@example.com", ...}
created, err := repo.Create(ctx, user)
if err != nil { t.Fatal(err) }
```

### 3. Verify Side Effects

Integration tests should verify database state:

```go
// Delete user
err := repo.Delete(ctx, userID)
require.NoError(t, err)

// Verify user is actually deleted
_, err = repo.FindByID(ctx, userID)
assert.Error(t, err)
assert.True(t, errors.IsNotFoundError(err))
```

### 4. Test Error Cases

Don't just test happy paths:

```go
// Test constraint violations
user := testutil.CreateTestUser(t, repo, "duplicate", "email@example.com")

// Try to create duplicate
_, err := repo.Create(ctx, testutil.NewTestUser("duplicate", "other@example.com"))
assert.Error(t, err)
assert.Contains(t, err.Error(), "constraint")
```

### 5. Table-Driven Tests

Use table-driven tests for comprehensive coverage:

```go
tests := []struct {
    name    string
    input   *entity.User
    wantErr bool
}{
    {name: "success", input: validUser, wantErr: false},
    {name: "duplicate", input: duplicateUser, wantErr: true},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test logic...
    })
}
```

## CI/CD Integration

Integration tests are designed to run in CI/CD pipelines:

```yaml
# Example GitHub Actions
- name: Run Unit Tests
  run: make test-unit

- name: Run Integration Tests
  run: make test-integration  # Uses SQLite, no setup needed

# Optional: Test against PostgreSQL
- name: Start PostgreSQL
  run: docker-compose up -d

- name: Run PostgreSQL Integration Tests
  run: make test-integration-postgres
```

## Troubleshooting

### Tests are slow

- Use `make test-unit` for quick feedback (no database)
- Use `make test-integration` (SQLite is faster than PostgreSQL)
- Run specific test: `go test -run TestName ./test/integration/...`

### "database is locked" (SQLite)

- Ensure you're using `cache=shared` in DSN
- Check that connections are properly closed
- Verify `defer cleanup()` is called

### "connection refused" (PostgreSQL)

- Ensure PostgreSQL is running: `docker-compose up -d`
- Check DSN is correct
- Verify database exists: `docker exec -it fusion-postgres psql -U postgres -l`

### Tests fail randomly

- Check for test isolation issues
- Ensure `defer cleanup()` is used in all tests
- Look for shared state or race conditions

## Contributing

When adding new features:

1. **Write unit tests first** (in `internal/`)
2. **Add integration tests** (in `test/integration/`)
3. **Run all tests**: `make test-all`
4. **Verify coverage**: `make test-coverage`

### Adding New Integration Tests

1. Repository tests go in `test/integration/repository/`
2. Service tests go in `test/integration/service/`
3. Use existing utilities in `test/integration/testutil/`
4. Follow table-driven test patterns
5. Test both success and error cases
6. Verify database state changes

## Resources

- [EntGO Testing](https://entgo.io/docs/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
