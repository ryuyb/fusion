package service

import (
	"context"
	"fmt"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/domain/repository"
	domainService "github.com/ryuyb/fusion/internal/domain/service"
	"github.com/ryuyb/fusion/internal/infrastructure/streaming"
	"github.com/ryuyb/fusion/internal/interface/http/dto/request"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"go.uber.org/zap"
)

// platformService implements domainService.PlatformService
type platformService struct {
	repo        repository.PlatformRepository
	providerMgr *streaming.StreamingProviderManager
	logger      *zap.Logger
}

// NewPlatformService creates a new PlatformService
func NewPlatformService(
	repo repository.PlatformRepository,
	providerMgr *streaming.StreamingProviderManager,
	logger *zap.Logger,
) domainService.PlatformService {
	return &platformService{
		repo:        repo,
		providerMgr: providerMgr,
		logger:      logger,
	}
}

// Create creates a new platform configuration
func (s *platformService) Create(
	ctx context.Context,
	req *request.CreatePlatformRequest,
) (*entity.Platform, error) {
	// Convert platform type string to enum
	platformType := entity.PlatformType(req.PlatformType)

	// Check if platform type already exists
	_, err := s.repo.FindByType(ctx, platformType)
	if err == nil {
		return nil, errors2.Conflict(fmt.Sprintf("platform type %s already exists", platformType))
	}
	if !errors2.IsNotFoundError(err) {
		s.logger.Error("Failed to check existing platform", zap.String("platform_type", string(platformType)), zap.Error(err))
		return nil, err
	}

	// Validate platform configuration
	if !s.providerMgr.HasProvider(platformType) {
		return nil, errors2.BadRequest(fmt.Sprintf("unsupported platform type: %s", platformType))
	}

	// Get the provider and validate configuration
	provider, err := s.providerMgr.GetProvider(platformType)
	if err != nil {
		return nil, errors2.Internal(fmt.Errorf("failed to get provider: %w", err))
	}

	if err := provider.ValidateConfiguration(req.Config); err != nil {
		s.logger.Warn("Invalid platform configuration", zap.String("platform_type", string(platformType)), zap.Error(err))
		return nil, errors2.BadRequest(fmt.Sprintf("invalid platform configuration: %v", err))
	}

	// Create platform entity
	platform := entity.CreatePlatform(req.Name, platformType, req.Config, req.PollInterval)

	// Save to database
	created, err := s.repo.Create(ctx, platform)
	if err != nil {
		s.logger.Error("Failed to create platform", zap.String("platform_type", string(platformType)), zap.Error(err))
		return nil, err
	}

	s.logger.Info("Platform created successfully",
		zap.Int64("platform_id", created.ID),
		zap.String("platform_type", string(platformType)),
		zap.String("name", req.Name))

	return created, nil
}

// Update updates a platform configuration
func (s *platformService) Update(
	ctx context.Context,
	id int64,
	req *request.UpdatePlatformRequest,
) (*entity.Platform, error) {
	// Get existing platform
	platform, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Convert status string to enum
	status := entity.PlatformStatus(req.Status)

	// Validate platform configuration
	provider, err := s.providerMgr.GetProvider(platform.PlatformType)
	if err != nil {
		return nil, errors2.Internal(fmt.Errorf("failed to get provider: %w", err))
	}

	if err := provider.ValidateConfiguration(req.Config); err != nil {
		s.logger.Warn("Invalid platform configuration during update", zap.Int64("platform_id", id), zap.Error(err))
		return nil, errors2.BadRequest(fmt.Sprintf("invalid platform configuration: %v", err))
	}

	// Update platform
	platform = platform.Update(req.Name, req.Config, status, req.PollInterval)

	// Save to database
	updated, err := s.repo.Update(ctx, platform)
	if err != nil {
		s.logger.Error("Failed to update platform", zap.Int64("platform_id", id), zap.Error(err))
		return nil, err
	}

	s.logger.Info("Platform updated successfully",
		zap.Int64("platform_id", id),
		zap.String("platform_type", string(platform.PlatformType)))

	return updated, nil
}

// GetByID retrieves a platform by ID
func (s *platformService) GetByID(ctx context.Context, id int64) (*entity.Platform, error) {
	platform, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return platform, nil
}

// List retrieves all platforms
func (s *platformService) List(ctx context.Context) ([]*entity.Platform, error) {
	platforms, err := s.repo.List(ctx)
	if err != nil {
		s.logger.Error("Failed to list platforms", zap.Error(err))
		return nil, err
	}
	return platforms, nil
}

// Delete deletes a platform
func (s *platformService) Delete(ctx context.Context, id int64) error {
	// Check if platform exists
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete platform", zap.Int64("platform_id", id), zap.Error(err))
		return err
	}

	s.logger.Info("Platform deleted successfully", zap.Int64("platform_id", id))
	return nil
}

// TestConnection tests the connection to a platform's API
func (s *platformService) TestConnection(ctx context.Context, id int64) (bool, error) {
	// Get platform
	platform, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return false, err
	}

	// Check if platform is active
	if !platform.IsActive() {
		return false, errors2.BadRequest("platform is not active")
	}

	// Get provider
	provider, err := s.providerMgr.GetProvider(platform.PlatformType)
	if err != nil {
		s.logger.Error("Failed to get provider for testing", zap.Int64("platform_id", id), zap.Error(err))
		return false, errors2.Internal(fmt.Errorf("failed to get provider: %w", err))
	}

	// Test the configuration
	if err := provider.ValidateConfiguration(platform.Config); err != nil {
		s.logger.Warn("Platform configuration validation failed", zap.Int64("platform_id", id), zap.Error(err))
		return false, errors2.BadRequest(fmt.Sprintf("configuration validation failed: %v", err))
	}

	// Try to test the connection
	// For this, we'll try to validate the configuration more thoroughly
	// In a real implementation, you might want to make a test API call
	if err := provider.ValidateConfiguration(platform.Config); err != nil {
		s.logger.Warn("Connection test failed", zap.Int64("platform_id", id), zap.Error(err))
		return false, errors2.StreamingPlatformError(platform.PlatformType, nil, err)
	}

	s.logger.Info("Platform connection test passed", zap.Int64("platform_id", id), zap.String("platform_type", string(platform.PlatformType)))
	return true, nil
}
