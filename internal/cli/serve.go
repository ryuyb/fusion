package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ryuyb/fusion/internal/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"

	_ "github.com/ryuyb/fusion/internal/infrastructure/database/ent/runtime"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	Long:  `Start the HTTP API server with all dependencies`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().String("host", "0.0.0.0", "The host to bind the API server")
	serveCmd.Flags().IntP("port", "p", 8080, "Server port")

	_ = viper.BindPFlag("server.host", serveCmd.Flags().Lookup("host"))
	_ = viper.BindPFlag("server.port", serveCmd.Flags().Lookup("port"))
}

func startServer() {
	fxApp := fx.New(
		fx.NopLogger,
		app.Module,

		fx.Invoke(func(lifecycle fx.Lifecycle, server *app.Server) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						if err := server.Start(); err != nil {
							fmt.Printf("Failed to start server: %v\n", err)
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return server.Shutdown(ctx)
				},
			})
		}),
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	startCtx, cancel := context.WithTimeout(context.Background(), fxApp.StartTimeout())
	defer cancel()

	if err := fxApp.Start(startCtx); err != nil {
		fmt.Printf("Failed to start application: %v\n", err)
		os.Exit(1)
	}

	<-quit
	fmt.Println("Shutting down server...")

	stopCtx, cancel := context.WithTimeout(context.Background(), fxApp.StopTimeout())
	defer cancel()

	if err := fxApp.Stop(stopCtx); err != nil {
		fmt.Printf("Failed to stop application: %v\n", err)
		os.Exit(1)
	}
}
