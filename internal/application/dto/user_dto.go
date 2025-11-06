package dto

import (
	"time"

	"github.com/ryuyb/fusion/internal/domain/entity"
)

type UserResponse struct {
	ID        int64             `json:"id"`
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	Status    entity.UserStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeleteAt  time.Time         `json:"delete_at,omitzero"`
}

type CreateUserRequest struct {
	Username       string            `json:"username" validate:"required,min=3,max=16"`
	Email          string            `json:"email" validate:"required,email"`
	Password       string            `json:"password" validate:"required,min=6,max=32"`
	RepeatPassword string            `json:"repeat_password" validate:"required,eqfield=Password"`
	Status         entity.UserStatus `json:"status" validate:"required,oneof=active inactive banned"`
}

type UpdateUserRequest struct {
	ID int64 `json:"id" validate:"required,gte=1"`
	CreateUserRequest
}
