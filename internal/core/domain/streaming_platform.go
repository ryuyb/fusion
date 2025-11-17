package domain

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/ryuyb/fusion/internal/pkg/errors"
)

// StreamingPlatformType defines the type of streaming platform
type StreamingPlatformType string

const (
	StreamingPlatformTypeDouyu    StreamingPlatformType = "douyu"
	StreamingPlatformTypeHuya     StreamingPlatformType = "huya"
	StreamingPlatformTypeBilibili StreamingPlatformType = "bilibili"
)

var supportedStreamingPlatformTypes = map[StreamingPlatformType]struct{}{
	StreamingPlatformTypeDouyu:    {},
	StreamingPlatformTypeHuya:     {},
	StreamingPlatformTypeBilibili: {},
}

// IsValid reports whether the platform type is supported by the domain layer
func (t StreamingPlatformType) IsValid() bool {
	_, ok := supportedStreamingPlatformTypes[t]
	return ok
}

// StreamingPlatform represents a supported live-streaming provider in the system.
type StreamingPlatform struct {
	ID          int64
	Type        StreamingPlatformType
	Name        string
	Description string
	BaseURL     string
	LogoURL     string
	Enabled     bool
	Priority    int
	Metadata    map[string]string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewStreamingPlatform validates and constructs a StreamingPlatform aggregate.
func NewStreamingPlatform(platformType StreamingPlatformType, name, baseURL string) (*StreamingPlatform, error) {
	if !platformType.IsValid() {
		return nil, errors.BadRequest("unsupported streaming platform type")
	}
	if strings.TrimSpace(name) == "" {
		return nil, errors.BadRequest("platform name is required")
	}
	if baseURL == "" {
		return nil, errors.BadRequest("platform base url is required")
	}
	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, errors.BadRequest(fmt.Sprintf("invalid platform base url: %v", err))
	}

	return &StreamingPlatform{
		Type:    platformType,
		Name:    strings.TrimSpace(name),
		BaseURL: baseURL,
		Enabled: true,
	}, nil
}

// UpdateMetadata refreshes human-facing metadata for the streaming platform.
func (p *StreamingPlatform) UpdateMetadata(name, description, baseURL, logoURL string, enabled bool, priority int, metadata map[string]string) error {
	if strings.TrimSpace(name) == "" {
		return errors.BadRequest("platform name is required")
	}
	if baseURL != "" {
		if _, err := url.ParseRequestURI(baseURL); err != nil {
			return errors.BadRequest(fmt.Sprintf("invalid platform base url: %v", err))
		}
	}

	p.Name = strings.TrimSpace(name)
	p.Description = description
	if baseURL != "" {
		p.BaseURL = baseURL
	}
	p.LogoURL = logoURL
	p.Enabled = enabled
	p.Priority = priority
	if metadata != nil {
		copied := make(map[string]string, len(metadata))
		for k, v := range metadata {
			copied[k] = v
		}
		p.Metadata = copied
	} else {
		p.Metadata = nil
	}
	return nil
}
