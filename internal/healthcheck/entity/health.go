package entity

// HealthResponse represents the health status of a service or component with its name and health condition.
type HealthResponse struct {
	Name      string `json:"name"`
	IsHealthy bool   `json:"is_healthy"`
}
