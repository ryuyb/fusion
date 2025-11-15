package service

import (
	"context"

	"github.com/ryuyb/fusion/internal/core/command"
	"github.com/ryuyb/fusion/internal/core/domain"
)

type UserService interface {
	Create(ctx context.Context, cmd *command.CreateUserCommand) (*domain.User, error)
}
