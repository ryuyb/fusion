package job

import (
	"context"
	"testing"
	"time"

	"github.com/ryuyb/fusion/internal/core/domain"
	coreExternal "github.com/ryuyb/fusion/internal/core/port/external"
	repoMocks "github.com/ryuyb/fusion/internal/core/port/repository"
	serviceMocks "github.com/ryuyb/fusion/internal/core/port/service"
	notificationInfra "github.com/ryuyb/fusion/internal/infrastructure/external/notification"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestBroadcastReminder_SendNotificationForLiveStreamer(t *testing.T) {
	ctx := context.Background()
	streamer := &domain.Streamer{
		ID:                 1,
		PlatformType:       domain.StreamingPlatformTypeBilibili,
		PlatformStreamerID: "1001",
		DisplayName:        "Streamer",
	}
	live := *streamer
	live.LiveStatus = domain.LiveStatusInfo{
		IsLive:    true,
		Title:     "Playing Game",
		StartTime: time.Now(),
	}
	live.RoomURL = "https://live.example/1001"

	streamerRepo := repoMocks.NewMockStreamerRepository(t)
	streamerRepo.EXPECT().
		List(mock.Anything, 0, streamerBatchSize).
		Return([]*domain.Streamer{streamer}, 1, nil).Once()

	follow := &domain.UserFollowedStreamer{
		ID:                     10,
		UserID:                 99,
		StreamerID:             streamer.ID,
		NotificationsEnabled:   true,
		NotificationChannelIDs: []int64{7},
	}
	followRepo := repoMocks.NewMockUserFollowedStreamerRepository(t)
	followRepo.EXPECT().
		ListByStreamerId(mock.Anything, streamer.ID, 0, followBatchSize).
		Return([]*domain.UserFollowedStreamer{follow}, 1, nil).Once()

	var updatedFollow *domain.UserFollowedStreamer
	followRepo.EXPECT().
		Update(mock.Anything, mock.MatchedBy(func(f *domain.UserFollowedStreamer) bool {
			updatedFollow = f
			return f.ID == follow.ID && f.LastNotificationSentAt != nil
		})).
		Return(follow, nil).Once()

	channel := &domain.NotificationChannel{
		ID:          7,
		UserID:      follow.UserID,
		ChannelType: domain.ChannelTypeBark,
		Name:        "bark",
		Enable:      true,
		Config: map[string]any{
			"device_key": "abc",
		},
	}
	channelRepo := repoMocks.NewMockNotificationChannelRepository(t)
	channelRepo.EXPECT().
		FindById(mock.Anything, channel.ID).
		Return(channel, nil).Once()

	streamerService := serviceMocks.NewMockStreamerService(t)
	streamerService.EXPECT().
		FindByPlatformStreamerId(mock.Anything, streamer.PlatformType, streamer.PlatformStreamerID, true).
		Return(&live, nil).Once()

	provider := coreExternal.NewMockNotificationProvider(t)
	provider.EXPECT().GetChannelType().Return(domain.ChannelTypeBark).Twice()
	provider.EXPECT().
		Send(mock.Anything, channel, mock.AnythingOfType("*external.NotificationData")).
		Return(nil).Once()

	manager := notificationInfra.NewNotificationProviderManager(
		[]coreExternal.NotificationProvider{provider},
		zap.NewNop(),
	)

	job := NewBroadcastReminder(
		zap.NewNop(),
		streamerRepo,
		followRepo,
		channelRepo,
		streamerService,
		manager,
	)

	err := job.Execute(ctx)
	require.NoError(t, err)
	require.NotNil(t, updatedFollow)
	require.NotNil(t, updatedFollow.LastNotificationSentAt)
}

func TestBroadcastReminder_FallbackToAllChannels(t *testing.T) {
	ctx := context.Background()
	streamer := &domain.Streamer{
		ID:                 2,
		PlatformType:       domain.StreamingPlatformTypeDouyu,
		PlatformStreamerID: "2002",
		DisplayName:        "Another",
	}
	live := *streamer
	live.LiveStatus = domain.LiveStatusInfo{
		IsLive:    true,
		Title:     "Just Chatting",
		Viewers:   1234,
		StartTime: time.Now(),
	}

	streamerRepo := repoMocks.NewMockStreamerRepository(t)
	streamerRepo.EXPECT().
		List(mock.Anything, 0, streamerBatchSize).
		Return([]*domain.Streamer{streamer}, 1, nil).Once()

	follow := &domain.UserFollowedStreamer{
		ID:                     20,
		UserID:                 77,
		StreamerID:             streamer.ID,
		NotificationsEnabled:   true,
		NotificationChannelIDs: nil,
	}
	followRepo := repoMocks.NewMockUserFollowedStreamerRepository(t)
	followRepo.EXPECT().
		ListByStreamerId(mock.Anything, streamer.ID, 0, followBatchSize).
		Return([]*domain.UserFollowedStreamer{follow}, 1, nil).Once()
	followRepo.EXPECT().
		Update(mock.Anything, mock.AnythingOfType("*domain.UserFollowedStreamer")).
		Return(follow, nil).Once()

	channels := []*domain.NotificationChannel{
		{
			ID:          8,
			UserID:      follow.UserID,
			ChannelType: domain.ChannelTypeBark,
			Name:        "bark",
			Enable:      true,
		},
		{
			ID:          9,
			UserID:      follow.UserID,
			ChannelType: domain.ChannelTypeTelegram,
			Name:        "telegram",
			Enable:      false,
		},
	}
	channelRepo := repoMocks.NewMockNotificationChannelRepository(t)
	channelRepo.EXPECT().
		ListByUserId(mock.Anything, follow.UserID, 0, channelBatchSize).
		Return(channels, len(channels), nil).Once()

	streamerService := serviceMocks.NewMockStreamerService(t)
	streamerService.EXPECT().
		FindByPlatformStreamerId(mock.Anything, streamer.PlatformType, streamer.PlatformStreamerID, true).
		Return(&live, nil).Once()

	provider := coreExternal.NewMockNotificationProvider(t)
	provider.EXPECT().GetChannelType().Return(domain.ChannelTypeBark).Twice()
	provider.EXPECT().
		Send(mock.Anything, channels[0], mock.AnythingOfType("*external.NotificationData")).
		Return(nil).Once()

	manager := notificationInfra.NewNotificationProviderManager(
		[]coreExternal.NotificationProvider{provider},
		zap.NewNop(),
	)

	job := NewBroadcastReminder(
		zap.NewNop(),
		streamerRepo,
		followRepo,
		channelRepo,
		streamerService,
		manager,
	)

	err := job.Execute(ctx)
	require.NoError(t, err)
	provider.AssertNumberOfCalls(t, "Send", 1)
}
