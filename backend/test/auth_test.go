package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/swapxs/LibMS/backend/controllers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	db := setupAuthTestDB()
	router := gin.Default()
	router.POST("/register", controllers.RegisterUser(db))

	input := controllers.RegisterInput{
		Name:          "New User",
		Email:         "newuser@example.com",
		Password:      "securepassword",
		ContactNumber: "1234567890",
		Role:          "Reader",
		LibraryID:     1,
	}
	jsonValue, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestLoginUser(t *testing.T) {
	db := setupAuthTestDB()
	router := gin.Default()
	router.POST("/login", controllers.Login(db))

	input := controllers.LoginInput{
		Email:    "test@example.com",
		Password: "password",
	}
	jsonValue, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
