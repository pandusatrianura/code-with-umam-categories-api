package service

import (
	"github.com/pandusatrianura/code-with-umam-categories-api/internal/healthcheck/entity"
)

// IHealthCheckService defines the contract for performing health checks and retrieving service health information.
type IHealthCheckService interface {
	API() entity.HealthResponse
}

// HealthCheckService provides functionality to perform health checks for system components or APIs.
type HealthCheckService struct{}

// NewHealthCheckService creates and returns a new instance of HealthCheckService or an error if the initialization fails.
func NewHealthCheckService() (*HealthCheckService, error) {
	return &HealthCheckService{}, nil
}

// API returns the health status of the Categories API as an entity.HealthResponse.
func (s *HealthCheckService) API() entity.HealthResponse {
	return entity.HealthResponse{
		Name:      "Categories API",
		IsHealthy: true,
	}
}
