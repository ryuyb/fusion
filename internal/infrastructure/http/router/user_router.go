package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/infrastructure/http/controller"
	"github.com/ryuyb/fusion/internal/infrastructure/http/middleware"
)

type UserRouter struct {
	userController *controller.UserController
	auth           *middleware.Auth
}

func (u *UserRouter) RegisterRouters(router fiber.Router) {
	group := router.Group("/api/v1/user")

	group.Post("/", u.userController.Create)
	group.Get("/list", u.userController.List)
	group.Put("/:id", u.userController.Update)
	group.Get("/:id", u.userController.GetByID)
	group.Delete("/:id", u.userController.DeleteByID)
}

func NewUserRouter(userController *controller.UserController, auth *middleware.Auth) *UserRouter {
	return &UserRouter{
		userController: userController,
		auth:           auth,
	}
}
