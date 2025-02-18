package controllers

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

// RaiseRequest allows a reader to raise an issue request.
// It checks if the requested book is available (i.e. available copies > 0)
// and uses the reader's ID from the token.
func RaiseRequest(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RaiseRequestInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get user info from token claims.
		claims, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		tokenClaims := claims.(jwt.MapClaims)
		readerID := uint(tokenClaims["id"].(float64))
		libraryID := uint(tokenClaims["library_id"].(float64))

		// Check if the book exists in the library and is available.
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
