package service

import (
	"context"
	"testing"

	"github.com/ryuyb/fusion/internal/core/domain"
	repoMocks "github.com/ryuyb/fusion/internal/core/port/repository"
	"github.com/ryuyb/fusion/internal/infrastructure/external/streaming"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestStreamerService_Create(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockStreamerRepository(t)
	spm := streaming.NewStreamingProviderManager(nil, zap.NewNop())
	svc := NewStreamerService(repo, spm, zap.NewNop())

	streamer := &domain.Streamer{ID: 1, PlatformType: domain.StreamingPlatformTypeBilibili, PlatformStreamerID: "123"}

	repo.EXPECT().ExistByPlatformStreamerId(ctx, streamer.PlatformType, streamer.PlatformStreamerID).Return(false, nil)
	repo.EXPECT().Create(ctx, streamer).Return(streamer, nil)

	created, err := svc.Create(ctx, streamer)
	require.NoError(t, err)
	require.Equal(t, streamer, created)
}

func TestStreamerService_CreateConflict(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockStreamerRepository(t)
	spm := streaming.NewStreamingProviderManager(nil, zap.NewNop())
	svc := NewStreamerService(repo, spm, zap.NewNop())

	streamer := &domain.Streamer{PlatformType: domain.StreamingPlatformTypeBilibili, PlatformStreamerID: "123"}

	repo.EXPECT().ExistByPlatformStreamerId(ctx, streamer.PlatformType, streamer.PlatformStreamerID).Return(true, nil)

	_, err := svc.Create(ctx, streamer)
	require.Error(t, err)
}

func TestStreamerService_Update(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockStreamerRepository(t)
	spm := streaming.NewStreamingProviderManager(nil, zap.NewNop())
	svc := NewStreamerService(repo, spm, zap.NewNop())

	existing := &domain.Streamer{ID: 1, PlatformType: domain.StreamingPlatformTypeBilibili, PlatformStreamerID: "123"}
	updated := &domain.Streamer{ID: 1, PlatformType: domain.StreamingPlatformTypeBilibili, PlatformStreamerID: "123", DisplayName: "New"}

	repo.EXPECT().FindById(ctx, updated.ID).Return(existing, nil)
	repo.EXPECT().Update(ctx, updated).Return(updated, nil)

	got, err := svc.Update(ctx, updated)
	require.NoError(t, err)
	require.Equal(t, updated, got)
}

func TestStreamerService_UpdateConflict(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockStreamerRepository(t)
	spm := streaming.NewStreamingProviderManager(nil, zap.NewNop())
	svc := NewStreamerService(repo, spm, zap.NewNop())

	current := &domain.Streamer{ID: 1, PlatformType: domain.StreamingPlatformTypeBilibili, PlatformStreamerID: "123"}
	updated := &domain.Streamer{ID: 1, PlatformType: domain.StreamingPlatformTypeBilibili, PlatformStreamerID: "456"}

	repo.EXPECT().FindById(ctx, updated.ID).Return(current, nil)
	repo.EXPECT().ExistByPlatformStreamerId(ctx, updated.PlatformType, updated.PlatformStreamerID).Return(true, nil)

	_, err := svc.Update(ctx, updated)
	require.Error(t, err)
}

func TestStreamerService_ListPaginationError(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockStreamerRepository(t)
	spm := streaming.NewStreamingProviderManager(nil, zap.NewNop())
	svc := NewStreamerService(repo, spm, zap.NewNop())

	_, _, err := svc.List(ctx, 0, 10)
	require.Error(t, err)
}
