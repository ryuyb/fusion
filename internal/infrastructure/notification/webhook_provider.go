package notification

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/domain/service"
	"github.com/ryuyb/fusion/internal/infrastructure/client"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// WebhookChannelProvider implements NotificationChannelProvider for Webhook notifications
type WebhookChannelProvider struct {
	client *client.RestyClient
	logger *zap.Logger
}

// NewWebhookProvider creates a new WebhookChannelProvider instance
func NewWebhookProvider(client *client.RestyClient, logger *zap.Logger) *WebhookChannelProvider {
	return &WebhookChannelProvider{
		client: client,
		logger: logger,
	}
}

// GetChannelType returns the channel type
func (p *WebhookChannelProvider) GetChannelType() entity.ChannelType {
	return entity.ChannelTypeWebhook
}

// Send sends a notification through webhook
func (p *WebhookChannelProvider) Send(ctx context.Context, channel *entity.NotificationChannel, notification *service.Notification) error {
	// Extract webhook URL from config
	webhookURL, ok := channel.Config["url"].(string)
	if !ok || webhookURL == "" {
		return fmt.Errorf("webhook URL not found in config")
	}

	// Build the payload
	payload := p.buildPayload(notification, channel)

	// Prepare the request
	req := p.client.NewRequest().
		SetContext(ctx).
		SetBody(payload).
		SetHeader("Content-Type", "application/json")

	// Add custom headers from config if present
	if headers, ok := channel.Config["headers"].(map[string]interface{}); ok {
		for k, v := range headers {
			if strVal, ok := v.(string); ok {
				req.SetHeader(k, strVal)
			}
		}
	}

	// Add custom method from config (default to POST)
	method := "POST"
	if methodVal, ok := channel.Config["method"].(string); ok && methodVal != "" {
		method = strings.ToUpper(methodVal)
	}

	var resp *resty.Response
	var err error

	switch method {
	case "POST":
		resp, err = req.Post(webhookURL)
	case "PUT":
		resp, err = req.Put(webhookURL)
	case "PATCH":
		resp, err = req.Patch(webhookURL)
	default:
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		p.logger.Error("Failed to send webhook notification", zap.Error(err), zap.String("url", webhookURL))
		return fmt.Errorf("failed to send webhook: %w", err)
	}

	// Check response status
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		body := resp.String()
		if body == "" {
			body = "no response body"
		}
		p.logger.Error("Webhook returned non-success status",
			zap.String("url", webhookURL),
			zap.Int("status_code", resp.StatusCode()),
			zap.String("body", body),
		)
		return fmt.Errorf("webhook returned status code %d", resp.StatusCode())
	}

	p.logger.Info("Successfully sent webhook notification",
		zap.String("url", webhookURL),
		zap.Int("status_code", resp.StatusCode()),
	)
	return nil
}

// ValidateConfiguration validates the webhook configuration
func (p *WebhookChannelProvider) ValidateConfiguration(config map[string]interface{}) error {
	// Check if URL is present
	webhookURL, ok := config["url"].(string)
	if !ok || webhookURL == "" {
		return fmt.Errorf("webhook URL is required")
	}

	// Validate URL format
	if !strings.HasPrefix(webhookURL, "http://") && !strings.HasPrefix(webhookURL, "https://") {
		return fmt.Errorf("webhook URL must start with http:// or https://")
	}

	// Validate method if present
	if method, ok := config["method"].(string); ok && method != "" {
		validMethods := []string{"POST", "PUT", "PATCH", "GET"}
		method = strings.ToUpper(method)
		isValid := false
		for _, valid := range validMethods {
			if method == valid {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("unsupported HTTP method: %s", method)
		}
	}

	return nil
}

// TestConnection tests the webhook connectivity
func (p *WebhookChannelProvider) TestConnection(ctx context.Context, config map[string]interface{}) error {
	// Validate config first
	if err := p.ValidateConfiguration(config); err != nil {
		return err
	}

	// Extract webhook URL
	webhookURL, _ := config["url"].(string)

	// Create a test notification
	testNotification := &service.Notification{
		Title:        "Test Notification",
		Content:      "This is a test notification from Fusion",
		StreamerName: "Test Streamer",
		RoomURL:      "https://example.com/room/123",
		ExtraData: map[string]interface{}{
			"test": true,
		},
	}

	// Create a temporary channel for testing
	tempChannel := &entity.NotificationChannel{
		Config: config,
	}

	// Send test notification
	err := p.Send(ctx, tempChannel, testNotification)
	if err != nil {
		p.logger.Error("Webhook test connection failed", zap.Error(err), zap.String("url", webhookURL))
		return fmt.Errorf("webhook test failed: %w", err)
	}

	p.logger.Info("Webhook test connection successful", zap.String("url", webhookURL))
	return nil
}

// buildPayload builds the webhook payload from notification data
func (p *WebhookChannelProvider) buildPayload(notification *service.Notification, channel *entity.NotificationChannel) map[string]interface{} {
	// Default payload structure
	payload := map[string]interface{}{
		"title":           notification.Title,
		"content":         notification.Content,
		"streamer_name":   notification.StreamerName,
		"streamer_avatar": notification.StreamerAvatar,
		"room_url":        notification.RoomURL,
		"cover_image":     notification.CoverImage,
		"timestamp":       time.Now().Unix(),
	}

	// Check if custom template is provided
	if template, ok := channel.Config["template"].(string); ok && template != "" {
		// Use custom template - simple string replacement
		payload["message"] = p.applyTemplate(template, notification)
	} else {
		// Default message format
		format := "🎮 %s is now LIVE!\n\nTitle: %s\nWatch now: %s"
		payload["message"] = fmt.Sprintf(format, notification.StreamerName, notification.Content, notification.RoomURL)
	}

	// Include extra data
	if len(notification.ExtraData) > 0 {
		payload["extra_data"] = notification.ExtraData
	}

	// Allow custom fields from config
	if customFields, ok := channel.Config["custom_fields"].(map[string]interface{}); ok {
		for k, v := range customFields {
			payload[k] = v
		}
	}

	return payload
}

// applyTemplate applies a simple string template with notification data
func (p *WebhookChannelProvider) applyTemplate(template string, notification *service.Notification) string {
	// Simple placeholder replacement
	replacements := map[string]string{
		"{title}":         notification.Title,
		"{content}":       notification.Content,
		"{streamer_name}": notification.StreamerName,
		"{room_url}":      notification.RoomURL,
		"{cover_image}":   notification.CoverImage,
	}

	result := template
	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}
