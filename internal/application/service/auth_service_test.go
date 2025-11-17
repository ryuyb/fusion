package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ryuyb/fusion/internal/core/command"
	"github.com/ryuyb/fusion/internal/core/domain"
	"github.com/ryuyb/fusion/internal/core/port/repository"
	"github.com/ryuyb/fusion/internal/infrastructure/provider/config"
	"github.com/ryuyb/fusion/internal/infrastructure/provider/jwt"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/ryuyb/fusion/internal/pkg/util"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAuthServiceLoginSuccess(t *testing.T) {
	t.Parallel()

	repo := repository.NewMockUserRepository(t)
	ctx := context.Background()
	password := "secret-password"
	hashed, err := util.HashPassword(password)
	require.NoError(t, err)
	user := &domain.User{
		ID:       1,
		Username: "alice",
		Email:    "alice@example.com",
		Password: hashed,
	}
	repo.EXPECT().FindByUsername(mock.Anything, "alice").Return(user, nil)

	cfg := &config.Config{
		App: config.AppConfig{Name: "fusion-test"},
		JWT: config.JWTConfig{Secret: "test-secret", Expiration: time.Minute},
	}
	service := NewAuthService(repo, jwt.NewJWTManager(cfg), zap.NewNop())

	token, expiresAt, err := service.Login(ctx, command.LoginCommand{Username: "alice", Password: password})
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.WithinDuration(t, time.Now().Add(cfg.JWT.Expiration), expiresAt, 2*time.Second)
}

func TestAuthServiceLoginInvalidPassword(t *testing.T) {
	t.Parallel()

	repo := repository.NewMockUserRepository(t)
	ctx := context.Background()
	hashed, err := util.HashPassword("correct-password")
	require.NoError(t, err)
	repo.EXPECT().FindByUsername(mock.Anything, "bob").Return(&domain.User{
		Username: "bob",
		Password: hashed,
	}, nil)

	cfg := &config.Config{App: config.AppConfig{Name: "fusion-test"}, JWT: config.JWTConfig{Secret: "secret", Expiration: time.Minute}}
	service := NewAuthService(repo, jwt.NewJWTManager(cfg), zap.NewNop())

	_, _, err = service.Login(ctx, command.LoginCommand{Username: "bob", Password: "wrong"})
	require.Error(t, err)
	appErr := errors2.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, errors2.ErrCodeUnauthorized, appErr.Code)
}

func TestAuthServiceRegisterValidationAndCreationFlow(t *testing.T) {
	t.Parallel()

	repo := repository.NewMockUserRepository(t)
	ctx := context.Background()
	cmd := command.RegisterCommand{Username: "new-user", Email: "new@example.com", Password: "strong-pass"}

	repo.EXPECT().ExistByUsername(mock.Anything, cmd.Username).Return(false, nil)
	repo.EXPECT().ExistByEmail(mock.Anything, cmd.Email).Return(false, nil)
	repo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Username == cmd.Username && u.Email == cmd.Email && u.Password != "" && u.Password != cmd.Password
	})).Return(&domain.User{}, nil)

	cfg := &config.Config{App: config.AppConfig{Name: "fusion-test"}, JWT: config.JWTConfig{Secret: "secret", Expiration: time.Minute}}
	svc := NewAuthService(repo, jwt.NewJWTManager(cfg), zap.NewNop())

	require.NoError(t, svc.Register(ctx, cmd))
}

func TestAuthServiceRegisterUsernameConflict(t *testing.T) {
	t.Parallel()

	repo := repository.NewMockUserRepository(t)
	repo.EXPECT().ExistByUsername(mock.Anything, "taken").Return(true, nil)

	cfg := &config.Config{App: config.AppConfig{Name: "fusion-test"}, JWT: config.JWTConfig{Secret: "secret", Expiration: time.Minute}}
	svc := NewAuthService(repo, jwt.NewJWTManager(cfg), zap.NewNop())

	err := svc.Register(context.Background(), command.RegisterCommand{Username: "taken", Email: "someone@example.com", Password: "pass"})
	require.Error(t, err)
	appErr := errors2.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, errors2.ErrCodeConflict, appErr.Code)
}

func TestAuthServiceRegisterCreateError(t *testing.T) {
	t.Parallel()

	repo := repository.NewMockUserRepository(t)
	repo.EXPECT().ExistByUsername(mock.Anything, "neo").Return(false, nil)
	repo.EXPECT().ExistByEmail(mock.Anything, "neo@example.com").Return(false, nil)
	expectedErr := errors.New("db error")
	repo.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, expectedErr)

	cfg := &config.Config{App: config.AppConfig{Name: "fusion-test"}, JWT: config.JWTConfig{Secret: "secret", Expiration: time.Minute}}
	svc := NewAuthService(repo, jwt.NewJWTManager(cfg), zap.NewNop())

	err := svc.Register(context.Background(), command.RegisterCommand{Username: "neo", Email: "neo@example.com", Password: "pass"})
	require.ErrorIs(t, err, expectedErr)
}
