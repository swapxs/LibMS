// /backend/test/user_handler_test.go
package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/swapxs/LibMS/backend/src/handlers"
	"github.com/swapxs/LibMS/backend/src/models"
	"github.com/golang-jwt/jwt/v4"
)

func TestGetUsers_Success(t *testing.T) {
	db := setupTestDB(t)

	// Seed users in the test database
	db.Create(&models.User{Name: "User1", Email: "user1@example.com", Role: "Reader", LibraryID: 1})
	db.Create(&models.User{Name: "User2", Email: "user2@example.com", Role: "LibraryAdmin", LibraryID: 1})

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Middleware to mock JWT authentication
	router.Use(func(c *gin.Context) {
		claims := jwt.MapClaims{
			"id":         1,
			"email":      "admin@example.com",
			"role":       "LibraryAdmin",
			"library_id": 1,
		}
		c.Set("user", claims)
		c.Next()
	})
	router.GET("/users", handlers.GetUsers(db))

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "user1@example.com")
	assert.Contains(t, w.Body.String(), "user2@example.com")
}

// TestGetUsers_Unauthorized ensures users cannot be fetched if no library ID is present.
func TestGetUsers_Unauthorized(t *testing.T) {
	db := setupTestDB(t)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Middleware without setting library_id in JWT claims
	router.Use(func(c *gin.Context) {
		claims := jwt.MapClaims{
			"id":    1,
			"email": "admin@example.com",
			"role":  "LibraryAdmin",
		}
		c.Set("user", claims)
		c.Next()
	})
	router.GET("/users", handlers.GetUsers(db))

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Expect Internal Server Error due to missing library_id
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "claim library_id not found")
}

// TestGetUserIssueInfo_Success ensures issue records are returned for a user.
func TestGetUserIssueInfo_Success(t *testing.T) {
	db := setupTestDB(t)

	// Seed user and issue registry records
	db.Create(&models.User{Name: "Reader One", Email: "reader@example.com", Role: "Reader", LibraryID: 1})
	db.Create(&models.IssueRegistry{ISBN: "12345", ReaderID: 1, IssueStatus: "Issued", LibraryID: 1})

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// FIX: Ensure `id` is stored as `float64` in JWT claims
	router.Use(func(c *gin.Context) {
		claims := jwt.MapClaims{
			"id":         float64(1),
			"email":      "reader@example.com",
			"role":       "Reader",
			"library_id": float64(1),
		}
		c.Set("user", claims)
		c.Next()
	})
	router.GET("/user/issues", handlers.GetUserIssueInfo(db))

	req, _ := http.NewRequest("GET", "/user/issues", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "12345") // ISBN should be in the response
}

// TestGetUserIssueInfo_NoIssues ensures an empty response when no issues exist.
func TestGetUserIssueInfo_NoIssues(t *testing.T) {
	db := setupTestDB(t)

	// Seed user without any issues
	db.Create(&models.User{Name: "Reader One", Email: "reader@example.com", Role: "Reader", LibraryID: 1})

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Middleware to mock JWT authentication
	router.Use(func(c *gin.Context) {
		claims := jwt.MapClaims{
			"id":         1,
			"email":      "reader@example.com",
			"role":       "Reader",
			"library_id": 1,
		}
		c.Set("user", claims)
		c.Next()
	})
	router.GET("/user/issues", handlers.GetUserIssueInfo(db))

	req, _ := http.NewRequest("GET", "/user/issues", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "[]") // Expect empty array
}

// TestGetUserIssueInfo_InvalidJWT ensures an error when JWT claims are incorrect.
func TestGetUserIssueInfo_InvalidJWT(t *testing.T) {
	db := setupTestDB(t)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Middleware with invalid JWT claim (missing "id")
	router.Use(func(c *gin.Context) {
		claims := jwt.MapClaims{
			"email":      "reader@example.com",
			"role":       "Reader",
			"library_id": float64(1),
		}
		c.Set("user", claims)
		c.Next()
	})
	router.GET("/user/issues", handlers.GetUserIssueInfo(db))

	req, _ := http.NewRequest("GET", "/user/issues", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ✅ Update the test expectation
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "claim id not found") // ✅ Expect this instead of panic
}
