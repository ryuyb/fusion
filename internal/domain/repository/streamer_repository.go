package repository

import (
	"context"
	"time"

	"github.com/ryuyb/fusion/internal/domain/entity"
)

type StreamerRepository interface {
	Create(ctx context.Context, streamer *entity.Streamer) (*entity.Streamer, error)

	Update(ctx context.Context, streamer *entity.Streamer) (*entity.Streamer, error)

	FindByID(ctx context.Context, id int64) (*entity.Streamer, error)

	FindByPlatformAndStreamerID(ctx context.Context, platformID int64, platformStreamerID string) (*entity.Streamer, error)

	FindByPlatform(ctx context.Context, platformID int64) ([]*entity.Streamer, error)

	// FindAllWithFollowers returns all streamers that have at least one follower
	FindAllWithFollowers(ctx context.Context) ([]*entity.Streamer, error)

	// UpdateLiveStatus updates the live status of a streamer
	UpdateLiveStatus(ctx context.Context, streamerID int64, isLive bool, lastLiveAt time.Time) error
}
