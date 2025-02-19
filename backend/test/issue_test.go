package test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/swapxs/LibMS/backend/controllers"
	"github.com/swapxs/LibMS/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupIssueTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Library{}, &models.User{}, &models.BookInventory{}, &models.RequestEvent{}, &models.IssueRegistry{})
	return db
}

func TestIssueBook(t *testing.T) {
	db := setupIssueTestDB()
	// Create dummy library, admin, reader, and request event.
	library := models.Library{Name: "Test Library"}
	db.Create(&library)
	admin := models.User{
		Name:      "Admin",
		Email:     "admin@example.com",
		Password:  "dummy",
		Role:      "LibraryAdmin",
		LibraryID: library.ID,
	}
	db.Create(&admin)
	reader := models.User{
		Name:      "Reader",
		Email:     "reader@example.com",
		Password:  "dummy",
		Role:      "Reader",
		LibraryID: library.ID,
	}
	db.Create(&reader)
	reqEvent := models.RequestEvent{
		BookID:      "1111111111",
		ReaderID:    reader.ID,
		RequestDate: time.Now(),
		RequestType: "Issue",
	}
	db.Create(&reqEvent)
	router := gin.Default()
	// Simulate JWT middleware for admin.
	router.POST("/issue/:id", func(c *gin.Context) {
		claims := map[string]interface{}{
			"id":         float64(admin.ID),
			"email":      admin.Email,
			"role":       admin.Role,
			"library_id": float64(library.ID),
		}
		c.Set("user", claims)
		controllers.IssueBook(db)(c)
	})
	req, _ := http.NewRequest("POST", "/issue/"+strconv.Itoa(int(reqEvent.ReqID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
