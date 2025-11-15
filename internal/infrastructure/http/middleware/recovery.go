package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/ryuyb/fusion/internal/app/errors"
	"go.uber.org/zap"
)

func Recovery(logger *zap.Logger) fiber.Handler {
	return func(c fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Panic recovered",
					zap.Any("error", r),
					zap.String("path", c.Path()),
					zap.String("method", c.Method()),
					zap.String("trace_id", requestid.FromContext(c)),
					zap.Stack("stack"),
				)
				err = errors.Internal(nil)
			}
		}()

		return c.Next()
	}
}
