package repository

import (
	"context"

	"github.com/ryuyb/fusion/internal/core/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
}
