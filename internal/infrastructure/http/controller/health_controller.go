package controller

import "github.com/gofiber/fiber/v3"

type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

// HealthCheck Health check
//
//	@Summary	Health check
//	@Tags		Health
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/health [get]
func (h *HealthController) HealthCheck(ctx fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"health": "ok",
	})
}
