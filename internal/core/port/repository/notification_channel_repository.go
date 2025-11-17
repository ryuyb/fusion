package repository

import (
	"context"

	"github.com/ryuyb/fusion/internal/core/domain"
)

type NotificationChannelRepository interface {
	Create(ctx context.Context, channel *domain.NotificationChannel) (*domain.NotificationChannel, error)

	Update(ctx context.Context, channel *domain.NotificationChannel) (*domain.NotificationChannel, error)

	Delete(ctx context.Context, id int64) error

	FindById(ctx context.Context, id int64) (*domain.NotificationChannel, error)

	ListByUserId(ctx context.Context, userID int64, offset, limit int) ([]*domain.NotificationChannel, int, error)

	ExistByName(ctx context.Context, userID int64, name string) (bool, error)
}
