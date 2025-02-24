// /backend/main.go
package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/swapxs/LibMS/backend/src/db"
	"github.com/swapxs/LibMS/backend/src/routes"
)

func main() {
	// Load environment variables from .env file, if available.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using system environment variables")
	}

	// Initialize the database.
	database := db.InitDB()

	// Set up the router with all endpoints.
	router := routes.SetupRouter(database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
