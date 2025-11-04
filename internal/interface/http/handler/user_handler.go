package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ryuyb/fusion/internal/domain/service"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService service.UserService
	logger      *zap.Logger
}

func NewUserHandler(userService service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	return nil
}
