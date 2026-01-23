package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pandusatrianura/code-with-umam-categories-api/constants"
	"github.com/pandusatrianura/code-with-umam-categories-api/internal/healthcheck/entity"
	"github.com/pandusatrianura/code-with-umam-categories-api/pkg/json_wrapper"
	"github.com/stretchr/testify/assert"
)

// MockHealthCheckService is a mock implementation of IHealthCheckService for testing.
type MockHealthCheckService struct {
	mockResponse entity.HealthResponse
}

func (m *MockHealthCheckService) API() entity.HealthResponse {
	return m.mockResponse
}

func TestHealthCheckHandler_API(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   entity.HealthResponse
		expectedCode   int
		expectedResult json_wrapper.APIResponse
	}{
		{
			name: "service is healthy",
			mockResponse: entity.HealthResponse{
				Name:      "Test Service",
				IsHealthy: true,
			},
			expectedCode: http.StatusOK,
			expectedResult: json_wrapper.APIResponse{
				Code:    constants.SuccessCode,
				Message: "Test Service is healthy",
			},
		},
		{
			name: "service is not healthy",
			mockResponse: entity.HealthResponse{
				Name:      "Test Service",
				IsHealthy: false,
			},
			expectedCode: http.StatusServiceUnavailable,
			expectedResult: json_wrapper.APIResponse{
				Code:    constants.ErrorCode,
				Message: "Test Service is not healthy",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockHealthCheckService{
				mockResponse: tt.mockResponse,
			}

			handler := &HealthCheckHandler{
				service: mockService,
			}

			req := httptest.NewRequest("GET", "/health", nil)
			rec := httptest.NewRecorder()

			handler.API(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			var actualResult json_wrapper.APIResponse
			err := json.NewDecoder(rec.Body).Decode(&actualResult)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, actualResult)
		})
	}
}

func TestNewHealthCheckHandler(t *testing.T) {
	mockService := &MockHealthCheckService{}
	handler, err := NewHealthCheckHandler(mockService)

	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
}
