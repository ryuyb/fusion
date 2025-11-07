package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/interface/http/handler"
	"github.com/ryuyb/fusion/internal/interface/http/middleware"
)

type UserRoute struct {
	userHandler *handler.UserHandler
	authMw      *middleware.Auth
}

func NewUserRoute(userHandler *handler.UserHandler, authMw *middleware.Auth) *UserRoute {
	return &UserRoute{
		userHandler: userHandler,
		authMw:      authMw,
	}
}

func (r *UserRoute) RegisterRouters(router fiber.Router) {
	group := router.Group("/api/v1/user")

	// All user endpoints require authentication
	group.Use(r.authMw.Handler())

	group.Post("/", r.userHandler.Create)
	group.Put("/", r.userHandler.Update)
	group.Get("/list", r.userHandler.List)
	group.Get("/:id", r.userHandler.GetByID)
	group.Delete("/:id", r.userHandler.DeleteByID)
}
