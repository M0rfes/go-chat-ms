package auth_test

import (
	"os"
	"testing"
	"time"

	"github.com/M0rfes/go-chat-ms/pkg/auth"
)

func TestSignToken(t *testing.T) {
	origSecret := os.Getenv("TOKEN_SECRET_KEY")
	defer os.Setenv("TOKEN_SECRET_KEY", origSecret)

	t.Run("panics when secret not set", func(t *testing.T) {
		os.Unsetenv("TOKEN_SECRET_KEY")
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic but got none")
			}
		}()
		_, _ = auth.SignToken(auth.Claims{UserID: "123"}, 300)
	})

	t.Run("signs token when secret is set", func(t *testing.T) {
		os.Setenv("TOKEN_SECRET_KEY", "test-secret")
		token, err := auth.SignToken(auth.Claims{
			UserID: "123",
		}, 300)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if token == "" {
			t.Error("Expected token to be non-empty")
		}
	})
}

func TestValidateToken(t *testing.T) {
	origSecret := os.Getenv("TOKEN_SECRET_KEY")
	defer os.Setenv("TOKEN_SECRET_KEY", origSecret)

	t.Run("validates token with correct secret", func(t *testing.T) {
		os.Setenv("TOKEN_SECRET_KEY", "test-secret")
		claims := auth.Claims{
			UserID: "123",
		}
		token, err := auth.SignToken(claims, 300)
		if err != nil {
			t.Fatalf("Failed to sign token: %v", err)
		}

		validatedClaims, err := auth.ValidateToken(token)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if validatedClaims.UserID != "123" {
			t.Errorf("Expected user_id 123, got %s", validatedClaims.UserID)
		}
	})

	t.Run("fails validation with wrong secret", func(t *testing.T) {
		os.Setenv("TOKEN_SECRET_KEY", "test-secret")
		claims := auth.Claims{
			UserID: "123",
		}
		token, err := auth.SignToken(claims, 300)
		if err != nil {
			t.Fatalf("Failed to sign token: %v", err)
		}

		os.Setenv("TOKEN_SECRET_KEY", "different-secret")
		_, err = auth.ValidateToken(token)
		if err == nil {
			t.Error("Expected error but got none")
		}
	})

	t.Run("fails validation with expired token", func(t *testing.T) {
		os.Setenv("TOKEN_SECRET_KEY", "test-secret")
		claims := auth.Claims{
			UserID: "123",
		}
		token, err := auth.SignToken(claims, 1)
		if err != nil {
			t.Fatalf("Failed to sign token: %v", err)
		}

		// Wait for token to expire
		time.Sleep(2 * time.Second)

		_, err = auth.ValidateToken(token)
		if err == nil {
			t.Error("Expected error for expired token but got none")
		}
	})
}
