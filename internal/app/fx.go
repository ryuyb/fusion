package app

import (
	"github.com/ryuyb/fusion/internal/application"
	"github.com/ryuyb/fusion/internal/infrastructure/database"
	"github.com/ryuyb/fusion/internal/infrastructure/http"
	"github.com/ryuyb/fusion/internal/infrastructure/provider/config"
	"github.com/ryuyb/fusion/internal/infrastructure/provider/logger"
	"github.com/ryuyb/fusion/internal/infrastructure/provider/validator"
	"go.uber.org/fx"
)

var AppModule = fx.Module("app",
	config.Module,
	logger.Module,
	validator.Module,

	database.Module,
	http.Module,

	application.Module,
)
