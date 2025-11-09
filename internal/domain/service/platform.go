package service

import (
	"context"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/interface/http/dto/request"
)

type PlatformService interface {
	Create(ctx context.Context, req *request.CreatePlatformRequest) (*entity.Platform, error)

	Update(ctx context.Context, id int64, req *request.UpdatePlatformRequest) (*entity.Platform, error)

	GetByID(ctx context.Context, id int64) (*entity.Platform, error)

	List(ctx context.Context) ([]*entity.Platform, error)

	Delete(ctx context.Context, id int64) error

	TestConnection(ctx context.Context, id int64) (bool, error)
}
