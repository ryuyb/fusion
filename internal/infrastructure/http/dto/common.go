package dto

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

type ErrorResponse struct {
	Code         string         `json:"code"`          // app error code
	Message      string         `json:"message"`       // error message
	Details      map[string]any `json:"details"`       // error details
	TraceID      string         `json:"trace_id"`      // trace ID
	OccurredTime time.Time      `json:"occurred_time"` // occurred time
}

func NewErrorResponse(c fiber.Ctx, code, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:         code,
		Message:      message,
		TraceID:      c.GetRespHeader(fiber.HeaderXRequestID),
		OccurredTime: time.Now(),
	}
}

func (e *ErrorResponse) WithDetails(details map[string]any) *ErrorResponse {
	e.Details = details
	return e
}

type PaginationResponse[T any] struct {
	Data       []T `json:"data"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}

func NewPaginationResponse[T any](data []T, total, page, pageSize int) *PaginationResponse[T] {
	return &PaginationResponse[T]{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (total + pageSize - 1) / pageSize,
	}
}
