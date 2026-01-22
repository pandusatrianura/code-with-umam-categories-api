package http

import (
	"fmt"
	"net/http"

	"github.com/pandusatrianura/code-with-umam-categories-api/constants"
	"github.com/pandusatrianura/code-with-umam-categories-api/internal/healthcheck/service"
	"github.com/pandusatrianura/code-with-umam-categories-api/pkg/json_wrapper"
)

// HealthCheckHandler is responsible for handling health check API requests by interacting with the health check service.
type HealthCheckHandler struct {
	service service.IHealthCheckService
}

// NewHealthCheckHandler initializes and returns a HealthCheckHandler with the provided IHealthCheckService implementation.
func NewHealthCheckHandler(service service.IHealthCheckService) (*HealthCheckHandler, error) {
	delegate := &HealthCheckHandler{
		service: service,
	}
	return delegate, nil
}

// API processes health check requests and returns a JSON response indicating the service's health status.
func (d *HealthCheckHandler) API(w http.ResponseWriter, r *http.Request) {
	var result json_wrapper.APIResponse
	svcHealthCheckResult := d.service.API()

	if svcHealthCheckResult.IsHealthy {
		result.Code = constants.SuccessCode
		result.Message = fmt.Sprintf("%s is healthy", svcHealthCheckResult.Name)
		json_wrapper.WriteJSONResponse(w, http.StatusOK, result)
	} else {
		result.Code = constants.ErrorCode
		result.Message = fmt.Sprintf("%s is not healthy", svcHealthCheckResult.Name)
		json_wrapper.WriteJSONResponse(w, http.StatusServiceUnavailable, result)
	}

	return
}
