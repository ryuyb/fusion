package testutil

import (
	"testing"

	serviceImpl "github.com/ryuyb/fusion/internal/application/service"
	"github.com/ryuyb/fusion/internal/domain/repository"
	"github.com/ryuyb/fusion/internal/domain/service"
	"github.com/ryuyb/fusion/internal/infrastructure/database"
	repoImpl "github.com/ryuyb/fusion/internal/infrastructure/repository"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// TestContainer holds all test dependencies for integration tests.
// It provides a centralized place to manage the dependency graph,
// making tests easier to write and maintain.
type TestContainer struct {
	// Infrastructure
	Client *database.Client
	Logger *zap.Logger

	// Repositories
	UserRepo repository.UserRepository

	// Services
	UserService service.UserService
}

// NewTestContainer creates a new test container with all dependencies wired.
// It automatically sets up the test database, logger, repositories, and services.
// Cleanup is automatically registered with t.Cleanup().
//
// Example usage:
//
//	func TestSomething(t *testing.T) {
//	    container := testutil.NewTestContainer(t)
//	    user, err := container.UserService.Create(ctx, req)
//	    // ...
//	}
func NewTestContainer(t *testing.T) *TestContainer {
	t.Helper()

	// Setup infrastructure
	client := SetupTestDB(t)
	logger := zaptest.NewLogger(t)

	// Setup repositories
	userRepo := repoImpl.NewUserRepository(client, logger)

	// Setup services
	userService := serviceImpl.NewUserService(userRepo, logger)

	// Register cleanup
	t.Cleanup(func() {
		CleanupDB(t, client)
	})

	return &TestContainer{
		Client:      client,
		Logger:      logger,
		UserRepo:    userRepo,
		UserService: userService,
	}
}
