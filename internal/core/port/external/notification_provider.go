package external

import (
	"context"

	"github.com/ryuyb/fusion/internal/core/domain"
)

type NotificationProvider interface {
	GetChannelType() domain.NotificationChannelType

	Send(ctx context.Context, channel *domain.NotificationChannel, data *NotificationData) error

	TestConnection(ctx context.Context, config map[string]any) error
}

type NotificationData struct {
	Title   string
	Content string
}
