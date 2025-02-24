package handlers_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/swapxs/LibMS/backend/src/handlers"
    "github.com/swapxs/LibMS/backend/src/models"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser_Success(t *testing.T) {
    db := setupTestDB(t)
    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.POST("/auth/register", handlers.RegisterUser(db))

    payload, _ := json.Marshal(map[string]any {
        "name":           "User 1",
        "email":          "user1@xenonstack.com",
        "password":       "123123123",
        "contact_number": "0912345678",
        "library_id":     1,
    })
    req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)

    var user models.User
    err := db.Where("email = ?", "user1@xenonstack.com").First(&user).Error
    assert.NoError(t, err)
    assert.Equal(t, "User 1", user.Name)
    assert.Equal(t, "Reader", user.Role)
}


func TestRegisterUser_UserExists(t *testing.T) {
    db := setupTestDB(t)
    // Seed an existing user
    existingUser := models.User{
        Name:          "Existing User",
        Email:         "existing.user@xenonstack.com",
        Password:      "passwd",
        ContactNumber: "1231231231",
        Role:          "Reader",
        LibraryID:     1,
    }
    db.Create(&existingUser)

    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.POST("/auth/register", handlers.RegisterUser(db))

    payload, _ := json.Marshal(map[string]any {
        "name":           "Jhon Doe",
        "email":          "existing.user@xenonstack.com",
        "password":       "pass123",
        "contact_number": "9999999999",
        "role":           "LibraryAdmin",
        "library_id":     1,
    })
    req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegisterUser_BadRequest(t *testing.T) {
    db := setupTestDB(t)

    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.POST("/auth/register", handlers.RegisterUser(db))

    // Missing required field 'email'
    payload, _ := json.Marshal(map[string]any {
        "name":           "User Without Email",
        "password":       "123123123",
        "contact_number": "01233443533",
    })
    req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegisterUser_InvalidEmailFormat(t *testing.T) {
    db := setupTestDB(t)

    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.POST("/auth/register", handlers.RegisterUser(db))

    // Provide an invalid email field
    payload, _ := json.Marshal(map[string]any {
        "name":           "XenonUser",
        "email":          "email",
        "password":       "password123",
        "contact_number": "1112223333",
    })
    req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code, "Invalid email should yield 400")
}

func TestRegisterUser_ShortPassword(t *testing.T) {
    db := setupTestDB(t)

    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.POST("/auth/register", handlers.RegisterUser(db))

    payload, _ := json.Marshal(map[string]any {
        "name":           "XenonUser",
        "email":          "shortpass@xenonstack.com",
        "password":       "abc",
        "contact_number": "9999999999",
    })
    req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code, "Short password should fail binding/validation")
}
