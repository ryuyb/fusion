package service

import (
	"context"

	"github.com/ryuyb/fusion/internal/core/command"
	"github.com/ryuyb/fusion/internal/core/domain"
	"github.com/ryuyb/fusion/internal/core/port/repository"
	"github.com/ryuyb/fusion/internal/core/port/service"
)

type userService struct {
	r repository.UserRepository
}

func (u *userService) Create(ctx context.Context, cmd *command.CreateUserCommand) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserService(r repository.UserRepository) service.UserService {
	return &userService{
		r: r,
	}
}
