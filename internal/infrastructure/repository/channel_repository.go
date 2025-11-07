package repository

import (
	"context"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/domain/repository"
	"github.com/ryuyb/fusion/internal/infrastructure/database"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent/notificationchannel"
	"github.com/ryuyb/fusion/internal/pkg/errors"
	"go.uber.org/zap"
)

type channelRepository struct {
	client *database.Client
	logger *zap.Logger
}

func (r *channelRepository) Create(ctx context.Context, channel *entity.NotificationChannel) (*entity.NotificationChannel, error) {
	created, err := r.client.NotificationChannel.
		Create().
		SetUserID(channel.UserID).
		SetChannelType(notificationchannel.ChannelType(channel.ChannelType)).
		SetName(channel.Name).
		SetConfig(channel.Config).
		SetIsEnabled(channel.IsEnabled).
		SetPriority(channel.Priority).
		Save(ctx)
	if err != nil {
		r.logger.Error("failed to create notification channel", zap.Error(err),
			zap.Int64("user_id", channel.UserID),
			zap.String("channel_type", string(channel.ChannelType)),
			zap.String("name", channel.Name))
		return nil, errors.ConvertDatabaseError(err, "NotificationChannel")
	}
	return r.toEntity(created), nil
}

func (r *channelRepository) Update(ctx context.Context, channel *entity.NotificationChannel) (*entity.NotificationChannel, error) {
	updated, err := r.client.NotificationChannel.
		UpdateOneID(channel.ID).
		SetName(channel.Name).
		SetConfig(channel.Config).
		SetIsEnabled(channel.IsEnabled).
		SetPriority(channel.Priority).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("NotificationChannel").WithDetail("id", channel.ID)
		}
		r.logger.Error("failed to update notification channel", zap.Error(err), zap.Int64("id", channel.ID))
		return nil, errors.ConvertDatabaseError(err, "NotificationChannel")
	}
	return r.toEntity(updated), nil
}

func (r *channelRepository) FindByID(ctx context.Context, id int64) (*entity.NotificationChannel, error) {
	ch, err := r.client.NotificationChannel.
		Query().
		Where(notificationchannel.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("NotificationChannel").WithDetail("id", id)
		}
		r.logger.Error("failed to find notification channel by id", zap.Error(err), zap.Int64("id", id))
		return nil, errors.ConvertDatabaseError(err, "NotificationChannel")
	}
	return r.toEntity(ch), nil
}

func (r *channelRepository) FindByUser(ctx context.Context, userID int64) ([]*entity.NotificationChannel, error) {
	channels, err := r.client.NotificationChannel.
		Query().
		Where(notificationchannel.UserID(userID)).
		Order(ent.Asc(notificationchannel.FieldPriority), ent.Desc(notificationchannel.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		r.logger.Error("failed to find notification channels by user", zap.Error(err), zap.Int64("user_id", userID))
		return nil, errors.ConvertDatabaseError(err, "NotificationChannel")
	}

	entities := make([]*entity.NotificationChannel, len(channels))
	for i, ch := range channels {
		entities[i] = r.toEntity(ch)
	}

	return entities, nil
}

func (r *channelRepository) FindEnabledByUser(ctx context.Context, userID int64) ([]*entity.NotificationChannel, error) {
	channels, err := r.client.NotificationChannel.
		Query().
		Where(
			notificationchannel.UserID(userID),
			notificationchannel.IsEnabled(true),
		).
		Order(ent.Asc(notificationchannel.FieldPriority), ent.Desc(notificationchannel.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		r.logger.Error("failed to find enabled notification channels by user", zap.Error(err), zap.Int64("user_id", userID))
		return nil, errors.ConvertDatabaseError(err, "NotificationChannel")
	}

	entities := make([]*entity.NotificationChannel, len(channels))
	for i, ch := range channels {
		entities[i] = r.toEntity(ch)
	}

	return entities, nil
}

func (r *channelRepository) Delete(ctx context.Context, id int64) error {
	err := r.client.NotificationChannel.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.NotFound("NotificationChannel").WithDetail("id", id)
		}
		r.logger.Error("failed to delete notification channel", zap.Error(err), zap.Int64("id", id))
		return errors.ConvertDatabaseError(err, "NotificationChannel")
	}
	return nil
}

func NewChannelRepository(client *database.Client, logger *zap.Logger) repository.ChannelRepository {
	return &channelRepository{
		client: client,
		logger: logger,
	}
}

func (r *channelRepository) toEntity(ch *ent.NotificationChannel) *entity.NotificationChannel {
	return &entity.NotificationChannel{
		ID:          ch.ID,
		UserID:      ch.UserID,
		ChannelType: entity.ChannelType(ch.ChannelType),
		Name:        ch.Name,
		Config:      ch.Config,
		IsEnabled:   ch.IsEnabled,
		Priority:    ch.Priority,
		CreatedAt:   ch.CreatedAt,
		UpdatedAt:   ch.UpdatedAt,
		DeleteAt:    ch.DeleteAt,
	}
}
