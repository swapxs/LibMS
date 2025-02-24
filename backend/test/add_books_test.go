// /backend/test/add_books_test.go
package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/swapxs/LibMS/backend/handlers"
	"github.com/swapxs/LibMS/backend/models"
	"github.com/stretchr/testify/assert"
)

func TestAddOrIncrementBook(t *testing.T) {
	db := setupTestDB(t)
	
	// Pre-seed a book record to simulate an existing book.
	book := models.BookInventory{
		ISBN:            "12345",
		LibraryID:       1,
		Title:           "Golang Book",
		Author:          "John Doe",
		Publisher:       "Some Publisher",
		Language:        "English",
		Version:         "v1",
		TotalCopies:     10,
		AvailableCopies: 10,
	}
	err := db.Create(&book).Error
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Inject test JWT claims.
	router.Use(func(c *gin.Context) {
		claims := jwt.MapClaims{
			"id":         1,
			"email":      "john@example.com",
			"role":       "LibraryAdmin",
			"library_id": 1,
		}
		c.Set("user", claims)
		c.Next()
	})
	router.POST("/books", handlers.AddOrIncrementBook(db))

	requestBody, _ := json.Marshal(map[string]interface{}{
		"isbn":     "12345",
		"title":    "Golang Book",
		"author":   "John Doe",
		"copies":   5,
		"language": "English",
	})
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Handler should return HTTP 201 Created.
	assert.Equal(t, http.StatusCreated, w.Code)

	// Verify that the book copies have been incremented.
	var updatedBook models.BookInventory
	err = db.First(&updatedBook, "isbn = ? AND library_id = ?", "12345", 1).Error
	assert.NoError(t, err)
	assert.Equal(t, 15, updatedBook.TotalCopies)
	assert.Equal(t, 15, updatedBook.AvailableCopies)
}

