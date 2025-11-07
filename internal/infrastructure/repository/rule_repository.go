package repository

import (
	"context"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/domain/repository"
	"github.com/ryuyb/fusion/internal/infrastructure/database"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent/notificationrule"
	"github.com/ryuyb/fusion/internal/pkg/errors"
	"go.uber.org/zap"
)

type ruleRepository struct {
	client *database.Client
	logger *zap.Logger
}

func NewRuleRepository(client *database.Client, logger *zap.Logger) repository.RuleRepository {
	return &ruleRepository{
		client: client,
		logger: logger,
	}
}

func (r *ruleRepository) Create(ctx context.Context, rule *entity.NotificationRule) (*entity.NotificationRule, error) {
	created, err := r.client.NotificationRule.
		Create().
		SetUserID(rule.UserID).
		SetRuleType(notificationrule.RuleType(rule.RuleType)).
		SetName(rule.Name).
		SetConfig(rule.Config).
		SetIsEnabled(rule.IsEnabled).
		Save(ctx)
	if err != nil {
		r.logger.Error("failed to create notification rule", zap.Error(err), zap.Int64("user_id", rule.UserID), zap.String("rule_type", string(rule.RuleType)))
		return nil, errors.ConvertDatabaseError(err, "NotificationRule")
	}
	return r.toEntity(created), nil
}

func (r *ruleRepository) Update(ctx context.Context, rule *entity.NotificationRule) (*entity.NotificationRule, error) {
	updated, err := r.client.NotificationRule.
		UpdateOneID(rule.ID).
		SetName(rule.Name).
		SetConfig(rule.Config).
		SetIsEnabled(rule.IsEnabled).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("NotificationRule").WithDetail("id", rule.ID)
		}
		r.logger.Error("failed to update notification rule", zap.Error(err), zap.Int64("id", rule.ID))
		return nil, errors.ConvertDatabaseError(err, "NotificationRule")
	}
	return r.toEntity(updated), nil
}

func (r *ruleRepository) FindByID(ctx context.Context, id int64) (*entity.NotificationRule, error) {
	rule, err := r.client.NotificationRule.
		Query().
		Where(notificationrule.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("NotificationRule").WithDetail("id", id)
		}
		r.logger.Error("failed to find notification rule by id", zap.Error(err), zap.Int64("id", id))
		return nil, errors.ConvertDatabaseError(err, "NotificationRule")
	}
	return r.toEntity(rule), nil
}

func (r *ruleRepository) FindByUser(ctx context.Context, userID int64) ([]*entity.NotificationRule, error) {
	rules, err := r.client.NotificationRule.
		Query().
		Where(notificationrule.UserID(userID)).
		Order(ent.Desc(notificationrule.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		r.logger.Error("failed to find notification rules by user", zap.Error(err), zap.Int64("user_id", userID))
		return nil, errors.DatabaseError(err)
	}

	entities := make([]*entity.NotificationRule, len(rules))
	for i, rule := range rules {
		entities[i] = r.toEntity(rule)
	}

	return entities, nil
}

func (r *ruleRepository) FindEnabledByUser(ctx context.Context, userID int64) ([]*entity.NotificationRule, error) {
	rules, err := r.client.NotificationRule.
		Query().
		Where(
			notificationrule.UserID(userID),
			notificationrule.IsEnabled(true),
		).
		Order(ent.Desc(notificationrule.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		r.logger.Error("failed to find enabled notification rules by user", zap.Error(err), zap.Int64("user_id", userID))
		return nil, errors.DatabaseError(err)
	}

	entities := make([]*entity.NotificationRule, len(rules))
	for i, rule := range rules {
		entities[i] = r.toEntity(rule)
	}

	return entities, nil
}

func (r *ruleRepository) Delete(ctx context.Context, id int64) error {
	err := r.client.NotificationRule.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.NotFound("NotificationRule").WithDetail("id", id)
		}
		r.logger.Error("failed to delete notification rule", zap.Error(err), zap.Int64("id", id))
		return errors.ConvertDatabaseError(err, "NotificationRule")
	}
	return nil
}

// toEntity converts an ent.NotificationRule to entity.NotificationRule
func (r *ruleRepository) toEntity(nr *ent.NotificationRule) *entity.NotificationRule {
	return &entity.NotificationRule{
		ID:        nr.ID,
		UserID:    nr.UserID,
		RuleType:  entity.RuleType(nr.RuleType),
		Name:      nr.Name,
		Config:    nr.Config,
		IsEnabled: nr.IsEnabled,
		CreatedAt: nr.CreatedAt,
		UpdatedAt: nr.UpdatedAt,
		DeleteAt:  nr.DeleteAt,
	}
}
