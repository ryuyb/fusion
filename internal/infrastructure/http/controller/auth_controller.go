package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/core/command"
	"github.com/ryuyb/fusion/internal/core/port/service"
	"github.com/ryuyb/fusion/internal/infrastructure/http/dto"
	"github.com/ryuyb/fusion/internal/infrastructure/provider/jwt"
	"github.com/ryuyb/fusion/internal/pkg/util"
)

type AuthController struct {
	authService service.AuthService
	jwtManager  *jwt.JWTManager
}

// Register New user registration
//
//	@Summary	Register
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.RegisterRequest	true	"Register info"
//	@Success	200
//	@Router		/auth/register [post]
func (a *AuthController) Register(ctx fiber.Ctx) error {
	req := new(dto.RegisterRequest)
	if err := util.ParseRequestJson(ctx, req); err != nil {
		return err
	}
	cmd := command.RegisterCommand{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	if err := a.authService.Register(ctx, cmd); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusCreated)
}

// Login User Login
//
//	@Summary	Login
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dto.LoginRequest	true	"Login info"
//	@Success	200		{object}	dto.LoginResponse
//	@Router		/auth/login [post]
func (a *AuthController) Login(ctx fiber.Ctx) error {
	req := new(dto.LoginRequest)
	if err := util.ParseRequestJson(ctx, req); err != nil {
		return err
	}
	cmd := command.LoginCommand{
		Username: req.Username,
		Password: req.Password,
	}
	token, expiresAt, err := a.authService.Login(ctx, cmd)
	if err != nil {
		return err
	}
	return ctx.JSON(&dto.LoginResponse{
		Token:    token,
		ExpireAt: expiresAt,
	})
}

func NewAuthController(authService service.AuthService, jwtManager *jwt.JWTManager) *AuthController {
	return &AuthController{
		authService: authService,
		jwtManager:  jwtManager,
	}
}
