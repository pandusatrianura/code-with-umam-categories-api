package entity

// Category represents a grouping of certain entities, containing an ID, name, and a description.
type Category struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// HealthResponse represents the health status of a service or component with its name and health condition.
type HealthResponse struct {
	Name      string `json:"name"`
	IsHealthy bool   `json:"is_healthy"`
}
