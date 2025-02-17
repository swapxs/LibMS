package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/swapxs/LibMS/backend/db"
	"github.com/swapxs/LibMS/backend/routes"
)

func main() {
	// Load environment variables from .env (if present)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using system environment variables")
	}

	// Initialize the DB connection and auto-migrate models.
	database := db.InitDB()

	// Set up the router (with public and protected endpoints)
	router := routes.SetupRouter(database)

	// Determine port (default to 5000 if not set)
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
