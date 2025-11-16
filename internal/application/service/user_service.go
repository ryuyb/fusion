package service

import (
	"context"

	"github.com/ryuyb/fusion/internal/app/errors"
	"github.com/ryuyb/fusion/internal/core/command"
	"github.com/ryuyb/fusion/internal/core/domain"
	"github.com/ryuyb/fusion/internal/core/port/repository"
	"github.com/ryuyb/fusion/internal/core/port/service"
	"go.uber.org/zap"
)

type userService struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

func (u *userService) Create(ctx context.Context, cmd *command.CreateUserCommand) (*domain.User, error) {
	existByUsername, err := u.repo.ExistByUsername(ctx, cmd.Username)
	if err != nil {
		return nil, err
	}
	if existByUsername {
		return nil, errors.Conflict("username already exists")
	}

	existByEmail, err := u.repo.ExistByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, err
	}
	if existByEmail {
		return nil, errors.Conflict("email already exists")
	}

	user, err := domain.CreateUser(cmd.Username, cmd.Email, cmd.Password)
	if err != nil {
		u.logger.Error("failed to create domain user", zap.Error(err))
		return nil, err
	}

	return u.repo.Create(ctx, user)
}

func (u *userService) Update(ctx context.Context, cmd *command.UpdateUserCommand) (*domain.User, error) {
	user, err := u.repo.FindById(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if cmd.Username != user.Username {
		u, err := u.repo.FindByUsername(ctx, cmd.Username)
		if err != nil && !errors.IsNotFoundError(err) {
			return nil, err
		}
		if u != nil {
			return nil, errors.Conflict("username already exists")
		}
	}
	if cmd.Email != user.Email {
		u, err := u.repo.FindByEmail(ctx, cmd.Email)
		if err != nil && !errors.IsNotFoundError(err) {
			return nil, err
		}
		if u != nil {
			return nil, errors.Conflict("email already exists")
		}
	}
	user, err = user.Update(cmd.Username, cmd.Email, cmd.Password)
	if err != nil {
		u.logger.Error("failed to update domain user", zap.Error(err))
		return nil, err
	}
	return u.repo.Update(ctx, user)
}

func (u *userService) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}

func (u *userService) FindById(ctx context.Context, id int64) (*domain.User, error) {
	return u.repo.FindById(ctx, id)
}

func (u *userService) List(ctx context.Context, page, pageSize int) ([]*domain.User, int, error) {
	if page < 1 {
		return nil, 0, errors.BadRequest("page must be greater than zero")
	}
	if pageSize < 1 || pageSize > 200 {
		return nil, 0, errors.BadRequest("page size must be between 1 and 200")
	}

	offset := (page - 1) * pageSize
	return u.repo.List(ctx, offset, pageSize)
}

func NewUserService(r repository.UserRepository, logger *zap.Logger) service.UserService {
	return &userService{
		repo:   r,
		logger: logger,
	}
}
