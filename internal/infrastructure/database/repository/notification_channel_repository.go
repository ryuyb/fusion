package repository

import (
	"context"

	"github.com/ryuyb/fusion/internal/core/domain"
	coreRepo "github.com/ryuyb/fusion/internal/core/port/repository"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent/notificationchannel"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type notificationChannelRepository struct {
	client *ent.Client
	logger *zap.Logger
}

func NewNotificationChannelRepository(client *ent.Client, logger *zap.Logger) coreRepo.NotificationChannelRepository {
	return &notificationChannelRepository{
		client: client,
		logger: logger,
	}
}

func (r *notificationChannelRepository) Create(ctx context.Context, channel *domain.NotificationChannel) (*domain.NotificationChannel, error) {
	builder := r.client.NotificationChannel.Create().
		SetUserID(channel.UserID).
		SetChannelType(string(channel.ChannelType)).
		SetName(channel.Name).
		SetEnable(channel.Enable).
		SetPriority(channel.Priority)

	if channel.Config != nil {
		builder.SetConfig(channel.Config)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		r.logger.Error("failed to create notification channel",
			zap.Error(err),
			zap.Int64("user_id", channel.UserID),
			zap.String("name", channel.Name),
		)
		return nil, errors2.ConvertDatabaseError(err, "NotificationChannel")
	}
	return r.toDomain(created), nil
}

func (r *notificationChannelRepository) Update(ctx context.Context, channel *domain.NotificationChannel) (*domain.NotificationChannel, error) {
	builder := r.client.NotificationChannel.UpdateOneID(channel.ID).
		SetChannelType(string(channel.ChannelType)).
		SetName(channel.Name).
		SetEnable(channel.Enable).
		SetPriority(channel.Priority)

	if channel.Config == nil {
		builder.ClearConfig()
	} else {
		builder.SetConfig(channel.Config)
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors2.NotFound("NotificationChannel").WithDetail("id", channel.ID)
		}
		r.logger.Error("failed to update notification channel",
			zap.Error(err),
			zap.Int64("id", channel.ID),
		)
		return nil, errors2.ConvertDatabaseError(err, "NotificationChannel")
	}
	return r.toDomain(updated), nil
}

func (r *notificationChannelRepository) Delete(ctx context.Context, id int64) error {
	err := r.client.NotificationChannel.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors2.NotFound("NotificationChannel").WithDetail("id", id)
		}
		r.logger.Error("failed to delete notification channel", zap.Error(err), zap.Int64("id", id))
		return errors2.ConvertDatabaseError(err, "NotificationChannel")
	}
	return nil
}

func (r *notificationChannelRepository) FindById(ctx context.Context, id int64) (*domain.NotificationChannel, error) {
	entity, err := r.client.NotificationChannel.
		Query().
		Where(notificationchannel.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors2.NotFound("NotificationChannel").WithDetail("id", id)
		}
		r.logger.Error("failed to find notification channel by id", zap.Error(err), zap.Int64("id", id))
		return nil, errors2.ConvertDatabaseError(err, "NotificationChannel")
	}
	return r.toDomain(entity), nil
}

func (r *notificationChannelRepository) ListByUserId(ctx context.Context, userID int64, offset, limit int) ([]*domain.NotificationChannel, int, error) {
	total, err := r.client.NotificationChannel.
		Query().
		Where(notificationchannel.UserIDEQ(userID)).
		Count(ctx)
	if err != nil {
		r.logger.Error("failed to count notification channels", zap.Error(err), zap.Int64("user_id", userID))
		return nil, 0, errors2.DatabaseError(err)
	}

	entities, err := r.client.NotificationChannel.
		Query().
		Where(notificationchannel.UserIDEQ(userID)).
		Offset(offset).
		Limit(limit).
		Order(
			ent.Desc(notificationchannel.FieldPriority),
			ent.Desc(notificationchannel.FieldCreatedAt),
		).
		All(ctx)
	if err != nil {
		r.logger.Error("failed to list notification channels", zap.Error(err), zap.Int64("user_id", userID))
		return nil, 0, errors2.DatabaseError(err)
	}

	results := make([]*domain.NotificationChannel, len(entities))
	for i, entity := range entities {
		results[i] = r.toDomain(entity)
	}
	return results, total, nil
}

func (r *notificationChannelRepository) ExistByName(ctx context.Context, userID int64, name string) (bool, error) {
	exist, err := r.client.NotificationChannel.
		Query().
		Where(
			notificationchannel.UserIDEQ(userID),
			notificationchannel.NameEQ(name),
		).
		Exist(ctx)
	if err != nil {
		r.logger.Error("failed to check notification channel name",
			zap.Error(err),
			zap.Int64("user_id", userID),
			zap.String("name", name),
		)
		return false, errors2.ConvertDatabaseError(err, "NotificationChannel")
	}
	return exist, nil
}

func (r *notificationChannelRepository) toDomain(entity *ent.NotificationChannel) *domain.NotificationChannel {
	var config map[string]any
	if entity.Config != nil {
		config = lo.Assign(map[string]any{}, entity.Config)
	}

	return &domain.NotificationChannel{
		ID:          entity.ID,
		UserID:      entity.UserID,
		ChannelType: domain.NotificationChannelType(entity.ChannelType),
		Name:        entity.Name,
		Config:      config,
		Enable:      entity.Enable,
		Priority:    entity.Priority,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
