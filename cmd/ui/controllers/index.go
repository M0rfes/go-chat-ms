package controllers

import (
	"net/http"

	"github.com/M0rfes/go-chat-ms/ui/services"
	"github.com/gin-gonic/gin"
)

type HTML func(code int, name string, obj any)

type IndexController interface {
	IndexPage(c *gin.Context, html HTML)
}

type indexController struct {
	userService services.UserService
}

func NewIndexController(userService services.UserService) IndexController {
	return &indexController{
		userService: userService,
	}
}

func (i *indexController) IndexPage(c *gin.Context, html HTML) {
	token, err := c.Cookie("token")

	if html == nil {
		html = c.HTML
	}
	if err != nil {
		html(http.StatusOK, "login", gin.H{
			"title":        "Login Page",
			"user_type":    "User",
			"auth_url":     "/auth/user/login",
			"redirect_url": "/chat-page",
		})
		return
	}
	refreshToken, err := c.Cookie("refresh-token")

	if err != nil {
		html(http.StatusOK, "login", gin.H{
			"title":        "Login Page",
			"user_type":    "User",
			"auth_url":     "/auth/user/login",
			"redirect_url": "/chat-page",
		})
		return
	}
	resp, err := i.userService.LoginCheck(token, refreshToken)

	if err != nil {
		html(http.StatusOK, "login", gin.H{
			"title":        "Login Page",
			"user_type":    "User",
			"auth_url":     "/auth/user/login",
			"redirect_url": "/chat-page",
		})
		return
	}
	c.SetCookie("token", resp.Token, 300, "/", "localhost", false, true)
	c.SetCookie("refresh_token", resp.RefreshToken, 30*24*60*60, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/chat-page")
}
