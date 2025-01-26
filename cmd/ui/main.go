package main

import (
	"github.com/M0rfes/go-chat-ms/ui/controllers"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

var port string

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

	// Routes
	router.GET("/health", controllers.Health)
	router.GET("/chat", controllers.ChatPage)
	router.GET("/admin", controllers.AdminPage)
	router.GET("/", controllers.IndexPage)

	// Start the server
	router.Run(":" + port)
}
