package repository

import (
	"context"

	"github.com/ryuyb/fusion/internal/core/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)

	Update(ctx context.Context, user *domain.User) (*domain.User, error)

	Delete(ctx context.Context, id int64) error

	FindById(ctx context.Context, id int64) (*domain.User, error)

	FindByUsername(ctx context.Context, username string) (*domain.User, error)

	ExistByUsername(ctx context.Context, username string) (bool, error)

	FindByEmail(ctx context.Context, email string) (*domain.User, error)

	ExistByEmail(ctx context.Context, email string) (bool, error)

	List(ctx context.Context, offset, limit int) ([]*domain.User, int, error)
}
