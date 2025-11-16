package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/infrastructure/http/dto"
	"github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ErrorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		appErr := errors.GetAppError(err)

		if appErr != nil {
			return handleAppError(c, appErr, logger)
		}

		if fiberErr, ok := lo.ErrorsAs[*fiber.Error](err); ok {
			return handleFiberError(c, fiberErr, logger)
		}

		return handleUnknownError(c, err, logger)
	}
}

func handleAppError(c fiber.Ctx, err *errors.AppError, logger *zap.Logger) error {
	logError(c, err, logger)

	errResp := dto.NewErrorResponse(c, string(err.Code), err.Message)

	if len(err.Details) > 0 {
		errResp.WithDetails(err.Details)
	}

	return c.Status(err.HTTPStatus).JSON(errResp)
}

func handleFiberError(c fiber.Ctx, err *fiber.Error, logger *zap.Logger) error {
	logger.Warn("Fiber error",
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.Int("status", err.Code),
		zap.String("ip", c.IP()),
		zap.String("trace_id", c.GetRespHeader(fiber.HeaderXRequestID)),
		zap.String("message", err.Message),
	)

	errResp := dto.NewErrorResponse(
		c,
		fmt.Sprintf("HTTP_%d", err.Code),
		err.Message,
	)

	return c.Status(err.Code).JSON(errResp)
}

func handleUnknownError(c fiber.Ctx, err error, logger *zap.Logger) error {
	logger.Error("Unknown error",
		zap.Error(err),
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.String("ip", c.IP()),
		zap.String("trace_id", c.GetRespHeader(fiber.HeaderXRequestID)),
	)

	errResp := dto.NewErrorResponse(
		c,
		"INTERNAL_ERROR",
		"An unexpected error occurred",
	)

	return c.Status(fiber.StatusInternalServerError).JSON(errResp)
}

func logError(c fiber.Ctx, err *errors.AppError, logger *zap.Logger) {
	logLevel := getLogLevel(err.HTTPStatus)

	fields := []zap.Field{
		zap.String("code", string(err.Code)),
		zap.String("message", err.Message),
		zap.Int("status", err.HTTPStatus),
		zap.String("path", c.Path()),
		zap.String("method", c.Method()),
		zap.String("ip", c.IP()),
		zap.String("trace_id", c.GetRespHeader(fiber.HeaderXRequestID)),
		zap.String("user_agent", c.Get(fiber.HeaderUserAgent)),
	}

	if err.Err != nil {
		fields = append(fields, zap.Error(err.Err))
	}

	if err.Details != nil {
		fields = append(fields, zap.Any("details", err.Details))
	}

	switch logLevel {
	case zap.ErrorLevel:
		logger.Error("Request error", fields...)
	case zap.WarnLevel:
		logger.Warn("Request warning", fields...)
	default:
		logger.Info("Request info", fields...)
	}
}

func getLogLevel(status int) zapcore.Level {
	switch {
	case status >= 500:
		return zap.ErrorLevel
	case status >= 400:
		return zap.WarnLevel
	default:
		return zap.InfoLevel
	}
}
