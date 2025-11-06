package service

import (
	"context"

	"github.com/ryuyb/fusion/internal/application/dto"
	"github.com/ryuyb/fusion/internal/interface/http/response"
)

type UserService interface {
	Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)

	Delete(ctx context.Context, id int64) error

	Update(ctx context.Context, req *dto.UpdateUserRequest) (*dto.UserResponse, error)

	GetByID(ctx context.Context, id int64) (*dto.UserResponse, error)

	List(ctx context.Context, page, pageSize int) (*response.PaginationResponse[*dto.UserResponse], error)
}
