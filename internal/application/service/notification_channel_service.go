package service

import (
	"context"

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

func (s *notificationChannelService) Create(ctx context.Context, channel *domain.NotificationChannel) (*domain.NotificationChannel, error) {
	exist, err := s.repo.ExistByName(ctx, channel.UserID, channel.Name)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errors.Conflict("notification channel already exists")
	}
	return s.repo.Create(ctx, channel)
}

func (s *notificationChannelService) Update(ctx context.Context, channel *domain.NotificationChannel) (*domain.NotificationChannel, error) {
	current, err := s.repo.FindById(ctx, channel.ID)
	if err != nil {
		return nil, err
	}
	if current.Name != channel.Name {
		exist, err := s.repo.ExistByName(ctx, current.UserID, channel.Name)
		if err != nil {
			return nil, err
		}
		if exist {
			return nil, errors.Conflict("notification channel already exists")
		}
	}
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
