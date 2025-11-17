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

type UserFollowedStreamerController struct {
	service coreService.UserFollowedStreamerService
}

func NewUserFollowedStreamerController(service coreService.UserFollowedStreamerService) *UserFollowedStreamerController {
	return &UserFollowedStreamerController{service: service}
}

// Create creates a follow relationship
//
//	@Summary	Create User Followed Streamer
//	@Tags		UserFollow
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.CreateUserFollowedStreamerRequest	true	"Follow data"
//	@Security	Bearer
//	@Success	201	{object}	dto.UserFollowedStreamerResponse
//	@Router		/follows [post]
func (c *UserFollowedStreamerController) Create(ctx fiber.Ctx) error {
	req := new(dto.CreateUserFollowedStreamerRequest)
	if err := ctx.Bind().JSON(req); err != nil {
		return errors.BadRequest("failed to parse request body").Wrap(err)
	}
	follow, err := domain.NewUserFollowedStreamer(req.UserID, req.StreamerID, req.Alias, req.Notes, req.NotificationChannelIDs)
	if err != nil {
		return err
	}
	follow.NotificationsEnabled = req.NotificationsEnabled

	created, err := c.service.Create(ctx, follow)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(c.toResponse(created))
}

// Update updates a follow relationship
//
//	@Summary	Update User Followed Streamer
//	@Tags		UserFollow
//	@Accept		json
//	@Produce	json
//	@Param		id		path	int										true	"Follow ID"
//	@Param		request	body	dto.UpdateUserFollowedStreamerRequest	true	"Follow data"
//	@Security	Bearer
//	@Success	200	{object}	dto.UserFollowedStreamerResponse
//	@Router		/follows/{id} [put]
func (c *UserFollowedStreamerController) Update(ctx fiber.Ctx) error {
	req := new(dto.UpdateUserFollowedStreamerRequest)
	if err := ctx.Bind().JSON(req); err != nil {
		return errors.BadRequest("failed to parse request body").Wrap(err)
	}
	follow := &domain.UserFollowedStreamer{
		ID:                     req.ID,
		Alias:                  req.Alias,
		Notes:                  req.Notes,
		NotificationsEnabled:   req.NotificationsEnabled,
		NotificationChannelIDs: req.NotificationChannelIDs,
	}
	updated, err := c.service.Update(ctx, follow)
	if err != nil {
		return err
	}
	return ctx.JSON(c.toResponse(updated))
}

// Delete deletes a follow relationship
//
//	@Summary	Delete User Followed Streamer
//	@Tags		UserFollow
//	@Produce	json
//	@Param		id	path	int	true	"Follow ID"
//	@Security	Bearer
//	@Success	200
//	@Router		/follows/{id} [delete]
func (c *UserFollowedStreamerController) Delete(ctx fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errors.BadRequest("invalid follow id").Wrap(err)
	}
	if err := c.service.Delete(ctx, id); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}

// GetByID fetches a follow relationship
//
//	@Summary	Get User Followed Streamer
//	@Tags		UserFollow
//	@Produce	json
//	@Param		id	path	int	true	"Follow ID"
//	@Security	Bearer
//	@Success	200	{object}	dto.UserFollowedStreamerResponse
//	@Router		/follows/{id} [get]
func (c *UserFollowedStreamerController) GetByID(ctx fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errors.BadRequest("invalid follow id").Wrap(err)
	}
	follow, err := c.service.FindById(ctx, id)
	if err != nil {
		return err
	}
	return ctx.JSON(c.toResponse(follow))
}

// ListByUser lists follows by user
//
//	@Summary	List Follows By User
//	@Tags		UserFollow
//	@Produce	json
//	@Param		user_id		path	int	true	"User ID"
//	@Param		page		query	int	false	"Page"		default(1)
//	@Param		page_size	query	int	false	"Page size"	default(10)
//	@Security	Bearer
//	@Success	200	{object}	dto.PaginationResponse[dto.UserFollowedStreamerResponse]
//	@Router		/follows/users/{user_id} [get]
func (c *UserFollowedStreamerController) ListByUser(ctx fiber.Ctx) error {
	userID, err := strconv.ParseInt(ctx.Params("user_id"), 10, 64)
	if err != nil {
		return errors.BadRequest("invalid user id").Wrap(err)
	}
	page, pageSize := util.ParsePagination(ctx)
	follows, total, err := c.service.ListByUserId(ctx, userID, page, pageSize)
	if err != nil {
		return err
	}
	items := make([]*dto.UserFollowedStreamerResponse, len(follows))
	for i, follow := range follows {
		items[i] = c.toResponse(follow)
	}
	return ctx.JSON(dto.NewPaginationResponse(items, total, page, pageSize))
}

// ListByStreamer lists follows by streamer
//
//	@Summary	List Follows By Streamer
//	@Tags		UserFollow
//	@Produce	json
//	@Param		streamer_id	path	int	true	"Streamer ID"
//	@Param		page		query	int	false	"Page"		default(1)
//	@Param		page_size	query	int	false	"Page size"	default(10)
//	@Security	Bearer
//	@Success	200	{object}	dto.PaginationResponse[dto.UserFollowedStreamerResponse]
//	@Router		/follows/streamers/{streamer_id} [get]
func (c *UserFollowedStreamerController) ListByStreamer(ctx fiber.Ctx) error {
	streamerID, err := strconv.ParseInt(ctx.Params("streamer_id"), 10, 64)
	if err != nil {
		return errors.BadRequest("invalid streamer id").Wrap(err)
	}
	page, pageSize := util.ParsePagination(ctx)
	follows, total, err := c.service.ListByStreamerId(ctx, streamerID, page, pageSize)
	if err != nil {
		return err
	}
	items := make([]*dto.UserFollowedStreamerResponse, len(follows))
	for i, follow := range follows {
		items[i] = c.toResponse(follow)
	}
	return ctx.JSON(dto.NewPaginationResponse(items, total, page, pageSize))
}

func (c *UserFollowedStreamerController) toResponse(follow *domain.UserFollowedStreamer) *dto.UserFollowedStreamerResponse {
	return &dto.UserFollowedStreamerResponse{
		ID:                     follow.ID,
		UserID:                 follow.UserID,
		StreamerID:             follow.StreamerID,
		Alias:                  follow.Alias,
		Notes:                  follow.Notes,
		NotificationsEnabled:   follow.NotificationsEnabled,
		NotificationChannelIDs: follow.NotificationChannelIDs,
	}
}
