package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/core/command"
	"github.com/ryuyb/fusion/internal/core/domain"
	coreService "github.com/ryuyb/fusion/internal/core/port/service"
	"github.com/ryuyb/fusion/internal/infrastructure/http/dto"
	"github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/ryuyb/fusion/internal/pkg/util"
)

type StreamerController struct {
	service coreService.StreamerService
}

func NewStreamerController(service coreService.StreamerService) *StreamerController {
	return &StreamerController{service: service}
}

// Create creates a streamer
//
//	@Summary	Create Streamer
//	@Tags		Streamer
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.CreateStreamerRequest	true	"Streamer data"
//	@Security	Bearer
//	@Success	201	{object}	dto.StreamerResponse
//	@Router		/streamers [post]
func (c *StreamerController) Create(ctx fiber.Ctx) error {
	req := new(dto.CreateStreamerRequest)
	if err := util.ParseRequestJson(ctx, req); err != nil {
		return err
	}
	cmd := &command.CreateStreamerCommand{
		PlatformType:       req.PlatformType,
		PlatformStreamerID: req.PlatformStreamerID,
		DisplayName:        req.DisplayName,
		AvatarURL:          req.AvatarURL,
		RoomURL:            req.RoomURL,
		Bio:                req.Bio,
		Tags:               req.Tags,
	}

	created, err := c.service.Create(ctx, cmd)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(c.toResponse(created))
}

// Update updates a streamer
//
//	@Summary	Update Streamer
//	@Tags		Streamer
//	@Accept		json
//	@Produce	json
//	@Param		id		path	int							true	"Streamer ID"
//	@Param		request	body	dto.UpdateStreamerRequest	true	"Streamer data"
//	@Security	Bearer
//	@Success	200	{object}	dto.StreamerResponse
//	@Router		/streamers/{id} [put]
func (c *StreamerController) Update(ctx fiber.Ctx) error {
	req := new(dto.UpdateStreamerRequest)
	if err := util.ParseRequestJson(ctx, req); err != nil {
		return err
	}
	cmd := &command.UpdateStreamerCommand{
		ID: req.ID,
		CreateStreamerCommand: &command.CreateStreamerCommand{
			PlatformType:       req.PlatformType,
			PlatformStreamerID: req.PlatformStreamerID,
			DisplayName:        req.DisplayName,
			AvatarURL:          req.AvatarURL,
			RoomURL:            req.RoomURL,
			Bio:                req.Bio,
			Tags:               req.Tags,
		},
	}

	updated, err := c.service.Update(ctx, cmd)
	if err != nil {
		return err
	}
	return ctx.JSON(c.toResponse(updated))
}

// Delete deletes a streamer
//
//	@Summary	Delete Streamer
//	@Tags		Streamer
//	@Produce	json
//	@Param		id	path	int	true	"Streamer ID"
//	@Security	Bearer
//	@Success	200
//	@Router		/streamers/{id} [delete]
func (c *StreamerController) Delete(ctx fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errors.BadRequest("invalid streamer id").Wrap(err)
	}
	if err := c.service.Delete(ctx, id); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}

// GetByID fetches a streamer by id
//
//	@Summary	Get Streamer
//	@Tags		Streamer
//	@Produce	json
//	@Param		id	path	int	true	"Streamer ID"
//	@Security	Bearer
//	@Success	200	{object}	dto.StreamerResponse
//	@Router		/streamers/{id} [get]
func (c *StreamerController) GetByID(ctx fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errors.BadRequest("invalid streamer id").Wrap(err)
	}
	streamer, err := c.service.FindById(ctx, id)
	if err != nil {
		return err
	}
	return ctx.JSON(c.toResponse(streamer))
}

// GetByPlatformStreamerID gets a streamer by platform type and platform streamer id
//
//	@Summary	Get Streamer By Platform
//	@Tags		Streamer
//	@Produce	json
//	@Param		platform_type			path	domain.StreamingPlatformType	true	"Platform type"
//	@Param		platform_streamer_id	path	string							true	"Platform streamer ID"
//	@Security	Bearer
//	@Param		refresh	query		bool	false	"Refresh from platform and update database"
//	@Success	200		{object}	dto.StreamerResponse
//	@Router		/streamers/{platform_type}/{platform_streamer_id} [get]
func (c *StreamerController) GetByPlatformStreamerID(ctx fiber.Ctx) error {
	platformType := domain.StreamingPlatformType(ctx.Params("platform_type"))
	platformStreamerID := ctx.Params("platform_streamer_id")
	refresh, _ := strconv.ParseBool(ctx.Query("refresh"))
	streamer, err := c.service.FindByPlatformStreamerId(ctx, platformType, platformStreamerID, refresh)
	if err != nil {
		return err
	}
	return ctx.JSON(c.toResponse(streamer))
}

// List lists streamers
//
//	@Summary	List Streamers
//	@Tags		Streamer
//	@Produce	json
//	@Param		page		query	int	false	"Page"		default(1)
//	@Param		page_size	query	int	false	"Page size"	default(10)
//	@Security	Bearer
//	@Success	200	{object}	dto.PaginationResponse[dto.StreamerResponse]
//	@Router		/streamers [get]
func (c *StreamerController) List(ctx fiber.Ctx) error {
	page, pageSize := util.ParsePagination(ctx)
	streamers, total, err := c.service.List(ctx, page, pageSize)
	if err != nil {
		return err
	}
	items := make([]*dto.StreamerResponse, len(streamers))
	for i, streamer := range streamers {
		items[i] = c.toResponse(streamer)
	}
	return ctx.JSON(dto.NewPaginationResponse(items, total, page, pageSize))
}

func (c *StreamerController) toResponse(streamer *domain.Streamer) *dto.StreamerResponse {
	return &dto.StreamerResponse{
		ID:                 streamer.ID,
		PlatformType:       string(streamer.PlatformType),
		PlatformStreamerID: streamer.PlatformStreamerID,
		DisplayName:        streamer.DisplayName,
		AvatarURL:          streamer.AvatarURL,
		RoomURL:            streamer.RoomURL,
		Bio:                streamer.Bio,
		Tags:               streamer.Tags,
	}
}
