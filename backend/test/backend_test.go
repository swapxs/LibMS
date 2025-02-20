package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/swapxs/LibMS/backend/handlers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/golang-jwt/jwt/v4"
)

// ✅ Include the setupTestDB function from test_utils.go
func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

// ✅ Test for User Login
func TestLogin(t *testing.T) {
	db, mock := setupTestDB(t)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/login", handlers.Login(db))

	mock.ExpectQuery(`SELECT \* FROM "book_inventories" WHERE \(isbn = \$1 AND library_id = \$2\) AND "book_inventories"."deleted_at" IS NULL ORDER BY "book_inventories"."id" LIMIT \$3`).
		WithArgs("12345", 1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"isbn", "library_id"}).
			AddRow("12345", 1))

	requestBody, _ := json.Marshal(map[string]string{
		"email":    "john@example.com",
		"password": "password123",
	})
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ✅ Test for Adding a Book
func TestAddOrIncrementBook(t *testing.T) {
    db, mock := setupTestDB(t)

    gin.SetMode(gin.TestMode)
    router := gin.New()

    // ✅ Mock JWT authentication with `jwt.MapClaims`
    router.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "id":         float64(1),
            "email":      "john@example.com",
            "role":       "LibraryAdmin",
            "library_id": float64(1),
        }
        c.Set("user", claims)
        c.Next()
    })

    router.POST("/books", handlers.AddOrIncrementBook(db))

    // ✅ Corrected SQL Mock Query
    mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT \$2`).
        WithArgs("john@example.com", 1).
        WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
            AddRow(1, "john@example.com", "$2a$10$hashedPassword"))

    mock.ExpectExec("INSERT INTO book_inventories").
        WillReturnResult(sqlmock.NewResult(1, 1))

    requestBody, _ := json.Marshal(map[string]interface{}{
        "isbn":     "12345",
        "title":    "Golang Book",
        "author":   "John Doe",
        "copies":   5,
        "language": "English",
    })
    req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(requestBody))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)
}

// ✅ Test for Raising an Issue Request
func TestRaiseRequest(t *testing.T) {
	db, mock := setupTestDB(t)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		claims := jwt.MapClaims{
			"id":         float64(1),
			"role":       "Reader",
			"library_id": float64(1),
		}
		c.Set("user", claims)
		c.Next()
	})
	router.POST("/requestEvents", handlers.RaiseRequest(db))

	mock.ExpectQuery("SELECT * FROM request_events WHERE reader_id = ?").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	mock.ExpectExec("INSERT INTO request_events").
		WillReturnResult(sqlmock.NewResult(1, 1))

	requestBody, _ := json.Marshal(map[string]interface{}{
		"bookID": "12345",
	})
	req, _ := http.NewRequest("POST", "/requestEvents", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

// ✅ Test for Creating a Library
func TestCreateLibrary(t *testing.T) {
	db, mock := setupTestDB(t)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/library", handlers.CreateLibrary(db))

	mock.ExpectQuery(`SELECT \* FROM "libraries" WHERE name = \$1 AND "libraries"."deleted_at" IS NULL ORDER BY "libraries"."id" LIMIT \$2`).
		WithArgs("New Library", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "New Library"))

	mock.ExpectExec("INSERT INTO libraries").WillReturnResult(sqlmock.NewResult(1, 1))

	requestBody, _ := json.Marshal(map[string]string{"name": "New Library"})
	req, _ := http.NewRequest("POST", "/library", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

