package main

import (
	"github.com/Divyshekhar/golang-coding-assessment/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDb()
	// initializers.Migrate()
}

func main() {
	router := gin.Default()

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Server is healthy"})
	})

	router.Run(":8080")
}
