package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/M0rfes/go-chat-ms/ui/controllers"
	"github.com/M0rfes/go-chat-ms/ui/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) LoginCheck(token, refreshToken string) (*services.LoginCheckResponse, error) {
	args := m.Called(token, refreshToken)
	return args.Get(0).(*services.LoginCheckResponse), args.Error(1)
}

func TestIndexPage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(MockUserService)

	indexController := controllers.NewIndexController(mockUserService)

	var html controllers.HTML

	t.Run("no token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
		html = func(code int, name string, obj any) {
			assert.Equal(t, http.StatusOK, code)
			assert.Equal(t, "login", name)
			assert.Equal(t, gin.H{
				"title":        "Login Page",
				"user_type":    "User",
				"auth_url":     "/auth/user/login",
				"redirect_url": "/chat-page",
			}, obj)
		}
		indexController.IndexPage(c, html)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("no refresh token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
		c.Request.AddCookie(&http.Cookie{
			Name:  "token",
			Value: "validToken",
		})
		html = func(code int, name string, obj any) {
			assert.Equal(t, http.StatusOK, code)
			assert.Equal(t, "login", name)
			assert.Equal(t, gin.H{
				"title":        "Login Page",
				"user_type":    "User",
				"auth_url":     "/auth/user/login",
				"redirect_url": "/chat-page",
			}, obj)
		}
		indexController.IndexPage(c, html)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("valid token and refresh token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
		c.Request.AddCookie(&http.Cookie{
			Name:  "token",
			Value: "validToken",
		})
		c.Request.AddCookie(&http.Cookie{
			Name:  "refresh-token",
			Value: "validRefreshToken",
		})

		mockUserService.On("LoginCheck", "validToken", "validRefreshToken").Return(&services.LoginCheckResponse{
			Token:        "validToken",
			RefreshToken: "validRefreshToken",
		}, nil)

		indexController.IndexPage(c, html)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Equal(t, "/chat-page", w.Header().Get("Location"))
		mockUserService.AssertExpectations(t)
	})

	t.Run("invalid token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
		c.Request.AddCookie(&http.Cookie{
			Name:  "token",
			Value: "invalidToken",
		})
		c.Request.AddCookie(&http.Cookie{
			Name:  "refresh-token",
			Value: "validRefreshToken",
		})

		mockUserService.On("LoginCheck", "invalidToken", "validRefreshToken").Return(&services.LoginCheckResponse{}, assert.AnError)

		html = func(code int, name string, obj any) {
			assert.Equal(t, http.StatusOK, code)
			assert.Equal(t, "login", name)
			assert.Equal(t, gin.H{
				"title":        "Login Page",
				"user_type":    "User",
				"auth_url":     "/auth/user/login",
				"redirect_url": "/chat-page-page",
			}, obj)
		}
		indexController.IndexPage(c, html)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUserService.AssertExpectations(t)
	})
}
