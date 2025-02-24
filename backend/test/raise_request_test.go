// /backend/test/raise_request_test.go
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

func TestRaiseRequest(t *testing.T) {
	db := setupTestDB(t)

	// Pre-seed a book record (needed for checking availability).
	book := models.BookInventory{
		ISBN:            "12345",
		LibraryID:       1,
		Title:           "Golang Book",
		Author:          "John Doe",
		Publisher:       "Some Publisher",
		Language:        "English",
		Version:         "v1",
		TotalCopies:     5,
		AvailableCopies: 5,
	}
	err := db.Create(&book).Error
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Inject test JWT claims for a Reader.
	router.Use(func(c *gin.Context) {
		claims := jwt.MapClaims{
			"id":         1,
			"role":       "Reader",
			"library_id": 1,
		}
		c.Set("user", claims)
		c.Next()
	})
	router.POST("/requestEvents", handlers.RaiseRequest(db))

	requestBody, _ := json.Marshal(map[string]interface{}{
		"bookID": "12345",
	})
	req, _ := http.NewRequest("POST", "/requestEvents", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Handler should return HTTP 201 Created.
	assert.Equal(t, http.StatusCreated, w.Code)

	// Verify that a request event is created in the database.
	var reqEvent models.RequestEvent
	err = db.First(&reqEvent, "book_id = ? AND reader_id = ?", "12345", 1).Error
	assert.NoError(t, err)
	assert.Equal(t, "Issue", reqEvent.RequestType)
}

