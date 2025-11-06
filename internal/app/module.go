package app

import (
	"github.com/ryuyb/fusion/internal/application"
	"github.com/ryuyb/fusion/internal/infrastructure/config"
	"github.com/ryuyb/fusion/internal/infrastructure/database"
	"github.com/ryuyb/fusion/internal/infrastructure/logger"
	"github.com/ryuyb/fusion/internal/infrastructure/repository"
	_interface "github.com/ryuyb/fusion/internal/interface"
	"github.com/ryuyb/fusion/internal/interface/http/route"
	"github.com/ryuyb/fusion/internal/pkg"
	"go.uber.org/fx"
)

var Module = fx.Module("app",
	config.Module,
	logger.Module,
	database.Module,
	repository.Module,
	pkg.Module,

	application.Module,

	_interface.Module,

	route.Module,

	fx.Provide(NewServer),
)
