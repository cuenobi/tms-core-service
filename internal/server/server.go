package server

import (
	"context"
	"fmt"
	"log"

	"tms-core-service/internal/config"

	"github.com/gofiber/fiber/v2"
)

// Server represents the HTTP server
type Server struct {
	config *config.ServerConfig
	app    *fiber.App
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.AppConfig) (*Server, error) {
	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "TMS Core Service",
		ReadTimeout:  cfg.Server.Timeout.Read,
		WriteTimeout: cfg.Server.Timeout.Write,
		IdleTimeout:  cfg.Server.Timeout.Idle,
	})

	// Apply global middleware
	ApplyMiddleware(app)

	// Wire dependencies and setup routes
	if err := WireDependencies(app, cfg); err != nil {
		return nil, fmt.Errorf("failed to wire dependencies: %w", err)
	}

	return &Server{
		config: &cfg.Server,
		app:    app,
	}, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	log.Printf("Server listening on %s\n", addr)
	if err := s.app.Listen(addr); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	return s.app.ShutdownWithContext(ctx)
}
