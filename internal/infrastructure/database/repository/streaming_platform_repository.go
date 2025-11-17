package repository

import (
	"context"

	"github.com/ryuyb/fusion/internal/core/domain"
	coreRepo "github.com/ryuyb/fusion/internal/core/port/repository"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent/streamingplatform"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type streamingPlatformRepository struct {
	client *ent.Client
	logger *zap.Logger
}

func NewStreamingPlatformRepository(client *ent.Client, logger *zap.Logger) coreRepo.StreamingPlatformRepository {
	return &streamingPlatformRepository{
		client: client,
		logger: logger,
	}
}

func (r *streamingPlatformRepository) Create(ctx context.Context, platform *domain.StreamingPlatform) (*domain.StreamingPlatform, error) {
	builder := r.client.StreamingPlatform.Create().
		SetType(string(platform.Type)).
		SetName(platform.Name).
		SetBaseURL(platform.BaseURL).
		SetEnabled(platform.Enabled).
		SetPriority(platform.Priority)

	if platform.Description != "" {
		builder.SetDescription(platform.Description)
	}
	if platform.LogoURL != "" {
		builder.SetLogoURL(platform.LogoURL)
	}
	if platform.Metadata != nil {
		builder.SetMetadata(platform.Metadata)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		r.logger.Error("failed to create streaming platform",
			zap.Error(err),
			zap.String("platform_type", string(platform.Type)),
			zap.String("name", platform.Name),
		)
		return nil, errors2.ConvertDatabaseError(err, "StreamingPlatform")
	}
	return r.toDomain(created), nil
}

func (r *streamingPlatformRepository) Update(ctx context.Context, platform *domain.StreamingPlatform) (*domain.StreamingPlatform, error) {
	builder := r.client.StreamingPlatform.UpdateOneID(platform.ID).
		SetType(string(platform.Type)).
		SetName(platform.Name).
		SetBaseURL(platform.BaseURL).
		SetEnabled(platform.Enabled).
		SetPriority(platform.Priority)

	if platform.Description == "" {
		builder.ClearDescription()
	} else {
		builder.SetDescription(platform.Description)
	}
	if platform.LogoURL == "" {
		builder.ClearLogoURL()
	} else {
		builder.SetLogoURL(platform.LogoURL)
	}
	if platform.Metadata == nil {
		builder.ClearMetadata()
	} else {
		builder.SetMetadata(platform.Metadata)
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors2.NotFound("StreamingPlatform").WithDetail("id", platform.ID)
		}
		r.logger.Error("failed to update streaming platform",
			zap.Error(err),
			zap.Int64("id", platform.ID),
		)
		return nil, errors2.ConvertDatabaseError(err, "StreamingPlatform")
	}
	return r.toDomain(updated), nil
}

func (r *streamingPlatformRepository) Delete(ctx context.Context, id int64) error {
	err := r.client.StreamingPlatform.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors2.NotFound("StreamingPlatform").WithDetail("id", id)
		}
		r.logger.Error("failed to delete streaming platform", zap.Error(err), zap.Int64("id", id))
		return errors2.ConvertDatabaseError(err, "StreamingPlatform")
	}
	return nil
}

func (r *streamingPlatformRepository) FindById(ctx context.Context, id int64) (*domain.StreamingPlatform, error) {
	entity, err := r.client.StreamingPlatform.
		Query().
		Where(streamingplatform.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors2.NotFound("StreamingPlatform").WithDetail("id", id)
		}
		r.logger.Error("failed to find streaming platform by id", zap.Error(err), zap.Int64("id", id))
		return nil, errors2.ConvertDatabaseError(err, "StreamingPlatform")
	}
	return r.toDomain(entity), nil
}

func (r *streamingPlatformRepository) FindByType(ctx context.Context, platformType domain.StreamingPlatformType) (*domain.StreamingPlatform, error) {
	entity, err := r.client.StreamingPlatform.
		Query().
		Where(streamingplatform.TypeEQ(string(platformType))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors2.NotFound("StreamingPlatform").WithDetail("type", platformType)
		}
		r.logger.Error("failed to find streaming platform by type", zap.Error(err), zap.String("type", string(platformType)))
		return nil, errors2.ConvertDatabaseError(err, "StreamingPlatform")
	}
	return r.toDomain(entity), nil
}

func (r *streamingPlatformRepository) FindByName(ctx context.Context, name string) (*domain.StreamingPlatform, error) {
	entity, err := r.client.StreamingPlatform.
		Query().
		Where(streamingplatform.NameEQ(name)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors2.NotFound("StreamingPlatform").WithDetail("name", name)
		}
		r.logger.Error("failed to find streaming platform by name", zap.Error(err), zap.String("name", name))
		return nil, errors2.ConvertDatabaseError(err, "StreamingPlatform")
	}
	return r.toDomain(entity), nil
}

func (r *streamingPlatformRepository) ExistByType(ctx context.Context, platformType domain.StreamingPlatformType) (bool, error) {
	exist, err := r.client.StreamingPlatform.
		Query().
		Where(streamingplatform.TypeEQ(string(platformType))).
		Exist(ctx)
	if err != nil {
		r.logger.Error("failed to check streaming platform by type", zap.Error(err), zap.String("type", string(platformType)))
		return false, errors2.ConvertDatabaseError(err, "StreamingPlatform")
	}
	return exist, nil
}

func (r *streamingPlatformRepository) ExistByName(ctx context.Context, name string) (bool, error) {
	exist, err := r.client.StreamingPlatform.
		Query().
		Where(streamingplatform.NameEQ(name)).
		Exist(ctx)
	if err != nil {
		r.logger.Error("failed to check streaming platform by name", zap.Error(err), zap.String("name", name))
		return false, errors2.ConvertDatabaseError(err, "StreamingPlatform")
	}
	return exist, nil
}

func (r *streamingPlatformRepository) List(ctx context.Context, offset, limit int) ([]*domain.StreamingPlatform, int, error) {
	total, err := r.client.StreamingPlatform.Query().Count(ctx)
	if err != nil {
		r.logger.Error("failed to count streaming platforms", zap.Error(err))
		return nil, 0, errors2.DatabaseError(err)
	}

	query := r.client.StreamingPlatform.
		Query().
		Offset(offset).
		Limit(limit).
		Order(ent.Desc(streamingplatform.FieldCreatedAt))

	entities, err := query.All(ctx)
	if err != nil {
		r.logger.Error("failed to list streaming platforms", zap.Error(err))
		return nil, 0, errors2.DatabaseError(err)
	}

	results := make([]*domain.StreamingPlatform, len(entities))
	for i, entity := range entities {
		results[i] = r.toDomain(entity)
	}
	return results, total, nil
}

func (r *streamingPlatformRepository) toDomain(entity *ent.StreamingPlatform) *domain.StreamingPlatform {
	var metadata map[string]string
	if entity.Metadata != nil {
		metadata = lo.Assign(map[string]string{}, entity.Metadata)
	}

	return &domain.StreamingPlatform{
		ID:          entity.ID,
		Type:        domain.StreamingPlatformType(entity.Type),
		Name:        entity.Name,
		Description: lo.FromPtr(entity.Description),
		BaseURL:     entity.BaseURL,
		LogoURL:     lo.FromPtr(entity.LogoURL),
		Enabled:     entity.Enabled,
		Priority:    entity.Priority,
		Metadata:    metadata,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
