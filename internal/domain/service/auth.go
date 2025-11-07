package service

import (
	"context"
	"time"

	"github.com/ryuyb/fusion/internal/interface/http/dto/request"
)

// AuthService defines the authentication business logic interface
type AuthService interface {
	// Login authenticates a user and returns a JWT token
	Login(ctx context.Context, req *request.LoginRequest) (string, time.Time, error)

	// Register creates a new user and returns a JWT token
	Register(ctx context.Context, req *request.RegisterRequest) (string, time.Time, error)
}
