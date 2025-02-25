// /backend/src/handlers/owner.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/swapxs/LibMS/backend/src/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RegisterOwnerInput is the payload for registering an owner.
type RegisterOwnerInput struct {
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Password      string `json:"password" binding:"required,min=6"`
	ContactNumber string `json:"contact_number" binding:"required"`
	LibraryName   string `json:"library_name" binding:"required"`
}

// RegisterLibraryOwner creates a new library and registers the owner.
func RegisterLibraryOwner(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RegisterOwnerInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var lib models.Library
		if err := db.Where("name = ?", input.LibraryName).First(&lib).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Library already exists with this name"})
			return
		}

		lib = models.Library{Name: input.LibraryName}
		if err := db.Create(&lib).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}

		owner := models.User{
			Name:          input.Name,
			Email:         input.Email,
			Password:      string(hashedPassword),
			ContactNumber: input.ContactNumber,
			Role:          "Owner",
			LibraryID:     lib.ID,
		}
		if err := db.Create(&owner).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Owner registered successfully", "library": lib, "owner": owner})
	}
}

// AssignAdminInput is the payload for assigning admin rights.
type AssignAdminInput struct {
	Email string `json:"email" binding:"required,email"`
}

func AssignAdmin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("user").(jwt.MapClaims)
		if claims["role"] != "Owner" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Only owner can assign admin rights"})
			return
		}
		ownerLibraryID, err := getUintFromClaim(claims, "library_id")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var input AssignAdminInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var user models.User
		if err := db.Where("email = ? AND library_id = ?", input.Email, ownerLibraryID).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found in your library"})
			return
		}
		user.Role = "LibraryAdmin"
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User promoted to LibraryAdmin", "user": user})
	}
}

// RevokeAdmin revokes admin rights (sets role to Reader).
func RevokeAdmin(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        claims := c.MustGet("user").(jwt.MapClaims)
        if claims["role"] != "Owner" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Only owner can revoke admin rights"})
            return
        }

        callerID, err := getUintFromClaim(claims, "id")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        ownerLibraryID, err := getUintFromClaim(claims, "library_id")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        var input AssignAdminInput
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        var user models.User
        if err := db.Where("email = ? AND library_id = ?", input.Email, ownerLibraryID).First(&user).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found in your library"})
            return
        }

        if user.ID == callerID {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Owner cannot revoke themselves"})
            return
        }

        user.Role = "Reader"
        if err := db.Save(&user).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "LibraryAdmin rights revoked", "user": user})
    }
}
