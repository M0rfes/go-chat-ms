package main

import (
	"net/http"
	"os"

	auth "github.com/M0rfes/go-chat-ms/pkg/auth"
	pkg "github.com/M0rfes/go-chat-ms/pkg/token"
	"github.com/M0rfes/go-chat-ms/ui/controllers"
	"github.com/M0rfes/go-chat-ms/ui/services"
	"github.com/gin-contrib/multitemplate"
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

func createRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("login", "templates/base.html", "templates/index.html")
	r.AddFromFiles("chat", "templates/base.html", "templates/chat.html")
	r.AddFromFiles("admin", "templates/base.html", "templates/admin.html")
	return r
}

func main() {
	router := gin.Default()

	// Set up templating engine
	router.HTMLRender = createRender()

	tokenService := pkg.NewTokenService(secret)

	userService := services.NewUserService(tokenService)
	indexController := controllers.NewIndexController(userService)
	chatController := controllers.NewChatController()
	authMiddleware := auth.NewAuthMiddleware(tokenService)
	cb := authMiddleware.TokenValidationMiddleware(
		func(ctx *gin.Context) {
			ctx.Redirect(http.StatusTemporaryRedirect, "/")
		},
	)
	// Routes
	router.GET("/health", controllers.Health)
	router.GET("/chat-page",
		cb, func(ctx *gin.Context) {
			chatController.ChatPage(ctx, nil)
		})
	router.GET("/admin", cb, controllers.AdminPage)
	router.GET("/", func(ctx *gin.Context) {
		indexController.IndexPage(ctx, nil)
	})

	// Start the server
	router.Run(":" + port)
}
