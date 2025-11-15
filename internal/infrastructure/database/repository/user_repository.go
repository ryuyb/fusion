package repository

import (
	"context"

	"github.com/ryuyb/fusion/internal/core/domain"
	"github.com/ryuyb/fusion/internal/core/port/repository"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent"
)

type userRepository struct {
	client *ent.Client
}

func (u *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserRepository(client *ent.Client) repository.UserRepository {
	return &userRepository{client: client}
}
