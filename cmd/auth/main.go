package main

import (
	"github.com/M0rfes/go-chat-ms/auth/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Routes
	router.POST("/auth/login", controllers.Login)
	router.POST("/auth/admin", controllers.AdminLogin)

	// Start the server
	router.Run(":8081")
}
