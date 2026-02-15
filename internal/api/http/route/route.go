package route

import (
	"tms-core-service/internal/api/http/handler/auth"
	"tms-core-service/internal/api/http/handler/healthcheck"
	"tms-core-service/internal/api/http/middleware"
	"tms-core-service/pkg/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// Dependencies holds all handler dependencies
type Dependencies struct {
	HealthCheckHandler *healthcheck.Handler
	AuthHandler        *auth.Handler
	JWTService         *jwt.JWTService
}

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, deps *Dependencies) {
	// Global Middleware
	app.Use(middleware.Trace())   // Generate trace ID first
	app.Use(middleware.Recover()) // Catch panics later

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Health check (no auth required)
	app.Get("/health", deps.HealthCheckHandler.Check)

	// API v1 routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Public routes (no auth required)
	authGroup := v1.Group("/auth")
	authGroup.Post("/register", deps.AuthHandler.Register)
	authGroup.Post("/login", deps.AuthHandler.Login)
	authGroup.Get("/google/login", deps.AuthHandler.GoogleLogin)
	authGroup.Get("/google/callback", deps.AuthHandler.GoogleCallback)
	authGroup.Get("/line/login", deps.AuthHandler.LineLogin)
	authGroup.Get("/line/callback", deps.AuthHandler.LineCallback)

	// Protected routes (JWT required)
	protected := v1.Group("", middleware.JWTAuth(deps.JWTService))
	protected.Get("/auth/me", deps.AuthHandler.GetProfile)
}
