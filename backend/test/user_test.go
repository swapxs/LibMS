package test

import (
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

func setupUserTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Library{}, &models.User{}, &models.IssueRegistry{})
	return db
}

func TestGetUsers(t *testing.T) {
	db := setupUserTestDB()
	library := models.Library{Name: "Test Library"}
	db.Create(&library)
	user1 := models.User{
		Name:      "User One",
		Email:     "one@example.com",
		Password:  "dummy",
		Role:      "Reader",
		LibraryID: library.ID,
	}
	db.Create(&user1)
	router := gin.Default()
	router.GET("/users", func(c *gin.Context) {
		claims := map[string]interface{}{
			"id":         float64(user1.ID),
			"email":      user1.Email,
			"role":       user1.Role,
			"library_id": float64(library.ID),
		}
		c.Set("user", claims)
		controllers.GetUsers(db)(c)
	})
	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

