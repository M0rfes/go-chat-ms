package main

import (
	"log"
	"net/http"
	"os"

	"github.com/M0rfes/go-chat-ms/chat/controllers"
	auth "github.com/M0rfes/go-chat-ms/pkg/auth"
	token "github.com/M0rfes/go-chat-ms/pkg/token"
	"github.com/gin-gonic/gin"
)

var (
	port   string = "8080"
	secret string
)

func init() {
	secret = os.Getenv("TOKEN_SECRET")
	if secret == "" {
		panic("TOKEN_SECRET is required")
	}
}

func main() {
	r := gin.Default()
	healthController := controllers.NewHealthController()
	messageController := controllers.NewMessagesController()
	tokenService := token.NewTokenService(secret)
	authMiddleware := auth.NewAuthMiddleware(tokenService)
	cb := authMiddleware.TokenValidationMiddleware(func(c *gin.Context) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	})
	r.GET("/ws", cb, messageController.HandleWebSocket)
	r.GET("/health", healthController.Health)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
