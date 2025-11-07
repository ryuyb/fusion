package entity

import (
	"time"
)

// Platform represents a streaming platform configuration
type Platform struct {
	ID           int64
	Name         string
	PlatformType PlatformType
	Config       map[string]interface{} // JSON configuration for platform API
	Status       PlatformStatus
	PollInterval int // Polling interval in seconds
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeleteAt     time.Time
}

// PlatformType defines the type of streaming platform
type PlatformType string

const (
	PlatformTypeDouyu    PlatformType = "douyu"
	PlatformTypeHuya     PlatformType = "huya"
	PlatformTypeBilibili PlatformType = "bilibili"
)

// PlatformStatus defines the status of a platform
type PlatformStatus string

const (
	PlatformStatusActive   PlatformStatus = "active"
	PlatformStatusInactive PlatformStatus = "inactive"
)

// CreatePlatform creates a new Platform instance
func CreatePlatform(name string, platformType PlatformType, config map[string]interface{}, pollInterval int) *Platform {
	return &Platform{
		Name:         name,
		PlatformType: platformType,
		Config:       config,
		Status:       PlatformStatusActive, // Default to active
		PollInterval: pollInterval,
	}
}

// Update updates the platform information
func (p *Platform) Update(name string, config map[string]interface{}, status PlatformStatus, pollInterval int) *Platform {
	p.Name = name
	p.Config = config
	p.Status = status
	p.PollInterval = pollInterval
	return p
}

// IsActive checks if the platform is active
func (p *Platform) IsActive() bool {
	return p.Status == PlatformStatusActive
}
