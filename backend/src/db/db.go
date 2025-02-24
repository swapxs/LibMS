// /backend/db/db.go
package db

import (
	"log"
	"os"

	"github.com/swapxs/LibMS/backend/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB opens a connection to the Postgres database and auto-migrates all models.
func InitDB() *gorm.DB {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN is not set in the environment")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate all models.
	err = db.AutoMigrate(
		&models.Library{},
		&models.User{},
		&models.BookInventory{},
		&models.RequestEvent{},
		&models.IssueRegistry{},
	)
	if err != nil {
		log.Fatalf("Auto-migration failed: %v", err)
	}
	return db
}
