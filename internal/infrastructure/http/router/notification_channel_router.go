package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/infrastructure/http/controller"
)

type NotificationChannelRouter struct {
	controller *controller.NotificationChannelController
}

func NewNotificationChannelRouter(controller *controller.NotificationChannelController) Router {
	return &NotificationChannelRouter{controller: controller}
}

func (r *NotificationChannelRouter) RegisterRouters(router fiber.Router) {
	group := router.Group("/notification-channels")
	group.Post("/", r.controller.Create)
	group.Put("/:id", r.controller.Update)
	group.Delete("/:id", r.controller.Delete)
	group.Get("/:id", r.controller.GetByID)
	group.Get("/users/:user_id", r.controller.ListByUser)
}
