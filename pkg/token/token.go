package token

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token interface {
	Sign(claims *Claims, ttl time.Duration) (string, error)
	Validate(tokenString string) (*Claims, error)
}

type token struct {
	secret []byte
}

// NewAuth creates a new instance of Token with a secret key.
// The secret key is retrieved using the getSecretKey function.
// Returns a Token interface.
func NewAuth() Token {
	return &token{
		secret: getSecretKey(),
	}
}

// getSecretKey retrieves the secret key from the environment variable "TOKEN_SECRET_KEY".
// If the environment variable is not set, it panics with an appropriate error message.
// Returns the secret key as a byte slice.
func getSecretKey() []byte {
	secretKey := os.Getenv("TOKEN_SECRET_KEY")
	if secretKey == "" {
		panic("SECRET_KEY environment variable is not set")
	}
	return []byte(secretKey)
}

type Claims struct {
	UserID string `json:"user_id"` // user ID
	jwt.RegisteredClaims
}

// Sign generates a signed JWT token with the given claims and time-to-live (ttl).
// The token is signed using the secret key associated with the token instance.
//
// Parameters:
//   - claims: The claims to be included in the JWT token.
//   - ttl: The time-to-live for the token (duration for which the token is valid).
//
// Returns:
//   - A signed JWT token as a string.
//   - An error if there was an issue signing the token.
func (a *token) Sign(claims *Claims, ttl time.Duration) (string, error) {
	secretKey := a.secret
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(ttl))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// Validate parses and validates a JWT token string. It returns the claims
// extracted from the token if it is valid, or an error if the token is invalid
// or if there is an issue during parsing.
//
// Parameters:
//   - tokenString: the JWT token string to be validated.
//
// Returns:
//   - *Claims: a pointer to the Claims struct containing the user ID if the token is valid.
//   - error: an error if the token is invalid or if there is an issue during parsing.
func (a *token) Validate(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return &Claims{
			UserID: userID,
		}, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
