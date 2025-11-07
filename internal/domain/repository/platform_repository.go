package repository

import (
	"context"

	"github.com/ryuyb/fusion/internal/domain/entity"
)

type PlatformRepository interface {
	Create(ctx context.Context, platform *entity.Platform) (*entity.Platform, error)

	Update(ctx context.Context, platform *entity.Platform) (*entity.Platform, error)

	FindByID(ctx context.Context, id int64) (*entity.Platform, error)

	FindByType(ctx context.Context, platformType entity.PlatformType) (*entity.Platform, error)

	List(ctx context.Context) ([]*entity.Platform, error)

	Delete(ctx context.Context, id int64) error
}
