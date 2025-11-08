package notification

import (
	"fmt"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/domain/service"
	"go.uber.org/zap"
)

// NotificationProviderManager manages all notification channel providers
type NotificationProviderManager struct {
	providers map[entity.ChannelType]service.NotificationChannelProvider
	logger    *zap.Logger
}

// NewNotificationProviderManager creates a new NotificationProviderManager
// It receives all providers registered via fx.Group("notification_providers")
func NewNotificationProviderManager(
	providers []service.NotificationChannelProvider,
	logger *zap.Logger,
) *NotificationProviderManager {
	pm := &NotificationProviderManager{
		providers: make(map[entity.ChannelType]service.NotificationChannelProvider),
		logger:    logger,
	}

	// Register all providers by their channel type
	for _, provider := range providers {
		channelType := provider.GetChannelType()
		pm.providers[channelType] = provider
		logger.Info("Registered notification channel provider",
			zap.String("channel_type", string(channelType)))
	}

	return pm
}

// GetProvider returns the provider for a specific channel type
func (pm *NotificationProviderManager) GetProvider(channelType entity.ChannelType) (service.NotificationChannelProvider, error) {
	provider, exists := pm.providers[channelType]
	if !exists {
		return nil, fmt.Errorf("provider not found for channel type: %s", channelType)
	}
	return provider, nil
}

// GetAllProviders returns all registered providers
func (pm *NotificationProviderManager) GetAllProviders() []service.NotificationChannelProvider {
	providers := make([]service.NotificationChannelProvider, 0, len(pm.providers))
	for _, provider := range pm.providers {
		providers = append(providers, provider)
	}
	return providers
}

// HasProvider checks if a provider exists for the given channel type
func (pm *NotificationProviderManager) HasProvider(channelType entity.ChannelType) bool {
	_, exists := pm.providers[channelType]
	return exists
}

// GetSupportedChannels returns a list of all supported channel types
func (pm *NotificationProviderManager) GetSupportedChannels() []entity.ChannelType {
	channels := make([]entity.ChannelType, 0, len(pm.providers))
	for channelType := range pm.providers {
		channels = append(channels, channelType)
	}
	return channels
}
