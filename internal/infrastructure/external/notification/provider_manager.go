package notification

import (
	"fmt"

	"github.com/ryuyb/fusion/internal/core/domain"
	"github.com/ryuyb/fusion/internal/core/port/external"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"go.uber.org/zap"
)

// NotificationProviderManager holds notification providers keyed by channel type.
type NotificationProviderManager struct {
	providers map[domain.NotificationChannelType]external.NotificationProvider
	logger    *zap.Logger
}

// NewNotificationProviderManager builds a manager with a lookup table for quick provider resolution.
func NewNotificationProviderManager(providers []external.NotificationProvider, logger *zap.Logger) *NotificationProviderManager {
	pm := &NotificationProviderManager{
		providers: make(map[domain.NotificationChannelType]external.NotificationProvider),
		logger:    logger,
	}

	for _, provider := range providers {
		pm.providers[provider.GetChannelType()] = provider
		logger.Info("registered streaming channel provider",
			zap.String("channel", string(provider.GetChannelType())))
	}

	return pm
}

func (pm *NotificationProviderManager) GetProvider(channelType domain.NotificationChannelType) (external.NotificationProvider, error) {
	provider, exists := pm.providers[channelType]
	if !exists {
		return nil, errors2.Internal(fmt.Errorf("provider not found for channel type: %s", channelType))
	}
	return provider, nil
}

// GetAllProviders returns all registered providers; useful for broadcast-style operations.
func (pm *NotificationProviderManager) GetAllProviders() []external.NotificationProvider {
	providers := make([]external.NotificationProvider, 0, len(pm.providers))
	for _, provider := range pm.providers {
		providers = append(providers, provider)
	}
	return providers
}

func (pm *NotificationProviderManager) HasProvider(channelType domain.NotificationChannelType) bool {
	_, exists := pm.providers[channelType]
	return exists
}

// GetSupportedChannels lists channel types for configuration or capability checks.
func (pm *NotificationProviderManager) GetSupportedChannels() []domain.NotificationChannelType {
	channelTypes := make([]domain.NotificationChannelType, 0, len(pm.providers))
	for channelType := range pm.providers {
		channelTypes = append(channelTypes, channelType)
	}
	return channelTypes
}
