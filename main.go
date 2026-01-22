package main

import (
	"log"

	"github.com/pandusatrianura/code-with-umam-categories-api/api"
)

// main starts the API server on the specified address and handles errors during its execution.
func main() {
	server := api.NewAPIServer(":8000")
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
