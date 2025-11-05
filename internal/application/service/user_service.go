package service

import (
	"context"

	"github.com/ryuyb/fusion/internal/application/dto"
	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/domain/repository"
	"github.com/ryuyb/fusion/internal/domain/service"
	"github.com/ryuyb/fusion/internal/interface/http/response"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
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

func (s *userService) Create(ctx context.Context, u *entity.User) (*entity.User, error) {
	return s.repo.Create(ctx, u)
}

func (s *userService) List(ctx context.Context, page, pageSize int) (*response.PaginationResponse[*dto.UserResponse], error) {
	if page < 1 {
		return nil, errors2.BadRequest("page must be greater than zero")
	}
	if pageSize < 1 || pageSize > 200 {
		return nil, errors2.BadRequest("page size must be between 1 and 200")
	}

	offset := (page - 1) * pageSize
	users, total, err := s.repo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, err
	}

	items := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		items[i] = s.toResponse(user)
	}

	return response.NewPaginationResponse[*dto.UserResponse](items, total, page, pageSize), nil
}

func (s *userService) toResponse(u *entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeleteAt:  u.DeleteAt,
	}
}
