package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/core/port/service"
)

type UserController struct {
	s service.UserService
}

func (u *UserController) Create(ctx fiber.Ctx) error {
	return nil
}

func (u *UserController) List(ctx fiber.Ctx) error {
	return nil
}

func NewUserController(s service.UserService) *UserController {
	return &UserController{s: s}
}
