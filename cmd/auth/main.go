package main

import (
	"os"

	"github.com/M0rfes/go-chat-ms/auth/controllers"
	"github.com/M0rfes/go-chat-ms/auth/services"
	pkg "github.com/M0rfes/go-chat-ms/pkg/token"
	"github.com/gin-gonic/gin"
)

var secret string

func init() {
	secret = os.Getenv("TOKEN_SECRET")
	if secret == "" {
		panic("TOKEN_SECRET is required")
	}
}

func main() {
	router := gin.Default()

	tokenService := pkg.NewAuth(secret)

	userService := services.NewUserService(tokenService)
	adminService := services.NewAdminService(tokenService)

	userController := controllers.NewUserController(userService)
	adminController := controllers.NewAdminController(adminService)

	// Routes
	router.GET("/health", controllers.Health)
	router.POST("/auth/user/login", userController.Login)
	router.GET("/auth/user/refresh", userController.Refresh)

	router.POST("/auth/admin/login", adminController.Login)
	router.GET("/auth/admin/refresh", adminController.Refresh)

	// Start the server
	router.Run(":8081")
}
