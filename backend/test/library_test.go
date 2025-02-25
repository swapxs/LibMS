// /backend/test/create_library_test.go
package handlers_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "github.com/stretchr/testify/assert"
    "github.com/swapxs/LibMS/backend/src/handlers"
    "github.com/swapxs/LibMS/backend/src/models"
)

// TestCreateLibrary_Owner tests that a user with role = "Owner" can create a library.
func TestCreateLibrary_Owner(t *testing.T) {
    db := setupTestDB(t)
    gin.SetMode(gin.TestMode)
    router := gin.New()

    router.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "role": "Owner",
            // You can add more claims if needed, e.g. "id", "library_id", etc.
        }
        c.Set("user", claims)
        c.Next()
    })

    // This route is now protected such that only an Owner can create
    router.POST("/library", handlers.CreateLibrary(db))

    payload, _ := json.Marshal(map[string]string{
        "name": "New Library",
    })
    req, _ := http.NewRequest("POST", "/library", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Handler should return HTTP 201 Created.
    assert.Equal(t, http.StatusCreated, w.Code)

    // Verify that the library is created in the database.
    var lib models.Library
    err := db.First(&lib, "name = ?", "New Library").Error
    assert.NoError(t, err)
    assert.Equal(t, "New Library", lib.Name)
}

// TestCreateLibrary_NonOwner tests that a user without the Owner role
// cannot create a library and should receive 401 Unauthorized.
func TestCreateLibrary_NonOwner(t *testing.T) {
    db := setupTestDB(t)
    gin.SetMode(gin.TestMode)
    router := gin.New()

    router.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "role": "Reader",
        }
        c.Set("user", claims)
        c.Next()
    })

    router.POST("/library", handlers.CreateLibrary(db))

    payload, _ := json.Marshal(map[string]string{
        "name": "Should Fail",
    })
    req, _ := http.NewRequest("POST", "/library", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Expect 401 Unauthorized
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestCreateLibrary_AlreadyExists verifies behavior when a duplicate library is attempted.
func TestCreateLibrary_AlreadyExists(t *testing.T) {
    db := setupTestDB(t)
    // Seed an existing library
    existingLibrary := models.Library{Name: "Duplicate Library"}
    db.Create(&existingLibrary)

    gin.SetMode(gin.TestMode)
    router := gin.New()

    router.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{"role": "Owner"}
        c.Set("user", claims)
        c.Next()
    })

    router.POST("/library", handlers.CreateLibrary(db))

    // Attempt to create the same library
    payload, _ := json.Marshal(map[string]string{
        "name": "Duplicate Library",
    })
    req, _ := http.NewRequest("POST", "/library", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Expect 400 when library already exists
    assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestCreateLibrary_BadRequest tests if missing 'name' yields 400 Bad Request
func TestCreateLibrary_BadRequest(t *testing.T) {
    db := setupTestDB(t)
    gin.SetMode(gin.TestMode)
    router := gin.New()

    // Mock Owner claims
    router.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{"role": "Owner"}
        c.Set("user", claims)
        c.Next()
    })

    router.POST("/library", handlers.CreateLibrary(db))

    // Missing required field "name"
    payload, _ := json.Marshal(map[string]string{})
    req, _ := http.NewRequest("POST", "/library", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestGetLibraries_EmptyDB checks retrieving libraries when DB has none.
func TestGetLibraries_EmptyDB(t *testing.T) {
    db := setupTestDB(t)
    gin.SetMode(gin.TestMode)

    router := gin.New()
    // GET /libraries is open/public (assuming your code allows that). If protected,
    // you'd also add JWT claims. Adjust as needed.
    router.GET("/libraries", handlers.GetLibraries(db))

    req, _ := http.NewRequest("GET", "/libraries", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 200 OK, but empty list
    assert.Equal(t, http.StatusOK, w.Code)

    var response map[string][]models.Library
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    libs, ok := response["libraries"]
    assert.True(t, ok, "Expected 'libraries' field in JSON response")
    assert.Empty(t, libs, "Expected an empty list of libraries")
}

// TestGetLibraries_WithSomeLibraries checks retrieving multiple libraries.
func TestGetLibraries_WithSomeLibraries(t *testing.T) {
    db := setupTestDB(t)
    // Seed multiple libraries
    libsToAdd := []models.Library{
        {Name: "Library A"},
        {Name: "Library B"},
        {Name: "Library C"},
    }
    db.Create(&libsToAdd)

    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.GET("/libraries", handlers.GetLibraries(db))

    req, _ := http.NewRequest("GET", "/libraries", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var response map[string][]models.Library
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    libs, ok := response["libraries"]
    assert.True(t, ok)
    assert.Len(t, libs, 3)
}
