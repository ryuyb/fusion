package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/ryuyb/fusion/internal/interface/http/response"
	"go.uber.org/zap"
)

func Recovery(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Path()),
					zap.String("method", c.Method()),
					zap.String("remote", c.IP()),
					zap.String("trace_id", c.GetRespHeader(fiber.HeaderXRequestID)),
					zap.Stack("stack"),
				)

				err = response.NewErrorResponse(
					c,
					http.StatusInternalServerError,
					http.StatusText(http.StatusInternalServerError),
					"Internal Server Error",
				)
			}
		}()

		return c.Next()
	}
}
