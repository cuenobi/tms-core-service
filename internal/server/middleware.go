package server

import (
	"tms-core-service/internal/api/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// ApplyMiddleware applies global middleware to the router
func ApplyMiddleware(app *fiber.App) {
	// Recovery middleware (recover from panics)
	app.Use(recover.New())

	// Logger middleware
	app.Use(logger.New())

	// CORS middleware
	app.Use(middleware.CORS())
}
