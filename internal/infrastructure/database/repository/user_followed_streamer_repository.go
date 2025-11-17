package repository

import (
	"context"
	"slices"
	"time"

	"github.com/ryuyb/fusion/internal/core/domain"
	coreRepo "github.com/ryuyb/fusion/internal/core/port/repository"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent/userfollowedstreamer"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type userFollowedStreamerRepository struct {
	client *ent.Client
	logger *zap.Logger
}

func NewUserFollowedStreamerRepository(client *ent.Client, logger *zap.Logger) coreRepo.UserFollowedStreamerRepository {
	return &userFollowedStreamerRepository{
		client: client,
		logger: logger,
	}
}

func (r *userFollowedStreamerRepository) Create(ctx context.Context, follow *domain.UserFollowedStreamer) (*domain.UserFollowedStreamer, error) {
	builder := r.client.UserFollowedStreamer.Create().
		SetUserID(follow.UserID).
		SetStreamerID(follow.StreamerID).
		SetNotificationsEnabled(follow.NotificationsEnabled)

	if follow.Alias != "" {
		builder.SetAlias(follow.Alias)
	}
	if follow.Notes != "" {
		builder.SetNotes(follow.Notes)
	}
	if len(follow.NotificationChannelIDs) > 0 {
		builder.SetNotificationChannelIds(follow.NotificationChannelIDs)
	}
	if follow.LastNotificationSentAt != nil {
		builder.SetLastNotificationSentAt(*follow.LastNotificationSentAt)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		r.logger.Error("failed to create user followed streamer",
			zap.Error(err),
			zap.Int64("user_id", follow.UserID),
			zap.Int64("streamer_id", follow.StreamerID),
		)
		return nil, errors2.ConvertDatabaseError(err, "UserFollowedStreamer")
	}
	return r.toDomain(created), nil
}

func (r *userFollowedStreamerRepository) Update(ctx context.Context, follow *domain.UserFollowedStreamer) (*domain.UserFollowedStreamer, error) {
	builder := r.client.UserFollowedStreamer.UpdateOneID(follow.ID).
		SetNotificationsEnabled(follow.NotificationsEnabled)

	if follow.Alias == "" {
		builder.ClearAlias()
	} else {
		builder.SetAlias(follow.Alias)
	}
	if follow.Notes == "" {
		builder.ClearNotes()
	} else {
		builder.SetNotes(follow.Notes)
	}
	builder.ClearNotificationChannelIds()
	if len(follow.NotificationChannelIDs) > 0 {
		builder.SetNotificationChannelIds(follow.NotificationChannelIDs)
	}
	if follow.LastNotificationSentAt == nil {
		builder.ClearLastNotificationSentAt()
	} else {
		builder.SetLastNotificationSentAt(*follow.LastNotificationSentAt)
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors2.NotFound("UserFollowedStreamer").WithDetail("id", follow.ID)
		}
		r.logger.Error("failed to update user followed streamer",
			zap.Error(err),
			zap.Int64("id", follow.ID),
		)
		return nil, errors2.ConvertDatabaseError(err, "UserFollowedStreamer")
	}
	return r.toDomain(updated), nil
}

func (r *userFollowedStreamerRepository) Delete(ctx context.Context, id int64) error {
	err := r.client.UserFollowedStreamer.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors2.NotFound("UserFollowedStreamer").WithDetail("id", id)
		}
		r.logger.Error("failed to delete user followed streamer", zap.Error(err), zap.Int64("id", id))
		return errors2.ConvertDatabaseError(err, "UserFollowedStreamer")
	}
	return nil
}

func (r *userFollowedStreamerRepository) FindById(ctx context.Context, id int64) (*domain.UserFollowedStreamer, error) {
	entity, err := r.client.UserFollowedStreamer.
		Query().
		Where(userfollowedstreamer.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors2.NotFound("UserFollowedStreamer").WithDetail("id", id)
		}
		r.logger.Error("failed to find user followed streamer by id", zap.Error(err), zap.Int64("id", id))
		return nil, errors2.ConvertDatabaseError(err, "UserFollowedStreamer")
	}
	return r.toDomain(entity), nil
}

func (r *userFollowedStreamerRepository) FindByUserAndStreamer(ctx context.Context, userID, streamerID int64) (*domain.UserFollowedStreamer, error) {
	entity, err := r.client.UserFollowedStreamer.
		Query().
		Where(
			userfollowedstreamer.UserIDEQ(userID),
			userfollowedstreamer.StreamerIDEQ(streamerID),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors2.NotFound("UserFollowedStreamer").
				WithDetail("user_id", userID).
				WithDetail("streamer_id", streamerID)
		}
		r.logger.Error("failed to find user followed streamer",
			zap.Error(err),
			zap.Int64("user_id", userID),
			zap.Int64("streamer_id", streamerID),
		)
		return nil, errors2.ConvertDatabaseError(err, "UserFollowedStreamer")
	}
	return r.toDomain(entity), nil
}

func (r *userFollowedStreamerRepository) ExistByUserAndStreamer(ctx context.Context, userID, streamerID int64) (bool, error) {
	exist, err := r.client.UserFollowedStreamer.
		Query().
		Where(
			userfollowedstreamer.UserIDEQ(userID),
			userfollowedstreamer.StreamerIDEQ(streamerID),
		).
		Exist(ctx)
	if err != nil {
		r.logger.Error("failed to check user followed streamer",
			zap.Error(err),
			zap.Int64("user_id", userID),
			zap.Int64("streamer_id", streamerID),
		)
		return false, errors2.ConvertDatabaseError(err, "UserFollowedStreamer")
	}
	return exist, nil
}

func (r *userFollowedStreamerRepository) ListByUserId(ctx context.Context, userID int64, offset, limit int) ([]*domain.UserFollowedStreamer, int, error) {
	total, err := r.client.UserFollowedStreamer.
		Query().
		Where(userfollowedstreamer.UserIDEQ(userID)).
		Count(ctx)
	if err != nil {
		r.logger.Error("failed to count follows by user", zap.Error(err), zap.Int64("user_id", userID))
		return nil, 0, errors2.DatabaseError(err)
	}

	entities, err := r.client.UserFollowedStreamer.
		Query().
		Where(userfollowedstreamer.UserIDEQ(userID)).
		Offset(offset).
		Limit(limit).
		Order(
			ent.Desc(userfollowedstreamer.FieldCreatedAt),
		).
		All(ctx)
	if err != nil {
		r.logger.Error("failed to list follows by user", zap.Error(err), zap.Int64("user_id", userID))
		return nil, 0, errors2.DatabaseError(err)
	}

	results := make([]*domain.UserFollowedStreamer, len(entities))
	for i, entity := range entities {
		results[i] = r.toDomain(entity)
	}
	return results, total, nil
}

func (r *userFollowedStreamerRepository) ListByStreamerId(ctx context.Context, streamerID int64, offset, limit int) ([]*domain.UserFollowedStreamer, int, error) {
	total, err := r.client.UserFollowedStreamer.
		Query().
		Where(userfollowedstreamer.StreamerIDEQ(streamerID)).
		Count(ctx)
	if err != nil {
		r.logger.Error("failed to count follows by streamer", zap.Error(err), zap.Int64("streamer_id", streamerID))
		return nil, 0, errors2.DatabaseError(err)
	}

	entities, err := r.client.UserFollowedStreamer.
		Query().
		Where(userfollowedstreamer.StreamerIDEQ(streamerID)).
		Offset(offset).
		Limit(limit).
		Order(
			ent.Desc(userfollowedstreamer.FieldCreatedAt),
		).
		All(ctx)
	if err != nil {
		r.logger.Error("failed to list follows by streamer", zap.Error(err), zap.Int64("streamer_id", streamerID))
		return nil, 0, errors2.DatabaseError(err)
	}

	results := make([]*domain.UserFollowedStreamer, len(entities))
	for i, entity := range entities {
		results[i] = r.toDomain(entity)
	}
	return results, total, nil
}

func (r *userFollowedStreamerRepository) toDomain(entity *ent.UserFollowedStreamer) *domain.UserFollowedStreamer {
	var lastNotification *time.Time
	if entity.LastNotificationSentAt != nil {
		clone := lo.FromPtr(entity.LastNotificationSentAt)
		lastNotification = &clone
	}

	return &domain.UserFollowedStreamer{
		ID:                     entity.ID,
		UserID:                 entity.UserID,
		StreamerID:             entity.StreamerID,
		Alias:                  lo.FromPtr(entity.Alias),
		Notes:                  lo.FromPtr(entity.Notes),
		NotificationsEnabled:   entity.NotificationsEnabled,
		NotificationChannelIDs: slices.Clone(entity.NotificationChannelIds),
		LastNotificationSentAt: lastNotification,
		CreatedAt:              entity.CreatedAt,
		UpdatedAt:              entity.UpdatedAt,
	}
}
