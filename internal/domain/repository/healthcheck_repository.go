package repository

import "context"

// HealthCheckRepository defines the interface for health check operations
type HealthCheckRepository interface {
	// Ping checks if the database connection is alive
	Ping(ctx context.Context) error
}
