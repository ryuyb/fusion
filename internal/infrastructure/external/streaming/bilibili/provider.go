package bilibili

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
	return domain.StreamingPlatformTypeBilibili
}

func (p *Provider) FetchStreamerInfo(ctx context.Context, platformStreamerId string) (*external.StreamerInfo, error) {
	roomID, err := strconv.ParseInt(platformStreamerId, 10, 64)
	if err != nil {
		return nil, errors2.StreamingPlatformError(string(p.GetPlatformType()), "invalid room id", err)
	}

	// Get room basic info
	var roomResp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			UID         int64  `json:"uid"`
			RoomID      int64  `json:"room_id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			UserCover   string `json:"user_cover"`
			LiveStatus  int    `json:"live_status"` // 0: offline, 1: live, 2: replay
		} `json:"data"`
	}

	resp, err := p.client.R().
		SetContext(ctx).
		SetQueryParam("room_id", platformStreamerId).
		SetResult(&roomResp).
		Get("https://api.live.bilibili.com/room/v1/Room/get_info")

	if err != nil {
		p.logger.Error("Failed to fetch Bilibili room info",
			zap.String("room_id", platformStreamerId),
			zap.Error(err))

		return nil, errors2.StreamingPlatformError(string(p.GetPlatformType()), "failed to fetch room info", err)
	}

	if resp.IsError() {
		p.logger.Error("Bilibili API returned error",
			zap.String("room_id", platformStreamerId),
			zap.Int("status_code", resp.StatusCode()))
		err = fmt.Errorf("API returned error status: %d", resp.StatusCode())
		return nil, errors2.StreamingPlatformError(string(p.GetPlatformType()), "", err)
	}

	if roomResp.Code != 0 {
		err := fmt.Errorf("bilibili API error: %s (code: %d)", roomResp.Message, roomResp.Code)
		return nil, errors2.StreamingPlatformError(string(p.GetPlatformType()), "", err)
	}

	// Get streamer info
	var streamerResp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Info struct {
				UID   int64  `json:"uid"`
				Uname string `json:"uname"`
				Face  string `json:"face"`
			} `json:"info"`
		} `json:"data"`
	}

	_, err = p.client.R().
		SetContext(ctx).
		SetQueryParam("uid", strconv.FormatInt(roomResp.Data.UID, 10)).
		SetResult(&streamerResp).
		Get("https://api.live.bilibili.com/live_user/v1/Master/info")

	if err != nil {
		p.logger.Error("Failed to fetch Bilibili streamer info",
			zap.Int64("uid", roomResp.Data.UID),
			zap.Error(err))
		return nil, errors2.StreamingPlatformError(string(p.GetPlatformType()), "failed to fetch streamer info", err)
	}

	if streamerResp.Code != 0 {
		err := fmt.Errorf("bilibili API error: %s (code: %d)", streamerResp.Message, streamerResp.Code)
		return nil, errors2.StreamingPlatformError(string(p.GetPlatformType()), "", err)
	}

	return &external.StreamerInfo{
		PlatformStreamerId: platformStreamerId,
		Name:               streamerResp.Data.Info.Uname,
		Avatar:             streamerResp.Data.Info.Face,
		Description:        roomResp.Data.Description,
		RoomURL:            fmt.Sprintf("https://live.bilibili.com/%d", roomID),
	}, nil
}

func (p *Provider) CheckLiveStatus(ctx context.Context, platformStreamerId string) (*external.LiveStatus, error) {
	roomID, err := strconv.ParseInt(platformStreamerId, 10, 64)
	if err != nil {
		return nil, errors2.StreamingPlatformError(string(p.GetPlatformType()), "invalid room id", err)
	}

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			RoomID     int64  `json:"room_id"`
			UID        int64  `json:"uid"`
			Title      string `json:"title"`
			LiveStatus int    `json:"live_status"` // 0: offline, 1: live, 2: replay
			LiveTime   string `json:"live_time"`   // "YYYY-MM-DD HH:mm:ss" or "0000-00-00 00:00:00"
			Online     int    `json:"online"`
			UserCover  string `json:"user_cover"`
			AreaName   string `json:"area_name"`
		} `json:"data"`
	}

	result, err := p.client.R().
		SetContext(ctx).
		SetQueryParam("room_id", strconv.FormatInt(roomID, 10)).
		SetResult(&resp).
		Get("https://api.live.bilibili.com/room/v1/Room/get_info")

	if err != nil {
		p.logger.Error("Failed to check Bilibili live status",
			zap.String("room_id", platformStreamerId),
			zap.Error(err))

		return nil, errors2.StreamingPlatformError(string(p.GetPlatformType()), "failed to check live status", err)
	}

	if result.IsError() {
		err := fmt.Errorf("API returned error status: %d", result.StatusCode())
		return nil, errors2.StreamingPlatformError(string(p.GetPlatformType()), "", err)
	}

	if resp.Code != 0 {
		err := fmt.Errorf("bilibili API error: %s (code: %d)", resp.Message, resp.Code)
		return nil, errors2.StreamingPlatformError(string(p.GetPlatformType()), "", err)
	}

	// Parse live time
	var startTime time.Time
	if resp.Data.LiveStatus == 1 && resp.Data.LiveTime != "0000-00-00 00:00:00" {
		startTime, _ = time.ParseInLocation("2006-01-02 15:04:05", resp.Data.LiveTime, time.Local)
	}

	return &external.LiveStatus{
		IsLive:     resp.Data.LiveStatus == 1, // Only count status 1 as live
		Title:      resp.Data.Title,
		GameName:   resp.Data.AreaName,
		StartTime:  startTime,
		Viewers:    resp.Data.Online,
		CoverImage: resp.Data.UserCover,
	}, nil
}

func (p *Provider) BatchCheckLiveStatus(ctx context.Context, platformStreamerIds []string) (map[string]*external.LiveStatus, error) {
	if len(platformStreamerIds) == 0 {
		return make(map[string]*external.LiveStatus), nil
	}

	results := make(map[string]*external.LiveStatus)

	for _, roomID := range platformStreamerIds {
		status, err := p.CheckLiveStatus(ctx, roomID)
		if err != nil {
			p.logger.Warn("Failed to check live status for room",
				zap.String("room_id", roomID),
				zap.Error(err))
			// Continue with other rooms even if one fails
			continue
		}
		results[roomID] = status
	}

	return results, nil
}
