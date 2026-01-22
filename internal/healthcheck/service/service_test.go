package service

import (
	"testing"

	"github.com/pandusatrianura/code-with-umam-categories-api/internal/healthcheck/entity"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckService_API(t *testing.T) {
	tests := []struct {
		name     string
		expected entity.HealthResponse
	}{
		{
			name: "returns healthy response",
			expected: entity.HealthResponse{
				Name:      "Categories API",
				IsHealthy: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &HealthCheckService{}
			result := service.API()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewHealthCheckService(t *testing.T) {
	tests := []struct {
		name      string
		expectErr bool
	}{
		{
			name:      "returns valid service instance without error",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewHealthCheckService()
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, service)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
			}
		})
	}
}
