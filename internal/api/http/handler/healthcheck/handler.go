package healthcheck

import (
	"tms-core-service/internal/api/http/dto"
	"tms-core-service/internal/usecase/healthcheck"
	"tms-core-service/internal/util/httpresponse"

	"github.com/gofiber/fiber/v2"
)

// Handler handles health check requests
type Handler struct {
	useCase *healthcheck.HealthCheckUseCase
}

// NewHandler creates a new health check handler
func NewHandler(useCase *healthcheck.HealthCheckUseCase) *Handler {
	return &Handler{useCase: useCase}
}

// Check godoc
// @Summary Health check
// @Description Check if the service and database are healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} httpresponse.Response{data=dto.HealthResponse}
// @Failure 500 {object} httpresponse.Response
// @Router /health [get]
func (h *Handler) Check(c *fiber.Ctx) error {
	if err := h.useCase.Check(c.Context()); err != nil {
		return httpresponse.Error(c, err)
	}

	return httpresponse.Success(c, dto.HealthResponse{
		Status:   "healthy",
		Database: "connected",
		Version:  "1.0.1",
	}, "Service is healthy")
}
