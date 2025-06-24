package routes

import (
	"github.com/Divyshekhar/golang-coding-assessment/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/user")
	{
		userGroup.POST("/login", controllers.Login)
		userGroup.POST("/signup", controllers.Signup)
	}
}
