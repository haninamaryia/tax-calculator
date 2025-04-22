package main

import (
	"log"

	"github.com/haninamaryia/tax-calculator/internal/handler"
	"github.com/haninamaryia/tax-calculator/internal/logger"
	"github.com/haninamaryia/tax-calculator/internal/service"
	"github.com/haninamaryia/tax-calculator/internal/storage"
)

func main() {
	// Initialize logger before anything else
	logger.InitLogger()

	// Initialize the storage client that talks to the API
	storageClient := storage.NewTaxAPIClient("http://localhost:5001")

	// Initialize the tax service with the storage client
	taxService := service.NewTaxService(storageClient)

	// Initialize the HTTP handler with the tax service
	taxHandler := handler.NewServer(8080, taxService)

	// Start the HTTP server
	log.Println("Starting server on :8080")
	if err := taxHandler.ListenAndServe(); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
