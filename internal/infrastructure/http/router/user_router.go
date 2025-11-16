package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/infrastructure/http/controller"
)

type UserRouter struct {
	userController *controller.UserController
}

func (u *UserRouter) RegisterRouters(router fiber.Router) {
	group := router.Group("/api/v1/user")

	group.Post("/", u.userController.Create)
	group.Put("/:id", u.userController.Update)
	group.Get("/:id", u.userController.GetByID)
	group.Delete("/:id", u.userController.DeleteByID)
	group.Get("/list", u.userController.List)
}

func NewUserRouter(userController *controller.UserController) *UserRouter {
	return &UserRouter{
		userController: userController,
	}
}
