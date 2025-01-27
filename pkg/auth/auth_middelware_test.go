package auth_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/M0rfes/go-chat-ms/pkg/auth"
	"github.com/M0rfes/go-chat-ms/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"

	pkg "github.com/M0rfes/go-chat-ms/pkg/token"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) Validate(token string) (*pkg.Claims, error) {
	args := m.Called(token)
	return args.Get(0).(*pkg.Claims), args.Error(1)
}

func (m *MockTokenService) Sign(claims *token.Claims, ttl time.Duration) (string, error) {
	args := m.Called(claims, ttl)
	return args.String(0), args.Error(1)
}

func TestTokenValidationMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockTokenService := new(MockTokenService)

	authMiddleware := auth.NewAuthMiddleware(mockTokenService)

	t.Run("valid token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
		c.Request.AddCookie(&http.Cookie{
			Name:  "token",
			Value: "validToken",
		})

		callbackCalled := 0
		mockTokenService.On("Validate", "validToken").Return(&pkg.Claims{
			UserID: "validUserID",
		}, nil)

		authMiddleware.TokenValidationMiddleware(func(c *gin.Context) {
			callbackCalled++
		})(c)
		assert.Equal(t, 0, callbackCalled)
		assert.Equal(t, http.StatusOK, w.Code)
		userID, ok := c.Get("userID")
		assert.Equal(t, true, ok)
		assert.Equal(t, "validUserID", userID.(string))

		mockTokenService.AssertExpectations(t)

	})

	t.Run("no token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

		callbackCalled := 0
		authMiddleware.TokenValidationMiddleware(func(c *gin.Context) {
			callbackCalled++
			assert.Equal(t, 1, callbackCalled)
			mockTokenService.AssertExpectations(t)
		})(c)
	})

	t.Run("invalid token without refresh token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
		c.Request.AddCookie(&http.Cookie{
			Name:  "token",
			Value: "validToken",
		})

		callbackCalled := 0
		mockTokenService.On("Validate", "validToken").Return(&pkg.Claims{}, fmt.Errorf("invalid token"))

		authMiddleware.TokenValidationMiddleware(func(c *gin.Context) {
			callbackCalled++
			assert.Equal(t, 1, callbackCalled)
			mockTokenService.AssertExpectations(t)
		})(c)

	})

	t.Run("when token and refresh both are invalid token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
		c.Request.AddCookie(&http.Cookie{
			Name:  "token",
			Value: "validToken",
		})

		callbackCalled := 0
		mockTokenService.On("Validate", "validToken").Return(&pkg.Claims{}, fmt.Errorf("invalid token"))

		authMiddleware.TokenValidationMiddleware(func(c *gin.Context) {
			callbackCalled++
			assert.Equal(t, 1, callbackCalled)
			mockTokenService.AssertExpectations(t)
		})(c)
	})

	t.Run("when token is invalid and refresh token is valid", func(t *testing.T) {
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

		callbackCalled := 0
		mockTokenService.Calls = nil
		mockTokenService.On("Validate", "validToken").Return(&pkg.Claims{}, fmt.Errorf("invalid token")).Once()
		mockTokenService.On("Validate", "validRefreshToken").Return(&pkg.Claims{
			UserID: "validUserID",
		}, nil).Once()
		mockTokenService.On("Sign", mock.Anything, 300*time.Second).Return("validToken", nil).Once()
		mockTokenService.On("Sign", mock.Anything, 30*24*time.Hour).Return("validRefreshToken", nil).Once()

		authMiddleware.TokenValidationMiddleware(func(c *gin.Context) {
			callbackCalled++
		})(c)
		assert.Equal(t, 0, callbackCalled)
		userID, ok := c.Get("userID")
		assert.Equal(t, true, ok)
		assert.Equal(t, "validUserID", userID.(string))
		// mockTokenService.AssertExpectations(t)
		// FIX later
	})

	t.Run("token signing error", func(t *testing.T) {
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

		callbackCalled := 0
		mockTokenService.On("Validate", "validToken").Return(&pkg.Claims{}, fmt.Errorf("invalid token")).Once()
		mockTokenService.On("Validate", "validRefreshToken").Return(&pkg.Claims{
			UserID: "validUserID",
		}, nil).Once()
		mockTokenService.On("Sign", mock.Anything, 300*time.Second).Return("", fmt.Errorf("signing error")).Once()

		authMiddleware.TokenValidationMiddleware(func(c *gin.Context) {
			callbackCalled++
			assert.Equal(t, 1, callbackCalled)
			mockTokenService.AssertExpectations(t)
		})(c)

	})

	t.Run("refresh token signing error", func(t *testing.T) {
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

		callbackCalled := 0
		mockTokenService.On("Validate", "validToken").Return(&pkg.Claims{}, fmt.Errorf("invalid token")).Once()
		mockTokenService.On("Validate", "validRefreshToken").Return(&pkg.Claims{
			UserID: "validUserID",
		}, nil).Once()
		mockTokenService.On("Sign", mock.Anything, 300*time.Second).Return("validToken", nil).Once()
		mockTokenService.On("Sign", mock.Anything, 30*24*time.Hour).Return("", fmt.Errorf("signing error")).Once()

		authMiddleware.TokenValidationMiddleware(func(c *gin.Context) {
			callbackCalled++
			assert.Equal(t, 1, callbackCalled)
			mockTokenService.AssertExpectations(t)
		})(c)
	})

}
