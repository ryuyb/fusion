package service

import (
	"context"

	"github.com/ryuyb/fusion/internal/core/command"
	"github.com/ryuyb/fusion/internal/core/domain"
)

type StreamerService interface {
	Create(ctx context.Context, cmd *command.CreateStreamerCommand) (*domain.Streamer, error)

	Update(ctx context.Context, cmd *command.UpdateStreamerCommand) (*domain.Streamer, error)

	Delete(ctx context.Context, id int64) error

	FindById(ctx context.Context, id int64) (*domain.Streamer, error)

	FindByPlatformStreamerId(ctx context.Context, platformType domain.StreamingPlatformType, platformStreamerID string, refresh bool) (*domain.Streamer, error)

	List(ctx context.Context, page, pageSize int) ([]*domain.Streamer, int, error)
}
