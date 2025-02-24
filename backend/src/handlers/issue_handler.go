// /backend/handlers/issue.go
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/swapxs/LibMS/backend/src/models"
	"gorm.io/gorm"
)

// CreateIssueRequest reuses RaiseRequest functionality.
func CreateIssueRequest(db *gorm.DB) gin.HandlerFunc {
	return RaiseRequest(db)
}

// IssueRequestDetail represents the joined issue request data.
type IssueRequestDetail struct {
    ReqID              uint       `json:"ReqID"`
    BookID             string     `json:"BookID"`
    BookName           string     `json:"BookName"`
    ReaderID           uint       `json:"ReaderID"`
    ReaderName         string     `json:"ReaderName"`
    RequestDate        *time.Time `json:"-"` // stored in DB query, hidden in JSON
    ApprovalDate       *time.Time `json:"-"` // same
    ApproverID         *uint      `json:"ApproverID"`
    IssueApproverEmail *string    `json:"IssueApproverEmail"`
    RequestType        string     `json:"RequestType"`
    IssueStatus        string     `json:"IssueStatus"`
    ReturnApproverEmail *string   `json:"ReturnApproverEmail"`
    ReturnStatus       string     `json:"ReturnStatus"`

    // Here are the string versions that we'll populate manually:
    RequestDateStr     string `json:"RequestDate"`
    ApprovalDateStr    string `json:"ApprovalDate"`
}

// GetIssueRequests returns issue requests with extra joined details.
func GetIssueRequests(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get library ID from token claims.
        claims := c.MustGet("user").(jwt.MapClaims)
        libraryID, err := getUintFromClaim(claims, "library_id")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // Query that retrieves raw timestamps, no to_char calls:
        rawQuery := `
            SELECT 
                re.req_id AS "ReqID",
                re.book_id AS "BookID",
                bi.title AS "BookName",
                re.reader_id AS "ReaderID",
                ru.name AS "ReaderName",
                re.request_date AS "RequestDate",
                re.approval_date AS "ApprovalDate",
                re.approver_id AS "ApproverID",
                ia.email AS "IssueApproverEmail",
                re.request_type AS "RequestType",
                CASE 
                    WHEN re.request_type = 'Approve' OR ir.issue_status = 'Approved' THEN 'Approved'
                    WHEN re.request_type = 'Issue' THEN 'Pending'
                    WHEN re.request_type = 'Reject' THEN 'Rejected'
                    ELSE 'Pending'
                END AS "IssueStatus",
                ret_ia.email AS "ReturnApproverEmail",
                CASE 
                    WHEN ir.return_date IS NOT NULL THEN 'Returned'
                    ELSE 'Not Returned'
                END AS "ReturnStatus"
            FROM request_events re
            JOIN book_inventories bi ON re.book_id = bi.isbn
            JOIN users ru ON re.reader_id = ru.id
            LEFT JOIN users ia ON re.approver_id = ia.id
            LEFT JOIN issue_registries ir ON re.book_id = ir.isbn AND re.reader_id = ir.reader_id
            LEFT JOIN users ret_ia ON ir.return_approver_id = ret_ia.id
            WHERE bi.library_id = ?
            ORDER BY re.req_id ASC
        `

        var details []IssueRequestDetail
        err = db.Raw(rawQuery, libraryID).Scan(&details).Error
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // Format the timestamps in Go (RFC3339).
        // If you'd prefer a custom layout, replace time.RFC3339 with e.g. "2006-01-02T15:04:05Z07:00"
        for i := range details {
            if details[i].RequestDate != nil {
                details[i].RequestDateStr = details[i].RequestDate.Format(time.RFC3339)
            }
            if details[i].ApprovalDate != nil {
                details[i].ApprovalDateStr = details[i].ApprovalDate.Format(time.RFC3339)
            }
        }

        c.JSON(http.StatusOK, gin.H{"success": true, "requests": details})
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

		// ─────────────────────────────────────────────────────────
		//  Fix: Reject unknown request types with a 400 Bad Request
		// ─────────────────────────────────────────────────────────
		if input.RequestType != "Approve" && input.RequestType != "Reject" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request type. Must be 'Approve' or 'Reject'."})
			return
		}

		// If we are approving, check book availability
		if input.RequestType == "Approve" {
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

		// For both Approve and Reject, we set approval date, type, and approver
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
