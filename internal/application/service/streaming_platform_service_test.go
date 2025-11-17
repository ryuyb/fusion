package service

import (
	"context"
	"testing"

	"github.com/ryuyb/fusion/internal/core/domain"
	repoMocks "github.com/ryuyb/fusion/internal/core/port/repository"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestStreamingPlatformService_Create(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockStreamingPlatformRepository(t)
	svc := NewStreamingPlatformService(repo, zap.NewNop())

	platform := &domain.StreamingPlatform{ID: 1, Name: "Bilibili"}

	repo.EXPECT().ExistByName(ctx, platform.Name).Return(false, nil)
	repo.EXPECT().Create(ctx, platform).Return(platform, nil)

	created, err := svc.Create(ctx, platform)
	require.NoError(t, err)
	require.Equal(t, platform, created)
}

func TestStreamingPlatformService_CreateDuplicate(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockStreamingPlatformRepository(t)
	svc := NewStreamingPlatformService(repo, zap.NewNop())

	platform := &domain.StreamingPlatform{Name: "Bilibili"}

	repo.EXPECT().ExistByName(ctx, platform.Name).Return(true, nil)

	_, err := svc.Create(ctx, platform)
	require.Error(t, err)
}

func TestStreamingPlatformService_Update(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockStreamingPlatformRepository(t)
	svc := NewStreamingPlatformService(repo, zap.NewNop())

	current := &domain.StreamingPlatform{ID: 1, Name: "Bilibili"}
	updated := &domain.StreamingPlatform{ID: 1, Name: "Bilibili"}

	repo.EXPECT().FindById(ctx, current.ID).Return(current, nil)
	repo.EXPECT().Update(ctx, updated).Return(updated, nil)

	got, err := svc.Update(ctx, updated)
	require.NoError(t, err)
	require.Equal(t, updated, got)
}

func TestStreamingPlatformService_UpdateConflict(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockStreamingPlatformRepository(t)
	svc := NewStreamingPlatformService(repo, zap.NewNop())

	current := &domain.StreamingPlatform{ID: 1, Name: "old"}
	updated := &domain.StreamingPlatform{ID: 1, Name: "new"}

	repo.EXPECT().FindById(ctx, updated.ID).Return(current, nil)
	repo.EXPECT().ExistByName(ctx, updated.Name).Return(true, nil)

	_, err := svc.Update(ctx, updated)
	require.Error(t, err)
}

func TestStreamingPlatformService_ListPaginationError(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockStreamingPlatformRepository(t)
	svc := NewStreamingPlatformService(repo, zap.NewNop())

	_, _, err := svc.List(ctx, 0, 10)
	require.Error(t, err)
}

func TestStreamingPlatformService_DeleteAndFind(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockStreamingPlatformRepository(t)
	svc := NewStreamingPlatformService(repo, zap.NewNop())

	repo.EXPECT().Delete(ctx, int64(1)).Return(nil)
	require.NoError(t, svc.Delete(ctx, 1))

	expected := &domain.StreamingPlatform{ID: 2}
	repo.EXPECT().FindById(ctx, int64(2)).Return(expected, nil)
	got, err := svc.FindById(ctx, 2)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}
