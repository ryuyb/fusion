package app

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ryuyb/fusion/internal/infrastructure/config"
	"go.uber.org/zap"

	_ "github.com/ryuyb/fusion/docs"
)

type Server struct {
	app    *fiber.App
	config *config.Config
	logger *zap.Logger
}

func NewServer(
	app *fiber.App,
	cfg *config.Config,
	logger *zap.Logger,
) *Server {
	return &Server{
		app:    app,
		config: cfg,
		logger: logger,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	s.logger.Info("Starting server", zap.String("address", addr))
	return s.app.Listen(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server gracefully")
	return s.app.ShutdownWithContext(ctx)
}
