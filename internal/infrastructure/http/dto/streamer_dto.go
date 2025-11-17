package dto

type CreateStreamerRequest struct {
	PlatformType       string   `json:"platform_type"`
	PlatformStreamerID string   `json:"platform_streamer_id"`
	DisplayName        string   `json:"display_name"`
	AvatarURL          string   `json:"avatar_url"`
	RoomURL            string   `json:"room_url"`
	Bio                string   `json:"bio"`
	Tags               []string `json:"tags"`
}

type UpdateStreamerRequest struct {
	ID                 int64    `json:"id"`
	PlatformType       string   `json:"platform_type"`
	PlatformStreamerID string   `json:"platform_streamer_id"`
	DisplayName        string   `json:"display_name"`
	AvatarURL          string   `json:"avatar_url"`
	RoomURL            string   `json:"room_url"`
	Bio                string   `json:"bio"`
	Tags               []string `json:"tags"`
}

type StreamerResponse struct {
	ID                 int64    `json:"id"`
	PlatformType       string   `json:"platform_type"`
	PlatformStreamerID string   `json:"platform_streamer_id"`
	DisplayName        string   `json:"display_name"`
	AvatarURL          string   `json:"avatar_url"`
	RoomURL            string   `json:"room_url"`
	Bio                string   `json:"bio"`
	Tags               []string `json:"tags"`
}
