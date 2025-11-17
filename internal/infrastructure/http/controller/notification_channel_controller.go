package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/core/domain"
	coreService "github.com/ryuyb/fusion/internal/core/port/service"
	"github.com/ryuyb/fusion/internal/infrastructure/http/dto"
	"github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/ryuyb/fusion/internal/pkg/util"
)

type NotificationChannelController struct {
	service coreService.NotificationChannelService
}

func NewNotificationChannelController(service coreService.NotificationChannelService) *NotificationChannelController {
	return &NotificationChannelController{service: service}
}

// Create creates a notification channel
//
//	@Summary	Create Notification Channel
//	@Tags		NotificationChannel
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.CreateNotificationChannelRequest	true	"Channel data"
//	@Security	Bearer
//	@Success	201	{object}	dto.NotificationChannelResponse
//	@Router		/notification-channels [post]
func (c *NotificationChannelController) Create(ctx fiber.Ctx) error {
	req := new(dto.CreateNotificationChannelRequest)
	if err := ctx.Bind().JSON(req); err != nil {
		return errors.BadRequest("failed to parse request body").Wrap(err)
	}
	channel := &domain.NotificationChannel{
		UserID:      req.UserID,
		ChannelType: domain.NotificationChannelType(req.ChannelType),
		Name:        req.Name,
		Config:      req.Config,
		Enable:      req.Enable,
		Priority:    req.Priority,
	}
	created, err := c.service.Create(ctx, channel)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(c.toResponse(created))
}

// Update updates a notification channel
//
//	@Summary	Update Notification Channel
//	@Tags		NotificationChannel
//	@Accept		json
//	@Produce	json
//	@Param		id		path	int										true	"Channel ID"
//	@Param		request	body	dto.UpdateNotificationChannelRequest	true	"Channel data"
//	@Security	Bearer
//	@Success	200	{object}	dto.NotificationChannelResponse
//	@Router		/notification-channels/{id} [put]
func (c *NotificationChannelController) Update(ctx fiber.Ctx) error {
	req := new(dto.UpdateNotificationChannelRequest)
	if err := ctx.Bind().JSON(req); err != nil {
		return errors.BadRequest("failed to parse request body").Wrap(err)
	}
	channel := &domain.NotificationChannel{
		ID:          req.ID,
		UserID:      req.UserID,
		ChannelType: domain.NotificationChannelType(req.ChannelType),
		Name:        req.Name,
		Config:      req.Config,
		Enable:      req.Enable,
		Priority:    req.Priority,
	}
	updated, err := c.service.Update(ctx, channel)
	if err != nil {
		return err
	}
	return ctx.JSON(c.toResponse(updated))
}

// Delete deletes a notification channel
//
//	@Summary	Delete Notification Channel
//	@Tags		NotificationChannel
//	@Produce	json
//	@Param		id	path	int	true	"Channel ID"
//	@Security	Bearer
//	@Success	200
//	@Router		/notification-channels/{id} [delete]
func (c *NotificationChannelController) Delete(ctx fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errors.BadRequest("invalid channel id").Wrap(err)
	}
	if err := c.service.Delete(ctx, id); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}

// GetByID fetches a notification channel
//
//	@Summary	Get Notification Channel
//	@Tags		NotificationChannel
//	@Produce	json
//	@Param		id	path	int	true	"Channel ID"
//	@Security	Bearer
//	@Success	200	{object}	dto.NotificationChannelResponse
//	@Router		/notification-channels/{id} [get]
func (c *NotificationChannelController) GetByID(ctx fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errors.BadRequest("invalid channel id").Wrap(err)
	}
	channel, err := c.service.FindById(ctx, id)
	if err != nil {
		return err
	}
	return ctx.JSON(c.toResponse(channel))
}

// ListByUser lists a user's notification channels
//
//	@Summary	List Notification Channels By User
//	@Tags		NotificationChannel
//	@Produce	json
//	@Param		user_id	path	int	true	"User ID"
//	@Param		page	query	int	false	"Page"			default(1)
//	@Param		page_size	query	int	false	"Page size"	default(10)
//	@Security	Bearer
//	@Success	200	{object}	dto.PaginationResponse[dto.NotificationChannelResponse]
//	@Router		/notification-channels/users/{user_id} [get]
func (c *NotificationChannelController) ListByUser(ctx fiber.Ctx) error {
	userID, err := strconv.ParseInt(ctx.Params("user_id"), 10, 64)
	if err != nil {
		return errors.BadRequest("invalid user id").Wrap(err)
	}
	page, pageSize := util.ParsePagination(ctx)
	channels, total, err := c.service.ListByUserId(ctx, userID, page, pageSize)
	if err != nil {
		return err
	}
	items := make([]*dto.NotificationChannelResponse, len(channels))
	for i, channel := range channels {
		items[i] = c.toResponse(channel)
	}
	return ctx.JSON(dto.NewPaginationResponse(items, total, page, pageSize))
}

func (c *NotificationChannelController) toResponse(channel *domain.NotificationChannel) *dto.NotificationChannelResponse {
	return &dto.NotificationChannelResponse{
		ID:          channel.ID,
		UserID:      channel.UserID,
		ChannelType: string(channel.ChannelType),
		Name:        channel.Name,
		Config:      channel.Config,
		Enable:      channel.Enable,
		Priority:    channel.Priority,
	}
}
