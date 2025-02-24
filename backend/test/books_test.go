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
	"github.com/swapxs/LibMS/backend/src/handlers"
	"github.com/swapxs/LibMS/backend/src/models"
	"github.com/stretchr/testify/assert"
)

func TestAddOrIncrementBook(t *testing.T) {
	db := setupTestDB(t)
	
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
	r := gin.New()

	r.Use(func(c *gin.Context) {
		claims := jwt.MapClaims{
			"id":         1,
			"email":      "john@xenonstack.com",
			"role":       "LibraryAdmin",
			"library_id": 1,
		}
		c.Set("user", claims)
		c.Next()
	})
	r.POST("/books", handlers.AddOrIncrementBook(db))

	payload, _ := json.Marshal(map[string]any {
		"isbn":     "12345",
		"title":    "Golang Book",
		"author":   "John Doe",
		"copies":   5,
		"language": "English",
	})
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var updatedBook models.BookInventory
	err = db.First(&updatedBook, "isbn = ? AND library_id = ?", "12345", 1).Error
	assert.NoError(t, err)
	assert.Equal(t, 15, updatedBook.TotalCopies)
	assert.Equal(t, 15, updatedBook.AvailableCopies)
}

func TestRemoveBook_ResultsInDeletion(t *testing.T) {
    db := setupTestDB(t)

    book := models.BookInventory{
        ISBN:            "rm-isbn",
        LibraryID:       1,
        Title:           "Delete Me",
        Author:          "Author",
        Publisher:       "Pub",
        Language:        "English",
        Version:         "1st",
        TotalCopies:     2,
        AvailableCopies: 2,
    }
    db.Create(&book)

    gin.SetMode(gin.TestMode)
    r := gin.New()

    r.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "id":         float64(99),
            "role":       "LibraryAdmin",
            "library_id": float64(1),
        }
        c.Set("user", claims)
        c.Next()
    })

    r.POST("/books/remove", handlers.RemoveBook(db))

    payload, _ := json.Marshal(map[string]any {
        "isbn":   "rm-isbn",
        "copies": 2,
    })
    req, _ := http.NewRequest("POST", "/books/remove", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var check models.BookInventory
    err := db.First(&check, "isbn = ?", "rm-isbn").Error
    assert.Error(t, err)
}

func TestUpdateBook_NotFound(t *testing.T) {
    db := setupTestDB(t)
    gin.SetMode(gin.TestMode)
    r := gin.New()

    r.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "id":         float64(101),
            "role":       "LibraryAdmin",
            "library_id": float64(1),
        }
        c.Set("user", claims)
        c.Next()
    })

    r.PUT("/books/:isbn", handlers.UpdateBook(db))

    payload, _ := json.Marshal(map[string]string{
        "author": "New Author",
    })
    req, _ := http.NewRequest("PUT", "/books/nonexistent-isbn", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateBook_PartialFields(t *testing.T) {
    db := setupTestDB(t)

    original := models.BookInventory{
        ISBN:            "upd-isbn",
        LibraryID:       1,
        Title:           "Original Title",
        Author:          "Original Author",
        Publisher:       "Orig Publisher",
        Language:        "English",
        Version:         "1st",
        TotalCopies:     5,
        AvailableCopies: 5,
    }
    db.Create(&original)

    gin.SetMode(gin.TestMode)
    r := gin.New()

    r.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "id":         float64(102),
            "role":       "LibraryAdmin",
            "library_id": float64(1),
        }
        c.Set("user", claims)
        c.Next()
    })

    r.PUT("/books/:isbn", handlers.UpdateBook(db))

    payload, _ := json.Marshal(map[string]string{
        "author": "New Author Only",
    })

    req, _ := http.NewRequest("PUT", "/books/upd-isbn", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    // Verify that only the Author changed
    var updated models.BookInventory
    err := db.First(&updated, "isbn = ?", "upd-isbn").Error
    assert.NoError(t, err)
    assert.Equal(t, "Original Title", updated.Title)
    assert.Equal(t, "New Author Only", updated.Author)
    assert.Equal(t, "Orig Publisher", updated.Publisher)
}

func TestAddOrIncrementBook_IncrementOnly_BookNotFound(t *testing.T) {
    db := setupTestDB(t)
    gin.SetMode(gin.TestMode)
    r := gin.New()

    r.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "id":         float64(103),
            "role":       "LibraryAdmin",
            "library_id": float64(1),
        }
        c.Set("user", claims)
        c.Next()
    })

    r.POST("/books", handlers.AddOrIncrementBook(db))

    payload, _ := json.Marshal(map[string]any {
        "isbn":           "no-such-book",
        "copies":         3,
        "increment_only": true,
    })

    req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}
