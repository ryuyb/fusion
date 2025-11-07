package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/interface/http/handler"
)

type HealthRoute struct {
	healthHandler *handler.HealthHandler
}

func NewHealthRoute(healthHandler *handler.HealthHandler) *HealthRoute {
	return &HealthRoute{
		healthHandler: healthHandler,
	}
}

func (r *HealthRoute) RegisterRouters(router fiber.Router) {
	router.Get("/api/v1/health", r.healthHandler.Health)
}
