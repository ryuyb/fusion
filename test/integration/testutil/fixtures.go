package testutil

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/domain/repository"
	"github.com/ryuyb/fusion/internal/pkg/utils"
)

var (
	// Default test password
	DefaultTestPassword = "test1234"
	// Default hashed test password
	DefaultTestPasswordHash string
)

func init() {
	// Pre-hash the default test password to avoid repeated hashing
	hash, err := utils.HashPassword(DefaultTestPassword)
	if err != nil {
		panic(fmt.Sprintf("failed to hash default test password: %v", err))
	}
	DefaultTestPasswordHash = hash
}

// NewTestUser creates a new test user entity with default values.
func NewTestUser(username, email string) *entity.User {
	return &entity.User{
		Username: username,
		Email:    email,
		Password: DefaultTestPasswordHash,
		Status:   entity.UserStatusActive,
	}
}

// NewTestUserWithPassword creates a new test user entity with a custom password.
func NewTestUserWithPassword(username, email, password string) (*entity.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return &entity.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Status:   entity.UserStatusActive,
	}, nil
}

// NewTestUserWithStatus creates a new test user entity with a specific status.
func NewTestUserWithStatus(username, email string, status entity.UserStatus) *entity.User {
	return &entity.User{
		Username: username,
		Email:    email,
		Password: DefaultTestPasswordHash,
		Status:   status,
	}
}

// RandomUser generates a user with random username and email.
// Useful for tests that don't care about specific values.
func RandomUser() *entity.User {
	rand.NewSource(time.Now().UnixNano())
	id := rand.Intn(1000000)

	return &entity.User{
		Username: fmt.Sprintf("user_%d", id),
		Email:    fmt.Sprintf("user_%d@example.com", id),
		Password: DefaultTestPasswordHash,
		Status:   entity.UserStatusActive,
	}
}

// RandomUsername generates a random username.
func RandomUsername() string {
	rand.NewSource(time.Now().UnixNano())
	return fmt.Sprintf("user_%d", rand.Intn(1000000))
}

// RandomEmail generates a random email address.
func RandomEmail() string {
	rand.NewSource(time.Now().UnixNano())
	return fmt.Sprintf("user_%d@example.com", rand.Intn(1000000))
}

// CreateTestUser creates a user in the database using the repository.
// Returns the created user or fails the test.
func CreateTestUser(t *testing.T, repo repository.UserRepository, username, email string) *entity.User {
	t.Helper()

	ctx := context.Background()
	user := NewTestUser(username, email)

	created, err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	return created
}

// CreateTestUserWithStatus creates a user with a specific status in the database.
func CreateTestUserWithStatus(t *testing.T, repo repository.UserRepository, username, email string, status entity.UserStatus) *entity.User {
	t.Helper()

	ctx := context.Background()
	user := NewTestUserWithStatus(username, email, status)

	created, err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	return created
}

// CreateTestUsers creates multiple test users in the database.
// Returns a slice of created users.
func CreateTestUsers(t *testing.T, repo repository.UserRepository, count int) []*entity.User {
	t.Helper()

	users := make([]*entity.User, count)
	ctx := context.Background()

	for i := 0; i < count; i++ {
		user := &entity.User{
			Username: fmt.Sprintf("testuser%d", i+1),
			Email:    fmt.Sprintf("testuser%d@example.com", i+1),
			Password: DefaultTestPasswordHash,
			Status:   entity.UserStatusActive,
		}

		created, err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("failed to create test user %d: %v", i+1, err)
		}

		users[i] = created
	}

	return users
}

// CreateRandomUsers creates multiple users with random data.
// Useful for pagination and listing tests.
func CreateRandomUsers(t *testing.T, repo repository.UserRepository, count int) []*entity.User {
	t.Helper()

	users := make([]*entity.User, count)
	ctx := context.Background()

	for i := 0; i < count; i++ {
		user := RandomUser()

		created, err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("failed to create random user %d: %v", i+1, err)
		}

		users[i] = created
	}

	return users
}

// MustHashPassword hashes a password or fails the test.
func MustHashPassword(t *testing.T, password string) string {
	t.Helper()

	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	return hash
}

// AssertUserEqual asserts that two users are equal (comparing relevant fields).
func AssertUserEqual(t *testing.T, expected, actual *entity.User) {
	t.Helper()

	if expected.Username != actual.Username {
		t.Errorf("username mismatch: expected %s, got %s", expected.Username, actual.Username)
	}
	if expected.Email != actual.Email {
		t.Errorf("email mismatch: expected %s, got %s", expected.Email, actual.Email)
	}
	if expected.Status != actual.Status {
		t.Errorf("status mismatch: expected %s, got %s", expected.Status, actual.Status)
	}
}
