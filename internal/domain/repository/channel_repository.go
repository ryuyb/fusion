package repository

import (
	"context"

	"github.com/ryuyb/fusion/internal/domain/entity"
)

type ChannelRepository interface {
	Create(ctx context.Context, channel *entity.NotificationChannel) (*entity.NotificationChannel, error)

	Update(ctx context.Context, channel *entity.NotificationChannel) (*entity.NotificationChannel, error)

	FindByID(ctx context.Context, id int64) (*entity.NotificationChannel, error)

	FindByUser(ctx context.Context, userID int64) ([]*entity.NotificationChannel, error)

	// FindEnabledByUser returns enabled channels for a user, sorted by priority (lower number = higher priority)
	FindEnabledByUser(ctx context.Context, userID int64) ([]*entity.NotificationChannel, error)

	Delete(ctx context.Context, id int64) error
}
