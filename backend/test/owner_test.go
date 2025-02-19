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

func setupOwnerTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Library{}, &models.User{})
	return db
}

func TestRegisterLibraryOwner(t *testing.T) {
	db := setupOwnerTestDB()
	router := gin.Default()
	router.POST("/register-owner", controllers.RegisterLibraryOwner(db))
	input := controllers.RegisterOwnerInput{
		Name:          "Owner Name",
		Email:         "owner@example.com",
		Password:      "password",
		ContactNumber: "1234567890",
		LibraryName:   "Owner Library",
	}
	jsonValue, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/register-owner", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}
