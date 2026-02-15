package middleware

import (
	"fmt"
	"log"
	"runtime/debug"

	"tms-core-service/internal/util/apierror"
	"tms-core-service/internal/util/httpresponse"

	"github.com/gofiber/fiber/v2"
)

// Recover middleware recovers from panics and returns a standardized error
func Recover() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}

				// Log the panic with stack trace
				log.Printf("[PANIC RECOVER] %v\n%s", err, debug.Stack())

				// Return standardized internal error
				apiErr := apierror.NewInternalError("An unexpected server error occurred")
				_ = httpresponse.Error(c, apiErr)
			}
		}()

		return c.Next()
	}
}
