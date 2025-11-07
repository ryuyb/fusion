package streaming

import (
	"context"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/domain/service"
	"github.com/ryuyb/fusion/internal/infrastructure/client"
	"go.uber.org/zap"
)

// HuyaProvider implements StreamingPlatformProvider for Huya Live
type HuyaProvider struct {
	client *client.RestyClient
	logger *zap.Logger
}

// NewHuyaProvider creates a new HuyaProvider instance
func NewHuyaProvider(client *client.RestyClient, logger *zap.Logger) *HuyaProvider {
	return &HuyaProvider{
		client: client,
		logger: logger,
	}
}

func (h HuyaProvider) GetPlatformType() entity.PlatformType {
	return entity.PlatformTypeHuya
}

func (h HuyaProvider) FetchStreamerInfo(ctx context.Context, platformStreamerId string) (*service.StreamerInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (h HuyaProvider) CheckLiveStatus(ctx context.Context, platformStreamerId string) (*service.LiveStatus, error) {
	//TODO implement me
	panic("implement me")
}

func (h HuyaProvider) BatchCheckLiveStatus(ctx context.Context, platformStreamerIds []string) (map[string]*service.LiveStatus, error) {
	//TODO implement me
	panic("implement me")
}

func (h HuyaProvider) ValidateConfiguration(config map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (h HuyaProvider) SearchStreamer(ctx context.Context, keyword string) ([]*service.StreamerInfo, error) {
	//TODO implement me
	panic("implement me")
}
