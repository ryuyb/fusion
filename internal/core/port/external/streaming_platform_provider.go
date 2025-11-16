package external

import (
	"context"
	"time"

	"github.com/ryuyb/fusion/internal/core/domain"
)

type StreamingPlatformProvider interface {
	// GetPlatformType returns the type of streaming platform
	GetPlatformType() domain.StreamingPlatformType

	// FetchStreamerInfo fetches detailed information about a streamer from the platform
	FetchStreamerInfo(ctx context.Context, platformStreamerId string) (*StreamerInfo, error)

	// CheckLiveStatus checks the live status of a single streamer
	CheckLiveStatus(ctx context.Context, platformStreamerId string) (*LiveStatus, error)

	// BatchCheckLiveStatus checks live status for multiple streamers in one call (performance optimization)
	BatchCheckLiveStatus(ctx context.Context, platformStreamerIds []string) (map[string]*LiveStatus, error)
}

// StreamerInfo contains basic information about a streamer
type StreamerInfo struct {
	PlatformStreamerId string // Unique streamer ID on the platform
	Name               string // Streamer's display name
	Avatar             string // Avatar image URL
	Description        string // Streamer description/bio
	RoomURL            string // Live room URL
}

// LiveStatus contains the current live status of a streamer
type LiveStatus struct {
	IsLive     bool      // Whether the streamer is currently live
	Title      string    // Live stream title
	GameName   string    // Game/category name
	StartTime  time.Time // Stream start time
	Viewers    int       // Current viewer count
	CoverImage string    // Cover/thumbnail image URL
}
