package repository

import (
	"context"

	"github.com/ryuyb/fusion/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)

	FindByID(ctx context.Context, id int64) (*entity.User, error)

	FindByUsername(ctx context.Context, username string) (*entity.User, error)

	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	Delete(ctx context.Context, id int64) error

	Update(ctx context.Context, u *entity.User) (*entity.User, error)

	List(ctx context.Context, offset, limit int) ([]*entity.User, int, error)
}
