package service_test

import (
	"context"
	"testing"

	"github.com/ryuyb/fusion/internal/domain/service"
	"github.com/ryuyb/fusion/internal/interface/http/dto/request"
	"github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/ryuyb/fusion/internal/pkg/utils"
	"github.com/ryuyb/fusion/test/integration/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestUserService_Create_Integration(t *testing.T) {
	tests := []struct {
		name        string
		request     *request.CreateUserRequest
		setup       func(service.UserService)
		wantErr     bool
		errContains string
		validate    func(*testing.T, service.UserService)
	}{
		{
			name: "success - creates user with hashed password",
			request: &request.CreateUserRequest{
				Username: "newuser",
				Email:    "newuser@example.com",
				Password: "password123",
				Status:   "active",
			},
			validate: func(t *testing.T, svc service.UserService) {
				// Verify user exists
				ctx := context.Background()
				users, _, err := svc.List(ctx, 1, 10)
				require.NoError(t, err)
				assert.Len(t, users, 1)

				// Verify password is hashed
				assert.NotEqual(t, "password123", users[0].Password)
				assert.True(t, utils.VerifyPassword("password123", users[0].Password))
			},
		},
		{
			name: "error - username conflict detected",
			request: &request.CreateUserRequest{
				Username: "existing",
				Email:    "new@example.com",
				Password: "password123",
				Status:   "active",
			},
			setup: func(svc service.UserService) {
				// Pre-create user with same username
				ctx := context.Background()
				_, _ = svc.Create(ctx, &request.CreateUserRequest{
					Username: "existing",
					Email:    "other@example.com",
					Password: "password123",
					Status:   "active",
				})
			},
			wantErr:     true,
			errContains: "username already exists",
		},
		{
			name: "error - email conflict detected",
			request: &request.CreateUserRequest{
				Username: "newuser",
				Email:    "existing@example.com",
				Password: "password123",
				Status:   "active",
			},
			setup: func(svc service.UserService) {
				// Pre-create user with same email
				ctx := context.Background()
				_, _ = svc.Create(ctx, &request.CreateUserRequest{
					Username: "other",
					Email:    "existing@example.com",
					Password: "password123",
					Status:   "active",
				})
			},
			wantErr:     true,
			errContains: "email already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := testutil.NewTestContainer(t)
			svc := container.UserService

			if tt.setup != nil {
				tt.setup(svc)
			}

			ctx := context.Background()
			result, err := svc.Create(ctx, tt.request)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.NotZero(t, result.ID)
				assert.Equal(t, tt.request.Username, result.Username)
				assert.Equal(t, tt.request.Email, result.Email)

				if tt.validate != nil {
					tt.validate(t, svc)
				}
			}
		})
	}
}

func TestUserService_Update_Integration(t *testing.T) {
	container := testutil.NewTestContainer(t)
	svc := container.UserService

	ctx := context.Background()

	// Create initial user
	created, err := svc.Create(ctx, &request.CreateUserRequest{
		Username: "original",
		Email:    "original@example.com",
		Password: "password123",
		Status:   "active",
	})
	require.NoError(t, err)

	tests := []struct {
		name        string
		request     func() *request.UpdateUserRequest
		setup       func()
		wantErr     bool
		errContains string
		validate    func(*testing.T)
	}{
		{
			name: "success - updates username and email",
			request: func() *request.UpdateUserRequest {
				return &request.UpdateUserRequest{
					ID: created.ID,
					CreateUserRequest: request.CreateUserRequest{
						Username: "updated",
						Email:    "updated@example.com",
						Password: "newpassword",
						Status:   "active",
					},
				}
			},
			validate: func(t *testing.T) {
				user, err := svc.GetByID(ctx, created.ID)
				require.NoError(t, err)
				assert.Equal(t, "updated", user.Username)
				assert.Equal(t, "updated@example.com", user.Email)
				assert.True(t, utils.VerifyPassword("newpassword", user.Password))
			},
		},
		{
			name: "success - changes status to inactive",
			request: func() *request.UpdateUserRequest {
				return &request.UpdateUserRequest{
					ID: created.ID,
					CreateUserRequest: request.CreateUserRequest{
						Username: "updated",
						Email:    "updated@example.com",
						Password: "newpassword",
						Status:   "inactive",
					},
				}
			},
			validate: func(t *testing.T) {
				user, err := svc.GetByID(ctx, created.ID)
				require.NoError(t, err)
				assert.Equal(t, "inactive", string(user.Status))
			},
		},
		{
			name: "error - update with conflicting username",
			setup: func() {
				// Create another user
				_, _ = svc.Create(ctx, &request.CreateUserRequest{
					Username: "taken",
					Email:    "taken@example.com",
					Password: "password123",
					Status:   "active",
				})
			},
			request: func() *request.UpdateUserRequest {
				return &request.UpdateUserRequest{
					ID: created.ID,
					CreateUserRequest: request.CreateUserRequest{
						Username: "taken",
						Email:    "newemail@example.com",
						Password: "password123",
						Status:   "active",
					},
				}
			},
			wantErr:     true,
			errContains: "username already exists",
		},
		{
			name: "error - user not found",
			request: func() *request.UpdateUserRequest {
				return &request.UpdateUserRequest{
					ID: 999999,
					CreateUserRequest: request.CreateUserRequest{
						Username: "notfound",
						Email:    "notfound@example.com",
						Password: "password123",
						Status:   "active",
					},
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			req := tt.request()
			result, err := svc.Update(ctx, req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)

				if tt.validate != nil {
					tt.validate(t)
				}
			}
		})
	}
}

func TestUserService_Delete_Integration(t *testing.T) {
	container := testutil.NewTestContainer(t)
	svc := container.UserService

	ctx := context.Background()

	// Create test user
	created, err := svc.Create(ctx, &request.CreateUserRequest{
		Username: "todelete",
		Email:    "delete@example.com",
		Password: "password123",
		Status:   "active",
	})
	require.NoError(t, err)

	tests := []struct {
		name        string
		userID      int64
		wantErr     bool
		errContains string
		validate    func(*testing.T)
	}{
		{
			name:   "success - soft deletes user",
			userID: created.ID,
			validate: func(t *testing.T) {
				// Verify user cannot be retrieved
				_, err := svc.GetByID(ctx, created.ID)
				assert.Error(t, err)
				assert.True(t, errors.IsNotFoundError(err))
			},
		},
		{
			name:        "error - user not found",
			userID:      999999,
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Delete(ctx, tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)

				if tt.validate != nil {
					tt.validate(t)
				}
			}
		})
	}
}

func TestUserService_GetByID_Integration(t *testing.T) {
	container := testutil.NewTestContainer(t)
	svc := container.UserService

	ctx := context.Background()

	// Create test user
	created, err := svc.Create(ctx, &request.CreateUserRequest{
		Username: "getbyid",
		Email:    "getbyid@example.com",
		Password: "password123",
		Status:   "active",
	})
	require.NoError(t, err)

	tests := []struct {
		name        string
		userID      int64
		wantErr     bool
		errContains string
	}{
		{
			name:   "success - retrieves existing user",
			userID: created.ID,
		},
		{
			name:        "error - user not found",
			userID:      999999,
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.GetByID(ctx, tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
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

func TestUserService_List_Integration(t *testing.T) {
	container := testutil.NewTestContainer(t)
	svc := container.UserService

	ctx := context.Background()

	// Create 25 test users
	for i := 1; i <= 25; i++ {
		_, err := svc.Create(ctx, &request.CreateUserRequest{
			Username: testutil.RandomUsername(),
			Email:    testutil.RandomEmail(),
			Password: "password123",
			Status:   "active",
		})
		require.NoError(t, err)
	}

	tests := []struct {
		name          string
		page          int
		pageSize      int
		wantErr       bool
		errContains   string
		expectedCount int
		expectedTotal int
	}{
		{
			name:          "success - first page with 10 items",
			page:          1,
			pageSize:      10,
			expectedCount: 10,
			expectedTotal: 25,
		},
		{
			name:          "success - second page with 10 items",
			page:          2,
			pageSize:      10,
			expectedCount: 10,
			expectedTotal: 25,
		},
		{
			name:          "success - third page with 5 items",
			page:          3,
			pageSize:      10,
			expectedCount: 5,
			expectedTotal: 25,
		},
		{
			name:          "success - page with all items",
			page:          1,
			pageSize:      100,
			expectedCount: 25,
			expectedTotal: 25,
		},
		{
			name:        "error - invalid page (zero)",
			page:        0,
			pageSize:    10,
			wantErr:     true,
			errContains: "page must be greater than zero",
		},
		{
			name:        "error - invalid page (negative)",
			page:        -1,
			pageSize:    10,
			wantErr:     true,
			errContains: "page must be greater than zero",
		},
		{
			name:        "error - invalid pageSize (zero)",
			page:        1,
			pageSize:    0,
			wantErr:     true,
			errContains: "page size must be between 1 and 200",
		},
		{
			name:        "error - invalid pageSize (too large)",
			page:        1,
			pageSize:    201,
			wantErr:     true,
			errContains: "page size must be between 1 and 200",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, total, err := svc.List(ctx, tt.page, tt.pageSize)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, users)
				assert.Equal(t, 0, total)
			} else {
				require.NoError(t, err)
				assert.Len(t, users, tt.expectedCount)
				assert.Equal(t, tt.expectedTotal, total)

				// Verify all users have valid data
				for _, user := range users {
					assert.NotZero(t, user.ID)
					assert.NotEmpty(t, user.Username)
					assert.NotEmpty(t, user.Email)
					assert.NotEmpty(t, user.Password)
					assert.NotZero(t, user.CreatedAt)
				}
			}
		})
	}
}

func TestUserService_CompleteWorkflow_Integration(t *testing.T) {
	container := testutil.NewTestContainer(t)
	svc := container.UserService

	ctx := context.Background()

	// 1. Create user
	created, err := svc.Create(ctx, &request.CreateUserRequest{
		Username: "workflow",
		Email:    "workflow@example.com",
		Password: "password123",
		Status:   "active",
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	assert.NotZero(t, created.ID)

	// 2. Retrieve user by ID
	retrieved, err := svc.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, "workflow", retrieved.Username)

	// 3. Update user
	updated, err := svc.Update(ctx, &request.UpdateUserRequest{
		ID: created.ID,
		CreateUserRequest: request.CreateUserRequest{
			Username: "workflow_updated",
			Email:    "workflow_updated@example.com",
			Password: "newpassword",
			Status:   "inactive",
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "workflow_updated", updated.Username)
	assert.Equal(t, "workflow_updated@example.com", updated.Email)
	assert.Equal(t, "inactive", string(updated.Status))

	// 4. Verify update persisted
	retrieved, err = svc.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, "workflow_updated", retrieved.Username)
	assert.True(t, utils.VerifyPassword("newpassword", retrieved.Password))

	// 5. List users (should include our user)
	users, total, err := svc.List(ctx, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, users, 1)
	assert.Equal(t, created.ID, users[0].ID)

	// 6. Delete user
	err = svc.Delete(ctx, created.ID)
	require.NoError(t, err)

	// 7. Verify user is deleted (not found)
	_, err = svc.GetByID(ctx, created.ID)
	require.Error(t, err)
	assert.True(t, errors.IsNotFoundError(err))

	// 8. List users (should be empty)
	users, total, err = svc.List(ctx, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Len(t, users, 0)
}
