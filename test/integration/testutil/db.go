package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/ryuyb/fusion/internal/infrastructure/database"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent/enttest"
	_ "modernc.org/sqlite"
)

const (
	// Environment variables for test database configuration
	EnvTestDBDriver = "TEST_DB_DRIVER"
	EnvTestDBDSN    = "TEST_DB_DSN"

	// Default values
	DefaultDriver = "sqlite"
	// modernc.org/sqlite requires _pragma parameter for foreign keys
	DefaultDSN = "file:ent?mode=memory&cache=shared&_pragma=foreign_keys(1)"
)

// SetupTestDB creates a test database client with automatic migration.
// By default, it uses SQLite in-memory database.
// Set TEST_DB_DRIVER and TEST_DB_DSN environment variables to use PostgreSQL.
//
// Example:
//
//	export TEST_DB_DRIVER=postgres
//	export TEST_DB_DSN="postgres://postgres:postgres@localhost:5432/fusion_test?sslmode=disable"
func SetupTestDB(t *testing.T) *database.Client {
	t.Helper()

	driver := os.Getenv(EnvTestDBDriver)
	dsn := os.Getenv(EnvTestDBDSN)

	if driver == "" {
		driver = DefaultDriver
	}
	if dsn == "" {
		dsn = DefaultDSN
	}

	return NewTestClient(t, driver, dsn)
}

// SetupSQLiteDB creates a SQLite in-memory test database.
// This is useful for fast, isolated tests.
func SetupSQLiteDB(t *testing.T) *database.Client {
	t.Helper()
	return NewTestClient(t, "sqlite", "file:ent?mode=memory&cache=shared&_pragma=foreign_keys(1)")
}

// SetupPostgresDB creates a PostgreSQL test database.
// Requires PostgreSQL to be running (e.g., via docker-compose).
// Default DSN can be overridden with TEST_DB_DSN environment variable.
func SetupPostgresDB(t *testing.T) *database.Client {
	t.Helper()

	dsn := os.Getenv(EnvTestDBDSN)
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/fusion_test?sslmode=disable"
	}

	return NewTestClient(t, "postgres", dsn)
}

// NewTestClient creates a new test database client with the given driver and DSN.
// The client is automatically closed when the test completes.
// Database schema is automatically migrated.
func NewTestClient(t *testing.T, driver, dsn string) *database.Client {
	t.Helper()

	// Map driver name to SQL driver and ent dialect
	var sqlDriver, dialectName string
	switch driver {
	case "sqlite", "sqlite3":
		sqlDriver = "sqlite" // modernc.org/sqlite registers as "sqlite"
		dialectName = dialect.SQLite
	case "postgres", "postgresql":
		sqlDriver = "postgres"
		dialectName = dialect.Postgres
	default:
		t.Fatalf("unsupported database driver: %s (use 'sqlite' or 'postgres')", driver)
	}

	// Open SQL connection with the correct driver
	db, err := sql.Open(sqlDriver, dsn)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// Create ent driver from sql.DB
	drv := entsql.OpenDB(dialectName, db)

	// Create test client with enttest
	entClient := enttest.NewClient(t, enttest.WithOptions(ent.Driver(drv)))

	// Wrap in our database.Client
	client := &database.Client{Client: *entClient}

	// Register cleanup
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Logf("failed to close test database: %v", err)
		}
	})

	return client
}

// CleanupDB deletes all data from the database.
// Useful for resetting database state between tests.
func CleanupDB(t *testing.T, client *database.Client) {
	t.Helper()

	ctx := context.Background()

	// Delete all users (this will cascade to related entities if any)
	if _, err := client.User.Delete().Exec(ctx); err != nil {
		t.Fatalf("failed to cleanup users: %v", err)
	}
}

// WithTransaction executes a test function within a database transaction.
// The transaction is automatically rolled back after the test, ensuring test isolation.
func WithTransaction(t *testing.T, client *database.Client, fn func(tx *ent.Tx)) {
	t.Helper()

	ctx := context.Background()
	tx, err := client.Tx(ctx)
	if err != nil {
		t.Fatalf("failed to start transaction: %v", err)
	}

	// Always rollback to ensure test isolation
	defer func() {
		if err := tx.Rollback(); err != nil {
			t.Logf("failed to rollback transaction: %v", err)
		}
	}()

	fn(tx)
}

// GetTestDSN returns the test database DSN for logging/debugging purposes.
func GetTestDSN() string {
	driver := os.Getenv(EnvTestDBDriver)
	dsn := os.Getenv(EnvTestDBDSN)

	if driver == "" {
		driver = DefaultDriver
	}
	if dsn == "" {
		dsn = DefaultDSN
	}

	return fmt.Sprintf("%s: %s", driver, dsn)
}
