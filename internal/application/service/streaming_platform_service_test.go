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

func TestStreamingPlatformService_Create(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockStreamingPlatformRepository(t)
	svc := NewStreamingPlatformService(repo, zap.NewNop())

	cmd := &command.CreateStreamingPlatformCommand{
		Type:        string(domain.StreamingPlatformTypeBilibili),
		Name:        "Bilibili",
		Description: "Video platform",
		BaseURL:     "https://www.bilibili.com",
		LogoURL:     "https://cdn/logo.png",
		Enabled:     true,
		Priority:    10,
		Metadata: map[string]string{
			"foo": "bar",
		},
	}
	expected := &domain.StreamingPlatform{ID: 1, Name: cmd.Name}

	repo.EXPECT().ExistByName(ctx, cmd.Name).Return(false, nil)
	repo.EXPECT().Create(ctx, mock.MatchedBy(func(platform *domain.StreamingPlatform) bool {
		return platform.Name == cmd.Name &&
			platform.Description == cmd.Description &&
			platform.BaseURL == cmd.BaseURL &&
			platform.LogoURL == cmd.LogoURL &&
			platform.Enabled == cmd.Enabled &&
			platform.Priority == cmd.Priority &&
			platform.Metadata["foo"] == "bar"
	})).Return(expected, nil)

	created, err := svc.Create(ctx, cmd)
	require.NoError(t, err)
	require.Equal(t, expected, created)
}

func TestStreamingPlatformService_CreateDuplicate(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockStreamingPlatformRepository(t)
	svc := NewStreamingPlatformService(repo, zap.NewNop())

	cmd := &command.CreateStreamingPlatformCommand{Name: "Bilibili"}

	repo.EXPECT().ExistByName(ctx, cmd.Name).Return(true, nil)

	_, err := svc.Create(ctx, cmd)
	require.Error(t, err)
}

func TestStreamingPlatformService_Update(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockStreamingPlatformRepository(t)
	svc := NewStreamingPlatformService(repo, zap.NewNop())

	current := &domain.StreamingPlatform{ID: 1, Name: "Bilibili"}
	cmd := &command.UpdateStreamingPlatformCommand{
		ID: current.ID,
		CreateStreamingPlatformCommand: &command.CreateStreamingPlatformCommand{
			Type:    string(domain.StreamingPlatformTypeBilibili),
			Name:    "Bilibili",
			BaseURL: "https://www.bilibili.com",
		},
	}
	expected := &domain.StreamingPlatform{ID: current.ID, Name: cmd.Name}

	repo.EXPECT().FindById(ctx, current.ID).Return(current, nil)
	repo.EXPECT().Update(ctx, mock.MatchedBy(func(platform *domain.StreamingPlatform) bool {
		return platform.ID == cmd.ID && platform.Name == cmd.Name
	})).Return(expected, nil)

	got, err := svc.Update(ctx, cmd)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestStreamingPlatformService_UpdateConflict(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockStreamingPlatformRepository(t)
	svc := NewStreamingPlatformService(repo, zap.NewNop())

	current := &domain.StreamingPlatform{ID: 1, Name: "old"}
	cmd := &command.UpdateStreamingPlatformCommand{
		ID: current.ID,
		CreateStreamingPlatformCommand: &command.CreateStreamingPlatformCommand{
			Name: "new",
		},
	}

	repo.EXPECT().FindById(ctx, cmd.ID).Return(current, nil)
	repo.EXPECT().ExistByName(ctx, cmd.Name).Return(true, nil)

	_, err := svc.Update(ctx, cmd)
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
