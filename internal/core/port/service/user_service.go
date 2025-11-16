package service

import (
	"context"

	"github.com/ryuyb/fusion/internal/core/command"
	"github.com/ryuyb/fusion/internal/core/domain"
)

type UserService interface {
	Create(ctx context.Context, cmd *command.CreateUserCommand) (*domain.User, error)

	Update(ctx context.Context, cmd *command.UpdateUserCommand) (*domain.User, error)

	Delete(ctx context.Context, id int64) error

	FindById(ctx context.Context, id int64) (*domain.User, error)

	List(ctx context.Context, page, pageSize int) ([]*domain.User, int, error)
}
