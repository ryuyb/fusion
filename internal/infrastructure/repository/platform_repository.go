package repository

import (
	"context"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/domain/repository"
	"github.com/ryuyb/fusion/internal/infrastructure/database"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent/platform"
	"github.com/ryuyb/fusion/internal/pkg/errors"
	"go.uber.org/zap"
)

type platformRepository struct {
	client *database.Client
	logger *zap.Logger
}

func NewPlatformRepository(client *database.Client, logger *zap.Logger) repository.PlatformRepository {
	return &platformRepository{
		client: client,
		logger: logger,
	}
}

func (r *platformRepository) Create(ctx context.Context, p *entity.Platform) (*entity.Platform, error) {
	created, err := r.client.Platform.
		Create().
		SetName(p.Name).
		SetPlatformType(platform.PlatformType(p.PlatformType)).
		SetConfig(p.Config).
		SetStatus(platform.Status(p.Status)).
		SetPollInterval(p.PollInterval).
		Save(ctx)
	if err != nil {
		r.logger.Error("failed to create platform", zap.Error(err), zap.String("name", p.Name), zap.String("type", string(p.PlatformType)))
		return nil, errors.ConvertDatabaseError(err, "Platform")
	}
	return r.toEntity(created), nil
}

func (r *platformRepository) Update(ctx context.Context, p *entity.Platform) (*entity.Platform, error) {
	updated, err := r.client.Platform.
		UpdateOneID(p.ID).
		SetName(p.Name).
		SetConfig(p.Config).
		SetStatus(platform.Status(p.Status)).
		SetPollInterval(p.PollInterval).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("Platform").WithDetail("id", p.ID)
		}
		r.logger.Error("failed to update platform", zap.Error(err), zap.Int64("id", p.ID))
		return nil, errors.ConvertDatabaseError(err, "Platform")
	}
	return r.toEntity(updated), nil
}

func (r *platformRepository) FindByID(ctx context.Context, id int64) (*entity.Platform, error) {
	p, err := r.client.Platform.
		Query().
		Where(platform.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("Platform").WithDetail("id", id)
		}
		r.logger.Error("failed to find platform by id", zap.Error(err), zap.Int64("id", id))
		return nil, errors.ConvertDatabaseError(err, "Platform")
	}
	return r.toEntity(p), nil
}

func (r *platformRepository) FindByType(ctx context.Context, platformType entity.PlatformType) (*entity.Platform, error) {
	p, err := r.client.Platform.
		Query().
		Where(platform.PlatformTypeEQ(platform.PlatformType(platformType))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("Platform").WithDetail("platform_type", platformType)
		}
		r.logger.Error("failed to find platform by type", zap.Error(err), zap.String("platform_type", string(platformType)))
		return nil, errors.ConvertDatabaseError(err, "Platform")
	}
	return r.toEntity(p), nil
}

func (r *platformRepository) List(ctx context.Context) ([]*entity.Platform, error) {
	platforms, err := r.client.Platform.
		Query().
		Order(ent.Desc(platform.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		r.logger.Error("failed to list platforms", zap.Error(err))
		return nil, errors.DatabaseError(err)
	}

	entities := make([]*entity.Platform, len(platforms))
	for i, p := range platforms {
		entities[i] = r.toEntity(p)
	}

	return entities, nil
}

func (r *platformRepository) Delete(ctx context.Context, id int64) error {
	err := r.client.Platform.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.NotFound("Platform").WithDetail("id", id)
		}
		r.logger.Error("failed to delete platform", zap.Error(err), zap.Int64("id", id))
		return errors.ConvertDatabaseError(err, "Platform")
	}
	return nil
}

// toEntity converts an ent.Platform to entity.Platform
func (r *platformRepository) toEntity(p *ent.Platform) *entity.Platform {
	return &entity.Platform{
		ID:           p.ID,
		Name:         p.Name,
		PlatformType: entity.PlatformType(p.PlatformType),
		Config:       p.Config,
		Status:       entity.PlatformStatus(p.Status),
		PollInterval: p.PollInterval,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
		DeleteAt:     p.DeleteAt,
	}
}
