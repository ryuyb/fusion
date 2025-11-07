package repository

import (
	"context"

	"github.com/ryuyb/fusion/internal/domain/entity"
)

type RuleRepository interface {
	Create(ctx context.Context, rule *entity.NotificationRule) (*entity.NotificationRule, error)

	Update(ctx context.Context, rule *entity.NotificationRule) (*entity.NotificationRule, error)

	FindByID(ctx context.Context, id int64) (*entity.NotificationRule, error)

	FindByUser(ctx context.Context, userID int64) ([]*entity.NotificationRule, error)

	// FindEnabledByUser returns only enabled rules for a user
	FindEnabledByUser(ctx context.Context, userID int64) ([]*entity.NotificationRule, error)

	Delete(ctx context.Context, id int64) error
}
