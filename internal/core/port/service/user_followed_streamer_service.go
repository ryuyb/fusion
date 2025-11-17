package service

import (
	"context"

	"github.com/ryuyb/fusion/internal/core/command"
	"github.com/ryuyb/fusion/internal/core/domain"
)

type UserFollowedStreamerService interface {
	Create(ctx context.Context, cmd *command.CreateUserFollowedStreamerCommand) (*domain.UserFollowedStreamer, error)

	Update(ctx context.Context, cmd *command.UpdateUserFollowedStreamerCommand) (*domain.UserFollowedStreamer, error)

	Delete(ctx context.Context, id int64) error

	FindById(ctx context.Context, id int64) (*domain.UserFollowedStreamer, error)

	FindByUserAndStreamer(ctx context.Context, userID, streamerID int64) (*domain.UserFollowedStreamer, error)

	ListByUserId(ctx context.Context, userID int64, page, pageSize int) ([]*domain.UserFollowedStreamer, int, error)

	ListByStreamerId(ctx context.Context, streamerID int64, page, pageSize int) ([]*domain.UserFollowedStreamer, int, error)
}
