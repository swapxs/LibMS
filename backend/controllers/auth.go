package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/swapxs/LibMS/backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// RegisterInput represents the payload for user registration.
type RegisterInput struct {
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Password      string `json:"password" binding:"required,min=6"`
	ContactNumber string `json:"contact_number" binding:"required"`
	// Role is now optional. If not provided, it will default to "Reader".
	Role      string `json:"role"`
	LibraryID uint   `json:"library_id"`
}

// RegisterUser registers a new user. If Role is not provided, it defaults to "Reader".
func RegisterUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RegisterInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Set default role to "Reader" if not provided.
		if input.Role == "" {
			input.Role = "Reader"
		}

		// Check if user already exists.
		var user models.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		}

		// Hash the password.
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
			return
		}

		user = models.User{
			Name:          input.Name,
			Email:         input.Email,
			Password:      string(hashedPassword),
			ContactNumber: input.ContactNumber,
			Role:          input.Role,
			LibraryID:     input.LibraryID,
		}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
	}
}

// LoginInput represents the payload for login.
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login authenticates a user and returns a JWT token.
func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Compare hashed password with the provided password.
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Create a JWT token with claims.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":         user.ID,
			"email":      user.Email,
			"role":       user.Role,
			"library_id": user.LibraryID,
			"exp":        time.Now().Add(72 * time.Hour).Unix(),
		})
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString, "role": user.Role, "library_id": user.LibraryID})
	}
}
