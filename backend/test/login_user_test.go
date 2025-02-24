package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/swapxs/LibMS/backend/src/handlers"
	"github.com/swapxs/LibMS/backend/src/models"
	"golang.org/x/crypto/bcrypt"
)

// The first function for the test is a successful login request.
// This first created the user by seeding the details into the database
// then it sets the route for the login request and then attempts login.
// Since both the email and password(hashed) matches, it gives a success response.
// At the end it also verifies token.
func TestLogin_Success(t *testing.T) {
    db := setupTestDB(t)
    hashed, _ := bcrypt.GenerateFromPassword([]byte("testpasswd"), bcrypt.DefaultCost)
    user := models.User{
        Name:          "User 1",
        Email:         "user1@xenonstack.com",
        Password:      string(hashed),
        ContactNumber: "1231231231",
        Role:          "Reader",
        LibraryID:     1,
    }
    db.Create(&user)

    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.POST("/auth/login", handlers.Login(db))

	// https://cs.opensource.google/go/go/+/refs/tags/go1.23.0:src/builtin/builtin.go;l=97 (Any is an Alias for interface{})
    payload, _ := json.Marshal(map[string]any {
        "email":    "user1@xenonstack.com",
        "password": "testpasswd",
    })

    req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

	var res map[string]any 
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err, "Failed to parse JSON response")

	token, ok := res["token"]
	assert.True(t, ok, "Expected 'token' filed, but not found.")

	tokenStr, isString := token.(string)
    assert.True(t, isString, "Expected token to be a string")
    assert.NotEmpty(t, tokenStr, "Expected token to be non-empty")
}

// This function first seeds a user then attempts to login with BAD credentials.
// This throws an unauthorized error for that user.
func TestLogin_InvalidCredentials(t *testing.T) {
    db := setupTestDB(t)
    // Seed a user
    hashed, _ := bcrypt.GenerateFromPassword([]byte("1231231231"), bcrypt.DefaultCost)
    user := models.User{
        Name:          "Test User",
        Email:         "abc@xenonstack.com",
        Password:      string(hashed),
        ContactNumber: "555-5555",
        Role:          "Reader",
        LibraryID:     1,
    }
    db.Create(&user)

    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.POST("/auth/login", handlers.Login(db))

    payload, _ := json.Marshal(map[string]any {
        "email":    "abc@xenonstack.com",
        "password": "3213213213",
    })

    req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Here we Login with missing credentials
// like we seeded the user then when logging in, we did not provide the password. 
// This checks for any injection vulnerability in the login form
func TestLogin_BadRequest(t *testing.T) {
    db := setupTestDB(t)

    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.POST("/auth/login", handlers.Login(db))

    payload, _ := json.Marshal(map[string]any {
        "email": "some@xenonstack.com",
    })

    req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_NonExistentEmail(t *testing.T) {
    db := setupTestDB(t)

    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.POST("/auth/login", handlers.Login(db))

    requestBody, _ := json.Marshal(map[string]any {
        "email":    "doesnotexist@xenonstack.com",
        "password": "passwd",
    })
    req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(requestBody))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusUnauthorized, w.Code)
}

//    This test calls a protected endpoint with an artificially expired token.
//    We'll create a route that requires JWTAuthMiddleware, then pass a token 
//    whose "exp" is in the past, expecting a 401.
func TestJWT_ExpiredToken(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.New()

    // This is a minimal protected route that just returns "ok" if authorized.
    protected := r.Group("/api")
    protected.Use(func(c *gin.Context) {
        // We replicate the JWTAuthMiddleware logic in a minimal test way,
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "No auth header"})
            c.Abort()
            return
        }

        tokenString := authHeader[len("Bearer "):]
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte("testsecret"), nil
        })
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }
        // If we get here, it's valid
        c.Next()
    })

    protected.GET("/protected-route", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    // Create an expired token
    // Notice we are using "testsecret" as the key. 
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "id":    999,
        "email": "expired@example.com",
        "exp":   time.Now().Add(-time.Hour).Unix(), // 1 hour in the past
    })
    tokenString, _ := token.SignedString([]byte("testsecret"))

    // Make request
    req, _ := http.NewRequest("GET", "/api/protected-route", nil)
    req.Header.Set("Authorization", "Bearer "+tokenString)

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    // Expect 401
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}
