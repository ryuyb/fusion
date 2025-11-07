package streaming

import (
	"fmt"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/domain/service"
	"go.uber.org/zap"
)

// StreamingProviderManager manages all streaming platform providers
type StreamingProviderManager struct {
	providers map[entity.PlatformType]service.StreamingPlatformProvider
	logger    *zap.Logger
}

// NewStreamingProviderManager creates a new StreamingProviderManager
// It receives all providers registered via fx.Group("streaming_providers")
func NewStreamingProviderManager(
	providers []service.StreamingPlatformProvider,
	logger *zap.Logger,
) *StreamingProviderManager {
	pm := &StreamingProviderManager{
		providers: make(map[entity.PlatformType]service.StreamingPlatformProvider),
		logger:    logger,
	}

	// Register all providers by their platform type
	for _, provider := range providers {
		platformType := provider.GetPlatformType()
		pm.providers[platformType] = provider
		logger.Info("Registered streaming platform provider",
			zap.String("platform_type", string(platformType)))
	}

	return pm
}

// GetProvider returns the provider for a specific platform type
func (pm *StreamingProviderManager) GetProvider(platformType entity.PlatformType) (service.StreamingPlatformProvider, error) {
	provider, exists := pm.providers[platformType]
	if !exists {
		return nil, fmt.Errorf("provider not found for platform type: %s", platformType)
	}
	return provider, nil
}

// GetAllProviders returns all registered providers
func (pm *StreamingProviderManager) GetAllProviders() []service.StreamingPlatformProvider {
	providers := make([]service.StreamingPlatformProvider, 0, len(pm.providers))
	for _, provider := range pm.providers {
		providers = append(providers, provider)
	}
	return providers
}

// HasProvider checks if a provider exists for the given platform type
func (pm *StreamingProviderManager) HasProvider(platformType entity.PlatformType) bool {
	_, exists := pm.providers[platformType]
	return exists
}

// GetSupportedPlatforms returns a list of all supported platform types
func (pm *StreamingProviderManager) GetSupportedPlatforms() []entity.PlatformType {
	platforms := make([]entity.PlatformType, 0, len(pm.providers))
	for platformType := range pm.providers {
		platforms = append(platforms, platformType)
	}
	return platforms
}
