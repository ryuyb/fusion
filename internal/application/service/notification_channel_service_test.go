package service

import (
	"context"
	"testing"

	"github.com/ryuyb/fusion/internal/core/command"
	"github.com/ryuyb/fusion/internal/core/domain"
	repoMocks "github.com/ryuyb/fusion/internal/core/port/repository"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNotificationChannelService_Create(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockNotificationChannelRepository(t)
	svc := NewNotificationChannelService(repo, zap.NewNop())

	cmd := &command.CreateNotificationChannelCommand{
		UserID:      10,
		ChannelType: string(domain.ChannelTypeEmail),
		Name:        "email",
		Config: map[string]any{
			"address": "user@example.com",
		},
		Enable:   true,
		Priority: 1,
	}
	expected := &domain.NotificationChannel{ID: 1, UserID: cmd.UserID, Name: cmd.Name}

	repo.EXPECT().ExistByName(ctx, cmd.UserID, cmd.Name).Return(false, nil)
	repo.EXPECT().Create(ctx, mock.MatchedBy(func(channel *domain.NotificationChannel) bool {
		return channel.UserID == cmd.UserID &&
			channel.Name == cmd.Name &&
			channel.ChannelType == domain.NotificationChannelType(cmd.ChannelType) &&
			channel.Config["address"] == cmd.Config["address"]
	})).Return(expected, nil)

	created, err := svc.Create(ctx, cmd)
	require.NoError(t, err)
	require.Equal(t, expected, created)
}

func TestNotificationChannelService_CreateConflict(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockNotificationChannelRepository(t)
	svc := NewNotificationChannelService(repo, zap.NewNop())

	cmd := &command.CreateNotificationChannelCommand{UserID: 10, Name: "email"}

	repo.EXPECT().ExistByName(ctx, cmd.UserID, cmd.Name).Return(true, nil)

	_, err := svc.Create(ctx, cmd)
	require.Error(t, err)
}

func TestNotificationChannelService_Update(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockNotificationChannelRepository(t)
	svc := NewNotificationChannelService(repo, zap.NewNop())

	current := &domain.NotificationChannel{ID: 1, UserID: 10, Name: "email"}
	cmd := &command.UpdateNotificationChannelCommand{
		ID: current.ID,
		CreateNotificationChannelCommand: &command.CreateNotificationChannelCommand{
			UserID: 10,
			Name:   "email",
		},
	}
	expected := &domain.NotificationChannel{ID: current.ID, UserID: current.UserID, Name: current.Name}

	repo.EXPECT().FindById(ctx, cmd.ID).Return(current, nil)
	repo.EXPECT().Update(ctx, mock.MatchedBy(func(channel *domain.NotificationChannel) bool {
		return channel.ID == cmd.ID && channel.Name == cmd.Name
	})).Return(expected, nil)

	got, err := svc.Update(ctx, cmd)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestNotificationChannelService_UpdateConflict(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockNotificationChannelRepository(t)
	svc := NewNotificationChannelService(repo, zap.NewNop())

	current := &domain.NotificationChannel{ID: 1, UserID: 10, Name: "email"}
	cmd := &command.UpdateNotificationChannelCommand{
		ID: current.ID,
		CreateNotificationChannelCommand: &command.CreateNotificationChannelCommand{
			UserID: 10,
			Name:   "push",
		},
	}

	repo.EXPECT().FindById(ctx, cmd.ID).Return(current, nil)
	repo.EXPECT().ExistByName(ctx, current.UserID, cmd.Name).Return(true, nil)

	_, err := svc.Update(ctx, cmd)
	require.Error(t, err)
}

func TestNotificationChannelService_ListInvalidPagination(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockNotificationChannelRepository(t)
	svc := NewNotificationChannelService(repo, zap.NewNop())

	_, _, err := svc.ListByUserId(ctx, 1, 0, 10)
	require.Error(t, err)
}
