package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/infrastructure/provider/jwt"
)

const (
	UserContextKey   = "auth_user"
	UserIdContextKey = "user_id"
)

func GetCurrentUser(c fiber.Ctx) (*jwt.UserClaims, bool) {
	user := c.Locals(UserContextKey)
	if user == nil {
		return nil, false
	}
	claims, ok := user.(*jwt.UserClaims)
	return claims, ok
}

func GetCurrentUserId(c fiber.Ctx) (int64, bool) {
	userID := c.Locals(UserIdContextKey)
	if userID == nil {
		return 0, false
	}
	id, ok := userID.(int64)
	return id, ok
}

func IsAuthenticated(c fiber.Ctx) bool {
	_, exists := GetCurrentUser(c)
	return exists
}
