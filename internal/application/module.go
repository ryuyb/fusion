package application

import (
	"github.com/ryuyb/fusion/internal/application/service"
	"go.uber.org/fx"
)

var Module = fx.Module("service",
	fx.Provide(
		service.NewUserService,
	),
)
