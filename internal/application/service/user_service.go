package service

import (
	"github.com/ryuyb/fusion/internal/domain/repository"
	"github.com/ryuyb/fusion/internal/domain/service"
	"go.uber.org/zap"
)

type userService struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

func NewUserService(repo repository.UserRepository, logger *zap.Logger) service.UserService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}
