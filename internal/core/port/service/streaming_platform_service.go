package service

import (
	"context"

	"github.com/ryuyb/fusion/internal/core/domain"
)

type StreamingPlatformService interface {
	Create(ctx context.Context, platform *domain.StreamingPlatform) (*domain.StreamingPlatform, error)

	Update(ctx context.Context, platform *domain.StreamingPlatform) (*domain.StreamingPlatform, error)

	Delete(ctx context.Context, id int64) error

	FindById(ctx context.Context, id int64) (*domain.StreamingPlatform, error)

	FindByType(ctx context.Context, platformType domain.StreamingPlatformType) (*domain.StreamingPlatform, error)

	List(ctx context.Context, page, pageSize int) ([]*domain.StreamingPlatform, int, error)
}
