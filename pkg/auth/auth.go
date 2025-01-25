package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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

// SignToken generates a signed JWT token with the given claims and time-to-live (ttl).
// The token is signed using the HS256 signing method and a secret key.
//
// Parameters:
//   - claims: The claims to be included in the JWT token.
//   - ttl: Seconds after which the token expires.
//
// Returns:
//   - A signed JWT token as a string.
//   - An error if there was an issue signing the token.
func SignToken(claims Claims, ttl int64) (string, error) {
	secretKey := getSecretKey()
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(ttl) * time.Second))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateToken validates a JWT token string and returns the claims if the token is valid.
// It takes a token string as input and returns a Claims struct and an error.
// If the token is invalid or the signing method is unexpected, it returns an error.
//
// Parameters:
//   - tokenString: the JWT token string to be validated.
//
// Returns:
//   - *Claims: a pointer to the Claims struct containing the user ID if the token is valid.
//   - error: an error if the token is invalid or the signing method is unexpected.
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getSecretKey(), nil
	})
	fmt.Println("=========>", err)
	if err != nil {
		return nil, err
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
