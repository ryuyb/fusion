package streaming

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/domain/service"
	"github.com/ryuyb/fusion/internal/infrastructure/client"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/samber/lo"
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

type BetardResponse struct {
	Room struct {
		Nickname    string `json:"nickname"`
		OwnerAvatar string `json:"owner_avatar"`
		Status      string `json:"status"`
		ShowStatus  int    `json:"show_status"`
		ShowDetails string `json:"show_details"`
		RoomName    string `json:"room_name"`
		RoomPic     string `json:"room_pic"`
		CoverSrc    string `json:"coverSrc"`
		ShowTime    int64  `json:"show_time"` // Unix seconds 1762482715
		Avatar      struct {
			Big    string `json:"big"`
			Middle string `json:"middle"`
			Small  string `json:"small"`
		} `json:"avatar"`
		CateName      string `json:"cate_name"`
		SecondLvlName string `json:"second_lvl_name"`
		RoomBizAll    struct {
			Hot string `json:"hot"`
		} `json:"room_biz_all"`
	} `json:"room"`
	Column struct {
		CateId   string `json:"cate_id"`
		CateName string `json:"cate_name"`
	} `json:"column"`
}

func (d DouyuProvider) GetPlatformType() entity.PlatformType {
	return entity.PlatformTypeDouyu
}

func (d DouyuProvider) FetchStreamerInfo(ctx context.Context, platformStreamerId string) (*service.StreamerInfo, error) {
	betardResp := &BetardResponse{}
	_, err := d.client.R().
		SetContext(ctx).
		SetPathParam("roomId", platformStreamerId).
		SetResult(betardResp).
		Get("https://www.douyu.com/betard/{roomId}")
	if err != nil {
		return nil, errors2.StreamingPlatformError(d.GetPlatformType(), lo.ToPtr("failed to fetch betard"), err)
	}
	return &service.StreamerInfo{
		PlatformStreamerId: platformStreamerId,
		Name:               betardResp.Room.Nickname,
		Avatar:             betardResp.Room.Avatar.Big,
		Description:        betardResp.Room.ShowDetails,
		RoomURL:            fmt.Sprintf("https://www.douyu.com/%s", platformStreamerId),
	}, nil
}

func (d DouyuProvider) CheckLiveStatus(ctx context.Context, platformStreamerId string) (*service.LiveStatus, error) {
	betardResp := &BetardResponse{}
	_, err := d.client.R().
		SetContext(ctx).
		SetPathParam("roomId", platformStreamerId).
		SetResult(betardResp).
		Get("https://www.douyu.com/betard/{roomId}")
	if err != nil {
		return nil, errors2.StreamingPlatformError(d.GetPlatformType(), lo.ToPtr("failed to fetch betard"), err)
	}
	viewers, err := strconv.Atoi(betardResp.Room.RoomBizAll.Hot)
	if err != nil {
		return nil, errors2.StreamingPlatformError(d.GetPlatformType(), lo.ToPtr("failed to parse hot"), err)
	}
	return &service.LiveStatus{
		IsLive:     betardResp.Room.ShowStatus == 1,
		Title:      betardResp.Room.RoomName,
		GameName:   betardResp.Room.SecondLvlName,
		StartTime:  time.Unix(betardResp.Room.ShowTime, 0),
		Viewers:    viewers,
		CoverImage: betardResp.Room.CoverSrc,
	}, nil
}

func (d DouyuProvider) BatchCheckLiveStatus(ctx context.Context, platformStreamerIds []string) (map[string]*service.LiveStatus, error) {
	results := make(map[string]*service.LiveStatus)

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

func (d DouyuProvider) ValidateConfiguration(config map[string]interface{}) error {
	return nil
}

func (d DouyuProvider) SearchStreamer(ctx context.Context, keyword string) ([]*service.StreamerInfo, error) {
	d.logger.Warn("Search streamer is not supported for Douyu", zap.String("keyword", keyword))
	return nil, errors2.StreamingPlatformError(d.GetPlatformType(), lo.ToPtr("Search streamer is not supported"), nil)
}
