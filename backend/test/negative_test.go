// /backend/test/negative_test.go
package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/swapxs/LibMS/backend/handlers"
	"github.com/swapxs/LibMS/backend/middleware"
)

// TestUnauthorizedAccess_NoToken calls a protected endpoint without providing any token.
// Expected result: 401 Unauthorized with an error message.
func TestUnauthorizedAccess_NoToken(t *testing.T) {
	db := setupTestDB(t) // from your db_setup_test.go
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.GET("/books", handlers.GetBooks(db))
	}

	req, _ := http.NewRequest("GET", "/api/books", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestUnauthorizedAccess_InvalidToken calls a protected endpoint with a malformed token.
// Expected result: 401 Unauthorized with an error message.
func TestUnauthorizedAccess_InvalidToken(t *testing.T) {
	db := setupTestDB(t)
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.GET("/books", handlers.GetBooks(db))
	}

	req, _ := http.NewRequest("GET", "/api/books", nil)
	req.Header.Set("Authorization", "Bearer invalid_token_here") // Malformed token
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestInvalidJSON_Login sends invalid/malformed JSON to /auth/login.
// Expected result: 400 Bad Request due to JSON parsing error.
func TestInvalidJSON_Login(t *testing.T) {
	db := setupTestDB(t)
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.POST("/auth/login", handlers.Login(db))

	// Malformed JSON (missing a quote)
	body := []byte(`{"email": "someone@example.com", "password": "abc}`)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestInvalidJSON_Book tries to create a book with completely invalid JSON structure.
// Expected result: 400 Bad Request from ShouldBindJSON.
func TestInvalidJSON_Book(t *testing.T) {
	db := setupTestDB(t)
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	// We'll inject JWT claims for LibraryAdmin so that the route passes authentication.
	r.Use(func(c *gin.Context) {
		c.Set("user", map[string]interface{}{
			"library_id": float64(1),
			"role":       "LibraryAdmin",
		})
		c.Next()
	})
	r.POST("/books", handlers.AddOrIncrementBook(db))

	// Intentionally broken JSON
	body := []byte(`{ "isbn": "99999", "copies": 3, `) // missing closing brace etc.
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestProtectedRoute_NonJSONBody tries to hit a protected route with the wrong Content-Type.
func TestProtectedRoute_NonJSONBody(t *testing.T) {
	db := setupTestDB(t)
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	// Provide valid claims so we pass the JWT check
	r.Use(func(c *gin.Context) {
		c.Set("user", map[string]interface{}{
			"library_id": float64(1),
			"role":       "LibraryAdmin",
		})
		c.Next()
	})
	r.POST("/books", handlers.AddOrIncrementBook(db))

	// Content-Type is plain text, not JSON
	body := []byte(`isbn=12345&copies=3`)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "text/plain")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Because ShouldBindJSON expects JSON, it will yield 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

