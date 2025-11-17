package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/infrastructure/http/controller"
)

type StreamingPlatformRouter struct {
	controller *controller.StreamingPlatformController
}

func NewStreamingPlatformRouter(controller *controller.StreamingPlatformController) Router {
	return &StreamingPlatformRouter{controller: controller}
}

func (r *StreamingPlatformRouter) RegisterRouters(router fiber.Router) {
	group := router.Group("/platforms")
	group.Post("/", r.controller.Create)
	group.Put("/:id", r.controller.Update)
	group.Delete("/:id", r.controller.Delete)
	group.Get("/:id", r.controller.GetByID)
	group.Get("/", r.controller.List)
}
