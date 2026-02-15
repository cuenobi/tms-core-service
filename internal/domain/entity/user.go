package entity

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system (Pure Domain Entity)
type User struct {
	ID           uuid.UUID
	Email        string
	PhoneNumber  string
	PasswordHash string
	FirstName    string
	LastName     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
