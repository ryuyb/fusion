package handler

import (
	"github.com/gofiber/fiber/v3"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health 健康检查
//
//	@Summary	健康检查
//	@Tags		Health
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/health [get]
func (h *HealthHandler) Health(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}
