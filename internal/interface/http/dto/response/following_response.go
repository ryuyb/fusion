package response

import (
	"time"

	"github.com/ryuyb/fusion/internal/domain/entity"
)

// StreamerInfo represents basic streamer information in following response
type StreamerInfo struct {
	ID                 int64  `json:"id"`
	PlatformID         int64  `json:"platform_id"`
	PlatformStreamerID string `json:"platform_streamer_id"`
	Name               string `json:"name"`
	Avatar             string `json:"avatar"`
	Description        string `json:"description"`
	RoomURL            string `json:"room_url"`
	IsLive             bool   `json:"is_live"`
}

// FollowingResponse represents a following relationship in the response
type FollowingResponse struct {
	ID                  int64        `json:"id"`
	UserID              int64        `json:"user_id"`
	Streamer            StreamerInfo `json:"streamer"`
	NotificationEnabled bool         `json:"notification_enabled"`
	LastNotifiedAt      *time.Time   `json:"last_notified_at,omitempty"`
	CreatedAt           time.Time    `json:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at"`
}

// ToFollowingResponse converts an entity.UserFollowing to FollowingResponse
func ToFollowingResponse(following *entity.UserFollowing) *FollowingResponse {
	var lastNotifiedAt *time.Time
	if !following.LastNotifiedAt.IsZero() {
		lastNotifiedAt = &following.LastNotifiedAt
	}

	// Note: Streamer information will be populated by the service layer
	// This is a placeholder for the structure
	streamer := StreamerInfo{}

	return &FollowingResponse{
		ID:                  following.ID,
		UserID:              following.UserID,
		Streamer:            streamer,
		NotificationEnabled: following.NotificationEnabled,
		LastNotifiedAt:      lastNotifiedAt,
		CreatedAt:           following.CreatedAt,
		UpdatedAt:           following.UpdatedAt,
	}
}

// ToFollowingResponseWithStreamer converts a UserFollowing with Streamer to FollowingResponse
func ToFollowingResponseWithStreamer(following *entity.UserFollowing, streamer *entity.Streamer) *FollowingResponse {
	var lastNotifiedAt *time.Time
	if !following.LastNotifiedAt.IsZero() {
		lastNotifiedAt = &following.LastNotifiedAt
	}

	streamerInfo := StreamerInfo{
		ID:                 streamer.ID,
		PlatformID:         streamer.PlatformID,
		PlatformStreamerID: streamer.PlatformStreamerID,
		Name:               streamer.Name,
		Avatar:             streamer.Avatar,
		Description:        streamer.Description,
		RoomURL:            streamer.RoomURL,
		IsLive:             streamer.IsLive,
	}

	return &FollowingResponse{
		ID:                  following.ID,
		UserID:              following.UserID,
		Streamer:            streamerInfo,
		NotificationEnabled: following.NotificationEnabled,
		LastNotifiedAt:      lastNotifiedAt,
		CreatedAt:           following.CreatedAt,
		UpdatedAt:           following.UpdatedAt,
	}
}

// ToFollowingResponseList converts a list of entities to response with pagination
func ToFollowingResponseList(followings []*entity.UserFollowing, total, page, pageSize int) *PaginationResponse[*FollowingResponse] {
	responses := make([]*FollowingResponse, 0, len(followings))
	for _, following := range followings {
		// Streamer will be populated by service layer when needed
		responses = append(responses, ToFollowingResponse(following))
	}

	return NewPaginationResponse[*FollowingResponse](responses, total, page, pageSize)
}
