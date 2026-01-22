package router

import (
	"net/http"

	categoriesHandler "github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/delivery/http"
	healthCheckHandler "github.com/pandusatrianura/code-with-umam-categories-api/internal/healthcheck/delivery/http"
)

// Router manages HTTP routing for various API endpoints.
// It integrates handlers for health checks and category-related operations.
type Router struct {
	healthCheck *healthCheckHandler.HealthCheckHandler
	categories  *categoriesHandler.CategoriesHandler
}

// NewRouter initializes a new Router with the given health check and categories handlers.
func NewRouter(healthCheckHandler *healthCheckHandler.HealthCheckHandler, categoriesHandler *categoriesHandler.CategoriesHandler) *Router {
	return &Router{
		healthCheck: healthCheckHandler,
		categories:  categoriesHandler,
	}
}

// RegisterRoutes initializes and registers all HTTP routes for health checks and category operations.
func (h *Router) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("GET /categories/health", h.healthCheck.API)
	r.HandleFunc("POST /categories", h.categories.InsertCategory)
	r.HandleFunc("GET /categories", h.categories.GetAllCategories)
	r.HandleFunc("GET /categories/{id}", h.categories.GetCategoryByID)
	r.HandleFunc("PUT /categories/{id}", h.categories.UpdateCategory)
	r.HandleFunc("DELETE /categories/{id}", h.categories.DeleteCategory)
	return r
}
