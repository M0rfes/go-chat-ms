package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/M0rfes/go-chat-ms/chat/controllers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	hc := controllers.NewHealthController()
	router.GET("/health", hc.Health)

	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status": "chat up"}`, w.Body.String())
}
