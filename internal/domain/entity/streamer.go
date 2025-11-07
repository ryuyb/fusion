package entity

import (
	"time"
)

// Streamer represents a streamer from a streaming platform
type Streamer struct {
	ID                 int64
	PlatformID         int64
	PlatformStreamerID string // Streamer ID on the platform
	Name               string
	Avatar             string
	Description        string
	RoomURL            string
	LastCheckedAt      time.Time
	IsLive             bool
	LastLiveAt         time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeleteAt           time.Time
}

// CreateStreamer creates a new Streamer instance
func CreateStreamer(platformID int64, platformStreamerID, name, avatar, description, roomURL string) *Streamer {
	return &Streamer{
		PlatformID:         platformID,
		PlatformStreamerID: platformStreamerID,
		Name:               name,
		Avatar:             avatar,
		Description:        description,
		RoomURL:            roomURL,
		IsLive:             false,
	}
}

// Update updates the streamer information
func (s *Streamer) Update(name, avatar, description, roomURL string) *Streamer {
	s.Name = name
	s.Avatar = avatar
	s.Description = description
	s.RoomURL = roomURL
	return s
}

// UpdateLiveStatus updates the live status of the streamer
func (s *Streamer) UpdateLiveStatus(isLive bool) *Streamer {
	s.IsLive = isLive
	s.LastCheckedAt = time.Now()
	if isLive {
		s.LastLiveAt = time.Now()
	}
	return s
}

// MarkAsChecked marks the streamer as checked
func (s *Streamer) MarkAsChecked() *Streamer {
	s.LastCheckedAt = time.Now()
	return s
}
