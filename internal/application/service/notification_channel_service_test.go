package service

import (
	"context"
	"testing"

	"github.com/ryuyb/fusion/internal/core/domain"
	repoMocks "github.com/ryuyb/fusion/internal/core/port/repository"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNotificationChannelService_Create(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockNotificationChannelRepository(t)
	svc := NewNotificationChannelService(repo, zap.NewNop())

	channel := &domain.NotificationChannel{ID: 1, UserID: 10, Name: "email"}

	repo.EXPECT().ExistByName(ctx, channel.UserID, channel.Name).Return(false, nil)
	repo.EXPECT().Create(ctx, channel).Return(channel, nil)

	created, err := svc.Create(ctx, channel)
	require.NoError(t, err)
	require.Equal(t, channel, created)
}

func TestNotificationChannelService_CreateConflict(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockNotificationChannelRepository(t)
	svc := NewNotificationChannelService(repo, zap.NewNop())

	channel := &domain.NotificationChannel{UserID: 10, Name: "email"}

	repo.EXPECT().ExistByName(ctx, channel.UserID, channel.Name).Return(true, nil)

	_, err := svc.Create(ctx, channel)
	require.Error(t, err)
}

func TestNotificationChannelService_Update(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockNotificationChannelRepository(t)
	svc := NewNotificationChannelService(repo, zap.NewNop())

	current := &domain.NotificationChannel{ID: 1, UserID: 10, Name: "email"}
	updated := &domain.NotificationChannel{ID: 1, UserID: 10, Name: "email"}

	repo.EXPECT().FindById(ctx, current.ID).Return(current, nil)
	repo.EXPECT().Update(ctx, updated).Return(updated, nil)

	got, err := svc.Update(ctx, updated)
	require.NoError(t, err)
	require.Equal(t, updated, got)
}

func TestNotificationChannelService_UpdateConflict(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockNotificationChannelRepository(t)
	svc := NewNotificationChannelService(repo, zap.NewNop())

	current := &domain.NotificationChannel{ID: 1, UserID: 10, Name: "email"}
	updated := &domain.NotificationChannel{ID: 1, UserID: 10, Name: "push"}

	repo.EXPECT().FindById(ctx, current.ID).Return(current, nil)
	repo.EXPECT().ExistByName(ctx, current.UserID, updated.Name).Return(true, nil)

	_, err := svc.Update(ctx, updated)
	require.Error(t, err)
}

func TestNotificationChannelService_ListInvalidPagination(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockNotificationChannelRepository(t)
	svc := NewNotificationChannelService(repo, zap.NewNop())

	_, _, err := svc.ListByUserId(ctx, 1, 0, 10)
	require.Error(t, err)
}
