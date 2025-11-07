package service

import (
	"context"
	"errors"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/domain/repository"
	"github.com/ryuyb/fusion/internal/domain/service"
	"github.com/ryuyb/fusion/internal/interface/http/dto/request"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/ryuyb/fusion/internal/pkg/utils"
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

func (s *userService) Create(ctx context.Context, req *request.CreateUserRequest) (*entity.User, error) {
	username, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil && !errors2.IsNotFoundError(err) {
		return nil, err
	}
	if username != nil {
		return nil, errors2.Conflict("username already exists")
	}
	email, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil && !errors2.IsNotFoundError(err) {
		return nil, err
	}
	if email != nil {
		return nil, errors2.Conflict("email already exists")
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors2.Internal(errors.New("hashing password failed")).Wrap(err)
	}
	user := entity.CreateUser(req.Username, hashedPassword, req.Email, req.Status)
	created, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return created, err
}

func (s *userService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) Update(ctx context.Context, req *request.UpdateUserRequest) (*entity.User, error) {
	user, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if req.Username != user.Username {
		existing, err := s.repo.FindByUsername(ctx, req.Username)
		if err != nil && !errors2.IsNotFoundError(err) {
			return nil, err
		}
		if existing != nil {
			return nil, errors2.Conflict("username already exists")
		}
	}
	if req.Email != user.Email {
		existing, err := s.repo.FindByEmail(ctx, req.Email)
		if err != nil && !errors2.IsNotFoundError(err) {
			return nil, err
		}
		if existing != nil {
			return nil, errors2.Conflict("email already exists")
		}
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors2.Internal(errors.New("hashing password failed")).Wrap(err)
	}
	user = user.Update(req.Username, hashedPassword, req.Email, req.Status)
	updated, err := s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	return updated, err
}

func (s *userService) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (s *userService) List(ctx context.Context, page, pageSize int) ([]*entity.User, int, error) {
	if page < 1 {
		return nil, 0, errors2.BadRequest("page must be greater than zero")
	}
	if pageSize < 1 || pageSize > 200 {
		return nil, 0, errors2.BadRequest("page size must be between 1 and 200")
	}

	offset := (page - 1) * pageSize
	return s.repo.List(ctx, offset, pageSize)
}
