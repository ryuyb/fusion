package request

// CreateChannelRequest represents the request to create a notification channel
type CreateChannelRequest struct {
	ChannelType string                 `json:"channel_type" validate:"required,oneof=email webhook telegram discord feishu"`
	Name        string                 `json:"name" validate:"required,min=2,max=50"`
	Config      map[string]interface{} `json:"config" validate:"required"`
	Priority    int                    `json:"priority" validate:"omitempty,gte=0,lte=100"`
}

// UpdateChannelRequest represents the request to update a notification channel
type UpdateChannelRequest struct {
	Name      string                 `json:"name" validate:"required,min=2,max=50"`
	Config    map[string]interface{} `json:"config" validate:"required"`
	Priority  int                    `json:"priority" validate:"omitempty,gte=0,lte=100"`
	IsEnabled *bool                  `json:"is_enabled" validate:"omitempty"`
}

// ToggleChannelRequest represents the request to toggle a channel enabled status
type ToggleChannelRequest struct {
	IsEnabled bool `json:"is_enabled" validate:"required"`
}

// ListChannelRequest represents the request to list notification channels
type ListChannelRequest struct {
	ChannelType string `json:"channel_type" validate:"omitempty,oneof=email webhook telegram discord feishu"`
	IsEnabled   *bool  `json:"is_enabled" validate:"omitempty"`
	Page        int    `json:"page" validate:"omitempty,gte=1"`
	PageSize    int    `json:"page_size" validate:"omitempty,gte=1,lte=100"`
}

// TestChannelRequest represents the request to test a channel
type TestChannelRequest struct {
	TestMessage string `json:"test_message" validate:"omitempty,max=200"`
}
