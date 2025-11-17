package service

import (
	"context"
	"testing"

	"github.com/ryuyb/fusion/internal/core/command"
	"github.com/ryuyb/fusion/internal/core/domain"
	repoMocks "github.com/ryuyb/fusion/internal/core/port/repository"
	"github.com/ryuyb/fusion/internal/infrastructure/external/streaming"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestStreamerService_Create(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockStreamerRepository(t)
	spm := streaming.NewStreamingProviderManager(nil, zap.NewNop())
	svc := NewStreamerService(repo, spm, zap.NewNop())

	cmd := &command.CreateStreamerCommand{
		PlatformType:       string(domain.StreamingPlatformTypeBilibili),
		PlatformStreamerID: "123",
		DisplayName:        "Neo",
		AvatarURL:          "https://cdn/avatar.png",
		RoomURL:            "https://room",
		Bio:                "bio",
		Tags:               []string{"tag1"},
	}
	expected := &domain.Streamer{ID: 1, PlatformType: domain.StreamingPlatformTypeBilibili, PlatformStreamerID: "123"}

	repo.EXPECT().ExistByPlatformStreamerId(ctx, domain.StreamingPlatformType(cmd.PlatformType), cmd.PlatformStreamerID).Return(false, nil)
	repo.EXPECT().Create(ctx, mock.MatchedBy(func(streamer *domain.Streamer) bool {
		return streamer.PlatformType == domain.StreamingPlatformType(cmd.PlatformType) &&
			streamer.PlatformStreamerID == cmd.PlatformStreamerID &&
			streamer.DisplayName == cmd.DisplayName &&
			streamer.AvatarURL == cmd.AvatarURL &&
			streamer.RoomURL == cmd.RoomURL &&
			streamer.Bio == cmd.Bio &&
			len(streamer.Tags) == len(cmd.Tags)
	})).Return(expected, nil)

	created, err := svc.Create(ctx, cmd)
	require.NoError(t, err)
	require.Equal(t, expected, created)
}

func TestStreamerService_CreateConflict(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockStreamerRepository(t)
	spm := streaming.NewStreamingProviderManager(nil, zap.NewNop())
	svc := NewStreamerService(repo, spm, zap.NewNop())

	cmd := &command.CreateStreamerCommand{PlatformType: string(domain.StreamingPlatformTypeBilibili), PlatformStreamerID: "123"}

	repo.EXPECT().ExistByPlatformStreamerId(ctx, domain.StreamingPlatformType(cmd.PlatformType), cmd.PlatformStreamerID).Return(true, nil)

	_, err := svc.Create(ctx, cmd)
	require.Error(t, err)
}

func TestStreamerService_Update(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockStreamerRepository(t)
	spm := streaming.NewStreamingProviderManager(nil, zap.NewNop())
	svc := NewStreamerService(repo, spm, zap.NewNop())

	existing := &domain.Streamer{ID: 1, PlatformType: domain.StreamingPlatformTypeBilibili, PlatformStreamerID: "123"}
	cmd := &command.UpdateStreamerCommand{
		ID: existing.ID,
		CreateStreamerCommand: &command.CreateStreamerCommand{
			PlatformType:       string(domain.StreamingPlatformTypeBilibili),
			PlatformStreamerID: "123",
			DisplayName:        "New",
		},
	}
	expected := &domain.Streamer{ID: existing.ID, PlatformType: existing.PlatformType, PlatformStreamerID: existing.PlatformStreamerID}

	repo.EXPECT().FindById(ctx, cmd.ID).Return(existing, nil)
	repo.EXPECT().Update(ctx, mock.MatchedBy(func(streamer *domain.Streamer) bool {
		return streamer.ID == cmd.ID && streamer.DisplayName == cmd.DisplayName
	})).Return(expected, nil)

	got, err := svc.Update(ctx, cmd)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestStreamerService_UpdateConflict(t *testing.T) {
	ctx := context.Background()
	repo := repoMocks.NewMockStreamerRepository(t)
	spm := streaming.NewStreamingProviderManager(nil, zap.NewNop())
	svc := NewStreamerService(repo, spm, zap.NewNop())

	current := &domain.Streamer{ID: 1, PlatformType: domain.StreamingPlatformTypeBilibili, PlatformStreamerID: "123"}
	cmd := &command.UpdateStreamerCommand{
		ID: current.ID,
		CreateStreamerCommand: &command.CreateStreamerCommand{
			PlatformType:       string(domain.StreamingPlatformTypeBilibili),
			PlatformStreamerID: "456",
		},
	}

	repo.EXPECT().FindById(ctx, cmd.ID).Return(current, nil)
	repo.EXPECT().ExistByPlatformStreamerId(ctx, current.PlatformType, cmd.PlatformStreamerID).Return(true, nil)

	_, err := svc.Update(ctx, cmd)
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
