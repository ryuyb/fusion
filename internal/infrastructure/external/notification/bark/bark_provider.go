package bark

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/core/domain"
	"github.com/ryuyb/fusion/internal/core/port/external"
	"github.com/ryuyb/fusion/internal/infrastructure/client"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"go.uber.org/zap"
	"resty.dev/v3"
)

const (
	DefaultUrl = "https://api.day.app/push"
)

type Provider struct {
	logger *zap.Logger
	client *resty.Client
}

func NewProvider(logger *zap.Logger) *Provider {
	return &Provider{
		logger: logger,
		client: client.NewRestyClient(logger),
	}
}

func (p *Provider) GetChannelType() domain.NotificationChannelType {
	return domain.ChannelTypeBark
}

func (p *Provider) Send(ctx context.Context, channel *domain.NotificationChannel, data *external.NotificationData) error {
	deviceKey, ok := channel.Config["device_key"].(string)
	if !ok {
		err := fmt.Errorf("device key is invalid: %s", channel.Config["device_key"])
		return errors2.Internal(err)
	}

	url, ok := channel.Config["url"].(string)
	if !ok {
		url = DefaultUrl
	}

	req := BarkRequest{
		Title:     data.Title,
		Body:      data.Content,
		DeviceKey: deviceKey,
	}

	response, err := p.client.R().
		SetContext(ctx).
		SetContentType(fiber.MIMEApplicationJSONCharsetUTF8).
		SetBody(req).
		Post(url)

	if err != nil {
		p.logger.Error("Failed to send bark notification", zap.Error(err))
		return errors2.Internal(err)
	}

	if response.StatusCode() < 200 || response.StatusCode() > 299 {
		body := response.String()
		if body == "" {
			body = "no response body"
		}
		p.logger.Error("Failed to send bark notification",
			zap.String("url", url),
			zap.Int("status", response.StatusCode()),
			zap.String("body", body),
		)
		return errors2.Internal(fmt.Errorf("bark returned status code %d", response.StatusCode()))
	}

	return nil
}

func (p *Provider) TestConnection(ctx context.Context, config map[string]any) error {
	data := &external.NotificationData{
		Title:   "Test Notification",
		Content: "This is a test notification from Fusion",
	}
	tempChannel := &domain.NotificationChannel{Config: config}

	err := p.Send(ctx, tempChannel, data)

	if err != nil {
		return errors2.Internal(err)
	}

	return nil
}
