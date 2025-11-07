package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/interface/http/handler"
)

type AuthRoute struct {
	authHandler *handler.AuthHandler
}

func NewAuthRoute(authHandler *handler.AuthHandler) *AuthRoute {
	return &AuthRoute{
		authHandler: authHandler,
	}
}

func (r *AuthRoute) RegisterRouters(router fiber.Router) {
	group := router.Group("/api/v1/auth")

	group.Post("/login", r.authHandler.Login)
	group.Post("/register", r.authHandler.Register)
}
