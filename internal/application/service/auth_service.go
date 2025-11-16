package service

import (
	"context"
	"time"

	"github.com/ryuyb/fusion/internal/core/command"
	"github.com/ryuyb/fusion/internal/core/domain"
	"github.com/ryuyb/fusion/internal/core/port/repository"
	"github.com/ryuyb/fusion/internal/core/port/service"
	"github.com/ryuyb/fusion/internal/infrastructure/provider/jwt"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/ryuyb/fusion/internal/pkg/util"
	"go.uber.org/zap"
)

type authService struct {
	userRepo   repository.UserRepository
	jwtManager *jwt.JWTManager
	logger     *zap.Logger
}

func NewAuthService(userRepo repository.UserRepository, jwtManager *jwt.JWTManager, logger *zap.Logger) service.AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

func (service *authService) Login(ctx context.Context, cmd command.LoginCommand) (string, time.Time, error) {
	user, err := service.userRepo.FindByUsername(ctx, cmd.Username)
	if err != nil {
		return "", time.Time{}, err
	}
	if !util.VerifyPassword(cmd.Password, user.Password) {
		return "", time.Time{}, errors2.Unauthorized("invalid username or password")
	}
	token, expiresAt, err := service.jwtManager.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		service.logger.Warn("failed to generate token", zap.Error(err))
		return "", time.Time{}, errors2.Internal(err)
	}
	return token, expiresAt, nil
}

func (service *authService) Register(ctx context.Context, cmd command.RegisterCommand) error {
	existByUsername, err := service.userRepo.ExistByUsername(ctx, cmd.Username)
	if err != nil {
		return err
	}
	if existByUsername {
		return errors2.Conflict("username already exist")
	}
	existByEmail, err := service.userRepo.ExistByEmail(ctx, cmd.Email)
	if err != nil {
		return err
	}
	if existByEmail {
		return errors2.Conflict("email already exist")
	}

	user, err := domain.CreateUser(cmd.Username, cmd.Email, cmd.Password)
	if err != nil {
		service.logger.Error("failed to create domain user", zap.Error(err))
		return err
	}

	if _, err := service.userRepo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}
