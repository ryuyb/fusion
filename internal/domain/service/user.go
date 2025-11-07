package service

import (
	"context"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/interface/http/dto/request"
)

type UserService interface {
	Create(ctx context.Context, req *request.CreateUserRequest) (*entity.User, error)

	Delete(ctx context.Context, id int64) error

	Update(ctx context.Context, req *request.UpdateUserRequest) (*entity.User, error)

	GetByID(ctx context.Context, id int64) (*entity.User, error)

	List(ctx context.Context, page, pageSize int) ([]*entity.User, int, error)
}
