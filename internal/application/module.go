package application

import (
	"github.com/ryuyb/fusion/internal/application/service"
	"go.uber.org/fx"
)

var Module = fx.Module("application",
	fx.Provide(
		service.NewUserService,
		service.NewAuthService,
	),
)
