package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/infrastructure/http/controller"
)

type AuthRouter struct {
	authController *controller.AuthController
}

func (u *AuthRouter) RegisterRouters(router fiber.Router) {
	group := router.Group("/api/v1/auth")

	group.Post("/register", u.authController.Register)
	group.Post("/login", u.authController.Login)
}

func NewAuthRouter(authController *controller.AuthController) *AuthRouter {
	return &AuthRouter{
		authController: authController,
	}
}
