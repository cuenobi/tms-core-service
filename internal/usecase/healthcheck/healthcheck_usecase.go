package healthcheck

import (
	"context"
	"fmt"

	"tms-core-service/internal/domain/repository"
)

// HealthCheckUseCase handles health check operations
type HealthCheckUseCase struct {
	repo repository.HealthCheckRepository
}

// NewHealthCheckUseCase creates a new health check use case
func NewHealthCheckUseCase(repo repository.HealthCheckRepository) *HealthCheckUseCase {
	return &HealthCheckUseCase{repo: repo}
}

// Check performs health check
func (uc *HealthCheckUseCase) Check(ctx context.Context) error {
	if err := uc.repo.Ping(ctx); err != nil {
		return fmt.Errorf("healthcheck repository: ping: %w", err)
	}
	return nil
}
