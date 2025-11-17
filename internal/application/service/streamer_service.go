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

type streamerService struct {
	repo   coreRepo.StreamerRepository
	logger *zap.Logger
}

func NewStreamerService(repo coreRepo.StreamerRepository, logger *zap.Logger) coreService.StreamerService {
	return &streamerService{
		repo:   repo,
		logger: logger,
	}
}

func (s *streamerService) Create(ctx context.Context, streamer *domain.Streamer) (*domain.Streamer, error) {
	exist, err := s.repo.ExistByPlatformStreamerId(ctx, streamer.PlatformType, streamer.PlatformStreamerID)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errors.Conflict("streamer already exists")
	}
	return s.repo.Create(ctx, streamer)
}

func (s *streamerService) Update(ctx context.Context, streamer *domain.Streamer) (*domain.Streamer, error) {
	current, err := s.repo.FindById(ctx, streamer.ID)
	if err != nil {
		return nil, err
	}
	if current.PlatformType != streamer.PlatformType || current.PlatformStreamerID != streamer.PlatformStreamerID {
		exist, err := s.repo.ExistByPlatformStreamerId(ctx, streamer.PlatformType, streamer.PlatformStreamerID)
		if err != nil {
			return nil, err
		}
		if exist {
			return nil, errors.Conflict("streamer already exists")
		}
	}
	return s.repo.Update(ctx, streamer)
}

func (s *streamerService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *streamerService) FindById(ctx context.Context, id int64) (*domain.Streamer, error) {
	return s.repo.FindById(ctx, id)
}

func (s *streamerService) FindByPlatformStreamerId(ctx context.Context, platformType domain.StreamingPlatformType, platformStreamerID string) (*domain.Streamer, error) {
	return s.repo.FindByPlatformStreamerId(ctx, platformType, platformStreamerID)
}

func (s *streamerService) List(ctx context.Context, page, pageSize int) ([]*domain.Streamer, int, error) {
	if err := util.ValidatePagination(page, pageSize); err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	return s.repo.List(ctx, offset, pageSize)
}
