package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/Divyshekhar/golang-coding-assessment/initializers"
	"github.com/Divyshekhar/golang-coding-assessment/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(ctx *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	var user models.User
	result := initializers.Db.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	if user.Role != "doctor" && user.Role != "receptionist" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized role"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not sign token"})
		return
	}
	ctx.SetCookie("jwt_token", tokenStr, 3600*24, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})

}
