package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/ryuyb/fusion/internal/infrastructure/provider/config"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewMigrator(cfg *config.Config, logger *zap.Logger) (*migrate.Migrate, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	m, err := migrate.New("", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to init migrator: %w", err)
	}
	return m, nil
}
