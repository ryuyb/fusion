package domain

import "time"

// NotificationChannelType defines the type of notification channel
type NotificationChannelType string

const (
	ChannelTypeEmail    NotificationChannelType = "email"
	ChannelTypeWebhook  NotificationChannelType = "webhook"
	ChannelTypeTelegram NotificationChannelType = "telegram"
	ChannelTypeDiscord  NotificationChannelType = "discord"
	ChannelTypeFeishu   NotificationChannelType = "feishu"
	ChannelTypeBark     NotificationChannelType = "bark"
)

type NotificationChannel struct {
	ID          int64
	UserID      int64
	ChannelType NotificationChannelType
	Name        string
	Config      map[string]any
	Enable      bool
	Priority    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
