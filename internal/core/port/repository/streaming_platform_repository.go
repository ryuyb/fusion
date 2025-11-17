package repository

import (
	"context"

	"github.com/ryuyb/fusion/internal/core/domain"
)

type StreamingPlatformRepository interface {
	Create(ctx context.Context, platform *domain.StreamingPlatform) (*domain.StreamingPlatform, error)

	Update(ctx context.Context, platform *domain.StreamingPlatform) (*domain.StreamingPlatform, error)

	Delete(ctx context.Context, id int64) error

	FindById(ctx context.Context, id int64) (*domain.StreamingPlatform, error)

	FindByType(ctx context.Context, platformType domain.StreamingPlatformType) (*domain.StreamingPlatform, error)

	FindByName(ctx context.Context, name string) (*domain.StreamingPlatform, error)

	ExistByType(ctx context.Context, platformType domain.StreamingPlatformType) (bool, error)

	ExistByName(ctx context.Context, name string) (bool, error)

	List(ctx context.Context, offset, limit int) ([]*domain.StreamingPlatform, int, error)
}
