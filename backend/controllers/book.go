package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/swapxs/LibMS/backend/models"
	"gorm.io/gorm"
)

// AddBookInput represents the payload for adding a new book.
type AddBookInput struct {
	ISBN          string `json:"isbn" binding:"required"`
	Title         string `json:"title" binding:"required"`
	Author        string `json:"author" binding:"required"`   // Changed field name to "Author"
	Publisher     string `json:"publisher"`
	Language      string `json:"language"`                      // New field for language
	Version       string `json:"version"`
	Copies        int    `json:"copies" binding:"required,gt=0"`
	IncrementOnly bool   `json:"increment_only"`
}

// AddOrIncrementBook adds a new book or increments copies if the book already exists.
// If an existing record has zero copies, it deletes that record before creating a new one.
func AddOrIncrementBook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input AddBookInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get library ID from JWT claims.
		claims := c.MustGet("user").(jwt.MapClaims)
		libraryID := uint(claims["library_id"].(float64))

		var book models.BookInventory
		err := db.Where("isbn = ? AND library_id = ?", input.ISBN, libraryID).First(&book).Error

		if err == nil {
			// Book record exists.
			if book.TotalCopies == 0 {
				// Delete record if total copies is zero.
				if err := db.Delete(&book).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			} else {
				// Increment the existing record's copies.
				book.TotalCopies += input.Copies
				book.AvailableCopies += input.Copies
				if err := db.Save(&book).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "Book copies incremented", "book": book})
				return
			}
		} else if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// At this point, either the book doesn't exist or was deleted.
		// Ensure Title, Author, and other required fields are provided.
		if input.Title == "" || input.Author == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Title and Author are required for a new book"})
			return
		}

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

// UpdateBookInput is the payload for updating a book.
type UpdateBookInput struct {
	Title       string `json:"title"`
	Author     string `json:"authors"`
	Publisher   string `json:"publisher"`
	Version     string `json:"version"`
	TotalCopies *int   `json:"total_copies"`
}

// UpdateBook updates details of a book.
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
		var input UpdateBookInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if input.Title != "" {
			book.Title = input.Title
		}
		if input.Author != "" {
			book.Author = input.Author
		}
		if input.Publisher != "" {
			book.Publisher = input.Publisher
		}
		if input.Version != "" {
			book.Version = input.Version
		}
		if input.TotalCopies != nil {
			diff := *input.TotalCopies - book.TotalCopies
			book.TotalCopies = *input.TotalCopies
			book.AvailableCopies += diff
			if book.AvailableCopies < 0 {
				book.AvailableCopies = 0
			}
		}
		if err := db.Save(&book).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Book updated", "book": book})
	}
}

// RemoveBook removes copies of a book (or deletes the record if copies become 0).
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
		issued := book.TotalCopies - book.AvailableCopies
		if payload.Copies > (book.TotalCopies - issued) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot remove issued copies"})
			return
		}
		if book.TotalCopies-payload.Copies == 0 {
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
