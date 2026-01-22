package api

import (
	"log"
	"net/http"

	route "github.com/pandusatrianura/code-with-umam-categories-api/api/router"

	HealthCheckHandler "github.com/pandusatrianura/code-with-umam-categories-api/internal/healthcheck/delivery/http"
	HealthCheckService "github.com/pandusatrianura/code-with-umam-categories-api/internal/healthcheck/service"

	CategoriesHandler "github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/delivery/http"
	CategoriesRepository "github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/repository"
	CategoriesService "github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/service"
)

// Server represents an HTTP server with an address for listening to incoming requests.
type Server struct {
	addr string
}

// NewAPIServer initializes and returns a new Server instance configured to listen on the specified address.
func NewAPIServer(addr string) *Server {
	return &Server{addr: addr}
}

// Run starts the server, initializes dependencies, registers routes, and listens for incoming HTTP requests.
func (s *Server) Run() error {

	healthCheckService, err := HealthCheckService.NewHealthCheckService()
	if err != nil {
		panic(err)
	}

	healthCheckHandler, err := HealthCheckHandler.NewHealthCheckHandler(healthCheckService)
	if err != nil {
		panic(err)
	}

	categoriesRepo, err := CategoriesRepository.NewCategoriesRepository()
	if err != nil {
		panic(err)
	}

	categoriesService, err := CategoriesService.NewCategoriesService(categoriesRepo)
	if err != nil {
		panic(err)
	}

	categoriesHandler, err := CategoriesHandler.NewCategoriesHandler(categoriesService)
	if err != nil {
		panic(err)
	}

	r := route.NewRouter(healthCheckHandler, categoriesHandler)
	routes := r.RegisterRoutes()
	router := http.NewServeMux()
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", routes))
	log.Println("Starting server on port", s.addr)
	return http.ListenAndServe(s.addr, router)
}
