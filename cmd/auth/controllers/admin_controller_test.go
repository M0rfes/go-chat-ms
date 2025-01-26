package controllers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/M0rfes/go-chat-ms/auth/controllers"
	"github.com/M0rfes/go-chat-ms/auth/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAdminService struct {
	mock.Mock
}

func (m *mockAdminService) Login(req services.LoginRequest) (*services.LoginResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*services.LoginResponse), args.Error(1)
}

func (m *mockAdminService) Refresh(refreshToken string) (*services.LoginResponse, error) {
	args := m.Called(refreshToken)
	return args.Get(0).(*services.LoginResponse), args.Error(1)
}

func TestAdminControllers_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockAdminService)
	adminController := controllers.NewAdminController(mockService)

	t.Run("successful login", func(t *testing.T) {
		mockService.On("Login", mock.Anything).Return(&services.LoginResponse{
			Token:        "admin-token",
			RefreshToken: "admin-refresh-token",
		}, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/admin/login", bytes.NewBufferString(`{"username":"admin","password":"admin"}`))
		c.Request.Header.Set("Content-Type", "application/json")

		adminController.Login(c)

		result := w.Result()
		cookies := result.Cookies()

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"token":"admin-token","refreshToken":"admin-refresh-token"}`, w.Body.String())

		assertCookie(t, cookies, "token", "admin-token", 60*60*24)
		assertCookie(t, cookies, "refresh-token", "admin-refresh-token", 60*60*24*30)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		mockService.On("Login", mock.Anything).Return(&services.LoginResponse{}, &services.InvalidCredentialsError{}).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/admin/login", bytes.NewBufferString(`{"username":"admin","password":"wrong"}`))
		c.Request.Header.Set("Content-Type", "application/json")

		adminController.Login(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, `{"error":"invalid credentials"}`, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/admin/login", bytes.NewBufferString(`{"username":"admin"`))
		c.Request.Header.Set("Content-Type", "application/json")

		adminController.Login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("server error", func(t *testing.T) {
		mockService.On("Login", mock.Anything).Return(&services.LoginResponse{}, errors.New("database error")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/admin/login", bytes.NewBufferString(`{"username":"admin","password":"admin"}`))
		c.Request.Header.Set("Content-Type", "application/json")

		adminController.Login(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, `{"error":"database error"}`, w.Body.String())
		mockService.AssertExpectations(t)
	})
}

func TestAdminControllers_Refresh(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockAdminService)
	adminController := controllers.NewAdminController(mockService)

	t.Run("successful refresh", func(t *testing.T) {
		mockService.On("Refresh", "valid-refresh-token").Return(&services.LoginResponse{
			Token:        "new-admin-token",
			RefreshToken: "new-admin-refresh-token",
		}, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/admin/refresh", nil)
		c.Request.Header.Set("Authorization", "valid-refresh-token")

		adminController.Refresh(c)

		result := w.Result()
		cookies := result.Cookies()

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"token":"new-admin-token","refreshToken":"new-admin-refresh-token"}`, w.Body.String())

		assertCookie(t, cookies, "token", "new-admin-token", 60*60*24)
		assertCookie(t, cookies, "refresh-token", "new-admin-refresh-token", 60*60*24*30)

		mockService.AssertExpectations(t)
	})

	t.Run("missing refresh token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/admin/refresh", nil)

		adminController.Refresh(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, `{"error":"missing refresh token"}`, w.Body.String())
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		mockService.On("Refresh", "invalid-token").Return(&services.LoginResponse{}, errors.New("invalid refresh token")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/admin/refresh", nil)
		c.Request.Header.Set("Authorization", "invalid-token")

		adminController.Refresh(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, `{"error":"invalid refresh token"}`, w.Body.String())
		mockService.AssertExpectations(t)
	})
}
