package douyu

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ryuyb/fusion/internal/core/domain"
	"github.com/ryuyb/fusion/internal/core/port/external"
	"github.com/ryuyb/fusion/internal/infrastructure/client"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"go.uber.org/zap"
	"resty.dev/v3"
)

type Provider struct {
	client *resty.Client
	logger *zap.Logger
}

func NewProvider(logger *zap.Logger) *Provider {
	return &Provider{
		client: client.NewRestyClient(logger),
		logger: logger,
	}
}

func (p *Provider) GetPlatformType() domain.StreamingPlatformType {
	return domain.StreamingPlatformTypeDouyu
}

func (d *Provider) FetchStreamerInfo(ctx context.Context, platformStreamerId string) (*external.StreamerInfo, error) {
	betardResp := &BetardResponse{}
	_, err := d.client.R().
		SetContext(ctx).
		SetPathParam("roomId", platformStreamerId).
		SetResult(betardResp).
		Get("https://www.douyu.com/betard/{roomId}")
	if err != nil {
		return nil, errors2.StreamingPlatformError(string(d.GetPlatformType()), "failed to fetch betard", err)
	}
	return &external.StreamerInfo{
		PlatformStreamerId: platformStreamerId,
		Name:               betardResp.Room.Nickname,
		Avatar:             betardResp.Room.Avatar.Big,
		Description:        betardResp.Room.ShowDetails,
		RoomURL:            fmt.Sprintf("https://www.douyu.com/%s", platformStreamerId),
	}, nil
}

func (d *Provider) CheckLiveStatus(ctx context.Context, platformStreamerId string) (*external.LiveStatus, error) {
	betardResp := &BetardResponse{}
	_, err := d.client.R().
		SetContext(ctx).
		SetPathParam("roomId", platformStreamerId).
		SetResult(betardResp).
		Get("https://www.douyu.com/betard/{roomId}")
	if err != nil {
		return nil, errors2.StreamingPlatformError(string(d.GetPlatformType()), "failed to fetch betard", err)
	}
	viewers, err := strconv.Atoi(betardResp.Room.RoomBizAll.Hot)
	if err != nil {
		return nil, errors2.StreamingPlatformError(string(d.GetPlatformType()), "failed to parse hot", err)
	}
	return &external.LiveStatus{
		IsLive:     betardResp.Room.ShowStatus == 1,
		Title:      betardResp.Room.RoomName,
		GameName:   betardResp.Room.SecondLvlName,
		StartTime:  time.Unix(betardResp.Room.ShowTime, 0),
		Viewers:    viewers,
		CoverImage: betardResp.Room.CoverSrc,
	}, nil
}

func (d *Provider) BatchCheckLiveStatus(ctx context.Context, platformStreamerIds []string) (map[string]*external.LiveStatus, error) {
	results := make(map[string]*external.LiveStatus)

	for _, platformStreamerId := range platformStreamerIds {
		liveStatus, err := d.CheckLiveStatus(ctx, platformStreamerId)
		if err != nil {
			d.logger.Warn("Failed to check live status for room",
				zap.String("room_id", platformStreamerId),
				zap.Error(err))
			continue
		}
		results[platformStreamerId] = liveStatus
	}
	return results, nil
}
