package database

import (
	"github.com/ryuyb/fusion/internal/infrastructure/database/repository"
	"go.uber.org/fx"
)

var Module = fx.Module("database",
	fx.Provide(NewClient),

	fx.Provide(
		repository.NewUserRepository,
	),
)

var MigrationModule = fx.Module("migration",
	fx.Provide(NewMigrator),
)
