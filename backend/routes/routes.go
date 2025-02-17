package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swapxs/LibMS/backend/controllers"
	"github.com/swapxs/LibMS/backend/middleware"
	"gorm.io/gorm"
)

// SetupRouter sets up all API routes with a complete CORS configuration.
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Apply custom CORS configuration.
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Change this to your frontend domain in production.
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Define all API endpoints under /api.
	api := r.Group("/api")
	{
		// Public endpoints.
		api.POST("/library", controllers.CreateLibrary(db))
		api.GET("/libraries", controllers.GetLibraries(db))
		api.POST("/owner/registration", controllers.RegisterLibraryOwner(db))
		api.POST("/auth/login", controllers.Login(db))
		api.POST("/auth/register", controllers.RegisterUser(db))

		// Protected endpoints.
		protected := api.Group("/")
		protected.Use(middleware.JWTAuthMiddleware())
		{
			protected.GET("/users", controllers.GetUsers(db))
			protected.GET("/auth/userIssueInfo", controllers.GetUserIssueInfo(db))

			// Book endpoints.
			books := protected.Group("/books")
			{
				books.POST("", controllers.AddOrIncrementBook(db))
				books.GET("", controllers.GetBooks(db))
				books.POST("/remove", controllers.RemoveBook(db))
				books.PUT("/:isbn", controllers.UpdateBook(db))
			}

			// Owner endpoints.
			owner := protected.Group("/owner")
			{
				owner.POST("/assign-admin", controllers.AssignAdmin(db))
				owner.POST("/revoke-admin", controllers.RevokeAdmin(db))
			}

			// Request events.
			protected.POST("/requestEvents", controllers.RaiseRequest(db))

			// Issue Request endpoints.
			issue := protected.Group("/issueRequests")
			{
				issue.POST("", controllers.CreateIssueRequest(db))
				issue.GET("", controllers.GetIssueRequests(db))
				issue.PUT("/:id", controllers.UpdateIssueRequestStatus(db))
			}

			// Issue Registry endpoints.
			protected.POST("/issueRegistry", controllers.IssueBook(db))
		}
	}

	return r
}
