package request

// CreatePlatformRequest represents the request to create a new platform
type CreatePlatformRequest struct {
	Name         string                 `json:"name" validate:"required,min=2,max=50"`
	PlatformType string                 `json:"platform_type" validate:"required,oneof=douyu huya bilibili"`
	Config       map[string]interface{} `json:"config" validate:"required"`
	PollInterval int                    `json:"poll_interval" validate:"required,min=10,max=3600"`
}

// UpdatePlatformRequest represents the request to update a platform
type UpdatePlatformRequest struct {
	Name         string                 `json:"name" validate:"required,min=2,max=50"`
	Config       map[string]interface{} `json:"config" validate:"required"`
	Status       string                 `json:"status" validate:"required,oneof=active inactive"`
	PollInterval int                    `json:"poll_interval" validate:"required,min=10,max=3600"`
}
