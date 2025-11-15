package database

import (
	"context"
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent"
	"github.com/ryuyb/fusion/internal/infrastructure/provider/config"
	"go.uber.org/fx"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewClient(cfg *config.Config, logger *zap.Logger, lc fx.Lifecycle) (*ent.Client, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	driver := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(driver), ent.Log(func(args ...any) {
		logger.Debug("ent query", zap.Any("query", args))
	}))

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Running database schema migrations...")
			if err = client.Schema.Create(ctx); err != nil {
				logger.Error("failed to create database schema", zap.Error(err))
				return fmt.Errorf("failed to create schema: %w", err)
			}
			logger.Info("Database schema migrations completed successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing database connection...")
			return client.Close()
		},
	})

	return client, nil
}
