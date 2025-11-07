package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ryuyb/fusion/internal/application/service/mocks"
	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/interface/http/dto/request"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

func TestUserService_Create(t *testing.T) {
	tests := []struct {
		name        string
		request     *request.CreateUserRequest
		setupMock   func(*mocks.MockUserRepository)
		wantErr     bool
		errContains string
		validate    func(*testing.T, *entity.User)
	}{
		{
			name: "success - creates user successfully",
			request: &request.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Status:   "active",
			},
			setupMock: func(repo *mocks.MockUserRepository) {
				// FindByUsername returns nil (user doesn't exist)
				repo.On("FindByUsername", mock.Anything, "testuser").
					Return(nil, errors2.NotFound("user"))

				// FindByEmail returns nil (email doesn't exist)
				repo.On("FindByEmail", mock.Anything, "test@example.com").
					Return(nil, errors2.NotFound("user"))

				// Create succeeds
				repo.On("Create", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
					return u.Username == "testuser" && u.Email == "test@example.com"
				})).Return(&entity.User{
					ID:        1,
					Username:  "testuser",
					Email:     "test@example.com",
					Status:    entity.UserStatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, user *entity.User) {
				assert.NotNil(t, user)
				assert.Equal(t, int64(1), user.ID)
				assert.Equal(t, "testuser", user.Username)
				assert.Equal(t, "test@example.com", user.Email)
				assert.Equal(t, entity.UserStatusActive, user.Status)
			},
		},
		{
			name: "error - username already exists",
			request: &request.CreateUserRequest{
				Username: "existing",
				Email:    "new@example.com",
				Password: "password123",
				Status:   "active",
			},
			setupMock: func(repo *mocks.MockUserRepository) {
				// FindByUsername returns existing user
				repo.On("FindByUsername", mock.Anything, "existing").
					Return(&entity.User{
						ID:       1,
						Username: "existing",
					}, nil)
			},
			wantErr:     true,
			errContains: "username already exists",
		},
		{
			name: "error - email already exists",
			request: &request.CreateUserRequest{
				Username: "newuser",
				Email:    "existing@example.com",
				Password: "password123",
				Status:   "active",
			},
			setupMock: func(repo *mocks.MockUserRepository) {
				// FindByUsername returns nil
				repo.On("FindByUsername", mock.Anything, "newuser").
					Return(nil, errors2.NotFound("user"))

				// FindByEmail returns existing user
				repo.On("FindByEmail", mock.Anything, "existing@example.com").
					Return(&entity.User{
						ID:    1,
						Email: "existing@example.com",
					}, nil)
			},
			wantErr:     true,
			errContains: "email already exists",
		},
		{
			name: "error - FindByUsername fails with non-NotFound error",
			request: &request.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Status:   "active",
			},
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.On("FindByUsername", mock.Anything, "testuser").
					Return(nil, errors.New("database error"))
			},
			wantErr:     true,
			errContains: "database error",
		},
		{
			name: "error - FindByEmail fails with non-NotFound error",
			request: &request.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Status:   "active",
			},
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.On("FindByUsername", mock.Anything, "testuser").
					Return(nil, errors2.NotFound("user"))
				repo.On("FindByEmail", mock.Anything, "test@example.com").
					Return(nil, errors.New("database error"))
			},
			wantErr:     true,
			errContains: "database error",
		},
		{
			name: "error - repository Create fails",
			request: &request.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Status:   "active",
			},
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.On("FindByUsername", mock.Anything, "testuser").
					Return(nil, errors2.NotFound("user"))
				repo.On("FindByEmail", mock.Anything, "test@example.com").
					Return(nil, errors2.NotFound("user"))
				repo.On("Create", mock.Anything, mock.Anything).
					Return(nil, errors.New("create failed"))
			},
			wantErr:     true,
			errContains: "create failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(mocks.MockUserRepository)
			tt.setupMock(mockRepo)

			logger := zaptest.NewLogger(t)
			svc := NewUserService(mockRepo, logger)

			// Execute
			result, err := svc.Create(context.Background(), tt.request)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Update(t *testing.T) {
	tests := []struct {
		name        string
		request     *request.UpdateUserRequest
		setupMock   func(*mocks.MockUserRepository)
		wantErr     bool
		errContains string
		validate    func(*testing.T, *entity.User)
	}{
		{
			name: "success - updates user successfully",
			request: &request.UpdateUserRequest{
				ID: 1,
				CreateUserRequest: request.CreateUserRequest{
					Username: "updated",
					Email:    "updated@example.com",
					Password: "newpassword123",
					Status:   "active",
				},
			},
			setupMock: func(repo *mocks.MockUserRepository) {
				existingUser := &entity.User{
					ID:       1,
					Username: "oldname",
					Email:    "old@example.com",
					Password: "oldhashedpass",
					Status:   entity.UserStatusActive,
				}

				// FindByID returns existing user
				repo.On("FindByID", mock.Anything, int64(1)).Return(existingUser, nil)

				// Check username availability
				repo.On("FindByUsername", mock.Anything, "updated").
					Return(nil, errors2.NotFound("user"))

				// Check email availability
				repo.On("FindByEmail", mock.Anything, "updated@example.com").
					Return(nil, errors2.NotFound("user"))

				// Update succeeds
				repo.On("Update", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
					return u.ID == 1 && u.Username == "updated"
				})).Return(&entity.User{
					ID:        1,
					Username:  "updated",
					Email:     "updated@example.com",
					Status:    entity.UserStatusActive,
					UpdatedAt: time.Now(),
				}, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, user *entity.User) {
				assert.NotNil(t, user)
				assert.Equal(t, int64(1), user.ID)
				assert.Equal(t, "updated", user.Username)
				assert.Equal(t, "updated@example.com", user.Email)
			},
		},
		{
			name: "success - same username and email (no conflict check)",
			request: &request.UpdateUserRequest{
				ID: 1,
				CreateUserRequest: request.CreateUserRequest{
					Username: "samename",
					Email:    "same@example.com",
					Password: "newpassword123",
					Status:   "active",
				},
			},
			setupMock: func(repo *mocks.MockUserRepository) {
				existingUser := &entity.User{
					ID:       1,
					Username: "samename",
					Email:    "same@example.com",
					Password: "oldhashedpass",
					Status:   entity.UserStatusActive,
				}

				repo.On("FindByID", mock.Anything, int64(1)).Return(existingUser, nil)

				// No FindByUsername or FindByEmail calls because they match existing

				repo.On("Update", mock.Anything, mock.Anything).Return(&entity.User{
					ID:       1,
					Username: "samename",
					Email:    "same@example.com",
					Status:   entity.UserStatusActive,
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error - user not found",
			request: &request.UpdateUserRequest{
				ID: 999,
				CreateUserRequest: request.CreateUserRequest{
					Username: "updated",
					Email:    "updated@example.com",
					Password: "newpassword123",
					Status:   "active",
				},
			},
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.On("FindByID", mock.Anything, int64(999)).
					Return(nil, errors2.NotFound("user"))
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "error - new username already exists",
			request: &request.UpdateUserRequest{
				ID: 1,
				CreateUserRequest: request.CreateUserRequest{
					Username: "taken",
					Email:    "new@example.com",
					Password: "newpassword123",
					Status:   "active",
				},
			},
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.On("FindByID", mock.Anything, int64(1)).Return(&entity.User{
					ID:       1,
					Username: "original",
					Email:    "original@example.com",
				}, nil)

				// New username is taken by another user
				repo.On("FindByUsername", mock.Anything, "taken").Return(&entity.User{
					ID:       2,
					Username: "taken",
				}, nil)
			},
			wantErr:     true,
			errContains: "username already exists",
		},
		{
			name: "error - new email already exists",
			request: &request.UpdateUserRequest{
				ID: 1,
				CreateUserRequest: request.CreateUserRequest{
					Username: "newname",
					Email:    "taken@example.com",
					Password: "newpassword123",
					Status:   "active",
				},
			},
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.On("FindByID", mock.Anything, int64(1)).Return(&entity.User{
					ID:       1,
					Username: "original",
					Email:    "original@example.com",
				}, nil)

				repo.On("FindByUsername", mock.Anything, "newname").
					Return(nil, errors2.NotFound("user"))

				// New email is taken by another user
				repo.On("FindByEmail", mock.Anything, "taken@example.com").Return(&entity.User{
					ID:    2,
					Email: "taken@example.com",
				}, nil)
			},
			wantErr:     true,
			errContains: "email already exists",
		},
		{
			name: "error - repository Update fails",
			request: &request.UpdateUserRequest{
				ID: 1,
				CreateUserRequest: request.CreateUserRequest{
					Username: "updated",
					Email:    "updated@example.com",
					Password: "newpassword123",
					Status:   "active",
				},
			},
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.On("FindByID", mock.Anything, int64(1)).Return(&entity.User{
					ID:       1,
					Username: "original",
					Email:    "original@example.com",
				}, nil)

				repo.On("FindByUsername", mock.Anything, "updated").
					Return(nil, errors2.NotFound("user"))
				repo.On("FindByEmail", mock.Anything, "updated@example.com").
					Return(nil, errors2.NotFound("user"))

				repo.On("Update", mock.Anything, mock.Anything).
					Return(nil, errors.New("update failed"))
			},
			wantErr:     true,
			errContains: "update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(mocks.MockUserRepository)
			tt.setupMock(mockRepo)

			logger := zaptest.NewLogger(t)
			svc := NewUserService(mockRepo, logger)

			// Execute
			result, err := svc.Update(context.Background(), tt.request)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		setupMock   func(*mocks.MockUserRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:   "success - deletes user successfully",
			userID: 1,
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.On("Delete", mock.Anything, int64(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "error - repository Delete fails",
			userID: 1,
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.On("Delete", mock.Anything, int64(1)).
					Return(errors.New("delete failed"))
			},
			wantErr:     true,
			errContains: "delete failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(mocks.MockUserRepository)
			tt.setupMock(mockRepo)

			logger := zaptest.NewLogger(t)
			svc := NewUserService(mockRepo, logger)

			// Execute
			err := svc.Delete(context.Background(), tt.userID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetByID(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		setupMock   func(*mocks.MockUserRepository)
		wantErr     bool
		errContains string
		validate    func(*testing.T, *entity.User)
	}{
		{
			name:   "success - retrieves user by ID",
			userID: 1,
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.On("FindByID", mock.Anything, int64(1)).Return(&entity.User{
					ID:       1,
					Username: "testuser",
					Email:    "test@example.com",
					Status:   entity.UserStatusActive,
				}, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, user *entity.User) {
				assert.NotNil(t, user)
				assert.Equal(t, int64(1), user.ID)
				assert.Equal(t, "testuser", user.Username)
				assert.Equal(t, "test@example.com", user.Email)
			},
		},
		{
			name:   "error - user not found",
			userID: 999,
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.On("FindByID", mock.Anything, int64(999)).
					Return(nil, errors2.NotFound("user"))
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name:   "error - repository FindByID fails",
			userID: 1,
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.On("FindByID", mock.Anything, int64(1)).
					Return(nil, errors.New("database error"))
			},
			wantErr:     true,
			errContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(mocks.MockUserRepository)
			tt.setupMock(mockRepo)

			logger := zaptest.NewLogger(t)
			svc := NewUserService(mockRepo, logger)

			// Execute
			result, err := svc.GetByID(context.Background(), tt.userID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_List(t *testing.T) {
	tests := []struct {
		name        string
		page        int
		pageSize    int
		setupMock   func(*mocks.MockUserRepository)
		wantErr     bool
		errContains string
		validate    func(*testing.T, []*entity.User, int)
	}{
		{
			name:     "success - lists users with pagination",
			page:     1,
			pageSize: 10,
			setupMock: func(repo *mocks.MockUserRepository) {
				users := []*entity.User{
					{ID: 1, Username: "user1", Email: "user1@example.com"},
					{ID: 2, Username: "user2", Email: "user2@example.com"},
				}
				repo.On("List", mock.Anything, 0, 10).Return(users, 2, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, users []*entity.User, total int) {
				assert.Len(t, users, 2)
				assert.Equal(t, 2, total)
				assert.Equal(t, "user1", users[0].Username)
				assert.Equal(t, "user2", users[1].Username)
			},
		},
		{
			name:     "success - page 2 with correct offset",
			page:     2,
			pageSize: 10,
			setupMock: func(repo *mocks.MockUserRepository) {
				users := []*entity.User{
					{ID: 11, Username: "user11", Email: "user11@example.com"},
				}
				repo.On("List", mock.Anything, 10, 10).Return(users, 15, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, users []*entity.User, total int) {
				assert.Len(t, users, 1)
				assert.Equal(t, 15, total)
			},
		},
		{
			name:     "error - page less than 1",
			page:     0,
			pageSize: 10,
			setupMock: func(repo *mocks.MockUserRepository) {
				// No repository calls expected
			},
			wantErr:     true,
			errContains: "page must be greater than zero",
		},
		{
			name:     "error - pageSize less than 1",
			page:     1,
			pageSize: 0,
			setupMock: func(repo *mocks.MockUserRepository) {
				// No repository calls expected
			},
			wantErr:     true,
			errContains: "page size must be between 1 and 200",
		},
		{
			name:     "error - pageSize greater than 200",
			page:     1,
			pageSize: 201,
			setupMock: func(repo *mocks.MockUserRepository) {
				// No repository calls expected
			},
			wantErr:     true,
			errContains: "page size must be between 1 and 200",
		},
		{
			name:     "error - repository List fails",
			page:     1,
			pageSize: 10,
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.On("List", mock.Anything, 0, 10).
					Return(nil, 0, errors.New("database error"))
			},
			wantErr:     true,
			errContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(mocks.MockUserRepository)
			tt.setupMock(mockRepo)

			logger := zaptest.NewLogger(t)
			svc := NewUserService(mockRepo, logger)

			// Execute
			users, total, err := svc.List(context.Background(), tt.page, tt.pageSize)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, users)
				assert.Equal(t, 0, total)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, users, total)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
