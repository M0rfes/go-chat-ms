package main

import (
	"os"

	"github.com/M0rfes/go-chat-ms/auth/controllers"
	"github.com/M0rfes/go-chat-ms/auth/services"
	pkg "github.com/M0rfes/go-chat-ms/pkg/token"
	"github.com/gin-gonic/gin"
)

var (
	secret string
	port   string
)

func init() {
	secret = os.Getenv("TOKEN_SECRET")
	if secret == "" {
		panic("TOKEN_SECRET is required")
	}
}

func main() {
	router := gin.Default()

	tokenService := pkg.NewTokenService(secret)

	userService := services.NewUserService(tokenService)
	adminService := services.NewAdminService(tokenService)

	userController := controllers.NewUserController(userService)
	adminController := controllers.NewAdminController(adminService)

	// Routes
	router.GET("/health", controllers.Health)
	router.POST("/user/login", userController.Login)
	router.GET("/user/refresh", userController.Refresh)

	router.POST("/admin/login", adminController.Login)
	router.GET("/admin/refresh", adminController.Refresh)

	// Start the server
	router.Run(":" + port)
}
