package middleware

import (
	"strings"

	"tms-core-service/internal/domain/errs"
	"tms-core-service/pkg/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
	userIDKey           = "user_id"
	userEmailKey        = "user_email"
	userStatusKey       = "user_status"
)

// JWTAuth creates a JWT authentication middleware
func JWTAuth(jwtService *jwt.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get(authorizationHeader)
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Authorization header is required",
			})
		}

		// Check Bearer prefix
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid authorization header format",
			})
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Token is required",
			})
		}

		// Validate token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   errs.ErrTokenInvalid.Error(),
			})
		}

		// Set user information in context
		c.Locals(userIDKey, claims.UserID)
		c.Locals(userStatusKey, claims.Status)

		return c.Next()
	}
}

// GetUserID gets user ID from context
func GetUserID(c *fiber.Ctx) (uuid.UUID, bool) {
	val := c.Locals(userIDKey)
	if val == nil {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}

// GetUserEmail gets user email from context
func GetUserEmail(c *fiber.Ctx) (string, bool) {
	email := c.Locals(userEmailKey)
	if email == nil {
		return "", false
	}
	emailStr, ok := email.(string)
	return emailStr, ok
}

// GetUserStatus gets user status from context
func GetUserStatus(c *fiber.Ctx) (string, bool) {
	status := c.Locals(userStatusKey)
	if status == nil {
		return "", false
	}
	statusStr, ok := status.(string)
	return statusStr, ok
}
