package entity

import (
	"time"
)

// NotificationChannel represents a user's notification channel configuration
type NotificationChannel struct {
	ID          int64
	UserID      int64
	ChannelType ChannelType
	Name        string
	Config      map[string]interface{} // JSON configuration for the channel
	IsEnabled   bool
	Priority    int // Lower number = higher priority
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeleteAt    time.Time
}

// ChannelType defines the type of notification channel
type ChannelType string

const (
	ChannelTypeEmail    ChannelType = "email"
	ChannelTypeWebhook  ChannelType = "webhook"
	ChannelTypeTelegram ChannelType = "telegram"
	ChannelTypeDiscord  ChannelType = "discord"
	ChannelTypeFeishu   ChannelType = "feishu"
)

// CreateChannel creates a new NotificationChannel instance
func CreateChannel(userID int64, channelType ChannelType, name string, config map[string]interface{}, priority int) *NotificationChannel {
	return &NotificationChannel{
		UserID:      userID,
		ChannelType: channelType,
		Name:        name,
		Config:      config,
		IsEnabled:   true, // Default to enabled
		Priority:    priority,
	}
}

// Update updates the channel information
func (nc *NotificationChannel) Update(name string, config map[string]interface{}, priority int) *NotificationChannel {
	nc.Name = name
	nc.Config = config
	nc.Priority = priority
	return nc
}

// Toggle toggles the enabled status of the channel
func (nc *NotificationChannel) Toggle(enabled bool) *NotificationChannel {
	nc.IsEnabled = enabled
	return nc
}

// IsActive checks if the channel is enabled
func (nc *NotificationChannel) IsActive() bool {
	return nc.IsEnabled
}
