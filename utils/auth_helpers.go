package utils

import (
	"net/http"

	"github.com/Divyshekhar/golang-coding-assessment/initializers"
	"github.com/Divyshekhar/golang-coding-assessment/models"
	"github.com/gin-gonic/gin"
)

func GetUserAndCheckRole(ctx *gin.Context, expectedRole string) (*models.User, bool) {
	userId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return nil, false
	}

	userID := userId.(uint)

	var user models.User
	if err := initializers.Db.Where("id = ?", userID).First(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not find the user"})
		return nil, false
	}

	if user.Role != expectedRole {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return nil, false
	}

	return &user, true
}
