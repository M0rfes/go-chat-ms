package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/M0rfes/go-chat-ms/auth/controllers"
	serves "github.com/M0rfes/go-chat-ms/auth/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) Login(req serves.LoginRequest) (*serves.LoginResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*serves.LoginResponse), args.Error(1)
}

func (m *mockUserService) Refresh(refreshToken string) (*serves.LoginResponse, error) {
	args := m.Called(refreshToken)
	return args.Get(0).(*serves.LoginResponse), args.Error(1)
}

func assertCookie(t *testing.T, cookies []*http.Cookie, name, expectedValue string, expectedMaxAge int) {
	var cookie *http.Cookie
	for _, c := range cookies {
		if c.Name == name {
			cookie = c
			break
		}
	}

	assert.NotNil(t, cookie, "Cookie %s should exist", name)
	if cookie != nil {
		assert.Equal(t, expectedValue, cookie.Value)
		assert.Equal(t, expectedMaxAge, cookie.MaxAge)
		assert.Equal(t, "/", cookie.Path)
		assert.Equal(t, "localhost", cookie.Domain)
		assert.False(t, cookie.Secure)
		assert.True(t, cookie.HttpOnly)
	}
}

func TestUserControllers_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockUserService)
	userController := controllers.NewUserController(mockService)

	t.Run("successful login", func(t *testing.T) {
		mockService.On("Login", mock.Anything).Return(&serves.LoginResponse{
			Token:        "access-token",
			RefreshToken: "refresh-token",
		}, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/login", bytes.NewBufferString(`{"username":"test","password":"test"}`))
		c.Request.Header.Set("Content-Type", "application/json")

		userController.Login(c)

		result := w.Result()
		cookies := result.Cookies()

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"token":"access-token","refreshToken":"refresh-token"}`, w.Body.String())

		assertCookie(t, cookies, "token", "access-token", 60*60*24)
		assertCookie(t, cookies, "refresh-token", "refresh-token", 60*60*24*30)

		mockService.AssertExpectations(t)
	})

	t.Run("missing request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/login", nil)

		userController.Login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, `{"error":"EOF"}`, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/login", bytes.NewBufferString(`{"username":"test"`))
		c.Request.Header.Set("Content-Type", "application/json")

		userController.Login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, `{"error":"unexpected EOF"}`, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("login error", func(t *testing.T) {
		mockService.On("Login", mock.Anything).Return(&serves.LoginResponse{}, errors.New("invalid credentials")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/login", bytes.NewBufferString(`{"username":"test","password":"test"}`))
		c.Request.Header.Set("Content-Type", "application/json")

		userController.Login(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid credentials", response["error"])
		assert.Empty(t, response["token"])
		assert.Empty(t, response["refreshToken"])
		mockService.AssertExpectations(t)
	})
}

func TestUserControllers_Refresh(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockUserService)
	userController := controllers.NewUserController(mockService)

	t.Run("successful refresh", func(t *testing.T) {
		mockService.On("Refresh", "refresh-token").Return(&serves.LoginResponse{
			Token:        "new-access-token",
			RefreshToken: "new-refresh-token",
		}, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/refresh", nil)
		c.Request.Header.Set("refresh-token", "refresh-token")

		userController.Refresh(c)

		result := w.Result()
		cookies := result.Cookies()

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"token":"new-access-token","refreshToken":"new-refresh-token"}`, w.Body.String())

		assertCookie(t, cookies, "token", "new-access-token", 60*60*24)

		mockService.AssertExpectations(t)
	})

	t.Run("missing refresh token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/refresh", nil)

		userController.Refresh(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, `{"error":"refresh token is required"}`, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("refresh error", func(t *testing.T) {
		mockService.On("Refresh", "refresh-token").Return(&serves.LoginResponse{}, errors.New("invalid refresh token")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/refresh", nil)
		c.Request.Header.Set("refresh-token", "refresh-token")

		userController.Refresh(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid refresh token", response["error"])
		assert.Empty(t, response["token"])
		assert.Empty(t, response["refreshToken"])
		mockService.AssertExpectations(t)
	})
}
