package service

import (
	"context"

	"github.com/ryuyb/fusion/internal/core/command"
	"github.com/ryuyb/fusion/internal/core/domain"
	coreRepo "github.com/ryuyb/fusion/internal/core/port/repository"
	coreService "github.com/ryuyb/fusion/internal/core/port/service"
	"github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/ryuyb/fusion/internal/pkg/util"
	"go.uber.org/zap"
)

type userFollowedStreamerService struct {
	repo   coreRepo.UserFollowedStreamerRepository
	logger *zap.Logger
}

func NewUserFollowedStreamerService(repo coreRepo.UserFollowedStreamerRepository, logger *zap.Logger) coreService.UserFollowedStreamerService {
	return &userFollowedStreamerService{
		repo:   repo,
		logger: logger,
	}
}

func (s *userFollowedStreamerService) Create(ctx context.Context, cmd *command.CreateUserFollowedStreamerCommand) (*domain.UserFollowedStreamer, error) {
	exist, err := s.repo.ExistByUserAndStreamer(ctx, cmd.UserID, cmd.StreamerID)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errors.Conflict("user already follows this streamer")
	}
	follow, err := domain.NewUserFollowedStreamer(cmd.UserID, cmd.StreamerID, cmd.Alias, cmd.Notes, cmd.NotificationChannelIDs)
	if err != nil {
		return nil, err
	}
	follow.NotificationsEnabled = cmd.NotificationsEnabled
	return s.repo.Create(ctx, follow)
}

func (s *userFollowedStreamerService) Update(ctx context.Context, cmd *command.UpdateUserFollowedStreamerCommand) (*domain.UserFollowedStreamer, error) {
	current, err := s.repo.FindById(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if err := current.UpdatePreferences(cmd.Alias, cmd.Notes, cmd.NotificationsEnabled, cmd.NotificationChannelIDs); err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, current)
}

func (s *userFollowedStreamerService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *userFollowedStreamerService) FindById(ctx context.Context, id int64) (*domain.UserFollowedStreamer, error) {
	return s.repo.FindById(ctx, id)
}

func (s *userFollowedStreamerService) FindByUserAndStreamer(ctx context.Context, userID, streamerID int64) (*domain.UserFollowedStreamer, error) {
	return s.repo.FindByUserAndStreamer(ctx, userID, streamerID)
}

func (s *userFollowedStreamerService) ListByUserId(ctx context.Context, userID int64, page, pageSize int) ([]*domain.UserFollowedStreamer, int, error) {
	return s.list(ctx, page, pageSize, func(ctx context.Context, offset, limit int) ([]*domain.UserFollowedStreamer, int, error) {
		return s.repo.ListByUserId(ctx, userID, offset, limit)
	})
}

func (s *userFollowedStreamerService) ListByStreamerId(ctx context.Context, streamerID int64, page, pageSize int) ([]*domain.UserFollowedStreamer, int, error) {
	return s.list(ctx, page, pageSize, func(ctx context.Context, offset, limit int) ([]*domain.UserFollowedStreamer, int, error) {
		return s.repo.ListByStreamerId(ctx, streamerID, offset, limit)
	})
}

func (s *userFollowedStreamerService) list(ctx context.Context, page, pageSize int, fn func(ctx context.Context, offset, limit int) ([]*domain.UserFollowedStreamer, int, error)) ([]*domain.UserFollowedStreamer, int, error) {
	if err := util.ValidatePagination(page, pageSize); err != nil {
		s.logger.Warn("invalid pagination parameters for user follow",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.Error(err),
		)
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	return fn(ctx, offset, pageSize)
}
