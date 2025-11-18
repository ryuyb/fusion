package douyu

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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
	resp, err := d.client.R().
		SetContext(ctx).
		SetPathParam("roomId", platformStreamerId).
		SetResult(betardResp).
		Get("https://www.douyu.com/betard/{roomId}")
	if err != nil {
		return nil, errors2.StreamingPlatformError(string(d.GetPlatformType()), "failed to fetch betard", err)
	}

	if isPromptHTML(resp) {
		if msg := extractPromptMessage(resp.String()); msg != "" {
			return nil, errors2.StreamingPlatformError(string(d.GetPlatformType()), msg, nil)
		}
		d.logger.Error("douyu betard returned prompt page", zap.String("body", resp.String()))
		return nil, errors2.StreamingPlatformError(string(d.GetPlatformType()), "douyu returned prompt page", nil)
	}

	if betardResp.Room.Nickname == "" {
		d.logger.Error("fetch streamer info failed", zap.String("platformStreamerId", platformStreamerId))
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

func isPromptHTML(resp *resty.Response) bool {
	if resp == nil {
		return false
	}
	contentType := resp.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return false
	}
	body := resp.String()
	return strings.Contains(body, "<title>提示信息 -斗鱼</title>")
}

func extractPromptMessage(body string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(doc.Find(".error > span > p").First().Text())
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
