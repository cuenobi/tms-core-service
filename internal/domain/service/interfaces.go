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

// TokenService defines the interface for token operations
type TokenService interface {
	GenerateToken(userID uuid.UUID, email string, expiry time.Duration) (string, error)
}
