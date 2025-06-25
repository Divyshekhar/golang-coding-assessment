package utils

import (
	"net/http"

	"github.com/Divyshekhar/golang-coding-assessment/initializers"
	"github.com/Divyshekhar/golang-coding-assessment/models"
	"github.com/gin-gonic/gin"
)

func GetUserAndCheckRole(ctx *gin.Context, expectedRole string) (*models.User, bool) {
	userIdVal, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return nil, false
	}

	var userID uint
	switch id := userIdVal.(type) {
	case float64:
		userID = uint(id)
	case int:
		userID = uint(id)
	case uint:
		userID = id
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID format"})
		return nil, false
	}

	var user models.User
	if err := initializers.Db.First(&user, userID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not find the user"})
		return nil, false
	}

	if user.Role != expectedRole {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return nil, false
	}

	return &user, true
}
