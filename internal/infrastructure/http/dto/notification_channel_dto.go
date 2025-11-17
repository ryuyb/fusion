package dto

type CreateNotificationChannelRequest struct {
	UserID      int64          `json:"user_id"`
	ChannelType string         `json:"channel_type"`
	Name        string         `json:"name"`
	Config      map[string]any `json:"config"`
	Enable      bool           `json:"enable"`
	Priority    int            `json:"priority"`
}

type UpdateNotificationChannelRequest struct {
	ID          int64          `json:"id"`
	UserID      int64          `json:"user_id"`
	ChannelType string         `json:"channel_type"`
	Name        string         `json:"name"`
	Config      map[string]any `json:"config"`
	Enable      bool           `json:"enable"`
	Priority    int            `json:"priority"`
}

type NotificationChannelResponse struct {
	ID          int64          `json:"id"`
	UserID      int64          `json:"user_id"`
	ChannelType string         `json:"channel_type"`
	Name        string         `json:"name"`
	Config      map[string]any `json:"config"`
	Enable      bool           `json:"enable"`
	Priority    int            `json:"priority"`
}
