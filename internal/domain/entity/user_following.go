package entity

import (
	"time"
)

// UserFollowing represents a user's following relationship with a streamer
type UserFollowing struct {
	ID                  int64
	UserID              int64
	StreamerID          int64
	NotificationEnabled bool
	LastNotifiedAt      time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeleteAt            time.Time
}

// CreateFollowing creates a new UserFollowing instance
func CreateFollowing(userID, streamerID int64) *UserFollowing {
	return &UserFollowing{
		UserID:              userID,
		StreamerID:          streamerID,
		NotificationEnabled: true, // Default to enabled
	}
}

// ToggleNotification toggles the notification enabled status
func (uf *UserFollowing) ToggleNotification(enabled bool) *UserFollowing {
	uf.NotificationEnabled = enabled
	return uf
}

// UpdateLastNotifiedAt updates the last notified timestamp
func (uf *UserFollowing) UpdateLastNotifiedAt() *UserFollowing {
	uf.LastNotifiedAt = time.Now()
	return uf
}

// ShouldNotify checks if notification should be sent based on enabled status
func (uf *UserFollowing) ShouldNotify() bool {
	return uf.NotificationEnabled
}
