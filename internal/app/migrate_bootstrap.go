package app

import (
	"context"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/ryuyb/fusion/internal/infrastructure/database"
	"github.com/ryuyb/fusion/internal/infrastructure/provider/config"
	"github.com/ryuyb/fusion/internal/infrastructure/provider/logger"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type MigrateDirection string

const (
	MigrateUp   MigrateDirection = "up"
	MigrateDown MigrateDirection = "down"
)

func RunMigrateApp(direction MigrateDirection) error {
	migrator, fxErr := getMigrator()
	if fxErr != nil {
		return fxErr
	}

	var err error
	log := zap.L()

	switch direction {
	case MigrateUp:
		log.Info("Running up migrations")
		err = migrator.Up()
	case MigrateDown:
		_, _, vErr := migrator.Version()
		if errors.Is(vErr, migrate.ErrNilVersion) {
			log.Info("No migrations to apply")
			return nil
		}
		log.Info("Running down migrations")
		err = migrator.Steps(-1)
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Error("Migration failed", zap.Error(err))
		return err
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Info("No migrations to apply")
	} else {
		log.Info("Migrations applied")
	}

	return nil
}

func RunMigrateVersionApp() error {
	migrator, fxErr := getMigrator()
	if fxErr != nil {
		return fxErr
	}
	log := zap.L()

	version, dirty, fxErr := migrator.Version()

	if errors.Is(fxErr, migrate.ErrNilVersion) {
		log.Info("No migrations to apply")
		return nil
	}
	if fxErr != nil {
		return fxErr
	}
	log.Info("Current migration version",
		zap.Uint("version", version),
		zap.Bool("dirty", dirty),
	)

	return nil
}

func getMigrator() (*migrate.Migrate, error) {
	var migrator *migrate.Migrate

	fxLogger := fx.NopLogger
	if viper.GetBool("logger.fx.enable") {
		fxLogger = fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger.Named("fx-migrate")}
		})
	}

	app := fx.New(
		config.Module,
		logger.Module,
		database.MigrationModule,

		fx.Populate(&migrator),

		fxLogger,
	)

	startCtx, cancel := context.WithTimeout(context.Background(), fx.DefaultTimeout)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		return nil, err
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), fx.DefaultTimeout)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		return nil, err
	}

	return migrator, nil
}
