package service

import (
	"context"
	"time"

	"github.com/ryuyb/fusion/internal/core/command"
)

type AuthService interface {
	Login(ctx context.Context, cmd command.LoginCommand) (string, time.Time, error)

	Register(ctx context.Context, cmd command.RegisterCommand) error
}
