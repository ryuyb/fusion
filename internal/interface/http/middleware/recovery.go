package middleware

import (
	"github.com/gofiber/fiber/v2"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"go.uber.org/zap"
)

func Recovery(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				logger.Error("Panic recovered",
					zap.Any("error", panicErr),
					zap.String("path", c.Path()),
					zap.String("method", c.Method()),
					zap.String("remote", c.IP()),
					zap.String("trace_id", c.GetRespHeader(fiber.HeaderXRequestID)),
					zap.Stack("stack"),
				)
				err = errors2.Internal(nil)
			}
		}()

		return c.Next()
	}
}
