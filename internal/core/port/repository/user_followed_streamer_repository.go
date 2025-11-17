package repository

import (
	"context"

	"github.com/ryuyb/fusion/internal/core/domain"
)

type UserFollowedStreamerRepository interface {
	Create(ctx context.Context, follow *domain.UserFollowedStreamer) (*domain.UserFollowedStreamer, error)

	Update(ctx context.Context, follow *domain.UserFollowedStreamer) (*domain.UserFollowedStreamer, error)

	Delete(ctx context.Context, id int64) error

	FindById(ctx context.Context, id int64) (*domain.UserFollowedStreamer, error)

	FindByUserAndStreamer(ctx context.Context, userID, streamerID int64) (*domain.UserFollowedStreamer, error)

	ExistByUserAndStreamer(ctx context.Context, userID, streamerID int64) (bool, error)

	ListByUserId(ctx context.Context, userID int64, offset, limit int) ([]*domain.UserFollowedStreamer, int, error)

	ListByStreamerId(ctx context.Context, streamerID int64, offset, limit int) ([]*domain.UserFollowedStreamer, int, error)
}
