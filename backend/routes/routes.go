// /backend/routes/routes.go
package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swapxs/LibMS/backend/handlers"
	"github.com/swapxs/LibMS/backend/middleware"
	"gorm.io/gorm"
)

// SetupRouter configures all routes and applies CORS.
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Custom CORS configuration.
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Change for production.
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api")
	{
		// Public endpoints.
		api.GET("/libraries", handlers.GetLibraries(db))
		api.POST("/owner/registration", handlers.RegisterLibraryOwner(db))
		api.POST("/auth/login", handlers.Login(db))
		api.POST("/auth/register", handlers.RegisterUser(db))

		// Protected endpoints.
		protected := api.Group("/")
		protected.Use(middleware.JWTAuthMiddleware())
		{
			protected.POST("/library", handlers.CreateLibrary(db))
			protected.GET("/users", handlers.GetUsers(db))
			protected.GET("/auth/userIssueInfo", handlers.GetUserIssueInfo(db))
			// Book endpoints.
			books := protected.Group("/books")
			{
				books.POST("", handlers.AddOrIncrementBook(db))
				books.GET("", handlers.GetBooks(db))
				books.POST("/remove", handlers.RemoveBook(db))
				books.PUT("/:isbn", handlers.UpdateBook(db))
			}
			// Owner endpoints.
			owner := protected.Group("/owner")
			{
				owner.POST("/assign-admin", handlers.AssignAdmin(db))
				owner.POST("/revoke-admin", handlers.RevokeAdmin(db))
			}
			// Request events.
			protected.POST("/requestEvents", handlers.RaiseRequest(db))
			// Issue Request endpoints.
			issue := protected.Group("/issueRequests")
			{
				issue.POST("", handlers.CreateIssueRequest(db))
				issue.GET("", handlers.GetIssueRequests(db))
				issue.PUT("/:id", handlers.UpdateIssueRequestStatus(db))
			}
			// Issue Registry endpoint.
			protected.POST("/issueRegistry", handlers.IssueBook(db))
		}
	}

	return r
}
