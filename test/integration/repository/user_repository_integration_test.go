package repository_test

import (
	"context"
	"testing"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/ryuyb/fusion/test/integration/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestUserRepository_Create_Integration(t *testing.T) {
	tests := []struct {
		name    string
		user    *entity.User
		wantErr bool
		errMsg  string
	}{
		{
			name: "success - creates user with all fields",
			user: testutil.NewTestUser("testuser", "test@example.com"),
		},
		{
			name: "success - creates inactive user",
			user: testutil.NewTestUserWithStatus("inactive_user", "inactive@example.com", entity.UserStatusInactive),
		},
		{
			name: "success - creates banned user",
			user: testutil.NewTestUserWithStatus("banned_user", "banned@example.com", entity.UserStatusBanned),
		},
		//{
		//	name:    "error - duplicate username",
		//	user:    testutil.NewTestUser("duplicate", "unique@example.com"),
		//	wantErr: true,
		//	errMsg:  "constraint",
		//},
		//{
		//	name:    "error - duplicate email",
		//	user:    testutil.NewTestUser("unique", "duplicate@example.com"),
		//	wantErr: true,
		//	errMsg:  "constraint",
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := testutil.NewTestContainer(t)
			repo := container.UserRepo

			ctx := context.Background()

			// Pre-create users for duplicate tests
			if tt.name == "error - duplicate username" {
				testutil.CreateTestUser(t, repo, "duplicate", "other@example.com")
			}
			if tt.name == "error - duplicate email" {
				testutil.CreateTestUser(t, repo, "other", "duplicate@example.com")
			}

			// Execute
			result, err := repo.Create(ctx, tt.user)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)

				// Verify all fields
				assert.NotZero(t, result.ID)
				assert.Equal(t, tt.user.Username, result.Username)
				assert.Equal(t, tt.user.Email, result.Email)
				assert.Equal(t, tt.user.Password, result.Password)
				assert.Equal(t, tt.user.Status, result.Status)
				assert.NotZero(t, result.CreatedAt)
				assert.NotZero(t, result.UpdatedAt)
				assert.True(t, result.DeleteAt.IsZero(), "DeleteAt should be zero for new users")
			}
		})
	}
}

func TestUserRepository_FindByID_Integration(t *testing.T) {
	container := testutil.NewTestContainer(t)
	repo := container.UserRepo

	ctx := context.Background()

	// Create test user
	created := testutil.CreateTestUser(t, repo, "findbyid", "findbyid@example.com")

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name: "success - finds existing user",
			id:   created.ID,
		},
		{
			name:    "error - user not found",
			id:      999999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.FindByID(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, errors.IsNotFoundError(err))
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, created.ID, result.ID)
				assert.Equal(t, created.Username, result.Username)
				assert.Equal(t, created.Email, result.Email)
			}
		})
	}
}

func TestUserRepository_FindByUsername_Integration(t *testing.T) {
	container := testutil.NewTestContainer(t)
	repo := container.UserRepo

	ctx := context.Background()

	// Create test users
	testutil.CreateTestUser(t, repo, "alice", "alice@example.com")
	testutil.CreateTestUser(t, repo, "bob", "bob@example.com")

	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			name:     "success - finds alice",
			username: "alice",
		},
		{
			name:     "success - finds bob",
			username: "bob",
		},
		{
			name:     "error - user not found",
			username: "nonexistent",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.FindByUsername(ctx, tt.username)

			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, errors.IsNotFoundError(err))
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.username, result.Username)
			}
		})
	}
}

func TestUserRepository_FindByEmail_Integration(t *testing.T) {
	container := testutil.NewTestContainer(t)
	repo := container.UserRepo

	ctx := context.Background()

	// Create test users
	testutil.CreateTestUser(t, repo, "user1", "user1@example.com")
	testutil.CreateTestUser(t, repo, "user2", "user2@test.com")

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:  "success - finds by example.com email",
			email: "user1@example.com",
		},
		{
			name:  "success - finds by test.com email",
			email: "user2@test.com",
		},
		{
			name:    "error - email not found",
			email:   "nonexistent@example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.FindByEmail(ctx, tt.email)

			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, errors.IsNotFoundError(err))
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.email, result.Email)
			}
		})
	}
}

func TestUserRepository_Update_Integration(t *testing.T) {
	container := testutil.NewTestContainer(t)
	repo := container.UserRepo

	ctx := context.Background()

	// Create test user
	original := testutil.CreateTestUser(t, repo, "original", "original@example.com")

	tests := []struct {
		name    string
		setup   func() *entity.User
		wantErr bool
		errMsg  string
	}{
		{
			name: "success - updates all fields",
			setup: func() *entity.User {
				updated := *original
				updated.Username = "updated"
				updated.Email = "updated@example.com"
				updated.Status = entity.UserStatusInactive
				updated.Password = testutil.MustHashPassword(t, "newpassword")
				return &updated
			},
		},
		{
			name: "success - updates only username",
			setup: func() *entity.User {
				updated := *original
				updated.Username = "newusername"
				return &updated
			},
		},
		{
			name: "success - updates only email",
			setup: func() *entity.User {
				updated := *original
				updated.Email = "newemail@example.com"
				return &updated
			},
		},
		{
			name: "success - updates only status",
			setup: func() *entity.User {
				updated := *original
				updated.Status = entity.UserStatusBanned
				return &updated
			},
		},
		{
			name: "error - user not found",
			setup: func() *entity.User {
				return &entity.User{
					ID:       999999,
					Username: "notfound",
					Email:    "notfound@example.com",
					Password: testutil.DefaultTestPasswordHash,
					Status:   entity.UserStatusActive,
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := tt.setup()

			result, err := repo.Update(ctx, user)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, user.ID, result.ID)
				assert.Equal(t, user.Username, result.Username)
				assert.Equal(t, user.Email, result.Email)
				assert.Equal(t, user.Status, result.Status)
				assert.NotZero(t, result.UpdatedAt)
			}
		})
	}
}

func TestUserRepository_Delete_Integration(t *testing.T) {
	container := testutil.NewTestContainer(t)
	repo := container.UserRepo

	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func() int64
		wantErr bool
	}{
		{
			name: "success - soft deletes existing user",
			setup: func() int64 {
				user := testutil.CreateTestUser(t, repo, "todelete", "delete@example.com")
				return user.ID
			},
		},
		{
			name: "error - user not found",
			setup: func() int64 {
				return 999999
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := tt.setup()

			err := repo.Delete(ctx, userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, errors.IsNotFoundError(err))
			} else {
				require.NoError(t, err)

				// Verify user is soft deleted (cannot be found)
				_, err := repo.FindByID(ctx, userID)
				assert.Error(t, err)
				assert.True(t, errors.IsNotFoundError(err), "Soft deleted user should not be found")
			}
		})
	}
}

func TestUserRepository_List_Integration(t *testing.T) {
	container := testutil.NewTestContainer(t)
	repo := container.UserRepo

	ctx := context.Background()

	// Create 15 test users
	testutil.CreateTestUsers(t, repo, 15)

	tests := []struct {
		name           string
		offset         int
		limit          int
		expectedCount  int
		expectedTotal  int
		validateResult func(*testing.T, []*entity.User)
	}{
		{
			name:          "success - first page (10 items)",
			offset:        0,
			limit:         10,
			expectedCount: 10,
			expectedTotal: 15,
		},
		{
			name:          "success - second page (5 items)",
			offset:        10,
			limit:         10,
			expectedCount: 5,
			expectedTotal: 15,
		},
		{
			name:          "success - limit larger than total",
			offset:        0,
			limit:         100,
			expectedCount: 15,
			expectedTotal: 15,
		},
		{
			name:          "success - offset at end",
			offset:        15,
			limit:         10,
			expectedCount: 0,
			expectedTotal: 15,
		},
		{
			name:          "success - small page size",
			offset:        0,
			limit:         3,
			expectedCount: 3,
			expectedTotal: 15,
			validateResult: func(t *testing.T, users []*entity.User) {
				// Verify ordering (should be DESC by created_at)
				// The most recently created users should come first
				for i := 0; i < len(users)-1; i++ {
					assert.True(t, users[i].CreatedAt.After(users[i+1].CreatedAt) || users[i].CreatedAt.Equal(users[i+1].CreatedAt),
						"Users should be ordered by created_at DESC")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, total, err := repo.List(ctx, tt.offset, tt.limit)

			require.NoError(t, err)
			assert.Len(t, users, tt.expectedCount)
			assert.Equal(t, tt.expectedTotal, total)

			// Verify all users have required fields
			for _, user := range users {
				assert.NotZero(t, user.ID)
				assert.NotEmpty(t, user.Username)
				assert.NotEmpty(t, user.Email)
				assert.NotZero(t, user.CreatedAt)
			}

			if tt.validateResult != nil {
				tt.validateResult(t, users)
			}
		})
	}
}

func TestUserRepository_List_EmptyDatabase_Integration(t *testing.T) {
	container := testutil.NewTestContainer(t)
	repo := container.UserRepo

	ctx := context.Background()

	// Don't create any users
	users, total, err := repo.List(ctx, 0, 10)

	require.NoError(t, err)
	assert.Empty(t, users)
	assert.Equal(t, 0, total)
}

func TestUserRepository_ConcurrentOperations_Integration(t *testing.T) {
	container := testutil.NewTestContainer(t)
	repo := container.UserRepo

	ctx := context.Background()

	// Create a test user
	user := testutil.CreateTestUser(t, repo, "concurrent", "concurrent@example.com")

	// Simulate concurrent reads
	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func() {
			_, err := repo.FindByID(ctx, user.ID)
			assert.NoError(t, err)
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		<-done
	}
}
