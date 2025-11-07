package repository

import (
	"context"
	"time"

	"github.com/ryuyb/fusion/internal/domain/entity"
)

// FollowingFilters represents the filter criteria for querying user followings
type FollowingFilters struct {
	PlatformType        *entity.PlatformType
	NotificationEnabled *bool
	Page                int
	PageSize            int
}

type FollowingRepository interface {
	Create(ctx context.Context, following *entity.UserFollowing) (*entity.UserFollowing, error)

	Update(ctx context.Context, following *entity.UserFollowing) (*entity.UserFollowing, error)

	FindByID(ctx context.Context, id int64) (*entity.UserFollowing, error)

	FindByUserAndStreamer(ctx context.Context, userID, streamerID int64) (*entity.UserFollowing, error)

	// FindByUser returns followings for a user with optional filters
	FindByUser(ctx context.Context, userID int64, filters *FollowingFilters) ([]*entity.UserFollowing, int, error)

	// FindByStreamer returns followers of a streamer, optionally filtered by notification_enabled
	FindByStreamer(ctx context.Context, streamerID int64, notificationEnabled *bool) ([]*entity.UserFollowing, error)

	Delete(ctx context.Context, id int64) error

	// UpdateLastNotifiedAt updates the last notified timestamp
	UpdateLastNotifiedAt(ctx context.Context, id int64, t time.Time) error
}
