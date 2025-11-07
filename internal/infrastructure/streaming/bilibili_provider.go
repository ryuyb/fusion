package streaming

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/domain/service"
	"go.uber.org/zap"
)

// BilibiliProvider implements StreamingPlatformProvider for Bilibili Live
type BilibiliProvider struct {
	client *RestyClient
	logger *zap.Logger
}

// NewBilibiliProvider creates a new BilibiliProvider instance
func NewBilibiliProvider(client *RestyClient, logger *zap.Logger) *BilibiliProvider {
	return &BilibiliProvider{
		client: client,
		logger: logger,
	}
}

// GetPlatformType returns the platform type
func (p *BilibiliProvider) GetPlatformType() entity.PlatformType {
	return entity.PlatformTypeBilibili
}

// FetchStreamerInfo fetches detailed information about a streamer
func (p *BilibiliProvider) FetchStreamerInfo(ctx context.Context, platformStreamerId string) (*service.StreamerInfo, error) {
	// First, get room info to get the room_id
	roomID, err := strconv.ParseInt(platformStreamerId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid room id: %w", err)
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
		return nil, fmt.Errorf("failed to fetch room info: %w", err)
	}

	if resp.IsError() {
		p.logger.Error("Bilibili API returned error",
			zap.String("room_id", platformStreamerId),
			zap.Int("status_code", resp.StatusCode()))
		return nil, fmt.Errorf("API returned error status: %d", resp.StatusCode())
	}

	if roomResp.Code != 0 {
		return nil, fmt.Errorf("bilibili API error: %s (code: %d)", roomResp.Message, roomResp.Code)
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
		return nil, fmt.Errorf("failed to fetch streamer info: %w", err)
	}

	if streamerResp.Code != 0 {
		return nil, fmt.Errorf("bilibili API error: %s (code: %d)", streamerResp.Message, streamerResp.Code)
	}

	return &service.StreamerInfo{
		PlatformStreamerId: platformStreamerId,
		Name:               streamerResp.Data.Info.Uname,
		Avatar:             streamerResp.Data.Info.Face,
		Description:        roomResp.Data.Description,
		RoomURL:            fmt.Sprintf("https://live.bilibili.com/%d", roomID),
	}, nil
}

// CheckLiveStatus checks the live status of a single streamer
func (p *BilibiliProvider) CheckLiveStatus(ctx context.Context, platformStreamerId string) (*service.LiveStatus, error) {
	roomID, err := strconv.ParseInt(platformStreamerId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid room id: %w", err)
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
		return nil, fmt.Errorf("failed to check live status: %w", err)
	}

	if result.IsError() {
		return nil, fmt.Errorf("API returned error status: %d", result.StatusCode())
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("bilibili API error: %s (code: %d)", resp.Message, resp.Code)
	}

	// Parse live time
	var startTime time.Time
	if resp.Data.LiveStatus == 1 && resp.Data.LiveTime != "0000-00-00 00:00:00" {
		startTime, _ = time.ParseInLocation("2006-01-02 15:04:05", resp.Data.LiveTime, time.Local)
	}

	return &service.LiveStatus{
		IsLive:     resp.Data.LiveStatus == 1, // Only count status 1 as live
		Title:      resp.Data.Title,
		GameName:   resp.Data.AreaName,
		StartTime:  startTime,
		Viewers:    resp.Data.Online,
		CoverImage: resp.Data.UserCover,
	}, nil
}

// BatchCheckLiveStatus checks live status for multiple streamers
func (p *BilibiliProvider) BatchCheckLiveStatus(ctx context.Context, platformStreamerIds []string) (map[string]*service.LiveStatus, error) {
	if len(platformStreamerIds) == 0 {
		return make(map[string]*service.LiveStatus), nil
	}

	// Convert room IDs to UIDs first (we need to get room info to get UIDs)
	// For batch checking, we'll use the get_status_info_by_uids API which requires UIDs
	// However, we have room IDs, so we need to first convert them to UIDs
	// For simplicity and to avoid rate limiting, we'll do sequential checks for now
	// In production, you might want to implement a cache mapping room_id -> uid

	results := make(map[string]*service.LiveStatus)

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

// ValidateConfiguration validates the platform configuration
func (p *BilibiliProvider) ValidateConfiguration(config map[string]interface{}) error {
	// Bilibili Live API doesn't require authentication for basic operations
	// No configuration validation needed
	return nil
}

// SearchStreamer searches for streamers by keyword
func (p *BilibiliProvider) SearchStreamer(ctx context.Context, keyword string) ([]*service.StreamerInfo, error) {
	// Note: Bilibili doesn't provide a public search API for live rooms
	// This would require web scraping or access to their internal API
	// For now, we return an error indicating this feature is not supported
	p.logger.Warn("Search streamer is not supported for Bilibili",
		zap.String("keyword", keyword))
	return nil, fmt.Errorf("search feature is not supported for Bilibili platform")
}
