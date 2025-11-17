package dto

type CreateStreamingPlatformRequest struct {
	Type        string            `json:"type"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	BaseURL     string            `json:"base_url"`
	LogoURL     string            `json:"logo_url"`
	Enabled     bool              `json:"enabled"`
	Priority    int               `json:"priority"`
	Metadata    map[string]string `json:"metadata"`
}

type UpdateStreamingPlatformRequest struct {
	ID          int64             `json:"id"`
	Type        string            `json:"type"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	BaseURL     string            `json:"base_url"`
	LogoURL     string            `json:"logo_url"`
	Enabled     bool              `json:"enabled"`
	Priority    int               `json:"priority"`
	Metadata    map[string]string `json:"metadata"`
}

type StreamingPlatformResponse struct {
	ID          int64             `json:"id"`
	Type        string            `json:"type"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	BaseURL     string            `json:"base_url"`
	LogoURL     string            `json:"logo_url"`
	Enabled     bool              `json:"enabled"`
	Priority    int               `json:"priority"`
	Metadata    map[string]string `json:"metadata"`
}
