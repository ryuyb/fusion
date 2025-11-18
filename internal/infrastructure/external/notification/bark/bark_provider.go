package bark

import (
	"context"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/core/domain"
	"github.com/ryuyb/fusion/internal/core/port/external"
	"github.com/ryuyb/fusion/internal/infrastructure/client"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"go.uber.org/zap"
	"resty.dev/v3"
)

const (
	DefaultURL = "https://api.day.app/push"
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
	endpointURL := resolveEndpoint(channel.Config)

	req, err := buildRequest(channel.Config, data)
	if err != nil {
		return err
	}

	response, err := p.client.R().
		SetContext(ctx).
		SetContentType(fiber.MIMEApplicationJSONCharsetUTF8).
		SetBody(req).
		SetResult(&barkResponse{}).
		Post(endpointURL)

	if err != nil {
		p.logger.Error("Failed to send bark notification", zap.Error(err))
		return errors2.Internal(err)
	}

	if err := p.checkHTTPStatus(response, endpointURL); err != nil {
		return err
	}

	if err := parseBarkResult(response); err != nil {
		return err
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
		return err
	}

	return nil
}

func buildRequest(cfg map[string]any, data *external.NotificationData) (BarkRequest, error) {
	deviceKey, ok := cfg["device_key"].(string)
	if !ok || strings.TrimSpace(deviceKey) == "" {
		return BarkRequest{}, errors2.BadRequest("device key is invalid").
			WithDetail("device_key", cfg["device_key"])
	}

	req := BarkRequest{
		Title:     data.Title,
		Body:      data.Content,
		DeviceKey: deviceKey,
	}

	setOptionalStrings(cfg, "subtitle", &req.Subtitle)
	setOptionalStrings(cfg, "call", &req.Call)
	setOptionalStrings(cfg, "autoCopy", &req.AutoCopy)
	setOptionalStrings(cfg, "copy", &req.Copy)
	setOptionalStrings(cfg, "sound", &req.Sound)
	setOptionalStrings(cfg, "group", &req.Group)
	setOptionalStrings(cfg, "action", &req.Action)

	if level, ok := cfg["level"].(string); ok {
		level = strings.TrimSpace(level)
		allowed := map[string]struct{}{
			"critical":      {},
			"active":        {},
			"timeSensitive": {},
			"passive":       {},
		}
		if _, exist := allowed[level]; !exist {
			return BarkRequest{}, errors2.BadRequest("level is invalid").WithDetail("level", level)
		}
		req.Level = level
	}

	if volume, ok := toInt(cfg["volume"]); ok {
		if volume < 0 || volume > 10 {
			return BarkRequest{}, errors2.BadRequest("volume must be between 0 and 10").
				WithDetail("volume", volume)
		}
		req.Volume = volume
	}

	if badge, ok := toInt(cfg["badge"]); ok {
		if badge < 0 {
			return BarkRequest{}, errors2.BadRequest("badge must be non-negative").
				WithDetail("badge", badge)
		}
		req.Badge = badge
	}

	if icon, ok := cfg["icon"].(string); ok && strings.TrimSpace(icon) != "" {
		if !isValidURL(icon) {
			return BarkRequest{}, errors2.BadRequest("icon must be a valid URL").
				WithDetail("icon", icon)
		}
		req.Icon = icon
	}

	if link, ok := cfg["link"].(string); ok && strings.TrimSpace(link) != "" {
		if !isValidURL(link) {
			return BarkRequest{}, errors2.BadRequest("link must be a valid URL").
				WithDetail("link", link)
		}
		req.Url = link
	}
	if req.Url == "" {
		if link, ok := cfg["open_url"].(string); ok && strings.TrimSpace(link) != "" {
			if !isValidURL(link) {
				return BarkRequest{}, errors2.BadRequest("open_url must be a valid URL").
					WithDetail("open_url", link)
			}
			req.Url = link
		}
	}

	return req, nil
}

func resolveEndpoint(cfg map[string]any) string {
	if endpoint, ok := cfg["url"].(string); ok && strings.TrimSpace(endpoint) != "" {
		return endpoint
	}
	return DefaultURL
}

func setOptionalStrings(cfg map[string]any, key string, target *string) {
	if v, ok := cfg[key].(string); ok && strings.TrimSpace(v) != "" {
		*target = v
	}
}

func (p *Provider) checkHTTPStatus(response *resty.Response, endpointURL string) error {
	if response.StatusCode() >= 200 && response.StatusCode() <= 299 {
		return nil
	}

	body := response.String()
	if body == "" {
		body = "no response body"
	}
	p.logger.Error("Failed to send bark notification",
		zap.String("url", endpointURL),
		zap.Int("status", response.StatusCode()),
		zap.String("body", body),
	)
	return errors2.BadRequest("bark returned non-success status").
		WithDetails(map[string]any{
			"status": response.StatusCode(),
			"body":   body,
		})
}

func parseBarkResult(response *resty.Response) error {
	if res, ok := response.Result().(*barkResponse); ok && res != nil {
		if res.Code != 0 && res.Code != 200 {
			msg := res.Message
			if msg == "" {
				msg = res.Msg
			}
			return errors2.BadRequest("bark returned error").
				WithDetails(map[string]any{
					"code":    res.Code,
					"message": msg,
				})
		}
	}
	return nil
}

func toInt(v any) (int, bool) {
	switch val := v.(type) {
	case int:
		return val, true
	case int32:
		return int(val), true
	case int64:
		return int(val), true
	case float32:
		return int(val), true
	case float64:
		return int(val), true
	default:
		return 0, false
	}
}

func isValidURL(raw string) bool {
	parsed, err := url.Parse(raw)
	if err != nil {
		return false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}
	return parsed.Host != ""
}
