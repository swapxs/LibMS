// backend/test/jwt_test.go
package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
    "os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/swapxs/LibMS/backend/src/middleware"
)

var jwtSecret = []byte("testsecret")

func generateTestToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func setupRouterWithMiddleware() *gin.Engine {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", string(jwtSecret))
	r := gin.Default()
	r.Use(middleware.JWTAuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Authorized"})
	})
	return r
}

func TestJWTAuthMiddleware_ExpiredToken(t *testing.T) {
	router := setupRouterWithMiddleware()
	expiredClaims := jwt.MapClaims{
		"id": 1,
		"email": "test@example.com",
		"role": "Reader",
		"exp": time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
	}
	tokenString, _ := generateTestToken(expiredClaims)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuthMiddleware_MissingToken(t *testing.T) {
	router := setupRouterWithMiddleware()

	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuthMiddleware_InvalidSignature(t *testing.T) {
	router := setupRouterWithMiddleware()

	claims := jwt.MapClaims{
		"id": 1,
		"email": "test@example.com",
		"role": "Reader",
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fakeTokenString, _ := token.SignedString([]byte("wrongsecret"))

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+fakeTokenString)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
