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

type StreamingPlatformController struct {
	service coreService.StreamingPlatformService
}

func NewStreamingPlatformController(service coreService.StreamingPlatformService) *StreamingPlatformController {
	return &StreamingPlatformController{service: service}
}

// Create creates a streaming platform
//
//	@Summary	Create Streaming Platform
//	@Tags		StreamingPlatform
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.CreateStreamingPlatformRequest	true	"Platform data"
//	@Security	Bearer
//	@Success	201	{object}	dto.StreamingPlatformResponse
//	@Router		/platforms [post]
func (c *StreamingPlatformController) Create(ctx fiber.Ctx) error {
	req := new(dto.CreateStreamingPlatformRequest)
	if err := ctx.Bind().JSON(req); err != nil {
		return errors.BadRequest("failed to parse request body").Wrap(err)
	}
	cmd := &command.CreateStreamingPlatformCommand{
		Type:        req.Type,
		Name:        req.Name,
		Description: req.Description,
		BaseURL:     req.BaseURL,
		LogoURL:     req.LogoURL,
		Enabled:     req.Enabled,
		Priority:    req.Priority,
		Metadata:    req.Metadata,
	}

	created, err := c.service.Create(ctx, cmd)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(c.toResponse(created))
}

// Update updates a streaming platform
//
//	@Summary	Update Streaming Platform
//	@Tags		StreamingPlatform
//	@Accept		json
//	@Produce	json
//	@Param		id		path	int									true	"Platform ID"
//	@Param		request	body	dto.UpdateStreamingPlatformRequest	true	"Platform data"
//	@Security	Bearer
//	@Success	200	{object}	dto.StreamingPlatformResponse
//	@Router		/platforms/{id} [put]
func (c *StreamingPlatformController) Update(ctx fiber.Ctx) error {
	req := new(dto.UpdateStreamingPlatformRequest)
	if err := ctx.Bind().JSON(req); err != nil {
		return errors.BadRequest("failed to parse request body").Wrap(err)
	}
	cmd := &command.UpdateStreamingPlatformCommand{
		ID: req.ID,
		CreateStreamingPlatformCommand: &command.CreateStreamingPlatformCommand{
			Type:        req.Type,
			Name:        req.Name,
			Description: req.Description,
			BaseURL:     req.BaseURL,
			LogoURL:     req.LogoURL,
			Enabled:     req.Enabled,
			Priority:    req.Priority,
			Metadata:    req.Metadata,
		},
	}

	updated, err := c.service.Update(ctx, cmd)
	if err != nil {
		return err
	}
	return ctx.JSON(c.toResponse(updated))
}

// Delete deletes a streaming platform
//
//	@Summary	Delete Streaming Platform
//	@Tags		StreamingPlatform
//	@Produce	json
//	@Param		id	path	int	true	"Platform ID"
//	@Security	Bearer
//	@Success	200
//	@Router		/platforms/{id} [delete]
func (c *StreamingPlatformController) Delete(ctx fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errors.BadRequest("invalid platform id").Wrap(err)
	}
	if err := c.service.Delete(ctx, id); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}

// GetByID fetches a streaming platform by id
//
//	@Summary	Get Streaming Platform
//	@Tags		StreamingPlatform
//	@Produce	json
//	@Param		id	path	int	true	"Platform ID"
//	@Security	Bearer
//	@Success	200	{object}	dto.StreamingPlatformResponse
//	@Router		/platforms/{id} [get]
func (c *StreamingPlatformController) GetByID(ctx fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errors.BadRequest("invalid platform id").Wrap(err)
	}
	platform, err := c.service.FindById(ctx, id)
	if err != nil {
		return err
	}
	return ctx.JSON(c.toResponse(platform))
}

// List lists streaming platforms
//
//	@Summary	List Streaming Platforms
//	@Tags		StreamingPlatform
//	@Produce	json
//	@Param		page		query	int	false	"Page"		default(1)
//	@Param		page_size	query	int	false	"Page size"	default(10)
//	@Security	Bearer
//	@Success	200	{object}	dto.PaginationResponse[dto.StreamingPlatformResponse]
//	@Router		/platforms [get]
func (c *StreamingPlatformController) List(ctx fiber.Ctx) error {
	page, pageSize := util.ParsePagination(ctx)
	platforms, total, err := c.service.List(ctx, page, pageSize)
	if err != nil {
		return err
	}
	items := make([]*dto.StreamingPlatformResponse, len(platforms))
	for i, platform := range platforms {
		items[i] = c.toResponse(platform)
	}
	return ctx.JSON(dto.NewPaginationResponse(items, total, page, pageSize))
}

func (c *StreamingPlatformController) toResponse(platform *domain.StreamingPlatform) *dto.StreamingPlatformResponse {
	return &dto.StreamingPlatformResponse{
		ID:          platform.ID,
		Type:        string(platform.Type),
		Name:        platform.Name,
		Description: platform.Description,
		BaseURL:     platform.BaseURL,
		LogoURL:     platform.LogoURL,
		Enabled:     platform.Enabled,
		Priority:    platform.Priority,
		Metadata:    platform.Metadata,
	}
}
