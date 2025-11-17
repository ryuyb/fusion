package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/infrastructure/http/controller"
)

type StreamerRouter struct {
	controller *controller.StreamerController
}

func NewStreamerRouter(controller *controller.StreamerController) Router {
	return &StreamerRouter{controller: controller}
}

func (r *StreamerRouter) RegisterRouters(router fiber.Router) {
	group := router.Group("/api/v1/streamers")
	group.Post("/", r.controller.Create)
	group.Put("/:id", r.controller.Update)
	group.Delete("/:id", r.controller.Delete)
	group.Get("/:id", r.controller.GetByID)
	group.Get("/:platform_type/:platform_streamer_id", r.controller.GetByPlatformStreamerID)
	group.Get("/", r.controller.List)
}
