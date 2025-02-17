package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swapxs/LibMS/backend/models"
	"gorm.io/gorm"
)

// CreateLibraryInput is the payload for creating a library.
type CreateLibraryInput struct {
	Name string `json:"name" binding:"required"`
}

// CreateLibrary creates a new library.
func CreateLibrary(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateLibraryInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Check for existing library.
		var lib models.Library
		if err := db.Where("name = ?", input.Name).First(&lib).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Library already exists"})
			return
		}
		lib = models.Library{Name: input.Name}
		if err := db.Create(&lib).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Library created", "library": lib})
	}
}

// GetLibraries returns all libraries.
func GetLibraries(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var libraries []models.Library
		if err := db.Find(&libraries).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"libraries": libraries})
	}
}
