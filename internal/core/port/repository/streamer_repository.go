package repository

import (
	"context"

	"github.com/ryuyb/fusion/internal/core/domain"
)

type StreamerRepository interface {
	Create(ctx context.Context, streamer *domain.Streamer) (*domain.Streamer, error)

	Update(ctx context.Context, streamer *domain.Streamer) (*domain.Streamer, error)

	Delete(ctx context.Context, id int64) error

	FindById(ctx context.Context, id int64) (*domain.Streamer, error)

	FindByPlatformStreamerId(ctx context.Context, platformType domain.StreamingPlatformType, platformStreamerID string) (*domain.Streamer, error)

	ExistByPlatformStreamerId(ctx context.Context, platformType domain.StreamingPlatformType, platformStreamerID string) (bool, error)

	List(ctx context.Context, offset, limit int) ([]*domain.Streamer, int, error)
}
