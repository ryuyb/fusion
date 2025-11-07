package service

import (
	"context"
	"time"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/domain/repository"
	"github.com/ryuyb/fusion/internal/domain/service"
	"github.com/ryuyb/fusion/internal/interface/http/dto/request"
	"github.com/ryuyb/fusion/internal/pkg/auth"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/ryuyb/fusion/internal/pkg/utils"
	"go.uber.org/zap"
)

type authService struct {
	userRepo   repository.UserRepository
	jwtManager *auth.JWTManager
	logger     *zap.Logger
}

// NewAuthService creates a new auth service instance
func NewAuthService(
	userRepo repository.UserRepository,
	jwtManager *auth.JWTManager,
	logger *zap.Logger,
) service.AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// Login authenticates a user with username and password, returns JWT token
func (s *authService) Login(ctx context.Context, req *request.LoginRequest) (string, time.Time, error) {
	// Find user by username
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors2.IsNotFoundError(err) {
			return "", time.Time{}, errors2.Unauthorized("invalid username or password")
		}
		return "", time.Time{}, err
	}

	// Verify password
	if !utils.VerifyPassword(req.Password, user.Password) {
		return "", time.Time{}, errors2.Unauthorized("invalid username or password")
	}

	// Check user status
	if user.Status != entity.UserStatusActive {
		return "", time.Time{}, errors2.Unauthorized("user account is not active")
	}

	// Generate JWT token
	token, expiresAt, err := s.jwtManager.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		s.logger.Error("failed to generate token", zap.Error(err))
		return "", time.Time{}, errors2.Internal(err)
	}

	return token, expiresAt, nil
}

// Register creates a new user and returns JWT token
func (s *authService) Register(ctx context.Context, req *request.RegisterRequest) (string, time.Time, error) {
	// Check if username already exists
	existingUser, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil && !errors2.IsNotFoundError(err) {
		return "", time.Time{}, err
	}
	if existingUser != nil {
		return "", time.Time{}, errors2.Conflict("username already exists")
	}

	// Check if email already exists
	existingEmail, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && !errors2.IsNotFoundError(err) {
		return "", time.Time{}, err
	}
	if existingEmail != nil {
		return "", time.Time{}, errors2.Conflict("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.logger.Error("failed to hash password", zap.Error(err))
		return "", time.Time{}, errors2.Internal(err)
	}

	// Create user entity (default status: active)
	user := entity.CreateUser(req.Username, hashedPassword, req.Email, string(entity.UserStatusActive))

	// Save user to database
	created, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return "", time.Time{}, err
	}

	// Generate JWT token
	token, expiresAt, err := s.jwtManager.GenerateToken(created.ID, created.Username, created.Email)
	if err != nil {
		s.logger.Error("failed to generate token", zap.Error(err))
		return "", time.Time{}, errors2.Internal(err)
	}

	return token, expiresAt, nil
}
