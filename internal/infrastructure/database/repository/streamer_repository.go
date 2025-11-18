package repository

import (
	"context"
	"slices"

	"github.com/ryuyb/fusion/internal/core/domain"
	coreRepo "github.com/ryuyb/fusion/internal/core/port/repository"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent/streamer"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type streamerRepository struct {
	client *ent.Client
	logger *zap.Logger
}

func NewStreamerRepository(client *ent.Client, logger *zap.Logger) coreRepo.StreamerRepository {
	return &streamerRepository{
		client: client,
		logger: logger,
	}
}

func (r *streamerRepository) Create(ctx context.Context, entity *domain.Streamer) (*domain.Streamer, error) {
	builder := r.client.Streamer.Create().
		SetPlatformType(string(entity.PlatformType)).
		SetPlatformStreamerID(entity.PlatformStreamerID).
		SetDisplayName(entity.DisplayName)

	if entity.AvatarURL != "" {
		builder.SetAvatarURL(entity.AvatarURL)
	}
	if entity.RoomURL != "" {
		builder.SetRoomURL(entity.RoomURL)
	}
	if entity.Bio != "" {
		builder.SetBio(entity.Bio)
	}
	if len(entity.Tags) > 0 {
		builder.SetTags(entity.Tags)
	}
	builder.SetIsLive(entity.LiveStatus.IsLive)
	if entity.LiveStatus.Title != "" {
		builder.SetLiveTitle(entity.LiveStatus.Title)
	}
	if entity.LiveStatus.GameName != "" {
		builder.SetLiveGameName(entity.LiveStatus.GameName)
	}
	if !entity.LiveStatus.StartTime.IsZero() {
		builder.SetLiveStartTime(entity.LiveStatus.StartTime)
	}
	builder.SetLiveViewers(entity.LiveStatus.Viewers)
	if entity.LiveStatus.CoverImage != "" {
		builder.SetLiveCoverImage(entity.LiveStatus.CoverImage)
	}
	if !entity.LastLiveSyncedAt.IsZero() {
		builder.SetLastLiveSyncedAt(entity.LastLiveSyncedAt)
	}
	if !entity.LastSyncedAt.IsZero() {
		builder.SetLastSyncedAt(entity.LastSyncedAt)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		r.logger.Error("failed to create streamer",
			zap.Error(err),
			zap.String("platform_type", string(entity.PlatformType)),
			zap.String("platform_streamer_id", entity.PlatformStreamerID),
		)
		return nil, errors2.ConvertDatabaseError(err, "Streamer")
	}
	return r.toDomain(created), nil
}

func (r *streamerRepository) Update(ctx context.Context, entity *domain.Streamer) (*domain.Streamer, error) {
	builder := r.client.Streamer.UpdateOneID(entity.ID).
		SetPlatformType(string(entity.PlatformType)).
		SetPlatformStreamerID(entity.PlatformStreamerID).
		SetDisplayName(entity.DisplayName)

	if entity.AvatarURL == "" {
		builder.ClearAvatarURL()
	} else {
		builder.SetAvatarURL(entity.AvatarURL)
	}
	if entity.RoomURL == "" {
		builder.ClearRoomURL()
	} else {
		builder.SetRoomURL(entity.RoomURL)
	}
	if entity.Bio == "" {
		builder.ClearBio()
	} else {
		builder.SetBio(entity.Bio)
	}
	builder.ClearTags()
	if len(entity.Tags) > 0 {
		builder.SetTags(entity.Tags)
	}
	builder.SetIsLive(entity.LiveStatus.IsLive)
	if entity.LiveStatus.Title == "" {
		builder.ClearLiveTitle()
	} else {
		builder.SetLiveTitle(entity.LiveStatus.Title)
	}
	if entity.LiveStatus.GameName == "" {
		builder.ClearLiveGameName()
	} else {
		builder.SetLiveGameName(entity.LiveStatus.GameName)
	}
	if entity.LiveStatus.StartTime.IsZero() {
		builder.ClearLiveStartTime()
	} else {
		builder.SetLiveStartTime(entity.LiveStatus.StartTime)
	}
	builder.SetLiveViewers(entity.LiveStatus.Viewers)
	if entity.LiveStatus.CoverImage == "" {
		builder.ClearLiveCoverImage()
	} else {
		builder.SetLiveCoverImage(entity.LiveStatus.CoverImage)
	}
	if entity.LastLiveSyncedAt.IsZero() {
		builder.ClearLastLiveSyncedAt()
	} else {
		builder.SetLastLiveSyncedAt(entity.LastLiveSyncedAt)
	}
	if entity.LastSyncedAt.IsZero() {
		builder.ClearLastSyncedAt()
	} else {
		builder.SetLastSyncedAt(entity.LastSyncedAt)
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors2.NotFound("Streamer").WithDetail("id", entity.ID)
		}
		r.logger.Error("failed to update streamer", zap.Error(err), zap.Int64("id", entity.ID))
		return nil, errors2.ConvertDatabaseError(err, "Streamer")
	}
	return r.toDomain(updated), nil
}

func (r *streamerRepository) Delete(ctx context.Context, id int64) error {
	err := r.client.Streamer.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors2.NotFound("Streamer").WithDetail("id", id)
		}
		r.logger.Error("failed to delete streamer", zap.Error(err), zap.Int64("id", id))
		return errors2.ConvertDatabaseError(err, "Streamer")
	}
	return nil
}

func (r *streamerRepository) FindById(ctx context.Context, id int64) (*domain.Streamer, error) {
	entity, err := r.client.Streamer.
		Query().
		Where(streamer.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors2.NotFound("Streamer").WithDetail("id", id)
		}
		r.logger.Error("failed to find streamer by id", zap.Error(err), zap.Int64("id", id))
		return nil, errors2.ConvertDatabaseError(err, "Streamer")
	}
	return r.toDomain(entity), nil
}

func (r *streamerRepository) FindByPlatformStreamerId(ctx context.Context, platformType domain.StreamingPlatformType, platformStreamerID string) (*domain.Streamer, error) {
	entity, err := r.client.Streamer.
		Query().
		Where(
			streamer.PlatformTypeEQ(string(platformType)),
			streamer.PlatformStreamerIDEQ(platformStreamerID),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors2.NotFound("Streamer").
				WithDetail("platform_type", platformType).
				WithDetail("platform_streamer_id", platformStreamerID)
		}
		r.logger.Error("failed to find streamer by platform id",
			zap.Error(err),
			zap.String("platform_type", string(platformType)),
			zap.String("platform_streamer_id", platformStreamerID),
		)
		return nil, errors2.ConvertDatabaseError(err, "Streamer")
	}
	return r.toDomain(entity), nil
}

func (r *streamerRepository) ExistByPlatformStreamerId(ctx context.Context, platformType domain.StreamingPlatformType, platformStreamerID string) (bool, error) {
	exist, err := r.client.Streamer.
		Query().
		Where(
			streamer.PlatformTypeEQ(string(platformType)),
			streamer.PlatformStreamerIDEQ(platformStreamerID),
		).
		Exist(ctx)
	if err != nil {
		r.logger.Error("failed to check streamer exist",
			zap.Error(err),
			zap.String("platform_type", string(platformType)),
			zap.String("platform_streamer_id", platformStreamerID),
		)
		return false, errors2.ConvertDatabaseError(err, "Streamer")
	}
	return exist, nil
}

func (r *streamerRepository) List(ctx context.Context, offset, limit int) ([]*domain.Streamer, int, error) {
	total, err := r.client.Streamer.Query().Count(ctx)
	if err != nil {
		r.logger.Error("failed to count streamers", zap.Error(err))
		return nil, 0, errors2.DatabaseError(err)
	}

	entities, err := r.client.Streamer.
		Query().
		Offset(offset).
		Limit(limit).
		Order(ent.Desc(streamer.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		r.logger.Error("failed to list streamers", zap.Error(err))
		return nil, 0, errors2.DatabaseError(err)
	}

	results := make([]*domain.Streamer, len(entities))
	for i, entity := range entities {
		results[i] = r.toDomain(entity)
	}
	return results, total, nil
}

func (r *streamerRepository) toDomain(entity *ent.Streamer) *domain.Streamer {
	return &domain.Streamer{
		ID:                 entity.ID,
		PlatformType:       domain.StreamingPlatformType(entity.PlatformType),
		PlatformStreamerID: entity.PlatformStreamerID,
		DisplayName:        entity.DisplayName,
		AvatarURL:          lo.FromPtr(entity.AvatarURL),
		RoomURL:            lo.FromPtr(entity.RoomURL),
		Bio:                lo.FromPtr(entity.Bio),
		Tags:               slices.Clone(entity.Tags),
		LiveStatus: domain.LiveStatusInfo{
			IsLive:     entity.IsLive,
			Title:      lo.FromPtr(entity.LiveTitle),
			GameName:   lo.FromPtr(entity.LiveGameName),
			StartTime:  lo.FromPtr(entity.LiveStartTime),
			Viewers:    lo.FromPtr(entity.LiveViewers),
			CoverImage: lo.FromPtr(entity.LiveCoverImage),
		},
		LastLiveSyncedAt: lo.FromPtr(entity.LastLiveSyncedAt),
		LastSyncedAt:     lo.FromPtr(entity.LastSyncedAt),
		CreatedAt:        entity.CreatedAt,
		UpdatedAt:        entity.UpdatedAt,
	}
}
