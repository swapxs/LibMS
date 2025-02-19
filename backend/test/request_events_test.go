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
)

func TestRaiseRequest(t *testing.T) {
	db := setupBookTestDB()
	router := gin.Default()

	var readerUser models.User
	db.Where("role = ?", "Reader").First(&readerUser)
	addJWTAuthMiddleware(router, readerUser)

	router.POST("/request", controllers.RaiseRequest(db))

	input := controllers.RaiseRequestInput{
		BookID: "9876543210",
	}
	jsonValue, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/request", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
