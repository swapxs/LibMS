package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/swapxs/LibMS/backend/controllers"
	"github.com/swapxs/LibMS/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	db := setupAuthTestDB()
	router := gin.Default()

	var adminUser models.User
	db.Where("role = ?", "LibraryAdmin").First(&adminUser)
	addJWTAuthMiddleware(router, adminUser)

	router.GET("/users", controllers.GetUsers(db))

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
