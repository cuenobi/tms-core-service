package repository

import (
	"context"

	"tms-core-service/internal/domain/entity"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// FindByID retrieves a user by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// FindByEmail retrieves a user by email
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// FindByPhoneNumber retrieves a user by phone number
	FindByPhoneNumber(ctx context.Context, phone string) (*entity.User, error)

	// FindByGoogleID retrieves a user by google ID
	FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error)

	// FindByLineID retrieves a user by LINE ID
	FindByLineID(ctx context.Context, lineID string) (*entity.User, error)

	// Create creates a new user
	Create(ctx context.Context, user *entity.User) error

	// Update updates an existing user
	Update(ctx context.Context, user *entity.User) error

	// Delete soft deletes a user
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves all users with pagination
	List(ctx context.Context, limit, offset int) ([]*entity.User, int64, error)
}
