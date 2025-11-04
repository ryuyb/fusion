package database

import (
	"context"
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	"github.com/ryuyb/fusion/internal/infrastructure/config"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent"
	"go.uber.org/fx"
	"go.uber.org/zap"

	entsql "entgo.io/ent/dialect/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Client struct {
	ent.Client
}

func NewClient(lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger) (*Client, error) {
	db, err := sql.Open("pgx", cfg.Database.DSN)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	driver := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(driver), ent.Log(func(args ...any) {
		logger.Debug("ent query", zap.Any("query", args))
	}))

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Running database migrations...")
			if err = client.Schema.Create(ctx); err != nil {
				logger.Error("failed to create database schema", zap.Error(err))
				return fmt.Errorf("failed to create schema: %w", err)
			}
			logger.Info("Database migrations completed successfully")
			return nil
		},
		OnStop: func(context.Context) error {
			logger.Info("Closing database connection...")
			return client.Close()
		},
	})

	return &Client{*client}, nil
}
