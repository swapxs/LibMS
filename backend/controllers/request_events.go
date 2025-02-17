package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/swapxs/LibMS/backend/models"
	"gorm.io/gorm"
)

// RaiseRequestInput is the payload for raising an issue request.
type RaiseRequestInput struct {
	ISBN string `json:"isbn" binding:"required"`
}

// RaiseRequest allows a Reader to raise an issue request.
func RaiseRequest(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RaiseRequestInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		claims := c.MustGet("user").(jwt.MapClaims)
		libraryID := uint(claims["library_id"].(float64))
		// Verify the book exists in the library.
		var book models.BookInventory
		if err := db.Where("isbn = ? AND library_id = ?", input.ISBN, libraryID).First(&book).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		if book.AvailableCopies <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Book not available for issue"})
			return
		}
		// Decrement available copies.
		book.AvailableCopies--
		if err := db.Save(&book).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Find the reader.
		var reader models.User
		if err := db.Where("email = ?", claims["email"]).First(&reader).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Reader not found"})
			return
		}
		reqEvent := models.RequestEvent{
			BookID:      input.ISBN,
			ReaderID:    reader.ID,
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
