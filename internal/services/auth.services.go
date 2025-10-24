package services

import (
	"context"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
)

type AuthService interface {
	AuthenticateWithGoogle(ctxt context.Context, code string) (sqlc.GetUserByGoogleIdRow, error)
}

type authService struct {
	googleAuth GoogleAuthService
}

// NewAuthService creates an AuthService that delegates Google authentication to the provided GoogleAuthService.
// The returned service uses the given googleAuth to perform AuthenticateWithGoogle calls.
func NewAuthService(googleAuth GoogleAuthService) AuthService {
	return &authService{
		googleAuth: googleAuth,
	}
}

func (a *authService) AuthenticateWithGoogle(ctxt context.Context, code string) (sqlc.GetUserByGoogleIdRow, error) {
	return a.googleAuth.Authenticate(ctxt, code)
}