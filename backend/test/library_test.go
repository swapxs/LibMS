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
	"gorm.io/gorm"
)

func setupLibraryTestDB() *gorm.DB {
	return setupAuthTestDB()
}

func TestCreateLibrary(t *testing.T) {
	db := setupLibraryTestDB()
	router := gin.Default()
	router.POST("/library", controllers.CreateLibrary(db))

	input := controllers.CreateLibraryInput{
		Name: "Another Library",
	}
	jsonValue, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/library", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
