package main

import (
	"log"
	"net/http"
	"os"

	"github.com/M0rfes/go-chat-ms/chat/controllers"
	auth "github.com/M0rfes/go-chat-ms/pkg/auth"
	producer "github.com/M0rfes/go-chat-ms/pkg/message-queue/producers"
	token "github.com/M0rfes/go-chat-ms/pkg/token"
	"github.com/gin-gonic/gin"
)

var (
	port     string = "8080"
	secret   string
	kafkaUrl string
)

func init() {
	secret = os.Getenv("TOKEN_SECRET")
	if secret == "" {
		panic("TOKEN_SECRET is required")
	}
	kafkaUrl = os.Getenv("KAFKA_URL")
	if kafkaUrl == "" {
		panic("KAFKA_URL is required")
	}
}

func main() {
	r := gin.Default()

	kafkaProducer, err := producer.NewMessage(kafkaUrl)

	if err != nil {
		log.Fatal("Failed to create kafka producer:", err)
	}

	healthController := controllers.NewHealthController()
	messageController := controllers.NewMessagesController(kafkaProducer)
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
