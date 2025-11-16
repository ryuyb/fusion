package streaming

import (
	"fmt"

	"github.com/ryuyb/fusion/internal/core/domain"
	"github.com/ryuyb/fusion/internal/core/port/external"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"go.uber.org/zap"
)

type StreamingProviderManager struct {
	providers map[domain.StreamingPlatformType]external.StreamingPlatformProvider
	logger    *zap.Logger
}

func NewStreamingProviderManager(providers []external.StreamingPlatformProvider, logger *zap.Logger) *StreamingProviderManager {
	pm := &StreamingProviderManager{
		providers: make(map[domain.StreamingPlatformType]external.StreamingPlatformProvider),
		logger:    logger,
	}

	for _, provider := range providers {
		pm.providers[provider.GetPlatformType()] = provider
		logger.Info("registered streaming platform provider",
			zap.String("platform", string(provider.GetPlatformType())))
	}

	return pm
}

func (pm *StreamingProviderManager) GetProvider(platformype domain.StreamingPlatformType) (external.StreamingPlatformProvider, error) {
	provider, exists := pm.providers[platformype]
	if !exists {
		return nil, errors2.Internal(fmt.Errorf("provider not found for platform type: %s", platformype))
	}
	return provider, nil
}

func (pm *StreamingProviderManager) GetAllProviders() []external.StreamingPlatformProvider {
	providers := make([]external.StreamingPlatformProvider, 0, len(pm.providers))
	for _, provider := range pm.providers {
		providers = append(providers, provider)
	}
	return providers
}

func (pm *StreamingProviderManager) HasProvider(platformType domain.StreamingPlatformType) bool {
	_, exists := pm.providers[platformType]
	return exists
}

func (pm StreamingProviderManager) GetSupportedPlatforms() []domain.StreamingPlatformType {
	platformTypes := make([]domain.StreamingPlatformType, 0, len(pm.providers))
	for platformType := range pm.providers {
		platformTypes = append(platformTypes, platformType)
	}
	return platformTypes
}
