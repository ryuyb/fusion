package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/infrastructure/provider/config"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func RunServeApp() error {
	fxLogger := fx.NopLogger
	if viper.GetBool("logger.fx.enable") {
		fxLogger = fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger.Named("fx")}
		})
	}

	app := fx.New(
		AppModule,

		fxLogger,

		fx.Invoke(func(lc fx.Lifecycle, app *fiber.App, cfg *config.Config, logger *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
						logger.Info("Starting server", zap.String("address", addr))
						if err := app.Listen(addr); err != nil {
							fmt.Printf("Failed to start server: %v\n", err)
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					logger.Info("Shutting down server gracefully")
					return app.ShutdownWithContext(ctx)
				},
			})
		}),
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	startCtx, cancel := context.WithTimeout(context.Background(), app.StartTimeout())
	defer cancel()

	if err := app.Start(startCtx); err != nil {
		fmt.Printf("Failed to start application: %v\n", err)
		os.Exit(1)
	}

	<-quit
	fmt.Println("Shutting down server...")

	stopCtx, cancel := context.WithTimeout(context.Background(), app.StopTimeout())
	defer cancel()

	return app.Stop(stopCtx)
}
