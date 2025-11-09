package response

import (
	"time"

	"github.com/ryuyb/fusion/internal/domain/entity"
)

// PlatformResponse represents a platform in the response
type PlatformResponse struct {
	ID           int64                  `json:"id"`
	Name         string                 `json:"name"`
	PlatformType string                 `json:"platform_type"`
	Config       map[string]interface{} `json:"config"`
	Status       string                 `json:"status"`
	PollInterval int                    `json:"poll_interval"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// ToPlatformResponse converts an entity.Platform to PlatformResponse
func ToPlatformResponse(platform *entity.Platform) *PlatformResponse {
	return &PlatformResponse{
		ID:           platform.ID,
		Name:         platform.Name,
		PlatformType: string(platform.PlatformType),
		Config:       platform.Config,
		Status:       string(platform.Status),
		PollInterval: platform.PollInterval,
		CreatedAt:    platform.CreatedAt,
		UpdatedAt:    platform.UpdatedAt,
	}
}

// ToPlatformResponseList converts a list of entities to response
func ToPlatformResponseList(platforms []*entity.Platform) []*PlatformResponse {
	responses := make([]*PlatformResponse, 0, len(platforms))
	for _, platform := range platforms {
		responses = append(responses, ToPlatformResponse(platform))
	}
	return responses
}
