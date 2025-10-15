package services

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userID string, expiresIn time.Duration, audience string) (string, error)
	ValidateToken(token string) (string, error)
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService() JWTService {
	secretKey := os.Getenv("JWT_SECRET")
	issuer := os.Getenv("JWT_ISSUER")
	if secretKey == "" {
		secretKey = "default_secret_key" // Use a default value or handle error appropriately
	}
	if issuer == "" {
		issuer = "my_app" // Use a default value or handle error appropriately
	}
	return &jwtService{
		secretKey: secretKey,
		issuer:    issuer,
	}
}

func (j *jwtService) GenerateToken(userID string, expiresIn time.Duration, audience string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud": audience,
		"iss": j.issuer,
		"sub": userID,
		"exp": time.Now().Add(expiresIn).Unix(),
		"iat": time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *jwtService) ValidateToken(tokenString string) (string, error) {
	// Expect "Bearer <token>"
	parts := strings.Split(tokenString, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid token format")
	}
	rawToken := strings.TrimSpace(parts[1])

	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return "", err
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"].(string); ok {
			return sub, nil
		}
		return "", fmt.Errorf("invalid subject claim")
	}

	return "", fmt.Errorf("invalid token claims")
}


