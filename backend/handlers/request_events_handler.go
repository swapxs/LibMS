// /backend/handlers/request_events.go
package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/swapxs/LibMS/backend/models"
	"gorm.io/gorm"
)

type RaiseRequestInput struct {
	BookID string `json:"bookID" binding:"required"`
}

// RaiseRequest allows a reader to raise an issue request with a limit of 4 active requests.
func RaiseRequest(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RaiseRequestInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		claims, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		tokenClaims := claims.(jwt.MapClaims)
		readerID, err := getUintFromClaim(tokenClaims, "id")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		libraryID, err := getUintFromClaim(tokenClaims, "library_id")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Count active requests (Pending or Approved) for the user.
		var activeRequests int64
		if err := db.Model(&models.RequestEvent{}).
			Where("reader_id = ? AND request_type IN (?)", readerID, []string{"Issue", "Approve"}).
			Count(&activeRequests).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count active requests"})
			return
		}

		if activeRequests >= 4 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Maximum of 4 active requests reached"})
			return
		}

		// Check if the book is available.
		var book models.BookInventory
		if err := db.Where("isbn = ? AND library_id = ?", input.BookID, libraryID).First(&book).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found in your library"})
			return
		}
		if book.AvailableCopies < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Book not available for issue"})
			return
		}

		// Create the request event.
		reqEvent := models.RequestEvent{
			BookID:      input.BookID,
			ReaderID:    readerID,
			RequestDate: time.Now(),
			RequestType: "Issue",
		}
		if err := db.Create(&reqEvent).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Issue request raised", "request": reqEvent})
	}
}
