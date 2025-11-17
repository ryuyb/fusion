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

func TestUserFollowedStreamerService_Create(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockUserFollowedStreamerRepository(t)
	svc := NewUserFollowedStreamerService(repo, zap.NewNop())

	cmd := &command.CreateUserFollowedStreamerCommand{
		UserID:               1,
		StreamerID:           2,
		Alias:                "fav",
		NotificationsEnabled: true,
		NotificationChannelIDs: []int64{
			10,
		},
	}
	follow := &domain.UserFollowedStreamer{ID: 1, UserID: cmd.UserID, StreamerID: cmd.StreamerID}

	repo.EXPECT().ExistByUserAndStreamer(ctx, cmd.UserID, cmd.StreamerID).Return(false, nil)
	repo.EXPECT().Create(ctx, mock.MatchedBy(func(f *domain.UserFollowedStreamer) bool {
		return f.UserID == cmd.UserID && f.StreamerID == cmd.StreamerID && f.Alias == cmd.Alias
	})).Return(follow, nil)

	created, err := svc.Create(ctx, cmd)
	require.NoError(t, err)
	require.Equal(t, follow, created)
}

func TestUserFollowedStreamerService_CreateConflict(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockUserFollowedStreamerRepository(t)
	svc := NewUserFollowedStreamerService(repo, zap.NewNop())

	cmd := &command.CreateUserFollowedStreamerCommand{UserID: 1, StreamerID: 2}

	repo.EXPECT().ExistByUserAndStreamer(ctx, cmd.UserID, cmd.StreamerID).Return(true, nil)

	_, err := svc.Create(ctx, cmd)
	require.Error(t, err)
}

func TestUserFollowedStreamerService_ListInvalid(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockUserFollowedStreamerRepository(t)
	svc := NewUserFollowedStreamerService(repo, zap.NewNop())

	_, _, err := svc.ListByUserId(ctx, 1, 0, 10)
	require.Error(t, err)
}
