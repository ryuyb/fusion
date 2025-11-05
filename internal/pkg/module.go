package pkg

import (
	"github.com/ryuyb/fusion/internal/pkg/auth"
	"github.com/ryuyb/fusion/internal/pkg/validator"
	"go.uber.org/fx"
)

var Module = fx.Module("route",
	fx.Provide(
		auth.NewJWTManager,

		validator.NewValidator,
	),
)
