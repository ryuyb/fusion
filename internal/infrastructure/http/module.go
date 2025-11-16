package http

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/swagger/v2"
	"github.com/ryuyb/fusion/internal/infrastructure/http/controller"
	"github.com/ryuyb/fusion/internal/infrastructure/http/middleware"
	"github.com/ryuyb/fusion/internal/infrastructure/http/router"
	"github.com/ryuyb/fusion/internal/infrastructure/provider/config"
	"github.com/ryuyb/fusion/internal/infrastructure/provider/validator"
	"go.uber.org/fx"
	"go.uber.org/zap"

	_ "github.com/ryuyb/fusion/docs/api"
)

var Module = fx.Module("http",
	fx.Provide(
		controller.NewUserController,
	),

	fx.Provide(NewFiberApp),

	router.Module,
)

func NewFiberApp(cfg *config.Config, logger *zap.Logger, routerRegistry *router.RouterRegistry, validate *validator.Validator) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:       cfg.App.Name,
		ServerHeader:  fmt.Sprintf("%s Server", cfg.App.Name),
		StrictRouting: false,
		CaseSensitive: false,
		ReadTimeout:   cfg.Server.ReadTimeout,
		WriteTimeout:  cfg.Server.WriteTimeout,
		ErrorHandler:  middleware.ErrorHandler(logger),
		StructValidator: &middleware.StructValidator{
			Validator: validate,
		},
	})

	app.Use(requestid.New())
	app.Use(middleware.Cors())
	app.Use(middleware.Recovery(logger.Named("recovery")))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use(middleware.Logger(logger))

	routerRegistry.RegisterAllRoutes(app)

	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}
