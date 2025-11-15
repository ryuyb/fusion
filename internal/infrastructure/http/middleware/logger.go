package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"go.uber.org/zap"
)

func Logger(logger *zap.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		path := c.Path()
		method := c.Method()

		err := c.Next()

		latency := time.Since(start)
		statusCode := c.Response().StatusCode()

		fields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get(fiber.HeaderUserAgent)),
			zap.String("trace_id", requestid.FromContext(c)),
		}
		logger.Info("HTTP Request", fields...)

		return err
	}
}
