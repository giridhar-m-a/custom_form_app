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
	"github.com/giridhar-m-a/custom_form_app/internal/services/templates"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
	"github.com/resend/resend-go/v3"
)

type AuthService interface {
	AuthenticateWithGoogle(ctxt context.Context, code string) (sqlc.GetUserByGoogleIdRow, error)
	AuthenticateWithEmailPassword(ctx context.Context, payload dto.EmailPasswordAuthRequest) (sqlc.GetUserByEmailRow, error)
	CreateUserWithEmailPassword(ctx context.Context, data dto.EmailPasswordRegisterRequest) (sqlc.User, error)
	GenerateTokens(userID string, audience string) (string, string, error)
	VerifyToken(token string) (string, error)
	RequestResetPassword(ctx context.Context, email string) (*resend.SendEmailResponse, error)
	ResetPassword(ctx context.Context, token string, newPassword string) error
}

type authService struct {
	googleAuth  GoogleAuthService
	userService UserService
	hashService BcryptService
	jwtService  JWTService
	mailService MailService
}

func NewAuthService(googleAuth GoogleAuthService, userService UserService) AuthService {
	bcryptService := NewBcryptService()
	jwtService := NewJWTService()
	mailService := NewMailService(utils.ResendClient)
	return &authService{
		googleAuth:  googleAuth,
		userService: userService,
		jwtService:  jwtService,
		hashService: bcryptService,
		mailService: mailService,
	}
}

func (a *authService) AuthenticateWithGoogle(ctxt context.Context, code string) (sqlc.GetUserByGoogleIdRow, error) {
	return a.googleAuth.Authenticate(ctxt, code)
}

func (a *authService) AuthenticateWithEmailPassword(ctx context.Context, payload dto.EmailPasswordAuthRequest) (sqlc.GetUserByEmailRow, error) {
	user, err := a.userService.GetUserDetailsByEmail(ctx, payload.Email)
	if err != nil {
		return sqlc.GetUserByEmailRow{}, errors.New("invalid credentials")
	}
	if !user.UserPassword.Valid {
		return sqlc.GetUserByEmailRow{}, errors.New("invalid credentials")
	}
	hashPassword := user.UserPassword.String
	if hashPassword == "" {
		return sqlc.GetUserByEmailRow{}, errors.New("invalid credentials")
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

func (a *authService) RequestResetPassword(
	ctx context.Context,
	email string,
) (*resend.SendEmailResponse, error) {

	frontendURL := utils.GetEnv("FRONTEND_URL", "")
	if frontendURL == "" {
		return nil, errors.New("frontend URL not set")
	}

	user, err := a.userService.GetUserDetailsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	verifyToken, err := a.jwtService.GenerateToken(
		user.UserID.String(),
		10*time.Minute,
		"reset",
	)
	if err != nil {
		return nil, err
	}

	// 1. Fixed scoping and initialization
	fromEmail := utils.GetEnv("SENDER_EMAIL_ADDRESS", "form-genius-no-reply@giridhar.dev")

	resetURL := fmt.Sprintf(
		"%s/reset-password?token=%s",
		frontendURL,
		verifyToken,
	)

	templateData := struct {
		PlatformName   string
		UserName       string
		ResetURL       string
		Year           int
		CompanyAddress string
	}{
		PlatformName:   "Form Genius",
		UserName:       user.UserFullName,
		ResetURL:       resetURL,
		Year:           time.Now().Year(),
		CompanyAddress: "Form Genius Inc",
	}

	// Note: Ensure the "templates" package is imported in your file
	templateService := templates.NewService()

	template, err := templateService.Render("password-reset.html", templateData)
	if err != nil {
		return nil, err
	}

	// 2. Correctly formatting the 'From' string
	sender := fmt.Sprintf("Form Genius <%s>", fromEmail)

	emailParams := resend.SendEmailRequest{
		From:    sender,
		To:      []string{user.UserEmail.String},
		Subject: "Password Reset Request",
		Html:    template,
		Tags: []resend.Tag{
			{
				Name:  "category",
				Value: "form-genius-password-reset",
			},
		},
	}

	return a.mailService.SendEmail(emailParams)
} // Added missing closing brace for the function

func (a *authService) ResetPassword(
	ctx context.Context,
	token string,
	newPassword string,
) error {

	userID, err := a.jwtService.ValidateToken(token)
	if err != nil {
		return err
	}

	hashedPassword, err := a.hashService.HashPassword(newPassword)
	if err != nil {
		return err
	}

	_, err = a.userService.UpdateUser(ctx, userID, dto.UserUpdateDTO{
		UserPassword: hashedPassword,
	})
	if err != nil {
		return err
	}

	return nil
} // Added missing closing brace for the function
