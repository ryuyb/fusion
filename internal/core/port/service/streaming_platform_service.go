package service

import (
	"context"

	"github.com/ryuyb/fusion/internal/core/command"
	"github.com/ryuyb/fusion/internal/core/domain"
)

type StreamingPlatformService interface {
	Create(ctx context.Context, cmd *command.CreateStreamingPlatformCommand) (*domain.StreamingPlatform, error)

	Update(ctx context.Context, cmd *command.UpdateStreamingPlatformCommand) (*domain.StreamingPlatform, error)

	Delete(ctx context.Context, id int64) error

	FindById(ctx context.Context, id int64) (*domain.StreamingPlatform, error)

	FindByType(ctx context.Context, platformType domain.StreamingPlatformType) (*domain.StreamingPlatform, error)

	List(ctx context.Context, page, pageSize int) ([]*domain.StreamingPlatform, int, error)
}
