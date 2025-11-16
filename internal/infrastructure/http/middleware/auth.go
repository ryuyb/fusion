package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/infrastructure/provider/jwt"
	"github.com/ryuyb/fusion/internal/pkg/auth"
	"github.com/ryuyb/fusion/internal/pkg/errors"
)

const (
	Bearer = "Bearer"
)

type Auth struct {
	jwtManager *jwt.JWTManager
}

func NewAuth(jwtManager *jwt.JWTManager) *Auth {
	return &Auth{
		jwtManager: jwtManager,
	}
}

func (a *Auth) Handler() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authHeader := ctx.Get(fiber.HeaderAuthorization)
		if authHeader == "" {
			return errors.Unauthorized("Missing Authorization header")
		}

		authHeaderParts := strings.SplitN(authHeader, " ", 2)
		if len(authHeaderParts) != 2 || authHeaderParts[0] != Bearer {
			return errors.Unauthorized("Authorization header format is invalid")
		}

		claims, err := a.jwtManager.ValidateToken(authHeaderParts[1])
		if err != nil {
			return err
		}

		ctx.Locals(auth.UserContextKey, claims)
		ctx.Locals(auth.UserIdContextKey, claims.ID)

		return ctx.Next()
	}
}

func (a *Auth) Optional() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authHeader := ctx.Get(fiber.HeaderAuthorization)
		if authHeader == "" {
			return ctx.Next()
		}

		authHeaderParts := strings.SplitN(authHeader, " ", 2)
		if len(authHeaderParts) == 2 && authHeaderParts[0] == Bearer {
			tokenString := authHeaderParts[1]
			if claims, err := a.jwtManager.ValidateToken(tokenString); err == nil {
				ctx.Locals(auth.UserContextKey, claims)
				ctx.Locals(auth.UserIdContextKey, claims.ID)
			}
		}
		return ctx.Next()
	}
}
