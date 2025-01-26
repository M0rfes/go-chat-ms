package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/M0rfes/go-chat-ms/ui/controllers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestChatPage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("renders chat page", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		var mockHTMLFunc controllers.HTML = func(code int, name string, obj interface{}) {
			assert.Equal(t, http.StatusOK, code)
			assert.Equal(t, "chat", name)
			assert.Equal(t, gin.H{
				"title": "Chat Page",
				"ws":    "ws://localhost:3000/chat/ws",
			}, obj)
		}

		controllers.NewChatController().ChatPage(c, mockHTMLFunc)
	})
}
