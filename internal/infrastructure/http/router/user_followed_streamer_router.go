package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/infrastructure/http/controller"
)

type UserFollowedStreamerRouter struct {
	controller *controller.UserFollowedStreamerController
}

func NewUserFollowedStreamerRouter(controller *controller.UserFollowedStreamerController) Router {
	return &UserFollowedStreamerRouter{controller: controller}
}

func (r *UserFollowedStreamerRouter) RegisterRouters(router fiber.Router) {
	group := router.Group("/api/v1/follows")
	group.Post("/", r.controller.Create)
	group.Put("/:id", r.controller.Update)
	group.Delete("/:id", r.controller.Delete)
	group.Get("/:id", r.controller.GetByID)
	group.Get("/users/:user_id", r.controller.ListByUser)
	group.Get("/streamers/:streamer_id", r.controller.ListByStreamer)
}
