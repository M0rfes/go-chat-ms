package main

import (
	"github.com/M0rfes/go-chat-ms/ui/controllers"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func createRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("login", "templates/base.html", "templates/index.html")
	return r
}

func main() {
	router := gin.Default()

	// Set up templating engine
	router.HTMLRender = createRender()

	// Routes
	router.GET("/health", controllers.Health)
	router.GET("/", controllers.IndexPage)

	// Start the server
	router.Run(":8080")
}
