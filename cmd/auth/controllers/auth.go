package controllers

import (
	"github.com/gin-gonic/gin"
)

// UserLogin handles the user login process.
// It expects a JSON payload with "username" and "password" fields.
// If the credentials are valid, it returns a signed JWT token and a refresh token.
// If the credentials are invalid, it returns an unauthorized error.
//
// @Summary User login
// @Description User login with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param admin body struct{ Username string `json:"username"`; Password string `json:"password"` } true "User credentials"
// @Success 200 {object} map[string]interface{} "User login successful"
// @Failure 400 {object} map[string]interface{} "Invalid JSON"
// @Failure 401 {object} map[string]interface{} "Invalid admin credentials"
// @Failure 500 {object} map[string]interface{} "Failed to sign token or refresh token"
// @Router /login/user [post]
func UserLogin(c *gin.Context) {

}

// AdminLogin handles the admin login process.
// It expects a JSON payload with "username" and "password" fields.
// If the credentials are valid, it returns a signed JWT token and a refresh token.
// If the credentials are invalid, it returns an unauthorized error.
//
// @Summary Admin login
// @Description Admin login with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param admin body struct{ Username string `json:"username"`; Password string `json:"password"` } true "Admin credentials"
// @Success 200 {object} map[string]interface{} "Admin login successful"
// @Failure 400 {object} map[string]interface{} "Invalid JSON"
// @Failure 401 {object} map[string]interface{} "Invalid admin credentials"
// @Failure 500 {object} map[string]interface{} "Failed to sign token or refresh token"
// @Router /login/admin [post]
func AdminLogin(c *gin.Context) {

}
