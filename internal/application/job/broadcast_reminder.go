package job

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ryuyb/fusion/internal/core/domain"
	coreExternal "github.com/ryuyb/fusion/internal/core/port/external"
	coreRepo "github.com/ryuyb/fusion/internal/core/port/repository"
	coreService "github.com/ryuyb/fusion/internal/core/port/service"
	notificationInfra "github.com/ryuyb/fusion/internal/infrastructure/external/notification"
	appErrors "github.com/ryuyb/fusion/internal/pkg/errors"
	"go.uber.org/zap"
)

const (
	BroadcastReminderJob = "broadcast_reminder"

	streamerBatchSize = 50
	followBatchSize   = 100
	channelBatchSize  = 100
)

type BroadcastReminder struct {
	logger                *zap.Logger
	streamerRepo          coreRepo.StreamerRepository
	followRepo            coreRepo.UserFollowedStreamerRepository
	channelRepo           coreRepo.NotificationChannelRepository
	streamerService       coreService.StreamerService
	notificationProviders *notificationInfra.NotificationProviderManager
}

func NewBroadcastReminder(
	logger *zap.Logger,
	streamerRepo coreRepo.StreamerRepository,
	followRepo coreRepo.UserFollowedStreamerRepository,
	channelRepo coreRepo.NotificationChannelRepository,
	streamerService coreService.StreamerService,
	notificationProviders *notificationInfra.NotificationProviderManager,
) *BroadcastReminder {
	return &BroadcastReminder{
		logger:                logger,
		streamerRepo:          streamerRepo,
		followRepo:            followRepo,
		channelRepo:           channelRepo,
		streamerService:       streamerService,
		notificationProviders: notificationProviders,
	}
}

func (j *BroadcastReminder) Name() string {
	return BroadcastReminderJob
}

func (j *BroadcastReminder) Execute(ctx context.Context) error {
	resolver := newChannelResolver(j.channelRepo)

	offset := 0
	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		streamers, total, err := j.streamerRepo.List(ctx, offset, streamerBatchSize)
		if err != nil {
			return err
		}
		if len(streamers) == 0 {
			break
		}

		for _, streamer := range streamers {
			if err := j.processStreamer(ctx, streamer, resolver); err != nil {
				j.logger.Warn("failed to process streamer for reminders",
					zap.Int64("streamer_id", streamer.ID),
					zap.Error(err))
			}
		}

		offset += len(streamers)
		if offset >= total {
			break
		}
	}

	return nil
}

func (j *BroadcastReminder) processStreamer(ctx context.Context, streamer *domain.Streamer, resolver *channelResolver) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	refreshed, err := j.streamerService.FindByPlatformStreamerId(ctx, streamer.PlatformType, streamer.PlatformStreamerID, true)
	if err != nil {
		return err
	}
	if refreshed == nil || !refreshed.LiveStatus.IsLive {
		return nil
	}

	follows, err := j.listFollowers(ctx, refreshed.ID)
	if err != nil {
		return err
	}

	for _, follow := range follows {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := j.processFollower(ctx, follow, refreshed, resolver); err != nil {
			j.logger.Warn("failed to process follower notification",
				zap.Int64("follow_id", follow.ID),
				zap.Int64("streamer_id", refreshed.ID),
				zap.Error(err))
		}
	}
	return nil
}

func (j *BroadcastReminder) listFollowers(ctx context.Context, streamerID int64) ([]*domain.UserFollowedStreamer, error) {
	var (
		results []*domain.UserFollowedStreamer
		offset  int
	)
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		follows, total, err := j.followRepo.ListByStreamerId(ctx, streamerID, offset, followBatchSize)
		if err != nil {
			return nil, err
		}
		if len(follows) == 0 {
			break
		}
		for _, follow := range follows {
			if follow.NotificationsEnabled {
				results = append(results, follow)
			}
		}
		offset += len(follows)
		if offset >= total {
			break
		}
	}
	return results, nil
}

func (j *BroadcastReminder) processFollower(ctx context.Context, follow *domain.UserFollowedStreamer, streamer *domain.Streamer, resolver *channelResolver) error {
	if !j.shouldSend(follow, streamer) {
		return nil
	}

	channels, err := resolver.Resolve(ctx, follow)
	if err != nil {
		return err
	}
	if len(channels) == 0 {
		return nil
	}

	data := buildNotificationData(follow, streamer)
	sent := false
	for _, channel := range channels {
		if !channel.Enable {
			continue
		}
		provider, err := j.notificationProviders.GetProvider(channel.ChannelType)
		if err != nil {
			j.logger.Warn("notification provider unavailable",
				zap.String("channel_type", string(channel.ChannelType)),
				zap.Error(err))
			continue
		}
		if err := provider.Send(ctx, channel, data); err != nil {
			j.logger.Warn("failed to send notification",
				zap.Int64("channel_id", channel.ID),
				zap.Int64("follow_id", follow.ID),
				zap.Error(err))
			continue
		}
		sent = true
	}

	if !sent {
		return nil
	}

	now := time.Now()
	follow.LastNotificationSentAt = &now
	_, err = j.followRepo.Update(ctx, follow)
	return err
}

func (j *BroadcastReminder) shouldSend(follow *domain.UserFollowedStreamer, streamer *domain.Streamer) bool {
	if !streamer.LiveStatus.IsLive {
		return false
	}
	if follow.LastNotificationSentAt == nil {
		return true
	}
	if !streamer.LiveStatus.StartTime.IsZero() {
		return streamer.LiveStatus.StartTime.After(*follow.LastNotificationSentAt)
	}
	return streamer.LastLiveSyncedAt.After(*follow.LastNotificationSentAt)
}

func buildNotificationData(follow *domain.UserFollowedStreamer, streamer *domain.Streamer) *coreExternal.NotificationData {
	displayName := streamer.DisplayName
	if strings.TrimSpace(follow.Alias) != "" {
		displayName = follow.Alias
	}

	title := fmt.Sprintf("%s is live now!", displayName)
	body := streamer.LiveStatus.Title
	if streamer.RoomURL != "" {
		if body != "" {
			body = body + "\n"
		}
		body += streamer.RoomURL
	}
	if body == "" {
		body = "Tune in now."
	}
	return &coreExternal.NotificationData{
		Title:   title,
		Content: body,
	}
}

type channelResolver struct {
	repo   coreRepo.NotificationChannelRepository
	byID   map[int64]*domain.NotificationChannel
	byUser map[int64][]*domain.NotificationChannel
}

func newChannelResolver(repo coreRepo.NotificationChannelRepository) *channelResolver {
	return &channelResolver{
		repo:   repo,
		byID:   make(map[int64]*domain.NotificationChannel),
		byUser: make(map[int64][]*domain.NotificationChannel),
	}
}

func (r *channelResolver) Resolve(ctx context.Context, follow *domain.UserFollowedStreamer) ([]*domain.NotificationChannel, error) {
	if len(follow.NotificationChannelIDs) > 0 {
		return r.channelsByIDs(ctx, follow.UserID, follow.NotificationChannelIDs)
	}
	return r.channelsByUser(ctx, follow.UserID)
}

func (r *channelResolver) channelsByIDs(ctx context.Context, userID int64, ids []int64) ([]*domain.NotificationChannel, error) {
	results := make([]*domain.NotificationChannel, 0, len(ids))
	for _, id := range ids {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		channel, err := r.channelByID(ctx, id)
		if err != nil {
			if appErrors.IsNotFoundError(err) {
				continue
			}
			return nil, err
		}
		if channel == nil {
			continue
		}
		if channel.UserID != userID || !channel.Enable {
			continue
		}
		results = append(results, channel)
	}
	return results, nil
}

func (r *channelResolver) channelsByUser(ctx context.Context, userID int64) ([]*domain.NotificationChannel, error) {
	if channels, ok := r.byUser[userID]; ok {
		return channels, nil
	}

	var (
		channels []*domain.NotificationChannel
		offset   int
		total    int
	)
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		batch, count, err := r.repo.ListByUserId(ctx, userID, offset, channelBatchSize)
		if err != nil {
			return nil, err
		}
		total = count
		if len(batch) == 0 {
			break
		}
		for _, channel := range batch {
			r.byID[channel.ID] = channel
			if channel.Enable {
				channels = append(channels, channel)
			}
		}
		offset += len(batch)
		if offset >= total {
			break
		}
	}

	r.byUser[userID] = channels
	return channels, nil
}

func (r *channelResolver) channelByID(ctx context.Context, id int64) (*domain.NotificationChannel, error) {
	if channel, ok := r.byID[id]; ok {
		return channel, nil
	}
	channel, err := r.repo.FindById(ctx, id)
	if err != nil {
		if appErrors.IsNotFoundError(err) {
			r.byID[id] = nil
			return nil, err
		}
		return nil, err
	}
	r.byID[id] = channel
	return channel, nil
}
