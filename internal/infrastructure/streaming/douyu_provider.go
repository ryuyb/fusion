package streaming

import (
	"context"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/domain/service"
	"github.com/ryuyb/fusion/internal/infrastructure/client"
	"go.uber.org/zap"
)

// DouyuProvider implements StreamingPlatformProvider for Douyu Live
type DouyuProvider struct {
	client *client.RestyClient
	logger *zap.Logger
}

// NewDouyuProvider creates a new DouyuProvider instance
func NewDouyuProvider(client *client.RestyClient, logger *zap.Logger) *DouyuProvider {
	return &DouyuProvider{
		client: client,
		logger: logger,
	}
}

func (d DouyuProvider) GetPlatformType() entity.PlatformType {
	return entity.PlatformTypeDouyu
}

func (d DouyuProvider) FetchStreamerInfo(ctx context.Context, platformStreamerId string) (*service.StreamerInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (d DouyuProvider) CheckLiveStatus(ctx context.Context, platformStreamerId string) (*service.LiveStatus, error) {
	//TODO implement me
	panic("implement me")
}

func (d DouyuProvider) BatchCheckLiveStatus(ctx context.Context, platformStreamerIds []string) (map[string]*service.LiveStatus, error) {
	//TODO implement me
	panic("implement me")
}

func (d DouyuProvider) ValidateConfiguration(config map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (d DouyuProvider) SearchStreamer(ctx context.Context, keyword string) ([]*service.StreamerInfo, error) {
	//TODO implement me
	panic("implement me")
}
