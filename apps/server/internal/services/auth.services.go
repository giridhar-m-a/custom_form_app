package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
)

type AuthService interface {
	AuthenticateWithGoogle(ctxt context.Context, code string) (sqlc.GetUserByGoogleIdRow, error)
	AuthenticateWithEmailPassword(ctx context.Context, payload dto.EmailPasswordAuthRequest) (sqlc.GetUserByEmailRow, error)
	CreateUserWithEmailPassword(ctx context.Context, data dto.EmailPasswordRegisterRequest) (sqlc.User, error)
	GenerateTokens(userID string, audience string) (string, string, error)
	VerifyToken(token string) (string, error)
}

type authService struct {
	googleAuth  GoogleAuthService
	userService UserService
	hashService BcryptService
	jwtService  JWTService
}

func NewAuthService(googleAuth GoogleAuthService, userService UserService) AuthService {
	bcryptService := NewBcryptService()
	jwtService := NewJWTService()
	return &authService{
		googleAuth:  googleAuth,
		userService: userService,
		jwtService:  jwtService,
		hashService: bcryptService,
	}
}

func (a *authService) AuthenticateWithGoogle(ctxt context.Context, code string) (sqlc.GetUserByGoogleIdRow, error) {
	return a.googleAuth.Authenticate(ctxt, code)
}

func (a *authService) AuthenticateWithEmailPassword(ctx context.Context, payload dto.EmailPasswordAuthRequest) (sqlc.GetUserByEmailRow, error) {
	user, err := a.userService.GetUserDetailsByEmail(ctx, payload.Email)
	if err != nil {
		return sqlc.GetUserByEmailRow{}, err
	}
	if !user.UserPassword.Valid {
		return sqlc.GetUserByEmailRow{}, errors.New("Invalid Credentials")
	}
	hashPassword := user.UserPassword.String
	if hashPassword == "" {
		return sqlc.GetUserByEmailRow{}, errors.New("Invalid Credentials")
	}
	isPasswordValid := a.hashService.ComparePassword(hashPassword, payload.Password)

	fmt.Printf("Password Valid: %v\n", isPasswordValid)

	if !isPasswordValid {
		return sqlc.GetUserByEmailRow{}, errors.New("invalid credentials")
	}

	return user, nil

}

func (a *authService) CreateUserWithEmailPassword(ctx context.Context, data dto.EmailPasswordRegisterRequest) (sqlc.User, error) {
	hashedPassword, err := a.hashService.HashPassword(data.Password)
	if err != nil {
		return sqlc.User{}, err
	}
	userData := map[string]any{
		"name":     data.Name,
		"email":    data.Email,
		"password": hashedPassword,
		"id":       "",
	}
	return a.userService.CreateUser(ctx, userData)
}

func (a *authService) GenerateTokens(userID string, audience string) (string, string, error) {
	jwtExpire, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	if err != nil {
		jwtExpire = 3600 // default to 1 hour if not set
	}
	refreshExpire, err := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRATION_HOURS"))
	if err != nil {
		refreshExpire = 86400 // default to 24 hours if not set
	}
	token, err := a.jwtService.GenerateToken(userID, time.Duration(jwtExpire)*1*time.Second, audience)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := a.jwtService.GenerateToken(userID, time.Duration(refreshExpire)*1*time.Second, audience)
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}

func (a *authService) VerifyToken(token string) (string, error) {
	return a.jwtService.ValidateToken(token)
}

