package service

import (
	"context"
	"testing"

	"github.com/ryuyb/fusion/internal/core/domain"
	repoMocks "github.com/ryuyb/fusion/internal/core/port/repository"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestUserFollowedStreamerService_Create(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockUserFollowedStreamerRepository(t)
	svc := NewUserFollowedStreamerService(repo, zap.NewNop())

	follow := &domain.UserFollowedStreamer{ID: 1, UserID: 1, StreamerID: 2}

	repo.EXPECT().ExistByUserAndStreamer(ctx, follow.UserID, follow.StreamerID).Return(false, nil)
	repo.EXPECT().Create(ctx, follow).Return(follow, nil)

	created, err := svc.Create(ctx, follow)
	require.NoError(t, err)
	require.Equal(t, follow, created)
}

func TestUserFollowedStreamerService_CreateConflict(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockUserFollowedStreamerRepository(t)
	svc := NewUserFollowedStreamerService(repo, zap.NewNop())

	follow := &domain.UserFollowedStreamer{UserID: 1, StreamerID: 2}

	repo.EXPECT().ExistByUserAndStreamer(ctx, follow.UserID, follow.StreamerID).Return(true, nil)

	_, err := svc.Create(ctx, follow)
	require.Error(t, err)
}

func TestUserFollowedStreamerService_ListInvalid(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockUserFollowedStreamerRepository(t)
	svc := NewUserFollowedStreamerService(repo, zap.NewNop())

	_, _, err := svc.ListByUserId(ctx, 1, 0, 10)
	require.Error(t, err)
}
