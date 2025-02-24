// /backend/handlers/user.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/swapxs/LibMS/backend/src/models"
	"gorm.io/gorm"
)

// GetUsers returns all users for the library.
func GetUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("user").(jwt.MapClaims)
		libraryID, err := getUintFromClaim(claims, "library_id")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var users []models.User
		if err := db.Where("library_id = ?", libraryID).Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": users})
	}
}

// GetUserIssueInfo returns issue registry records for the logged-in user.
func GetUserIssueInfo(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("user").(jwt.MapClaims)
        userID, err := getUintFromClaim(claims, "id")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var issues []models.IssueRegistry

		if err := db.Where("reader_id = ?", userID).Find(&issues).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
            "success": true,
            "data": issues,
        })
	}
}
