// /backend/test/issue_requests_test.go
package handlers_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strconv"
    "testing"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "github.com/stretchr/testify/assert"
    "github.com/swapxs/LibMS/backend/handlers"
    "github.com/swapxs/LibMS/backend/models"
)

// TestCreateIssueRequest_Success checks the scenario where a user successfully creates an issue request.
// "CreateIssueRequest" internally reuses RaiseRequest logic.
func TestCreateIssueRequest_Success(t *testing.T) {
    db := setupTestDB(t)

    // Seed a user.
    user := models.User{
        Name:          "IssueRequestUser",
        Email:         "issuereq@example.com",
        Password:      "hashed",
        ContactNumber: "123",
        Role:          "Reader",
        LibraryID:     1,
    }
    db.Create(&user)

    // Seed a book that is available.
    book := models.BookInventory{
        ISBN:            "issue-req-isbn",
        LibraryID:       1,
        Title:           "Issue Request Book",
        Author:          "Author",
        Publisher:       "Publisher",
        Language:        "English",
        Version:         "v1",
        TotalCopies:     5,
        AvailableCopies: 5,
    }
    db.Create(&book)

    gin.SetMode(gin.TestMode)
    router := gin.New()

    // Middleware to inject Reader claims
    router.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "id":         float64(user.ID),
            "role":       user.Role,
            "library_id": float64(user.LibraryID),
        }
        c.Set("user", claims)
        c.Next()
    })

    // We call the "CreateIssueRequest" handler, which reuses RaiseRequest
    router.POST("/issueRequests", handlers.CreateIssueRequest(db))

    requestBody, _ := json.Marshal(map[string]string{
        "bookID": "issue-req-isbn",
    })
    req, _ := http.NewRequest("POST", "/issueRequests", bytes.NewBuffer(requestBody))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)

    // Verify a RequestEvent was created with "Issue"
    var reqEvent models.RequestEvent
    err := db.First(&reqEvent, "book_id = ? AND reader_id = ?", "issue-req-isbn", user.ID).Error
    assert.NoError(t, err)
    assert.Equal(t, "Issue", reqEvent.RequestType)
}

// TestGetIssueRequests_Success seeds a few request events and checks if we can retrieve them.
func TestGetIssueRequests_Success(t *testing.T) {
    db := setupTestDB(t)

    // Seed a user with library_id=10
    user := models.User{
        Name:          "IssueRequestsUser",
        Email:         "issuerequests@example.com",
        Password:      "hashed",
        ContactNumber: "999",
        Role:          "LibraryAdmin",
        LibraryID:     10,
    }
    db.Create(&user)

    // Seed two book records, each with library_id=10
    book1 := models.BookInventory{
        ISBN:            "isbn1",
        LibraryID:       10,
        Title:           "Book 1",
        Author:          "Author 1",
        Publisher:       "Pub 1",
        Language:        "English",
        Version:         "v1",
        TotalCopies:     5,
        AvailableCopies: 5,
    }
    db.Create(&book1)

    book2 := models.BookInventory{
        ISBN:            "isbn2",
        LibraryID:       10,
        Title:           "Book 2",
        Author:          "Author 2",
        Publisher:       "Pub 2",
        Language:        "English",
        Version:         "v1",
        TotalCopies:     5,
        AvailableCopies: 5,
    }
    db.Create(&book2)

    // Now seed two request events referencing these books + user.
    re1 := models.RequestEvent{
        BookID:      "isbn1",
        ReaderID:    user.ID,
        RequestDate: time.Now(),
        RequestType: "Issue",
    }
    re2 := models.RequestEvent{
        BookID:      "isbn2",
        ReaderID:    user.ID,
        RequestDate: time.Now(),
        RequestType: "Approve",
    }
    db.Create(&re1)
    db.Create(&re2)

    gin.SetMode(gin.TestMode)
    router := gin.New()

    // Inject JWT claims referencing the user
    router.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "id":         float64(user.ID),
            "role":       user.Role,
            "library_id": float64(user.LibraryID),
        }
        c.Set("user", claims)
        c.Next()
    })

    // The actual endpoint under test
    router.GET("/issueRequests", handlers.GetIssueRequests(db))

    req, _ := http.NewRequest("GET", "/issueRequests", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // We expect 200 OK since we have 2 request events for library_id=10
    assert.Equal(t, http.StatusOK, w.Code)

    // Optionally parse JSON response to ensure the data is present
    var resp struct {
        Success  bool                           `json:"success"`
        Requests []handlers.IssueRequestDetail  `json:"requests"`
    }
    json.Unmarshal(w.Body.Bytes(), &resp)

    assert.True(t, resp.Success)
    assert.Len(t, resp.Requests, 2, "Expected 2 request events in response")
}

// TestIssueBook_Success checks a scenario in which we successfully create an IssueRegistry record.
func TestIssueBook_Success(t *testing.T) {
    db := setupTestDB(t)
    gin.SetMode(gin.TestMode)

    router := gin.New()
    // JWT claims as an admin or owner
    router.Use(func(c *gin.Context) {
        c.Set("user", jwt.MapClaims{
            "id":         float64(999),
            "role":       "LibraryAdmin",
            "library_id": float64(1),
        })
        c.Next()
    })
    router.POST("/issueRegistry", handlers.IssueBook(db))

    // Provide all required fields, including expected_return_date
    futureDate := time.Now().Add(48 * time.Hour)
    requestBody, _ := json.Marshal(map[string]interface{}{
        "isbn":               "issuebookisbn",
        "reader_id":          1,
        "issue_approver_id":  999,
        "issue_status":       "Issued",
        "expected_return_date": futureDate,
        "library_id":         1,
    })

    req, _ := http.NewRequest("POST", "/issueRegistry", bytes.NewBuffer(requestBody))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    // Check that the record was created
    var record models.IssueRegistry
    err := db.First(&record, "isbn = ? AND reader_id = ?", "issuebookisbn", 1).Error
    assert.NoError(t, err)
    assert.Equal(t, "Issued", record.IssueStatus)
    assert.WithinDuration(t, futureDate, record.ExpectedReturnDate, time.Second)
}

// TestIssueBook_MissingFields omits expected_return_date => 400
func TestIssueBook_MissingFields(t *testing.T) {
    db := setupTestDB(t)
    gin.SetMode(gin.TestMode)
    router := gin.New()

    // JWT claims
    router.Use(func(c *gin.Context) {
        c.Set("user", jwt.MapClaims{
            "id":         float64(100),
            "role":       "LibraryAdmin",
            "library_id": float64(1),
        })
        c.Next()
    })
    router.POST("/issueRegistry", handlers.IssueBook(db))

    // Missing "expected_return_date"
    requestBody, _ := json.Marshal(map[string]interface{}{
        "isbn":             "whatever",
        "reader_id":        2,
        "issue_approver_id": 100,
        "issue_status":     "Issued",
        "library_id":       1,
    })

    req, _ := http.NewRequest("POST", "/issueRegistry", bytes.NewBuffer(requestBody))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestUpdateIssueRequestStatus_ApproveWithZeroCopies seeds a request + a book with 0 copies => 400
func TestUpdateIssueRequestStatus_ApproveWithZeroCopies(t *testing.T) {
    db := setupTestDB(t)

    // Seed a book with 0 available copies
    book := models.BookInventory{
        ISBN:            "approvezero",
        LibraryID:       1,
        Title:           "No Copies Book",
        Author:          "No Author",
        Publisher:       "No Publisher",
        Language:        "English",
        Version:         "1st",
        TotalCopies:     3,
        AvailableCopies: 0, // none available
    }
    db.Create(&book)

    // Seed a request event
    reqEvent := models.RequestEvent{
        BookID:      "approvezero",
        ReaderID:    2,
        RequestDate: time.Now(),
        RequestType: "Issue", // means "waiting for approval"
    }
    db.Create(&reqEvent)

    gin.SetMode(gin.TestMode)
    router := gin.New()

    // JWT claims
    router.Use(func(c *gin.Context) {
        c.Set("user", jwt.MapClaims{
            "id":         float64(999),
            "role":       "LibraryAdmin",
            "library_id": float64(1),
        })
        c.Next()
    })
    router.PUT("/issueRequests/:id", handlers.UpdateIssueRequestStatus(db))

    // Approve request
    requestBody, _ := json.Marshal(map[string]interface{}{
        "request_type": "Approve",
    })
    url := "/issueRequests/" + strconv.Itoa(int(reqEvent.ReqID))
    req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestUpdateIssueRequestStatus_NotFound tries to update a non-existent ReqID => 404
func TestUpdateIssueRequestStatus_NotFound(t *testing.T) {
    db := setupTestDB(t)
    gin.SetMode(gin.TestMode)
    router := gin.New()

    // JWT claims
    router.Use(func(c *gin.Context) {
        claims := jwt.MapClaims{
            "id":         float64(55),
            "role":       "LibraryAdmin",
            "library_id": float64(2),
        }
        c.Set("user", claims)
        c.Next()
    })
    router.PUT("/issueRequests/:id", handlers.UpdateIssueRequestStatus(db))

    // We haven't created any request events; ID "999" won't exist
    requestBody, _ := json.Marshal(map[string]interface{}{
        "request_type": "Approve",
    })
    req, _ := http.NewRequest("PUT", "/issueRequests/999", bytes.NewBuffer(requestBody))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestUpdateIssueRequestStatus_UnknownType tries a "weird" request_type => code might just set it and do 200.
func TestUpdateIssueRequestStatus_UnknownType(t *testing.T) {
    db := setupTestDB(t)

    // Seed a request event
    reqEvent := models.RequestEvent{
        BookID:      "someisbn",
        ReaderID:    10,
        RequestDate: time.Now(),
        RequestType: "Issue",
    }
    db.Create(&reqEvent)

    gin.SetMode(gin.TestMode)
    router := gin.New()

    // JWT claims
    router.Use(func(c *gin.Context) {
        c.Set("user", jwt.MapClaims{
            "id":         float64(111),
            "role":       "LibraryAdmin",
            "library_id": float64(1),
        })
        c.Next()
    })
    router.PUT("/issueRequests/:id", handlers.UpdateIssueRequestStatus(db))

    // Provide an unknown request type
    requestBody, _ := json.Marshal(map[string]string{
        "request_type": "WeirdStatus",
    })
    url := "/issueRequests/" + strconv.Itoa(int(reqEvent.ReqID))
    req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}

