package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ryuyb/fusion/internal/application/dto"
	"github.com/ryuyb/fusion/internal/domain/service"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/ryuyb/fusion/internal/pkg/validator"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService service.UserService
	validate    *validator.Validator
	logger      *zap.Logger
}

func NewUserHandler(userService service.UserService, validate *validator.Validator, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validate,
		logger:      logger,
	}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return errors2.BadRequest("failed to parse request body")
	}
	if err := h.validate.Validate(req); err != nil {
		errs := h.validate.TranslateErrorsAuto(err, c.Get(fiber.HeaderAcceptLanguage))
		return errors2.PkgValidationError(errs)
	}

	return nil
}
