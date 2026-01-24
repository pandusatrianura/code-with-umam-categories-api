package router

import (
	"fmt"
	"net/http"

	categoriesHandler "github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/delivery/http"
	"github.com/pandusatrianura/code-with-umam-categories-api/pkg/scalar"
)

// Router manages HTTP routing for various API endpoints.
// It integrates handlers for health checks and category-related operations.
type Router struct {
	categories *categoriesHandler.CategoriesHandler
}

// NewRouter initializes a new Router with the given health check and categories handlers.
func NewRouter(categoriesHandler *categoriesHandler.CategoriesHandler) *Router {
	return &Router{
		categories: categoriesHandler,
	}
}

// RegisterRoutes initializes and registers all HTTP routes for health checks and category operations.
func (h *Router) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("GET /categories/health", h.categories.API)
	r.HandleFunc("POST /categories", h.categories.InsertCategory)
	r.HandleFunc("GET /categories", h.categories.GetAllCategories)
	r.HandleFunc("GET /categories/{id}", h.categories.GetCategoryByID)
	r.HandleFunc("PUT /categories/{id}", h.categories.UpdateCategory)
	r.HandleFunc("DELETE /categories/{id}", h.categories.DeleteCategory)
	r.HandleFunc("GET /categories/docs", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Test Categories API",
			},
			DarkMode: true,
		})

		if err != nil {
			fmt.Printf("%v", err)
		}

		_, err = fmt.Fprintln(w, htmlContent)
		if err != nil {
			return
		}
	})
	return r
}
