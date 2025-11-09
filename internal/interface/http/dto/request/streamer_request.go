package request

// SearchStreamerRequest represents the request to search for streamers
type SearchStreamerRequest struct {
	PlatformType string `json:"platform_type" validate:"required,oneof=douyu huya bilibili"`
	Keyword      string `json:"keyword" validate:"required,min=1,max=100"`
	Page         int    `json:"page" validate:"omitempty,gte=1"`
	PageSize     int    `json:"page_size" validate:"omitempty,gte=1,lte=100"`
}
