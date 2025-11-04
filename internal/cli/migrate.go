package cli

import (
	"context"
	"log"

	"github.com/ryuyb/fusion/internal/infrastructure/config"
	"github.com/ryuyb/fusion/internal/infrastructure/database"
	"github.com/ryuyb/fusion/internal/infrastructure/logger"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Execute database schema migrations using EntGO`,
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run migrations up",
	Run: func(cmd *cobra.Command, args []string) {
		runMigrationUp()
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback migrations",
	Run: func(cmd *cobra.Command, args []string) {
		runMigrationDown()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateUpCmd, migrateDownCmd)
}

func runMigrationUp() {
	app := fx.New(
		fx.NopLogger,
		logger.Module,
		config.Module,
		database.Module,

		fx.Invoke(func(client *database.Client) {
		}),
	)

	if err := app.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func runMigrationDown() {
	app := fx.New(
		fx.NopLogger,
		logger.Module,
		config.Module,

		fx.Invoke(func(zapLogger *zap.Logger) {
			zapLogger.Info("Rollback not implemented")
		}),
	)

	if err := app.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}
