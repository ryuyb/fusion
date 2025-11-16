package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/ryuyb/fusion/internal/infrastructure/provider/config"
	"go.uber.org/zap"

	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewMigrator(cfg *config.Config, logger *zap.Logger) (*migrate.Migrate, error) {
	dsn := fmt.Sprintf(
		"pgx5://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)
	logger.Info("migrating database", zap.String("dsn", dsn))

	m, err := migrate.New("file://./migrations", dsn)
	if err != nil {
		logger.Error("failed to create migrator", zap.Error(err))
		return nil, fmt.Errorf("failed to init migrator: %w", err)
	}
	return m, nil
}
