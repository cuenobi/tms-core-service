package middleware

import (
	"tms-core-service/internal/util/apierror"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Trace middleware adds a unique Trace ID to each request
func Trace() fiber.Handler {
	return func(c *fiber.Ctx) error {
		traceID := c.Get(apierror.HeaderXTraceID)
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Update response header
		c.Set(apierror.HeaderXTraceID, traceID)

		// Store in locals for later retrieval
		c.Locals(apierror.CtxTraceID, traceID)

		return c.Next()
	}
}
