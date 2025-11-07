package repository

import (
	"context"
	"time"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/domain/repository"
	"github.com/ryuyb/fusion/internal/infrastructure/database"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent/platform"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent/userfollowing"
	"github.com/ryuyb/fusion/internal/pkg/errors"
	"go.uber.org/zap"
)

type followingRepository struct {
	client *database.Client
	logger *zap.Logger
}

func (r *followingRepository) Create(ctx context.Context, following *entity.UserFollowing) (*entity.UserFollowing, error) {
	create := r.client.UserFollowing.
		Create().
		SetUserID(following.UserID).
		SetStreamerID(following.StreamerID).
		SetNotificationEnabled(following.NotificationEnabled)

	// LastNotifiedAt is nullable, only set if not zero
	if !following.LastNotifiedAt.IsZero() {
		create.SetLastNotifiedAt(following.LastNotifiedAt)
	}

	created, err := create.Save(ctx)
	if err != nil {
		r.logger.Error("failed to create user following", zap.Error(err),
			zap.Int64("user_id", following.UserID),
			zap.Int64("streamer_id", following.StreamerID),
			zap.Bool("notification_enabled", following.NotificationEnabled))
		return nil, errors.ConvertDatabaseError(err, "UserFollowing")
	}
	return r.toEntity(created), nil
}

func (r *followingRepository) Update(ctx context.Context, following *entity.UserFollowing) (*entity.UserFollowing, error) {
	update := r.client.UserFollowing.
		UpdateOneID(following.ID).
		SetNotificationEnabled(following.NotificationEnabled)

	// Update LastNotifiedAt if provided
	if !following.LastNotifiedAt.IsZero() {
		update.SetLastNotifiedAt(following.LastNotifiedAt)
	}

	updated, err := update.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("UserFollowing").WithDetail("id", following.ID)
		}
		r.logger.Error("failed to update user following", zap.Error(err), zap.Int64("id", following.ID))
		return nil, errors.ConvertDatabaseError(err, "UserFollowing")
	}
	return r.toEntity(updated), nil
}

func (r *followingRepository) FindByID(ctx context.Context, id int64) (*entity.UserFollowing, error) {
	following, err := r.client.UserFollowing.
		Query().
		Where(userfollowing.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("UserFollowing").WithDetail("id", id)
		}
		r.logger.Error("failed to find user following by id", zap.Error(err), zap.Int64("id", id))
		return nil, errors.ConvertDatabaseError(err, "UserFollowing")
	}
	return r.toEntity(following), nil
}

func (r *followingRepository) FindByUserAndStreamer(ctx context.Context, userID, streamerID int64) (*entity.UserFollowing, error) {
	following, err := r.client.UserFollowing.
		Query().
		Where(
			userfollowing.UserID(userID),
			userfollowing.StreamerID(streamerID),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("UserFollowing").WithDetail("user_id", userID).WithDetail("streamer_id", streamerID)
		}
		r.logger.Error("failed to find user following by user and streamer", zap.Error(err),
			zap.Int64("user_id", userID), zap.Int64("streamer_id", streamerID))
		return nil, errors.ConvertDatabaseError(err, "UserFollowing")
	}
	return r.toEntity(following), nil
}

func (r *followingRepository) FindByUser(ctx context.Context, userID int64, filters *repository.FollowingFilters) ([]*entity.UserFollowing, int, error) {
	query := r.client.UserFollowing.
		Query().
		Where(userfollowing.UserID(userID))

	// Apply filters
	if filters != nil {
		if filters.NotificationEnabled != nil {
			query = query.Where(userfollowing.NotificationEnabled(*filters.NotificationEnabled))
		}

		// Join with streamer to filter by platform type
		if filters.PlatformType != nil {
			query.WithStreamer(func(streamerQuery *ent.StreamerQuery) {
				streamerQuery.WithPlatform(func(platformQuery *ent.PlatformQuery) {
					platformQuery.Where(platform.PlatformTypeEQ(platform.PlatformType(*filters.PlatformType)))
				})
			})
		}
	}

	// Get total count
	total, err := query.Count(ctx)
	if err != nil {
		r.logger.Error("failed to count user followings", zap.Error(err), zap.Int64("user_id", userID))
		return nil, 0, errors.ConvertDatabaseError(err, "UserFollowing")
	}

	// Apply pagination
	offset, limit := 0, 10 // default
	if filters != nil {
		if filters.Page > 0 && filters.PageSize > 0 {
			offset = (filters.Page - 1) * filters.PageSize
			limit = filters.PageSize
		}
	}

	// Order by created_at descending (most recent first)
	query = query.
		Offset(offset).
		Limit(limit).
		Order(ent.Desc(userfollowing.FieldCreatedAt))

	followings, err := query.All(ctx)
	if err != nil {
		r.logger.Error("failed to find user followings", zap.Error(err), zap.Int64("user_id", userID))
		return nil, 0, errors.ConvertDatabaseError(err, "UserFollowing")
	}

	entities := make([]*entity.UserFollowing, len(followings))
	for i, f := range followings {
		entities[i] = r.toEntity(f)
	}

	return entities, total, nil
}

func (r *followingRepository) FindByStreamer(ctx context.Context, streamerID int64, notificationEnabled *bool) ([]*entity.UserFollowing, error) {
	query := r.client.UserFollowing.
		Query().
		Where(userfollowing.StreamerID(streamerID))

	// Apply notification enabled filter if provided
	if notificationEnabled != nil {
		query = query.Where(userfollowing.NotificationEnabled(*notificationEnabled))
	}

	// Order by created_at descending
	query = query.Order(ent.Desc(userfollowing.FieldCreatedAt))

	followings, err := query.All(ctx)
	if err != nil {
		r.logger.Error("failed to find followers by streamer", zap.Error(err), zap.Int64("streamer_id", streamerID))
		return nil, errors.ConvertDatabaseError(err, "UserFollowing")
	}

	entities := make([]*entity.UserFollowing, len(followings))
	for i, f := range followings {
		entities[i] = r.toEntity(f)
	}

	return entities, nil
}

func (r *followingRepository) Delete(ctx context.Context, id int64) error {
	err := r.client.UserFollowing.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.NotFound("UserFollowing").WithDetail("id", id)
		}
		r.logger.Error("failed to delete user following", zap.Error(err), zap.Int64("id", id))
		return errors.ConvertDatabaseError(err, "UserFollowing")
	}
	return nil
}

func (r *followingRepository) UpdateLastNotifiedAt(ctx context.Context, id int64, t time.Time) error {
	err := r.client.UserFollowing.
		UpdateOneID(id).
		SetLastNotifiedAt(t).
		Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.NotFound("UserFollowing").WithDetail("id", id)
		}
		r.logger.Error("failed to update last notified at", zap.Error(err), zap.Int64("id", id), zap.Time("time", t))
		return errors.ConvertDatabaseError(err, "UserFollowing")
	}
	return nil
}

func NewFollowingRepository(client *database.Client, logger *zap.Logger) repository.FollowingRepository {
	return &followingRepository{
		client: client,
		logger: logger,
	}
}

func (r *followingRepository) toEntity(f *ent.UserFollowing) *entity.UserFollowing {
	var lastNotifiedAt time.Time
	if f.LastNotifiedAt != nil {
		lastNotifiedAt = *f.LastNotifiedAt
	}

	return &entity.UserFollowing{
		ID:                  f.ID,
		UserID:              f.UserID,
		StreamerID:          f.StreamerID,
		NotificationEnabled: f.NotificationEnabled,
		LastNotifiedAt:      lastNotifiedAt,
		CreatedAt:           f.CreatedAt,
		UpdatedAt:           f.UpdatedAt,
		DeleteAt:            f.DeleteAt,
	}
}
