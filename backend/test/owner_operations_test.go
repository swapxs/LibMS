// /backend/test/owner_operations_test.go
package handlers_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/swapxs/LibMS/backend/handlers"
    "github.com/swapxs/LibMS/backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/golang-jwt/jwt/v4"
)


// This function tests the owner creation process.
// This will create a new owner and their library alongside 
// This test is going to return success as this is the base operation which is
// to be true all the time.
func TestRegisterLibraryOwner_Success(t *testing.T) {
    db := setupTestDB(t)
    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.POST("/owner/registration", handlers.RegisterLibraryOwner(db))

    payload, _ := json.Marshal(map[string]any {
        "name":           "Owner1",
        "email":          "owner.lib@xenonstack.com",
        "password":       "passwd",
        "contact_number": "2223334444",
        "library_name":   "Test Library",
    })
    req, _ := http.NewRequest("POST", "/owner/registration", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)

    // Check that both the library and the owner were created
    var lib models.Library
    err := db.Where("name = ?", "Test Library").First(&lib).Error
    assert.NoError(t, err)

    var owner models.User
    err = db.Where("email = ?", "owner.lib@xenonstack.com").First(&owner).Error
    assert.NoError(t, err)
    assert.Equal(t, "Owner", owner.Role)
    assert.Equal(t, lib.ID, owner.LibraryID)
}

// Registration with the library already existing
func TestRegisterLibraryOwner_LibraryExists(t *testing.T) {
    db := setupTestDB(t)
    lib := models.Library{Name: "Library"}
    db.Create(&lib)

    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.POST("/owner/registration", handlers.RegisterLibraryOwner(db))

    payload, _ := json.Marshal(map[string]any {
        "name":           "Owner2",
        "email":          "owner2@xenonstack.com",
        "password":       "ownerpass",
        "contact_number": "2223334444",
        "library_name":   "Library",
    })
    req, _ := http.NewRequest("POST", "/owner/registration", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}

// We'll test assigning admin rights (successful)
func TestAssignAdmin_Success(t *testing.T) {
    db := setupTestDB(t)
    owner := models.User{
        Name:          "Test Owner",
        Email:         "testowner@xenonstack.com",
        Password:      "passwd",
        ContactNumber: "1231231231",
        Role:          "Owner",
        LibraryID:     1,
    }
    db.Create(&owner)
    // Seed a normal user (Reader)
    normalUser := models.User{
        Name:          "Test Reader",
        Email:         "reader@xenonstack.com",
        Password:      "hashedpw2",
        ContactNumber: "456",
        Role:          "Reader",
        LibraryID:     1,
    }
    db.Create(&normalUser)

    gin.SetMode(gin.TestMode)
    r := gin.New()

    r.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "id":         float64(owner.ID),
            "email":      owner.Email,
            "role":       owner.Role,
            "library_id": float64(owner.LibraryID),
        }
        c.Set("user", claims)
        c.Next()
    })

    r.POST("/owner/assign-admin", handlers.AssignAdmin(db))

    payload, _ := json.Marshal(map[string]any {
        "email": "reader@xenonstack.com",
    })
    req, _ := http.NewRequest("POST", "/owner/assign-admin", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var updatedUser models.User
    err := db.Where("email = ?", "reader@xenonstack.com").First(&updatedUser).Error
    assert.NoError(t, err)
    assert.Equal(t, "LibraryAdmin", updatedUser.Role)
}

// Assigning Admin as an unauthorized user or just a normal user
func TestAssignAdmin_Unauthorized(t *testing.T) {
    db := setupTestDB(t)

    malicious := models.User{
        Name:          "Not Owner",
        Email:         "notowner@xenonstack.com",
        Password:      "hashedpw",
        ContactNumber: "123",
        Role:          "Reader", 
        LibraryID:     1,
    }

    db.Create(&malicious)

    gin.SetMode(gin.TestMode)

    r := gin.New()
    r.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "id":         float64(malicious.ID),
            "email":      malicious.Email,
            "role":       malicious.Role,
            "library_id": float64(malicious.LibraryID),
        }
        c.Set("user", claims)
        c.Next()
    })

    r.POST("/owner/assign-admin", handlers.AssignAdmin(db))

    payload, _ := json.Marshal(map[string]any {
        "email": "any@xenonstack.com",
    })

    req, _ := http.NewRequest("POST", "/owner/assign-admin", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusUnauthorized, w.Code) //returns 401
}

// Assigning admin but the user is not present
func TestAssignAdmin_UserNotFound(t *testing.T) {
    db := setupTestDB(t)

    // Seed an Owner
    owner := models.User{
        Name:          "Owner Person",
        Email:         "realowner@xenonstack.com",
        Password:      "dummy",
        ContactNumber: "abc",
        Role:          "Owner",
        LibraryID:     2,
    }

    db.Create(&owner)

    gin.SetMode(gin.TestMode)

    r := gin.New()

    r.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "id":         float64(owner.ID),
            "email":      owner.Email,
            "role":       owner.Role,
            "library_id": float64(owner.LibraryID),
        }
        c.Set("user", claims)
        c.Next()
    })

    r.POST("/owner/assign-admin", handlers.AssignAdmin(db))

    payload, _ := json.Marshal(map[string]string{
        "email": "nonexistent@xenonstack.com",
    })

    req, _ := http.NewRequest("POST", "/owner/assign-admin", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestRevokeAdmin_Success seeds an Owner and a LibraryAdmin user, then revokes admin rights.
func TestRevokeAdmin_Success(t *testing.T) {
    db := setupTestDB(t)

    // Seed an Owner user
    owner := models.User{
        Name:          "Test Owner",
        Email:         "owner@xenonstack.com",
        Password:      "hashedownerpw",
        ContactNumber: "123456",
        Role:          "Owner",
        LibraryID:     1,
    }

    db.Create(&owner)

    // Seed a user with LibraryAdmin role
    adminUser := models.User{
        Name:          "Admin User",
        Email:         "admin@xenonstack.com",
        Password:      "hashedadminpw",
        ContactNumber: "555555",
        Role:          "LibraryAdmin",
        LibraryID:     1,
    }

    db.Create(&adminUser)

    gin.SetMode(gin.TestMode)

    r := gin.New()
    // Mock middleware to inject Owner JWT claims
    r.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "id":         float64(owner.ID),
            "email":      owner.Email,
            "role":       owner.Role,
            "library_id": float64(owner.LibraryID),
        }
        c.Set("user", claims)
        c.Next()
    })

    r.POST("/owner/revoke-admin", handlers.RevokeAdmin(db))

    payload, _ := json.Marshal(map[string]string{
        "email": "admin@xenonstack.com",
    })
    req, _ := http.NewRequest("POST", "/owner/revoke-admin", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var updatedUser models.User
    err := db.Where("email = ?", "admin@xenonstack.com").First(&updatedUser).Error
    assert.NoError(t, err)
    assert.Equal(t, "Reader", updatedUser.Role)
}

// TestRevokeAdmin_Unauthorized seeds a non-owner user who tries to revoke admin.
func TestRevokeAdmin_Unauthorized(t *testing.T) {
    db := setupTestDB(t)

    malicious := models.User{
        Name:          "Admin",
        Email:         "hacker@xenonstack.com",
        Password:      "fakehash",
        ContactNumber: "123",
        Role:          "LibraryAdmin",
        LibraryID:     1,
    }
    db.Create(&malicious)

    gin.SetMode(gin.TestMode)
    r := gin.New()

    // Mock middleware to inject non-owner JWT claims
    r.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "id":         float64(malicious.ID),
            "email":      malicious.Email,
            "role":       malicious.Role, // "LibraryAdmin"
            "library_id": float64(malicious.LibraryID),
        }
        c.Set("user", claims)
        c.Next()
    })

    r.POST("/owner/revoke-admin", handlers.RevokeAdmin(db))

    payload, _ := json.Marshal(map[string]string{
        "email": "someadmin@xenonstack.com", 
    })
    req, _ := http.NewRequest("POST", "/owner/revoke-admin", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    // Non-owner user => Expect 401 Unauthorized
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestRevokeAdmin_UserNotFound seeds an Owner but tries to revoke a nonexistent user.
func TestRevokeAdmin_UserNotFound(t *testing.T) {
    db := setupTestDB(t)

    owner := models.User{
        Name:          "Actual Owner",
        Email:         "actualowner@xenonstack.com",
        Password:      "fakepw",
        ContactNumber: "7890",
        Role:          "Owner",
        LibraryID:     5,
    }
    db.Create(&owner)

    gin.SetMode(gin.TestMode)

    r := gin.New()
    r.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "id":         float64(owner.ID),
            "email":      owner.Email,
            "role":       owner.Role, 
            "library_id": float64(owner.LibraryID),
        }
        c.Set("user", claims)
        c.Next()
    })

    r.POST("/owner/revoke-admin", handlers.RevokeAdmin(db))

    payload, _ := json.Marshal(map[string]string{
        "email": "nonexistent@xenonstack.com",
    })

    req, _ := http.NewRequest("POST", "/owner/revoke-admin", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    // The user is not found in the DB => 404
    assert.Equal(t, http.StatusNotFound, w.Code)
}
