package pkg

import (
	"github.com/ryuyb/fusion/internal/pkg/auth"
	"go.uber.org/fx"
)

var Module = fx.Module("route",
	fx.Provide(
		auth.NewJWTManager,
	),
)
