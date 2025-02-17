package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/swapxs/LibMS/backend/controllers"
	"github.com/swapxs/LibMS/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.User{}, &models.Library{}, &models.BookInventory{}, &models.RequestEvent{}, &models.IssueRegistry{})
	return db
}

func TestRegisterUser(t *testing.T) {
	testDB := setupTestDB()
	router := gin.Default()
	router.POST("/register", controllers.RegisterUser(testDB))

	input := controllers.RegisterInput{
		Name:          "Test User",
		Email:         "test@example.com",
		Password:      "password",
		ContactNumber: "1234567890",
		Role:          "Reader",
		LibraryID:     0,
	}
	jsonValue, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}
