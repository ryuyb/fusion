package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ryuyb/fusion/internal/interface/http/handler"
)

type UserRoute struct {
	userHandler *handler.UserHandler
}

func NewUserRoute(userHandler *handler.UserHandler) *UserRoute {
	return &UserRoute{
		userHandler: userHandler,
	}
}

func (r *UserRoute) RegisterRouters(router fiber.Router) {
	group := router.Group("/api/v1/user")

	group.Post("/", r.userHandler.Create)
	group.Delete("/:id", r.userHandler.DeleteByID)
}
