package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/domain/service"
	"github.com/ryuyb/fusion/internal/interface/http/dto/request"
	"github.com/ryuyb/fusion/internal/interface/http/dto/response"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/ryuyb/fusion/internal/pkg/validator"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService service.AuthService
	validate    *validator.Validator
	logger      *zap.Logger
}

func NewAuthHandler(authService service.AuthService, validate *validator.Validator, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validate,
		logger:      logger,
	}
}

// Login authenticates a user
//
//	@Summary	User login
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		request	body		request.LoginRequest	true	"Login credentials"
//	@Success	200		{object}	response.TokenResponse
//	@Failure	400		{object}	response.ErrorResponse
//	@Failure	401		{object}	response.ErrorResponse
//	@Failure	422		{object}	response.ErrorResponse
//	@Router		/auth/login [post]
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req request.LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return errors2.BadRequest("failed to parse request body").Wrap(err)
	}

	if err := h.validate.Validate(req); err != nil {
		errs := h.validate.TranslateErrorsAuto(err, c.Get(fiber.HeaderAcceptLanguage))
		return errors2.PkgValidationError(errs)
	}

	token, expiresAt, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.JSON(&response.TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	})
}

// Register creates a new user account
//
//	@Summary	User registration
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		request	body		request.RegisterRequest	true	"Registration data"
//	@Success	201		{object}	response.TokenResponse
//	@Failure	400		{object}	response.ErrorResponse
//	@Failure	409		{object}	response.ErrorResponse
//	@Failure	422		{object}	response.ErrorResponse
//	@Router		/auth/register [post]
func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req request.RegisterRequest
	if err := c.Bind().JSON(&req); err != nil {
		return errors2.BadRequest("failed to parse request body").Wrap(err)
	}

	if err := h.validate.Validate(req); err != nil {
		errs := h.validate.TranslateErrorsAuto(err, c.Get(fiber.HeaderAcceptLanguage))
		return errors2.PkgValidationError(errs)
	}

	token, expiresAt, err := h.authService.Register(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(&response.TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	})
}
