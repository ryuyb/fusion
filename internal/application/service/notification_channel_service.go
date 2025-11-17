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

type notificationChannelService struct {
	repo   coreRepo.NotificationChannelRepository
	logger *zap.Logger
}

func NewNotificationChannelService(repo coreRepo.NotificationChannelRepository, logger *zap.Logger) coreService.NotificationChannelService {
	return &notificationChannelService{
		repo:   repo,
		logger: logger,
	}
}

func (s *notificationChannelService) Create(ctx context.Context, cmd *command.CreateNotificationChannelCommand) (*domain.NotificationChannel, error) {
	exist, err := s.repo.ExistByName(ctx, cmd.UserID, cmd.Name)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errors.Conflict("notification channel already exists")
	}
	channel, err := buildNotificationChannelFromCommand(cmd)
	if err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, channel)
}

func (s *notificationChannelService) Update(ctx context.Context, cmd *command.UpdateNotificationChannelCommand) (*domain.NotificationChannel, error) {
	current, err := s.repo.FindById(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if current.Name != cmd.Name {
		exist, err := s.repo.ExistByName(ctx, current.UserID, cmd.Name)
		if err != nil {
			return nil, err
		}
		if exist {
			return nil, errors.Conflict("notification channel already exists")
		}
	}
	channel, err := buildNotificationChannelFromCommand(cmd.CreateNotificationChannelCommand)
	if err != nil {
		return nil, err
	}
	channel.ID = cmd.ID
	return s.repo.Update(ctx, channel)
}

func (s *notificationChannelService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *notificationChannelService) FindById(ctx context.Context, id int64) (*domain.NotificationChannel, error) {
	return s.repo.FindById(ctx, id)
}

func (s *notificationChannelService) ListByUserId(ctx context.Context, userID int64, page, pageSize int) ([]*domain.NotificationChannel, int, error) {
	if err := util.ValidatePagination(page, pageSize); err != nil {
		s.logger.Warn("invalid pagination parameters for notification channel",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.Error(err),
		)
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	return s.repo.ListByUserId(ctx, userID, offset, pageSize)
}

func buildNotificationChannelFromCommand(cmd *command.CreateNotificationChannelCommand) (*domain.NotificationChannel, error) {
	if cmd == nil {
		return nil, errors.BadRequest("notification channel command is required")
	}

	var config map[string]any
	if cmd.Config != nil {
		config = make(map[string]any, len(cmd.Config))
		for k, v := range cmd.Config {
			config[k] = v
		}
	}

	return &domain.NotificationChannel{
		UserID:      cmd.UserID,
		ChannelType: domain.NotificationChannelType(cmd.ChannelType),
		Name:        cmd.Name,
		Config:      config,
		Enable:      cmd.Enable,
		Priority:    cmd.Priority,
	}, nil
}
