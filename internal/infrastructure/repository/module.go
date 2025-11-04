package repository

import (
	//domainRepo "github.com/ryuyb/fusion/internal/domain/repository"
	"go.uber.org/fx"
)

var Module = fx.Module("repository",
	fx.Provide(
		NewUserRepository,
	),
)
