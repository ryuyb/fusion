package notification

import (
	"context"
	"testing"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/domain/service"
	"go.uber.org/zap"
)

// mockProvider is a simple mock implementation of NotificationChannelProvider for testing
type mockProvider struct {
	channelType entity.ChannelType
}

func (m *mockProvider) GetChannelType() entity.ChannelType {
	return m.channelType
}

func (m *mockProvider) Send(_ context.Context, _ *entity.NotificationChannel, _ *service.Notification) error {
	return nil
}

func (m *mockProvider) ValidateConfiguration(_ map[string]interface{}) error {
	return nil
}

func (m *mockProvider) TestConnection(_ context.Context, _ map[string]interface{}) error {
	return nil
}

func TestNewNotificationProviderManager(t *testing.T) {
	logger := zap.NewNop()

	// Create mock providers
	providers := []service.NotificationChannelProvider{
		&mockProvider{channelType: entity.ChannelTypeEmail},
		&mockProvider{channelType: entity.ChannelTypeWebhook},
		&mockProvider{channelType: entity.ChannelTypeTelegram},
	}

	// Create manager
	manager := NewNotificationProviderManager(providers, logger)

	// Test GetAllProviders
	allProviders := manager.GetAllProviders()
	if len(allProviders) != 3 {
		t.Errorf("Expected 3 providers, got %d", len(allProviders))
	}

	// Test GetProvider
	emailProvider, err := manager.GetProvider(entity.ChannelTypeEmail)
	if err != nil {
		t.Errorf("Expected to get email provider, got error: %v", err)
	}
	if emailProvider.GetChannelType() != entity.ChannelTypeEmail {
		t.Errorf("Expected email provider type, got %s", emailProvider.GetChannelType())
	}

	// Test HasProvider
	if !manager.HasProvider(entity.ChannelTypeWebhook) {
		t.Errorf("Expected to find webhook provider")
	}

	if manager.HasProvider(entity.ChannelTypeDiscord) {
		t.Errorf("Expected not to find discord provider")
	}

	// Test GetSupportedChannels
	supportedChannels := manager.GetSupportedChannels()
	if len(supportedChannels) != 3 {
		t.Errorf("Expected 3 supported channels, got %d", len(supportedChannels))
	}

	// Test GetProvider not found
	_, err = manager.GetProvider(entity.ChannelTypeDiscord)
	if err == nil {
		t.Errorf("Expected error when getting non-existent provider")
	}
}

func TestNotificationProviderManager_GetProvider(t *testing.T) {
	logger := zap.NewNop()

	providers := []service.NotificationChannelProvider{
		&mockProvider{channelType: entity.ChannelTypeEmail},
	}

	manager := NewNotificationProviderManager(providers, logger)

	// Test successful get
	provider, err := manager.GetProvider(entity.ChannelTypeEmail)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if provider == nil {
		t.Errorf("Expected provider, got nil")
	}

	// Test non-existent provider
	_, err = manager.GetProvider(entity.ChannelTypeFeishu)
	if err == nil {
		t.Errorf("Expected error for non-existent provider")
	}
}
