package domain

import (
	"strings"
	"time"

	"github.com/ryuyb/fusion/internal/pkg/errors"
)

// Streamer captures an individual content creator hosted on a streaming platform.
type Streamer struct {
	ID                 int64
	PlatformType       StreamingPlatformType
	PlatformStreamerID string
	DisplayName        string
	AvatarURL          string
	RoomURL            string
	Bio                string
	Tags               []string
	LastSyncedAt       time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// NewStreamer validates input and constructs a Streamer aggregate.
func NewStreamer(platformType StreamingPlatformType, platformStreamerID, displayName string) (*Streamer, error) {
	if !platformType.IsValid() {
		return nil, errors.BadRequest("unsupported streaming platform type")
	}
	if strings.TrimSpace(platformStreamerID) == "" {
		return nil, errors.BadRequest("platform streamer id is required")
	}
	if strings.TrimSpace(displayName) == "" {
		return nil, errors.BadRequest("streamer display name is required")
	}

	return &Streamer{
		PlatformType:       platformType,
		PlatformStreamerID: strings.TrimSpace(platformStreamerID),
		DisplayName:        strings.TrimSpace(displayName),
	}, nil
}

// UpdateProfile refreshes basic profile data from the upstream platform.
func (s *Streamer) UpdateProfile(displayName, avatarURL, roomURL, bio string, tags []string) error {
	if strings.TrimSpace(displayName) == "" {
		return errors.BadRequest("streamer display name is required")
	}
	s.DisplayName = strings.TrimSpace(displayName)
	s.AvatarURL = avatarURL
	s.RoomURL = roomURL
	s.Bio = bio
	s.Tags = copyStringSlice(tags)
	return nil
}

func copyStringSlice(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	dup := make([]string, len(values))
	copy(dup, values)
	return dup
}
