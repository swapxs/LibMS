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

func setupBookTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.User{}, &models.Library{}, &models.BookInventory{})
	return db
}

func TestAddOrIncrementBook(t *testing.T) {
	db := setupBookTestDB()
	// Create a dummy library and admin user.
	library := models.Library{Name: "Test Library"}
	db.Create(&library)
	user := models.User{
		Name:      "Admin",
		Email:     "admin@example.com",
		Password:  "dummy",
		Role:      "LibraryAdmin",
		LibraryID: library.ID,
	}
	db.Create(&user)
	router := gin.Default()
	// Simulate JWT middleware by setting a dummy user context.
	router.POST("/books", func(c *gin.Context) {
		claims := map[string]interface{}{
			"id":         float64(user.ID),
			"email":      user.Email,
			"role":       user.Role,
			"library_id": float64(library.ID),
		}
		c.Set("user", claims)
		controllers.AddOrIncrementBook(db)(c)
	})
	input := controllers.AddBookInput{
		ISBN:    "1234567890",
		Title:   "Test Book",
		Authors: "Author One",
		Publisher: "Test Publisher",
		Version: "1st",
		Copies:  5,
	}
	jsonValue, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

