package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Authenticate user (mock logic)
	if user.Username == "user" && user.Password == "password" {
		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": "user-token"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

func AdminLogin(c *gin.Context) {
	var admin struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Authenticate admin (mock logic)
	if admin.Username == "admin" && admin.Password == "admin" {
		c.JSON(http.StatusOK, gin.H{"message": "Admin login successful", "token": "admin-token"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid admin credentials"})
	}
}
