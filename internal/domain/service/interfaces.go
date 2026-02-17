package service

import (
	"time"

	"github.com/google/uuid"
)

// HashService defines the interface for password hashing
type HashService interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, hash string) bool
}

// TokenClaims represents the claims in a JWT token
type TokenClaims struct {
	UserID uuid.UUID
	Email  string
}

// TokenService defines the interface for token operations
type TokenService interface {
	GenerateToken(userID uuid.UUID, email string, expiry time.Duration) (string, error)
	ValidateToken(tokenString string) (*TokenClaims, error)
}
