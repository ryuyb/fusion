package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/ryuyb/fusion/internal/infrastructure/config"
	"github.com/ryuyb/fusion/internal/interface/http/middleware"
	"go.uber.org/zap"
)

func NewFiberApp(cfg *config.Config, logger *zap.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "Fusion",
		ServerHeader:          "Fusion Server",
		StrictRouting:         false,
		CaseSensitive:         false,
		DisableStartupMessage: false,
		ReadTimeout:           cfg.Server.ReadTimeout * time.Second,
		WriteTimeout:          cfg.Server.WriteTimeout * time.Second,
		ErrorHandler:          middleware.ErrorHandler(logger),
	})

	app.Use(requestid.New())
	app.Use(middleware.Cors())
	app.Use(middleware.Recovery(logger))
	app.Use(middleware.Logger(logger))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	return app
}
