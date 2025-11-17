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

type streamingPlatformService struct {
	repo   coreRepo.StreamingPlatformRepository
	logger *zap.Logger
}

func NewStreamingPlatformService(repo coreRepo.StreamingPlatformRepository, logger *zap.Logger) coreService.StreamingPlatformService {
	return &streamingPlatformService{
		repo:   repo,
		logger: logger,
	}
}

func (s *streamingPlatformService) Create(ctx context.Context, platform *domain.StreamingPlatform) (*domain.StreamingPlatform, error) {
	exist, err := s.repo.ExistByName(ctx, platform.Name)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errors.Conflict("streaming platform already exists")
	}

	return s.repo.Create(ctx, platform)
}

func (s *streamingPlatformService) Update(ctx context.Context, platform *domain.StreamingPlatform) (*domain.StreamingPlatform, error) {
	current, err := s.repo.FindById(ctx, platform.ID)
	if err != nil {
		return nil, err
	}
	if platform.Name != current.Name {
		exist, err := s.repo.ExistByName(ctx, platform.Name)
		if err != nil {
			return nil, err
		}
		if exist {
			return nil, errors.Conflict("streaming platform already exists")
		}
	}

	return s.repo.Update(ctx, platform)
}

func (s *streamingPlatformService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *streamingPlatformService) FindById(ctx context.Context, id int64) (*domain.StreamingPlatform, error) {
	return s.repo.FindById(ctx, id)
}

func (s *streamingPlatformService) FindByType(ctx context.Context, platformType domain.StreamingPlatformType) (*domain.StreamingPlatform, error) {
	return s.repo.FindByType(ctx, platformType)
}

func (s *streamingPlatformService) List(ctx context.Context, page, pageSize int) ([]*domain.StreamingPlatform, int, error) {
	if err := util.ValidatePagination(page, pageSize); err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	return s.repo.List(ctx, offset, pageSize)
}
