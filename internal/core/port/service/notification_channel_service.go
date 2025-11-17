package service

import (
	"context"

	"github.com/ryuyb/fusion/internal/core/domain"
)

type NotificationChannelService interface {
	Create(ctx context.Context, channel *domain.NotificationChannel) (*domain.NotificationChannel, error)

	Update(ctx context.Context, channel *domain.NotificationChannel) (*domain.NotificationChannel, error)

	Delete(ctx context.Context, id int64) error

	FindById(ctx context.Context, id int64) (*domain.NotificationChannel, error)

	ListByUserId(ctx context.Context, userID int64, page, pageSize int) ([]*domain.NotificationChannel, int, error)
}
