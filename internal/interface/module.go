package _interface

import (
	"github.com/ryuyb/fusion/internal/interface/http"
	"github.com/ryuyb/fusion/internal/interface/http/handler"
	"github.com/ryuyb/fusion/internal/interface/http/middleware"
	"go.uber.org/fx"
)

var Module = fx.Module("http",
	fx.Provide(http.NewFiberApp),

	fx.Provide(middleware.NewAuth),

	fx.Provide(
		handler.NewHealthHandler,
		handler.NewUserHandler,
		handler.NewAuthHandler,
	),
)
