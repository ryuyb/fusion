package application

import (
	"github.com/ryuyb/fusion/internal/application/service"
	domainService "github.com/ryuyb/fusion/internal/domain/service"
	"go.uber.org/fx"
)

var Module = fx.Module("service",
	fx.Provide(
		fx.Annotate(service.NewUserService, fx.As(new(domainService.UserService))),
		fx.Annotate(service.NewAuthService, fx.As(new(domainService.AuthService))),
		fx.Annotate(service.NewPlatformService, fx.As(new(domainService.PlatformService))),
	),
)
