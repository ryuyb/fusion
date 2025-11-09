package response

import (
	"time"

	"github.com/ryuyb/fusion/internal/domain/entity"
)

// ConfigResponse represents sanitized channel configuration in response
type ConfigResponse map[string]interface{}

// ChannelResponse represents a notification channel in the response
type ChannelResponse struct {
	ID          int64          `json:"id"`
	UserID      int64          `json:"user_id"`
	ChannelType string         `json:"channel_type"`
	Name        string         `json:"name"`
	Config      ConfigResponse `json:"config"`
	IsEnabled   bool           `json:"is_enabled"`
	Priority    int            `json:"priority"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// ToChannelResponse converts an entity.NotificationChannel to ChannelResponse
func ToChannelResponse(channel *entity.NotificationChannel) *ChannelResponse {
	return &ChannelResponse{
		ID:          channel.ID,
		UserID:      channel.UserID,
		ChannelType: string(channel.ChannelType),
		Name:        channel.Name,
		Config:      channel.Config,
		IsEnabled:   channel.IsEnabled,
		Priority:    channel.Priority,
		CreatedAt:   channel.CreatedAt,
		UpdatedAt:   channel.UpdatedAt,
	}
}

// ToChannelResponseList converts a list of entities to response with pagination
func ToChannelResponseList(channels []*entity.NotificationChannel, total, page, pageSize int) *PaginationResponse[*ChannelResponse] {
	responses := make([]*ChannelResponse, 0, len(channels))
	for _, channel := range channels {
		responses = append(responses, ToChannelResponse(channel))
	}

	return NewPaginationResponse[*ChannelResponse](responses, total, page, pageSize)
}
