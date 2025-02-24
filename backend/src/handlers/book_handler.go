// /backend/handlers/book.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/swapxs/LibMS/backend/src/models"
	"gorm.io/gorm"
)

// AddBookInput represents the payload for adding a book.
type AddBookInput struct {
	ISBN          string `json:"isbn" binding:"required"`
	Title         string `json:"title"`  // Required only when creating a new book
	Author        string `json:"author"` // Required only when creating a new book
	Publisher     string `json:"publisher"`
	Language      string `json:"language"` // Required only when creating a new book
	Version       string `json:"version"`
	Copies        int    `json:"copies" binding:"required,gt=0"`
	IncrementOnly bool   `json:"increment_only"`
}

// AddOrIncrementBook adds a new book or increments copies if the book already exists.
// If the book exists, it ignores Title, Author, and Language.
// If no record is found and IncrementOnly is true, it returns an error.
func AddOrIncrementBook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input AddBookInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		claims := c.MustGet("user").(jwt.MapClaims)
		libraryID, err := getUintFromClaim(claims, "library_id")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}


		var book models.BookInventory
		err = db.Where("isbn = ? AND library_id = ?", input.ISBN, libraryID).First(&book).Error

		if err == nil {
			// Book record exists: Increment copies.
			book.TotalCopies += input.Copies
			book.AvailableCopies += input.Copies
			if err := db.Save(&book).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"message": "Book copies incremented", "book": book})
			return
		} else if err != gorm.ErrRecordNotFound {
			// Some other error occurred.
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Record not found.
		if input.IncrementOnly {
			// Cannot increment copies if record does not exist.
			c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found, cannot increment copies. Please add new book details."})
			return
		}

		// For a new book, require Title, Author, and Language.
		if input.Title == "" || input.Author == "" || input.Language == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Title, Author, and Language are required for a new book"})
			return
		}

		// Create a new book record.
		newBook := models.BookInventory{
			ISBN:            input.ISBN,
			LibraryID:       libraryID,
			Title:           input.Title,
			Author:          input.Author,
			Publisher:       input.Publisher,
			Language:        input.Language,
			Version:         input.Version,
			TotalCopies:     input.Copies,
			AvailableCopies: input.Copies,
		}
		if err := db.Create(&newBook).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Book added successfully", "book": newBook})
	}
}

// GetBooks returns all books for the library.
func GetBooks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("user").(jwt.MapClaims)
		libraryID := uint(claims["library_id"].(float64))
		var books []models.BookInventory
		if err := db.Where("library_id = ?", libraryID).Find(&books).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"books": books})
	}
}

// RemoveBook removes copies of a book. If removal makes copies 0, delete the record.
func RemoveBook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("user").(jwt.MapClaims)
		libraryID := uint(claims["library_id"].(float64))
		var payload struct {
			ISBN   string `json:"isbn" binding:"required"`
			Copies int    `json:"copies" binding:"required,gt=0"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var book models.BookInventory
		if err := db.Where("isbn = ? AND library_id = ?", payload.ISBN, libraryID).First(&book).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		// Check if removal is valid (can't remove more than available non-issued copies)
		issued := book.TotalCopies - book.AvailableCopies
		if payload.Copies > (book.TotalCopies - issued) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot remove issued copies"})
			return
		}
		if book.TotalCopies-payload.Copies == 0 {
			// Delete record if resulting copies is zero.
			if err := db.Delete(&book).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Book record deleted"})
			return
		}
		book.TotalCopies -= payload.Copies
		if book.AvailableCopies >= payload.Copies {
			book.AvailableCopies -= payload.Copies
		} else {
			book.AvailableCopies = 0
		}
		if err := db.Save(&book).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Book copies removed", "book": book})
	}
}

// UpdateBook updates book details.
func UpdateBook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		isbn := c.Param("isbn")
		claims := c.MustGet("user").(jwt.MapClaims)
		libraryID := uint(claims["library_id"].(float64))
		var book models.BookInventory
		if err := db.Where("isbn = ? AND library_id = ?", isbn, libraryID).First(&book).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		var input map[string]interface{}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Update only provided fields.
		if err := db.Model(&book).Updates(input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Book updated", "book": book})
	}
}
