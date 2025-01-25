package main

import (
	"github.com/M0rfes/go-chat-ms/auth/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Routes
	router.GET("/health", controllers.Health)
	router.POST("/auth/login", controllers.UserLogin)

	// Start the server
	router.Run(":8081")
}
