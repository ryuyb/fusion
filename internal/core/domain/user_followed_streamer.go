package domain

import (
	"strings"
	"time"

	"github.com/ryuyb/fusion/internal/pkg/errors"
)

// UserFollowedStreamer links a user with a streamer they follow together with notification preferences.
type UserFollowedStreamer struct {
	ID                     int64
	UserID                 int64
	StreamerID             int64
	Alias                  string
	Notes                  string
	NotificationsEnabled   bool
	NotificationChannelIDs []int64
	LastNotificationSentAt *time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// NewUserFollowedStreamer validates inputs and builds a follow relationship with default notification settings.
func NewUserFollowedStreamer(userID, streamerID int64, alias, notes string, channelIDs []int64) (*UserFollowedStreamer, error) {
	if userID <= 0 {
		return nil, errors.BadRequest("user id must be greater than zero")
	}
	if streamerID <= 0 {
		return nil, errors.BadRequest("streamer id must be greater than zero")
	}

	sanitizedAlias := strings.TrimSpace(alias)
	sanitizedNotes := strings.TrimSpace(notes)
	normalizedChannelIDs, err := normalizeNotificationChannelIDs(channelIDs)
	if err != nil {
		return nil, err
	}

	return &UserFollowedStreamer{
		UserID:                 userID,
		StreamerID:             streamerID,
		Alias:                  sanitizedAlias,
		Notes:                  sanitizedNotes,
		NotificationsEnabled:   true,
		NotificationChannelIDs: normalizedChannelIDs,
	}, nil
}

// UpdatePreferences refreshes notification settings and user facing metadata.
func (f *UserFollowedStreamer) UpdatePreferences(alias, notes string, notificationsEnabled bool, channelIDs []int64) error {
	sanitizedAlias := strings.TrimSpace(alias)
	sanitizedNotes := strings.TrimSpace(notes)
	normalizedChannelIDs, err := normalizeNotificationChannelIDs(channelIDs)
	if err != nil {
		return err
	}

	f.Alias = sanitizedAlias
	f.Notes = sanitizedNotes
	f.NotificationsEnabled = notificationsEnabled
	f.NotificationChannelIDs = normalizedChannelIDs
	return nil
}

func normalizeNotificationChannelIDs(channelIDs []int64) ([]int64, error) {
	if len(channelIDs) == 0 {
		return nil, nil
	}

	unique := make(map[int64]struct{}, len(channelIDs))
	result := make([]int64, 0, len(channelIDs))
	for _, id := range channelIDs {
		if id <= 0 {
			return nil, errors.BadRequest("notification channel id must be greater than zero")
		}
		if _, exists := unique[id]; exists {
			continue
		}
		unique[id] = struct{}{}
		result = append(result, id)
	}
	return result, nil
}
