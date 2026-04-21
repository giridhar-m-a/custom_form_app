package services

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userID string, expiresIn time.Duration, audience string) (string, error)
	GenerateInvitationToken(invitationID string, formID string, expiresIn time.Duration) (string, error)
	ValidateToken(token string) (string, error)
	ValidateInvitationToken(token string) (*InvitationClaims, error)
	GenerateAnonymousInvitationToken(formId string, expiresIn time.Duration) (string, error)
}

type jwtService struct {
	secretKey string
	issuer    string
}

type InvitationClaims struct {
	InvitationID *string `json:"invitationId"`
	FormID       string  `json:"formId"`
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

func (j *jwtService) GenerateInvitationToken(invitationID string, formID string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":    j.issuer,
		"sub":    invitationID,
		"formId": formID,
		"exp":    time.Now().Add(expiresIn).Unix(),
		"iat":    time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *jwtService) ValidateToken(tokenString string) (string, error) {
	if tokenString == "" {
		return "", fmt.Errorf("token is empty")
	}
	// Expect "Bearer <token>"
	tokenString = strings.ReplaceAll(tokenString, "Bearer", "")
	tokenString = strings.TrimSpace(tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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

func (j *jwtService) ValidateInvitationToken(tokenString string) (*InvitationClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}

	tokenString = strings.ReplaceAll(tokenString, "Bearer", "")
	tokenString = strings.TrimSpace(tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var invitationClaims InvitationClaims

		// ✅ optional invitationId (from "sub")
		if invitationID, ok := claims["sub"].(string); ok && invitationID != "" {
			invitationClaims.InvitationID = &invitationID
		}

		// ✅ required formId
		if formID, ok := claims["formId"].(string); ok {
			invitationClaims.FormID = formID
		} else {
			return nil, errors.New("formId missing in token")
		}

		return &invitationClaims, nil
	}

	return nil, errors.New("invalid token claims")
}

func (j *jwtService) GenerateAnonymousInvitationToken(formId string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":    j.issuer,
		"formId": formId,
		"exp":    time.Now().Add(expiresIn).Unix(),
		"iat":    time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
