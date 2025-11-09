package response

import (
	"time"

	"github.com/ryuyb/fusion/internal/domain/entity"
)

// StreamerResponse represents a streamer in the response
type StreamerResponse struct {
	ID                 int64      `json:"id"`
	PlatformID         int64      `json:"platform_id"`
	PlatformStreamerID string     `json:"platform_streamer_id"`
	Name               string     `json:"name"`
	Avatar             string     `json:"avatar"`
	Description        string     `json:"description"`
	RoomURL            string     `json:"room_url"`
	IsLive             bool       `json:"is_live"`
	LastLiveAt         *time.Time `json:"last_live_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// LiveStatusResponse represents the live status of a streamer
type LiveStatusResponse struct {
	IsLive     bool       `json:"is_live"`
	Title      string     `json:"title"`
	GameName   string     `json:"game_name"`
	StartTime  *time.Time `json:"start_time,omitempty"`
	Viewers    int        `json:"viewers"`
	CoverImage string     `json:"cover_image"`
}

// ToStreamerResponse converts an entity.Streamer to StreamerResponse
func ToStreamerResponse(streamer *entity.Streamer) *StreamerResponse {
	var lastLiveAt *time.Time
	if !streamer.LastLiveAt.IsZero() {
		lastLiveAt = &streamer.LastLiveAt
	}

	return &StreamerResponse{
		ID:                 streamer.ID,
		PlatformID:         streamer.PlatformID,
		PlatformStreamerID: streamer.PlatformStreamerID,
		Name:               streamer.Name,
		Avatar:             streamer.Avatar,
		Description:        streamer.Description,
		RoomURL:            streamer.RoomURL,
		IsLive:             streamer.IsLive,
		LastLiveAt:         lastLiveAt,
		CreatedAt:          streamer.CreatedAt,
		UpdatedAt:          streamer.UpdatedAt,
	}
}

// ToStreamerResponseList converts a list of entities to response
func ToStreamerResponseList(streamers []*entity.Streamer) []*StreamerResponse {
	responses := make([]*StreamerResponse, 0, len(streamers))
	for _, streamer := range streamers {
		responses = append(responses, ToStreamerResponse(streamer))
	}
	return responses
}
