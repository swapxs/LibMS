package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/swapxs/LibMS/backend/models"
	"gorm.io/gorm"
)

// GetUsers returns all users in the same library.
func GetUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("user").(jwt.MapClaims)
		libraryID := uint(claims["library_id"].(float64))
		var users []models.User
		if err := db.Where("library_id = ?", libraryID).Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}

// GetUserIssueInfo returns the issue registry records for the logged-in user.
func GetUserIssueInfo(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("user").(jwt.MapClaims)
		userID := uint(claims["id"].(float64))
		var issues []models.IssueRegistry
		if err := db.Where("reader_id = ?", userID).Find(&issues).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"issues": issues})
	}
}
