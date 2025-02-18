// /backend/controllers/issue.go
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

// CreateIssueRequest reuses RaiseRequest functionality.
func CreateIssueRequest(db *gorm.DB) gin.HandlerFunc {
	return RaiseRequest(db)
}

// IssueRequestDetail represents the joined issue request data.
type IssueRequestDetail struct {
	ReqID              uint   `json:"ReqID"`
	BookID             string `json:"BookID"`
	BookName           string `json:"BookName"`
	ReaderID           uint   `json:"ReaderID"`
	ReaderName         string `json:"ReaderName"`
	RequestDate        string `json:"RequestDate"`
	ApprovalDate       string `json:"ApprovalDate"`
	ApproverID         uint   `json:"ApproverID"`
	IssueApproverEmail string `json:"IssueApproverEmail"`
	RequestType        string `json:"RequestType"`
	// Since request_events doesn't store return data, we'll default these:
	ReturnApproverEmail string `json:"ReturnApproverEmail"`
	ReturnStatus        string `json:"ReturnStatus"`
}

// GetIssueRequests returns issue requests with extra joined details.
func GetIssueRequests(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get library id from token claims.
		claims := c.MustGet("user").(jwt.MapClaims)
		libraryID := uint(claims["library_id"].(float64))

		var details []IssueRequestDetail
		// Join request_events (re), book_inventories (bi), and users (for reader and issue approver)
		err := db.Raw(`
			SELECT 
				re.req_id as "ReqID",
				re.book_id as "BookID",
				bi.title as "BookName",
				re.reader_id as "ReaderID",
				ru.name as "ReaderName",
				to_char(re.request_date, 'YYYY-MM-DD"T"HH24:MI:SSZ') as "RequestDate",
				COALESCE(to_char(re.approval_date, 'YYYY-MM-DD"T"HH24:MI:SSZ'), '') as "ApprovalDate",
				re.approver_id as "ApproverID",
				ia.email as "IssueApproverEmail",
				re.request_type as "RequestType",
				'N/A' as "ReturnApproverEmail",
				'N/A' as "ReturnStatus"
			FROM request_events re
			JOIN book_inventories bi ON re.book_id = bi.isbn
			JOIN users ru ON re.reader_id = ru.id
			LEFT JOIN users ia ON re.approver_id = ia.id
			WHERE bi.library_id = ?
			ORDER BY re.req_id ASC
		`, libraryID).Scan(&details).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"requests": details})
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

		var input struct {
			RequestType        string     `json:"request_type" binding:"required"`
			ExpectedReturnDate *time.Time `json:"expected_return_date"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		claims := c.MustGet("user").(jwt.MapClaims)
		approverID := uint(claims["id"].(float64))

		if input.RequestType == "Approve" {
			// Reduce available copies by 1.
			var book models.BookInventory
			if err := db.Where("isbn = ?", reqEvent.BookID).First(&book).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
				return
			}

			if book.AvailableCopies < 1 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "No available copies left to issue"})
				return
			}

			book.AvailableCopies -= 1
			if err := db.Save(&book).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book availability"})
				return
			}
		}

		now := time.Now()
		reqEvent.ApprovalDate = &now
		reqEvent.RequestType = input.RequestType
		reqEvent.ApproverID = &approverID

		if input.ExpectedReturnDate != nil {
			reqEvent.ApprovalDate = input.ExpectedReturnDate
		}

		if err := db.Save(&reqEvent).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Issue request status updated and available copies adjusted", "request": reqEvent})
	}
}

// IssueBook creates an issuance record in the issue_registries table.
func IssueBook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.IssueRegistry
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Set issue date to now.
		payload.IssueDate = time.Now()
		// Check that expected_return_date is not zero.
		if payload.ExpectedReturnDate.IsZero() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Expected return date is required"})
			return
		}
		if err := db.Create(&payload).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Book issued", "issue": payload})
	}
}
