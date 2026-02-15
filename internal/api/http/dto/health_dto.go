package dto

// HealthResponse represents a health check response
type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}
