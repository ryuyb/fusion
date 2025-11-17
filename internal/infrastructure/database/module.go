package database

import (
	"github.com/ryuyb/fusion/internal/infrastructure/database/repository"
	"go.uber.org/fx"
)

var Module = fx.Module("database",
	fx.Provide(NewClient),

	fx.Provide(
		repository.NewUserRepository,
		repository.NewStreamingPlatformRepository,
		repository.NewStreamerRepository,
		repository.NewNotificationChannelRepository,
		repository.NewUserFollowedStreamerRepository,
	),
)

var MigrationModule = fx.Module("migration",
	fx.Provide(NewMigrator),
)
