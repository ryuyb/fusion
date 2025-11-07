package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/domain/service"
	"github.com/ryuyb/fusion/internal/interface/http/dto/request"
	"github.com/ryuyb/fusion/internal/interface/http/dto/response"
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
//	@Param		request	body		request.CreateUserRequest	true	"用户信息"
//	@Success	200		{object}	response.UserResponse
//	@Router		/user [post]
func (h *UserHandler) Create(c fiber.Ctx) error {
	var req request.CreateUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return errors2.BadRequest("failed to parse request body").Wrap(err)
	}
	if err := h.validate.Validate(req); err != nil {
		errs := h.validate.TranslateErrorsAuto(err, c.Get(fiber.HeaderAcceptLanguage))
		return errors2.PkgValidationError(errs)
	}
	created, err := h.userService.Create(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(h.toResponse(created))
}

// Update 更新用户
//
//	@Summary	更新用户
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		request	body		request.UpdateUserRequest	true	"用户信息"
//	@Success	200		{object}	response.UserResponse
//	@Router		/user [put]
func (h *UserHandler) Update(c fiber.Ctx) error {
	var req request.UpdateUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return errors2.BadRequest("failed to parse request body").Wrap(err)
	}
	if err := h.validate.Validate(req); err != nil {
		errs := h.validate.TranslateErrorsAuto(err, c.Get(fiber.HeaderAcceptLanguage))
		return errors2.PkgValidationError(errs)
	}
	updated, err := h.userService.Update(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.JSON(h.toResponse(updated))
}

// GetByID 获取用户详情
//
//	@Summary	获取用户详情
//	@Tags		User
//	@Produce	json
//	@Param		id	path		int	true	"用户ID"
//	@Success	200	{object}	response.UserResponse
//	@Router		/user/{id} [get]
func (h *UserHandler) GetByID(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errors2.BadRequest("failed to parse id as integer").Wrap(err)
	}
	user, err := h.userService.GetByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(h.toResponse(user))
}

// List 获取用户列表
//
//	@Summary	获取用户列表
//	@Tags		User
//	@Produce	json
//	@Param		page		query		int	false	"页码"	default(1)
//	@Param		page_size	query		int	false	"每页数量"	default(10)
//	@Success	200			{object}	response.PaginationResponse[response.UserResponse]
//	@Router		/user/list [get]
func (h *UserHandler) List(c fiber.Ctx) error {
	page := 1
	pageSize := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	users, total, err := h.userService.List(c.Context(), page, pageSize)
	if err != nil {
		return err
	}
	items := make([]*response.UserResponse, len(users))
	for i, user := range users {
		items[i] = h.toResponse(user)
	}
	resp := response.NewPaginationResponse[*response.UserResponse](items, total, page, pageSize)
	return c.JSON(resp)
}

// DeleteByID 删除用户
//
//	@Summary	删除用户
//	@Tags		User
//	@Produce	json
//	@Param		id	path		int	true	"用户ID"
//	@Success	200	{object}	nil
//	@Router		/user/{id} [delete]
func (h *UserHandler) DeleteByID(c fiber.Ctx) error {
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

func (h *UserHandler) toResponse(u *entity.User) *response.UserResponse {
	return &response.UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Status:    string(u.Status),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeleteAt:  u.DeleteAt,
	}
}
