package healthcheck

import (
	"context"

	"tms-core-service/internal/domain/repository"

	"gorm.io/gorm"
)

type healthCheckRepo struct {
	db *gorm.DB
}

// NewHealthCheckRepository creates a new health check repository
func NewHealthCheckRepository(db *gorm.DB) repository.HealthCheckRepository {
	return &healthCheckRepo{db: db}
}

// Ping checks if the database connection is alive
func (r *healthCheckRepo) Ping(ctx context.Context) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
