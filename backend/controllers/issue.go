package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/swapxs/LibMS/backend/models"
	"gorm.io/gorm"
)

// CreateIssueRequest creates an issue request record (reusing RaiseRequest functionality).
func CreateIssueRequest(db *gorm.DB) gin.HandlerFunc {
	return RaiseRequest(db)
}

// GetIssueRequests returns all issue requests for the library.
func GetIssueRequests(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("user").(jwt.MapClaims)
		libraryID := uint(claims["library_id"].(float64))
		var requests []models.RequestEvent
		if err := db.Table("request_events").
			Select("request_events.*").
			Joins("JOIN book_inventories ON request_events.book_id = book_inventories.isbn").
			Where("book_inventories.library_id = ?", libraryID).
			Find(&requests).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"requests": requests})
	}
}

// UpdateIssueRequestStatus updates the status of an issue request.
func UpdateIssueRequestStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		reqID, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
			return
		}
		var reqEvent models.RequestEvent
		if err := db.First(&reqEvent, reqID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
			return
		}
		type StatusInput struct {
			RequestType string `json:"request_type" binding:"required"`
		}
		var input StatusInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		now := time.Now()
		reqEvent.ApprovalDate = &now
		reqEvent.RequestType = input.RequestType
		claims := c.MustGet("user").(jwt.MapClaims)
		approverID := uint(claims["id"].(float64))
		reqEvent.ApproverID = &approverID
		if err := db.Save(&reqEvent).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Issue request status updated", "request": reqEvent})
	}
}

// IssueBook approves an issue request and creates an issue registry entry.
func IssueBook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("user").(jwt.MapClaims)
		role := claims["role"].(string)
		if role != "LibraryAdmin" && role != "Owner" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		idParam := c.Param("id")
		reqID, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
			return
		}
		var reqEvent models.RequestEvent
		if err := db.First(&reqEvent, reqID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
			return
		}
		now := time.Now()
		approverID := uint(claims["id"].(float64))
		reqEvent.ApprovalDate = &now
		reqEvent.ApproverID = &approverID
		reqEvent.RequestType = "Approved"
		if err := db.Save(&reqEvent).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		issue := models.IssueRegistry{
			ISBN:               reqEvent.BookID,
			ReaderID:           reqEvent.ReaderID,
			IssueApproverID:    approverID,
			IssueStatus:        "Issued",
			IssueDate:          now,
			ExpectedReturnDate: now.AddDate(0, 0, 14),
			LibraryID:          uint(claims["library_id"].(float64)),
		}
		if err := db.Create(&issue).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Book issued", "issue": issue})
	}
}
