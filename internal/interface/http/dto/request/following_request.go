package request

// FollowRequest represents the request to follow a streamer
type FollowRequest struct {
	PlatformType       string `json:"platform_type" validate:"required,oneof=douyu huya bilibili"`
	PlatformStreamerID string `json:"platform_streamer_id" validate:"required,min=1,max=100"`
}

// UpdateNotificationRequest represents the request to update notification settings
type UpdateNotificationRequest struct {
	Enabled bool `json:"enabled" validate:"required"`
}

// ListFollowingRequest represents the request to list user's followings
type ListFollowingRequest struct {
	PlatformType        string `json:"platform_type" validate:"omitempty,oneof=douyu huya bilibili"`
	NotificationEnabled *bool  `json:"notification_enabled" validate:"omitempty"`
	Page                int    `json:"page" validate:"omitempty,gte=1"`
	PageSize            int    `json:"page_size" validate:"omitempty,gte=1,lte=100"`
}
