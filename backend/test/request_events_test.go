package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/swapxs/LibMS/backend/controllers"
	"github.com/swapxs/LibMS/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRequestTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Library{}, &models.User{}, &models.BookInventory{}, &models.RequestEvent{})
	return db
}

func TestRaiseRequest(t *testing.T) {
	db := setupRequestTestDB()
	// Create dummy library, reader, and book.
	library := models.Library{Name: "Test Library"}
	db.Create(&library)
	reader := models.User{
		Name:      "Reader",
		Email:     "reader@example.com",
		Password:  "dummy",
		Role:      "Reader",
		LibraryID: library.ID,
	}
	db.Create(&reader)
	book := models.BookInventory{
		ISBN:            "9876543210",
		LibraryID:       library.ID,
		Title:           "Test Book",
		Authors:         "Author",
		TotalCopies:     5,
		AvailableCopies: 5,
	}
	db.Create(&book)
	router := gin.Default()
	// Simulate JWT middleware.
	router.POST("/request", func(c *gin.Context) {
		claims := map[string]interface{}{
			"id":         float64(reader.ID),
			"email":      reader.Email,
			"role":       reader.Role,
			"library_id": float64(library.ID),
		}
		c.Set("user", claims)
		controllers.RaiseRequest(db)(c)
	})
	input := controllers.RaiseRequestInput{
		ISBN: "9876543210",
	}
	jsonValue, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/request", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

