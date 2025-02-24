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
	r := gin.New()

	// Inject test JWT claims for a Reader.
	r.Use(func(c *gin.Context) {
		claims := jwt.MapClaims{
			"id":         1,
			"role":       "Reader",
			"library_id": 1,
		}
		c.Set("user", claims)
		c.Next()
	})
	r.POST("/requestEvents", handlers.RaiseRequest(db))

	payload, _ := json.Marshal(map[string]any {
		"bookID": "12345",
	})
	req, _ := http.NewRequest("POST", "/requestEvents", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Handler should return HTTP 201 Created.
	assert.Equal(t, http.StatusCreated, w.Code)

	// Verify that a request event is created in the database.
	var reqEvent models.RequestEvent
	err = db.First(&reqEvent, "book_id = ? AND reader_id = ?", "12345", 1).Error
	assert.NoError(t, err)
	assert.Equal(t, "Issue", reqEvent.RequestType)
}

// TestRaiseRequest_TooManyExistingRequests ensures a user cannot exceed 4 active requests.
func TestRaiseRequest_TooManyExistingRequests(t *testing.T) {
    db := setupTestDB(t)

    // Seed a user
    userID := uint(1)
    user := models.User{
        Name:          "RequestUser",
        Email:         "requser@xenonstack.com",
        Password:      "hashedpw",
        ContactNumber: "12345",
        Role:          "Reader",
        LibraryID:     1,
    }
    db.Create(&user)

    // Seed a book
    book := models.BookInventory{
        ISBN:            "testisbn",
        LibraryID:       1,
        Title:           "Test Book",
        Author:          "Test Author",
        Publisher:       "Test Publisher",
        Language:        "English",
        Version:         "v1",
        TotalCopies:     5,
        AvailableCopies: 5,
    }
    db.Create(&book)

    // Create 4 active requests for that user
    // request_type = "Issue" or "Approve" indicates active
    for i := 0; i < 4; i++ {
        re := models.RequestEvent{
            BookID:      "testisbn",
            ReaderID:    userID,
            RequestDate: book.Model.CreatedAt, // arbitrary date
            RequestType: "Issue",              // active
        }
        db.Create(&re)
    }

    gin.SetMode(gin.TestMode)
    r := gin.New()

    // Inject Reader claims
    r.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "id":         float64(userID),
            "role":       "Reader",
            "library_id": float64(1),
        }
        c.Set("user", claims)
        c.Next()
    })

    r.POST("/requestEvents", handlers.RaiseRequest(db))

    // Attempt to raise a new request
    payload, _ := json.Marshal(map[string]string{
        "bookID": "testisbn",
    })
    req, _ := http.NewRequest("POST", "/requestEvents", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    // Because activeRequests >= 4, we should see a 403 Forbidden
    assert.Equal(t, http.StatusForbidden, w.Code)
}

// TestRaiseRequest_BookNotAvailable sets AvailableCopies=0 and tries to raise a new request.
func TestRaiseRequest_BookNotAvailable(t *testing.T) {
    db := setupTestDB(t)

    // Seed a user
    user := models.User{
        Name:          "NotAvailableUser",
        Email:         "notavail@xenonstack.com",
        Password:      "dummy",
        ContactNumber: "111",
        Role:          "Reader",
        LibraryID:     1,
    }
    db.Create(&user)

    // Seed a book with 0 available copies
    book := models.BookInventory{
        ISBN:            "noavail",
        LibraryID:       1,
        Title:           "Zero Copies",
        Author:          "Test Author",
        Publisher:       "Test Pub",
        Language:        "English",
        Version:         "v1",
        TotalCopies:     5,
        AvailableCopies: 0, // not available
    }
    db.Create(&book)

    gin.SetMode(gin.TestMode)
    r := gin.New()
    // JWT claims
    r.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "id":         float64(user.ID),
            "role":       user.Role,
            "library_id": float64(user.LibraryID),
        }
        c.Set("user", claims)
        c.Next()
    })
    r.POST("/requestEvents", handlers.RaiseRequest(db))

    payload, _ := json.Marshal(map[string]string{
        "bookID": "noavail",
    })
    req, _ := http.NewRequest("POST", "/requestEvents", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    // Should 400 "Book not available for issue"
    assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestRaiseRequest_MissingBookID tries to raise a request with no bookID in the JSON.
func TestRaiseRequest_MissingBookID(t *testing.T) {
    db := setupTestDB(t)

    // Seed a user
    user := models.User{
        Name:          "Missing BookID",
        Email:         "missingbk@xenonstack.com",
        Password:      "pw",
        ContactNumber: "333",
        Role:          "Reader",
        LibraryID:     1,
    }
    db.Create(&user)

    gin.SetMode(gin.TestMode)
    r := gin.New()
    // JWT claims
    r.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "id":         float64(user.ID),
            "role":       user.Role,
            "library_id": float64(user.LibraryID),
        }
        c.Set("user", claims)
        c.Next()
    })
    r.POST("/requestEvents", handlers.RaiseRequest(db))

    // No bookID in the payload
    payload, _ := json.Marshal(map[string]string{})
    req, _ := http.NewRequest("POST", "/requestEvents", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    // Should 400 because bookID is required
    assert.Equal(t, http.StatusBadRequest, w.Code)
}

