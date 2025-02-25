// /backend/test/books_test.go
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
	db.Create(&book)

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

	payload, _ := json.Marshal(map[string]any{
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
	err := db.First(&updatedBook, "isbn = ? AND library_id = ?", "12345", 1).Error
	assert.NoError(t, err)
	assert.Equal(t, 15, updatedBook.TotalCopies)
	assert.Equal(t, 15, updatedBook.AvailableCopies)
}

// Test retrieving a book by ISBN
func TestGetBookByISBN_Success(t *testing.T) {
    db := setupTestDB(t)

    // Seed a test book
    book := models.BookInventory{
        ISBN:            "get-isbn",
        LibraryID:       1,
        Title:           "Fetch Me",
        Author:          "Test Author",
        Publisher:       "Test Publisher",
        Language:        "English",
        Version:         "1st",
        TotalCopies:     5,
        AvailableCopies: 5,
    }
    db.Create(&book)

    gin.SetMode(gin.TestMode)
    r := gin.New()

    //  Inject user claims so the handler doesn't panic
    r.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "id":         float64(101),
            "email":      "user@example.com",
            "role":       "Reader",
            "library_id": float64(1),
        }
        c.Set("user", claims)
        c.Next()
    })

    r.GET("/books/:isbn", handlers.GetBooks(db))

    req, _ := http.NewRequest("GET", "/books/get-isbn", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    //  Print raw response body for debugging
    t.Logf("Response Body: %s", w.Body.String())

    //  Assert HTTP response
    assert.Equal(t, http.StatusOK, w.Code)

    //  Parse JSON response
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)

    //  Verify response structure
    books, exists := response["books"].([]interface{})
    assert.True(t, exists, "Expected 'books' array in response")
    assert.GreaterOrEqual(t, len(books), 1, "Expected at least one book in response")

    //  Extract the first book from the array
    bookData, ok := books[0].(map[string]interface{})
    assert.True(t, ok, "Expected book object inside 'books' array")

    //  Ensure 'title' field exists and is correct
    title, titleExists := bookData["Title"].(string)
    assert.True(t, titleExists, "Expected 'Title' field in book response")
    assert.Equal(t, "Fetch Me", title)
}

//  Test retrieving all books
func TestGetAllBooks_Success(t *testing.T) {
    db := setupTestDB(t)

    // Seed multiple books
    db.Create(&models.BookInventory{ISBN: "bk-1", LibraryID: 1, Title: "Book 1"})
    db.Create(&models.BookInventory{ISBN: "bk-2", LibraryID: 1, Title: "Book 2"})

    gin.SetMode(gin.TestMode)
    r := gin.New()

    //  Inject user claims so the handler doesn't panic
    r.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "id":         float64(101),
            "email":      "user@example.com",
            "role":       "Reader", // Adjust role as needed
            "library_id": float64(1),
        }
        c.Set("user", claims)
        c.Next()
    })

    r.GET("/books", handlers.GetBooks(db))

    req, _ := http.NewRequest("GET", "/books", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    // Assert HTTP response
    assert.Equal(t, http.StatusOK, w.Code)

    // Parse JSON response
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)

	// Check if books exist in response
    books, ok := response["books"].([]interface{})
    assert.True(t, ok, "Expected books array in response")
    assert.GreaterOrEqual(t, len(books), 2, "Expected at least 2 books")
}

// Test removing a book that has active issues (should fail)
func TestRemoveBook_WithActiveIssues_Fails(t *testing.T) {
	db := setupTestDB(t)

	book := models.BookInventory{
		ISBN:            "issued-isbn",
		LibraryID:       1,
		Title:           "Cannot Delete",
		TotalCopies:     2,
		AvailableCopies: 1, // 1 book is issued
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

	payload, _ := json.Marshal(map[string]any{
		"isbn":   "issued-isbn",
		"copies": 2,
	})
	req, _ := http.NewRequest("POST", "/books/remove", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Should fail since some books are issued
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test updating a book with an invalid ISBN (should fail)
func TestUpdateBook_InvalidISBN_Fails(t *testing.T) {
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
		"author": "Updated Author",
	})
	req, _ := http.NewRequest("PUT", "/books/invalid-isbn", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Test adding a book without required fields (should fail)
func TestAddBook_MissingFields_Fails(t *testing.T) {
	db := setupTestDB(t)
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

	r.POST("/books", handlers.AddOrIncrementBook(db))

	payload, _ := json.Marshal(map[string]any{
		"isbn":    "no-title-book",
		"copies":  3,
	})
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
