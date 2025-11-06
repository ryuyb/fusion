package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/swagger/v2"
)

type SwaggerRoute struct{}

func NewSwaggerRoute() *SwaggerRoute {
	return &SwaggerRoute{}
}

func (r *SwaggerRoute) RegisterRouters(router fiber.Router) {
	router.Get("/swagger/*", swagger.HandlerDefault)
}
