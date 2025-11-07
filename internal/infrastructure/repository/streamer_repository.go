package repository

import (
	"context"
	"time"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/domain/repository"
	"github.com/ryuyb/fusion/internal/infrastructure/database"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent/streamer"
	"github.com/ryuyb/fusion/internal/pkg/errors"
	"go.uber.org/zap"
)

type streamerRepository struct {
	client *database.Client
	logger *zap.Logger
}

func NewStreamerRepository(client *database.Client, logger *zap.Logger) repository.StreamerRepository {
	return &streamerRepository{
		client: client,
		logger: logger,
	}
}

func (r *streamerRepository) Create(ctx context.Context, streamer *entity.Streamer) (*entity.Streamer, error) {
	created, err := r.client.Streamer.
		Create().
		SetPlatformID(streamer.PlatformID).
		SetPlatformStreamerID(streamer.PlatformStreamerID).
		SetName(streamer.Name).
		SetAvatar(streamer.Avatar).
		SetDescription(streamer.Description).
		SetRoomURL(streamer.RoomURL).
		SetIsLive(streamer.IsLive).
		Save(ctx)
	if err != nil {
		r.logger.Error("failed to create streamer", zap.Error(err), zap.Int64("platform_id", streamer.PlatformID), zap.String("platform_streamer_id", streamer.PlatformStreamerID))
		return nil, errors.ConvertDatabaseError(err, "Streamer")
	}
	return r.toEntity(created), nil
}

func (r *streamerRepository) Update(ctx context.Context, streamer *entity.Streamer) (*entity.Streamer, error) {
	updated, err := r.client.Streamer.
		UpdateOneID(streamer.ID).
		SetName(streamer.Name).
		SetAvatar(streamer.Avatar).
		SetDescription(streamer.Description).
		SetRoomURL(streamer.RoomURL).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("Streamer").WithDetail("id", streamer.ID)
		}
		r.logger.Error("failed to update streamer", zap.Error(err), zap.Int64("id", streamer.ID))
		return nil, errors.ConvertDatabaseError(err, "Streamer")
	}
	return r.toEntity(updated), nil
}

func (r *streamerRepository) FindByID(ctx context.Context, id int64) (*entity.Streamer, error) {
	s, err := r.client.Streamer.
		Query().
		Where(streamer.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("Streamer").WithDetail("id", id)
		}
		r.logger.Error("failed to find streamer by id", zap.Error(err), zap.Int64("id", id))
		return nil, errors.ConvertDatabaseError(err, "Streamer")
	}
	return r.toEntity(s), nil
}

func (r *streamerRepository) FindByPlatformAndStreamerID(ctx context.Context, platformID int64, platformStreamerID string) (*entity.Streamer, error) {
	s, err := r.client.Streamer.
		Query().
		Where(
			streamer.PlatformID(platformID),
			streamer.PlatformStreamerID(platformStreamerID),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("Streamer").WithDetail("platform_id", platformID).WithDetail("platform_streamer_id", platformStreamerID)
		}
		r.logger.Error("failed to find streamer by platform and platform_streamer_id", zap.Error(err), zap.Int64("platform_id", platformID), zap.String("platform_streamer_id", platformStreamerID))
		return nil, errors.ConvertDatabaseError(err, "Streamer")
	}
	return r.toEntity(s), nil
}

func (r *streamerRepository) FindByPlatform(ctx context.Context, platformID int64) ([]*entity.Streamer, error) {
	streamers, err := r.client.Streamer.
		Query().
		Where(streamer.PlatformID(platformID)).
		Order(ent.Desc(streamer.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		r.logger.Error("failed to find streamers by platform", zap.Error(err), zap.Int64("platform_id", platformID))
		return nil, errors.DatabaseError(err)
	}

	entities := make([]*entity.Streamer, len(streamers))
	for i, s := range streamers {
		entities[i] = r.toEntity(s)
	}

	return entities, nil
}

func (r *streamerRepository) FindAllWithFollowers(ctx context.Context) ([]*entity.Streamer, error) {
	// Query streamers that have at least one following (follower)
	streamers, err := r.client.Streamer.
		Query().
		Where(streamer.HasFollowings()).
		Order(ent.Desc(streamer.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		r.logger.Error("failed to find streamers with followers", zap.Error(err))
		return nil, errors.DatabaseError(err)
	}

	entities := make([]*entity.Streamer, len(streamers))
	for i, s := range streamers {
		entities[i] = r.toEntity(s)
	}

	return entities, nil
}

func (r *streamerRepository) UpdateLiveStatus(ctx context.Context, streamerID int64, isLive bool, lastLiveAt time.Time) error {
	query := r.client.Streamer.
		UpdateOneID(streamerID).
		SetIsLive(isLive).
		SetLastCheckedAt(time.Now())

	// Set last_live_at if isLive is true
	if isLive {
		query = query.SetLastLiveAt(lastLiveAt)
	}

	err := query.Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.NotFound("Streamer").WithDetail("id", streamerID)
		}
		r.logger.Error("failed to update live status", zap.Error(err), zap.Int64("streamer_id", streamerID), zap.Bool("is_live", isLive))
		return errors.ConvertDatabaseError(err, "Streamer")
	}
	return nil
}

// toEntity converts an ent.Streamer to entity.Streamer
func (r *streamerRepository) toEntity(s *ent.Streamer) *entity.Streamer {
	var lastCheckedAt time.Time
	if s.LastCheckedAt != nil {
		lastCheckedAt = *s.LastCheckedAt
	}

	var lastLiveAt time.Time
	if s.LastLiveAt != nil {
		lastLiveAt = *s.LastLiveAt
	}

	return &entity.Streamer{
		ID:                 s.ID,
		PlatformID:         s.PlatformID,
		PlatformStreamerID: s.PlatformStreamerID,
		Name:               s.Name,
		Avatar:             s.Avatar,
		Description:        s.Description,
		RoomURL:            s.RoomURL,
		LastCheckedAt:      lastCheckedAt,
		IsLive:             s.IsLive,
		LastLiveAt:         lastLiveAt,
		CreatedAt:          s.CreatedAt,
		UpdatedAt:          s.UpdatedAt,
		DeleteAt:           s.DeleteAt,
	}
}
