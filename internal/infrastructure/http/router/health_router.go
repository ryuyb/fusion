package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/infrastructure/http/controller"
)

type HealthRouter struct {
	controller *controller.HealthController
}

func (u *HealthRouter) RegisterRouters(router fiber.Router) {
	router.Get("/api/v1/health", u.controller.HealthCheck)
}

func NewHealthRouter(controller *controller.HealthController) *HealthRouter {
	return &HealthRouter{
		controller: controller,
	}
}
