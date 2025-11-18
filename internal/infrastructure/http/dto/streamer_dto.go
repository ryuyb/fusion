package dto

import "time"

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
	ID                 int64               `json:"id"`
	PlatformType       string              `json:"platform_type"`
	PlatformStreamerID string              `json:"platform_streamer_id"`
	DisplayName        string              `json:"display_name"`
	AvatarURL          string              `json:"avatar_url"`
	RoomURL            string              `json:"room_url"`
	Bio                string              `json:"bio"`
	Tags               []string            `json:"tags"`
	LiveStatus         *LiveStatusResponse `json:"live_status,omitempty"`
	LastLiveSyncedAt   *time.Time          `json:"last_live_synced_at,omitempty"`
	LastProfileSynced  *time.Time          `json:"last_profile_synced_at,omitempty"`
}

type LiveStatusResponse struct {
	IsLive     bool       `json:"is_live"`
	Title      string     `json:"title"`
	GameName   string     `json:"game_name"`
	StartTime  *time.Time `json:"start_time,omitempty"`
	Viewers    int        `json:"viewers"`
	CoverImage string     `json:"cover_image"`
}
