package dto

type CreateUserFollowedStreamerRequest struct {
	UserID                 int64   `json:"user_id"`
	StreamerID             int64   `json:"streamer_id"`
	Alias                  string  `json:"alias"`
	Notes                  string  `json:"notes"`
	NotificationsEnabled   bool    `json:"notifications_enabled"`
	NotificationChannelIDs []int64 `json:"notification_channel_ids"`
}

type UpdateUserFollowedStreamerRequest struct {
	ID                     int64   `json:"id"`
	Alias                  string  `json:"alias"`
	Notes                  string  `json:"notes"`
	NotificationsEnabled   bool    `json:"notifications_enabled"`
	NotificationChannelIDs []int64 `json:"notification_channel_ids"`
}

type UserFollowedStreamerResponse struct {
	ID                     int64   `json:"id"`
	UserID                 int64   `json:"user_id"`
	StreamerID             int64   `json:"streamer_id"`
	Alias                  string  `json:"alias"`
	Notes                  string  `json:"notes"`
	NotificationsEnabled   bool    `json:"notifications_enabled"`
	NotificationChannelIDs []int64 `json:"notification_channel_ids"`
}
