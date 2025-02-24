// /backend/test/db_setup_test.go
package handlers_test

import (
	"testing"

	"github.com/swapxs/LibMS/backend/src/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite DB: %v", err)
	}
	err = db.AutoMigrate(
		&models.Library{},
		&models.User{},
		&models.BookInventory{},
		&models.RequestEvent{},
		&models.IssueRegistry{},
	)
	if err != nil {
		t.Fatalf("failed to auto-migrate models: %v", err)
	}
	return db
}
