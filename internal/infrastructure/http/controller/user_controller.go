package controller

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/app/errors"
	"github.com/ryuyb/fusion/internal/core/command"
	"github.com/ryuyb/fusion/internal/core/domain"
	"github.com/ryuyb/fusion/internal/core/port/service"
	"github.com/ryuyb/fusion/internal/infrastructure/http/dto"
	validator2 "github.com/ryuyb/fusion/internal/infrastructure/provider/validator"
	"github.com/samber/lo"
)

type UserController struct {
	service service.UserService
}

// Create Create new user
//
//	@Summary	Create user
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.CreateUserRequest	true	"User info"
//	@Security	Bearer
//	@Success	200	{object}	dto.UserResponse
//	@Router		/user [post]
func (u *UserController) Create(ctx fiber.Ctx) error {
	req := new(dto.CreateUserRequest)
	if err := ctx.Bind().JSON(req); err != nil {
		if errs, ok := lo.ErrorsAs[validator.ValidationErrors](err); ok {
			validationErrors := validator2.VALIDATOR.TranslateErrorsAuto(errs, ctx.Get(fiber.HeaderAcceptLanguage))
			return errors.CustomValidationError(validationErrors)
		}
		return errors.BadRequest("failed to parse request body").Wrap(err)
	}
	cmd := &command.CreateUserCommand{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	create, err := u.service.Create(ctx, cmd)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(u.toResponse(create))
}

// Update Update user info
//
//	@Summary	Update user
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		id		path	int						true	"User ID"
//	@Param		request	body	dto.UpdateUserRequest	true	"User info"
//	@Security	Bearer
//	@Success	200	{object}	dto.UserResponse
//	@Router		/user/{id} [put]
func (u *UserController) Update(ctx fiber.Ctx) error {
	req := new(dto.UpdateUserRequest)
	if err := ctx.Bind().JSON(req); err != nil {
		if errs, ok := lo.ErrorsAs[validator.ValidationErrors](err); ok {
			validationErrors := validator2.VALIDATOR.TranslateErrorsAuto(errs, ctx.Get(fiber.HeaderAcceptLanguage))
			return errors.CustomValidationError(validationErrors)
		}
		return errors.BadRequest("failed to parse request body").Wrap(err)
	}
	cmd := &command.UpdateUserCommand{
		ID: req.ID,
		CreateUserCommand: &command.CreateUserCommand{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		},
	}
	update, err := u.service.Update(ctx, cmd)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(u.toResponse(update))
}

// GetByID Get user by id
//
//	@Summary	Get user By ID
//	@Tags		User
//	@Produce	json
//	@Param		id	path	int	true	"User ID"
//	@Security	Bearer
//	@Success	200	{object}	dto.UserResponse
//	@Router		/user/{id} [get]
func (u *UserController) GetByID(ctx fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errors.BadRequest("invalid user id").Wrap(err)
	}
	user, err := u.service.FindById(ctx, id)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(u.toResponse(user))
}

// DeleteByID Delete user by user id
//
//	@Summary	Delete user by user ID
//	@Tags		User
//	@Produce	json
//	@Param		id	path	int	true	"User ID"
//	@Security	Bearer
//	@Success	200	{object}	nil
//	@Router		/user/{id} [delete]
func (u *UserController) DeleteByID(ctx fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errors.BadRequest("invalid user id").Wrap(err)
	}
	if err = u.service.Delete(ctx, id); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}

// List List all users by page
//
//	@Summary	List all users by page
//	@Tags		User
//	@Produce	json
//	@Param		page		query	int	false	"page"		default(1)
//	@Param		page_size	query	int	false	"page size"	default(10)
//	@Security	Bearer
//	@Success	200	{object}	dto.PaginationResponse[dto.UserResponse]
//	@Router		/user/list [get]
func (u *UserController) List(ctx fiber.Ctx) error {
	page, pageSize := 1, 10

	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if pageSizeStr := ctx.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps >= 0 {
			pageSize = ps
		}
	}

	users, total, err := u.service.List(ctx, page, pageSize)
	if err != nil {
		return err
	}

	items := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		items[i] = u.toResponse(user)
	}
	response := dto.NewPaginationResponse[*dto.UserResponse](items, total, page, pageSize)
	return ctx.JSON(response)
}

func (u *UserController) toResponse(user *domain.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{service: service}
}
