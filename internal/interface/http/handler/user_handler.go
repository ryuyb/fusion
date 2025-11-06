package handler

import (
	"strconv"

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

// Create 创建用户
//
//	@Summary	创建用户
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dto.CreateUserRequest	true	"用户信息"
//	@Success	200		{object}	dto.CreateUserRequest
//	@Router		/api/v1/user/create [post]
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return errors2.BadRequest("failed to parse request body")
	}
	if err := h.validate.Validate(req); err != nil {
		errs := h.validate.TranslateErrorsAuto(err, c.Get(fiber.HeaderAcceptLanguage))
		return errors2.PkgValidationError(errs)
	}
	created, err := h.userService.Create(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h *UserHandler) DeleteByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errors2.BadRequest("failed to parse id as integer")
	}
	err = h.userService.Delete(c.Context(), id)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}
