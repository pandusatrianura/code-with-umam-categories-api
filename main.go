package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pandusatrianura/code-with-umam-categories-api/api"
)

// @title Categories API
// @version 1.0
// @host pandusatrianura-categories-api-production.up.railway.app/
// @BasePath /

// main starts the API server on the specified address and handles errors during its execution.
func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	port := os.Getenv("PORT")

	fmt.Println("Server started successfully on PORT ", port)
	server := api.NewAPIServer(fmt.Sprintf(":%s", port))
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
