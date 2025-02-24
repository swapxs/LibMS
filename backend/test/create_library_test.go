// /backend/test/create_library_test.go
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
)


func TestCreateLibrary(t *testing.T) {
	db := setupTestDB(t)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/library", handlers.CreateLibrary(db))

	requestBody, _ := json.Marshal(map[string]string{
		"name": "New Library",
	})
	req, _ := http.NewRequest("POST", "/library", bytes.NewBuffer(requestBody))
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
