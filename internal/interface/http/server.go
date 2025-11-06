package http

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/ryuyb/fusion/internal/infrastructure/config"
	"github.com/ryuyb/fusion/internal/interface/http/middleware"
	"go.uber.org/zap"
)

func NewFiberApp(cfg *config.Config, logger *zap.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:       "Fusion",
		ServerHeader:  "Fusion Server",
		StrictRouting: false,
		CaseSensitive: false,
		ReadTimeout:   cfg.Server.ReadTimeout * time.Second,
		WriteTimeout:  cfg.Server.WriteTimeout * time.Second,
		ErrorHandler:  middleware.ErrorHandler(logger),
	})

	app.Use(requestid.New())
	app.Use(middleware.Cors())
	app.Use(middleware.Recovery(logger))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use(middleware.Logger(logger))

	return app
}
