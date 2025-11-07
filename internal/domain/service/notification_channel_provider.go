package service

import (
	"context"

	"github.com/ryuyb/fusion/internal/domain/entity"
)

// NotificationChannelProvider defines the interface for notification channel integrations
// Each notification channel (Email, Webhook, Telegram, Discord, Feishu, etc.) must implement this interface
type NotificationChannelProvider interface {
	// GetChannelType returns the type of notification channel
	GetChannelType() entity.ChannelType

	// Send sends a notification through this channel
	Send(ctx context.Context, channel *entity.NotificationChannel, notification *Notification) error

	// ValidateConfiguration validates the channel configuration
	ValidateConfiguration(config map[string]interface{}) error

	// TestConnection tests the channel connectivity
	TestConnection(ctx context.Context, config map[string]interface{}) error
}

// Notification contains the notification data to be sent
type Notification struct {
	Title          string                 // Notification title
	Content        string                 // Notification content/message
	StreamerName   string                 // Streamer's name
	StreamerAvatar string                 // Streamer's avatar URL
	RoomURL        string                 // Live room URL
	CoverImage     string                 // Stream cover/thumbnail URL
	ExtraData      map[string]interface{} // Additional data for specific channels
}
