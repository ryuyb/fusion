package service

import (
	"context"
	"testing"

	"github.com/ryuyb/fusion/internal/core/command"
	"github.com/ryuyb/fusion/internal/core/domain"
	"github.com/ryuyb/fusion/internal/core/port/repository"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/ryuyb/fusion/internal/pkg/util"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestUserServiceCreateSuccess(t *testing.T) {
	t.Parallel()
	repo := repository.NewMockUserRepository(t)
	ctx := context.Background()
	cmd := &command.CreateUserCommand{Username: "neo", Email: "neo@example.com", Password: "secret"}

	repo.EXPECT().ExistByUsername(mock.Anything, cmd.Username).Return(false, nil)
	repo.EXPECT().ExistByEmail(mock.Anything, cmd.Email).Return(false, nil)
	repo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Username == cmd.Username && u.Email == cmd.Email && u.Password != cmd.Password && util.VerifyPassword(cmd.Password, u.Password)
	})).Return(&domain.User{ID: 1}, nil)

	svc := NewUserService(repo, zap.NewNop())
	user, err := svc.Create(ctx, cmd)
	require.NoError(t, err)
	require.Equal(t, int64(1), user.ID)
}

func TestUserServiceCreateUsernameConflict(t *testing.T) {
	t.Parallel()
	repo := repository.NewMockUserRepository(t)
	repo.EXPECT().ExistByUsername(mock.Anything, "taken").Return(true, nil)

	svc := NewUserService(repo, zap.NewNop())
	_, err := svc.Create(context.Background(), &command.CreateUserCommand{Username: "taken", Email: "taken@example.com", Password: "secret"})
	require.Error(t, err)
	appErr := errors2.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, errors2.ErrCodeConflict, appErr.Code)
}

func TestUserServiceUpdateUsernameConflict(t *testing.T) {
	t.Parallel()
	repo := repository.NewMockUserRepository(t)
	existing := &domain.User{ID: 7, Username: "old", Email: "old@example.com"}
	repo.EXPECT().FindById(mock.Anything, existing.ID).Return(existing, nil)
	repo.EXPECT().FindByUsername(mock.Anything, "new").Return(&domain.User{ID: 9, Username: "new"}, nil)

	svc := NewUserService(repo, zap.NewNop())
	cmd := &command.UpdateUserCommand{CreateUserCommand: &command.CreateUserCommand{Username: "new", Email: existing.Email, Password: "pwd"}, ID: existing.ID}
	_, err := svc.Update(context.Background(), cmd)
	require.Error(t, err)
	appErr := errors2.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, errors2.ErrCodeConflict, appErr.Code)
}

func TestUserServiceUpdateSuccess(t *testing.T) {
	t.Parallel()
	repo := repository.NewMockUserRepository(t)
	existing := &domain.User{ID: 5, Username: "old", Email: "old@example.com", Password: "hashed"}
	repo.EXPECT().FindById(mock.Anything, existing.ID).Return(existing, nil)
	repo.EXPECT().FindByUsername(mock.Anything, "new").Return(nil, errors2.NotFound("user"))
	repo.EXPECT().FindByEmail(mock.Anything, "new@example.com").Return(nil, errors2.NotFound("user"))
	repo.EXPECT().Update(mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.ID == existing.ID && u.Username == "new" && u.Email == "new@example.com" && u.Password != "newpass" && util.VerifyPassword("newpass", u.Password)
	})).Return(&domain.User{ID: existing.ID, Username: "new", Email: "new@example.com"}, nil)

	svc := NewUserService(repo, zap.NewNop())
	cmd := &command.UpdateUserCommand{CreateUserCommand: &command.CreateUserCommand{Username: "new", Email: "new@example.com", Password: "newpass"}, ID: existing.ID}
	updated, err := svc.Update(context.Background(), cmd)
	require.NoError(t, err)
	require.Equal(t, "new", updated.Username)
}

func TestUserServiceListValidation(t *testing.T) {
	t.Parallel()
	repo := repository.NewMockUserRepository(t)
	svc := NewUserService(repo, zap.NewNop())

	_, _, err := svc.List(context.Background(), 0, 10)
	require.Equal(t, errors2.ErrCodeBadRequest, errors2.GetAppError(err).Code)

	_, _, err = svc.List(context.Background(), 1, 0)
	require.Equal(t, errors2.ErrCodeBadRequest, errors2.GetAppError(err).Code)

	_, _, err = svc.List(context.Background(), 1, 201)
	require.Equal(t, errors2.ErrCodeBadRequest, errors2.GetAppError(err).Code)
}

func TestUserServiceListSuccess(t *testing.T) {
	t.Parallel()
	repo := repository.NewMockUserRepository(t)
	repo.EXPECT().List(mock.Anything, 10, 5).Return([]*domain.User{{ID: 1}}, 1, nil)

	svc := NewUserService(repo, zap.NewNop())
	users, total, err := svc.List(context.Background(), 3, 5)
	require.NoError(t, err)
	require.Len(t, users, 1)
	require.Equal(t, 1, total)
}

func TestUserServiceDeleteDelegates(t *testing.T) {
	t.Parallel()
	repo := repository.NewMockUserRepository(t)
	repo.EXPECT().Delete(mock.Anything, int64(99)).Return(nil)

	svc := NewUserService(repo, zap.NewNop())
	require.NoError(t, svc.Delete(context.Background(), 99))
}
