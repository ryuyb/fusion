package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ryuyb/fusion/internal/pkg/auth"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
)

type Auth struct {
	jwtManager *auth.JWTManager
}

func NewAuth(jwtManager *auth.JWTManager) *Auth {
	return &Auth{jwtManager: jwtManager}
}

func (a *Auth) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return errors2.Unauthorized("Authorization header is empty")
		}

		authHeaderParts := strings.SplitN(authHeader, " ", 2)
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			return errors2.Unauthorized("Authorization header format is invalid")
		}

		tokenString := authHeaderParts[1]

		claims, err := a.jwtManager.ValidateToken(tokenString)
		if err != nil {
			return err
		}

		c.Locals(auth.UserContextKey, claims)
		c.Locals(auth.UserIdContextKey, claims.UserID)

		return c.Next()
	}
}

func (a *Auth) Optional() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		authHeaderParts := strings.SplitN(authHeader, " ", 2)
		if len(authHeaderParts) == 2 && authHeaderParts[0] == "Bearer" {
			tokenString := authHeaderParts[1]
			if claims, err := a.jwtManager.ValidateToken(tokenString); err == nil {
				c.Locals(auth.UserContextKey, claims)
				c.Locals(auth.UserIdContextKey, claims.UserID)
			}
		}

		return c.Next()
	}
}
