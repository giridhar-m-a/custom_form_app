package services

import (
	"context"
	"log"

	"github.com/giridhar-m-a/custom_form_app/configs"
	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
)

type GoogleAuthService interface {
	Authenticate(ctx context.Context, code string) (sqlc.GetUserByGoogleIdRow, error)
}

type googleAuthService struct {
	service UserService
}

func NewGoogleAuthService(service UserService) GoogleAuthService {
	return &googleAuthService{service: service}
}

func (s *googleAuthService) Authenticate(ctx context.Context, code string) (sqlc.GetUserByGoogleIdRow, error) {
	token, err := configs.GoogleOauthConfig.Exchange(ctx, code)
	if err != nil {
		log.Printf("GoogleAuthService: token exchange error: %v", err)
		return sqlc.GetUserByGoogleIdRow{}, err
	}

	userInfo, err := utils.GetUserDetails(token.AccessToken)
	if err != nil {
		log.Printf("GoogleAuthService: fetch user info error: %v", err)
		return sqlc.GetUserByGoogleIdRow{}, err
	}

	existingUser, err := s.service.GetUserDetailsByGoogleId(ctx, userInfo["id"].(string))
	if err == nil {
		return existingUser, nil
	}

	existingMailUser, err := s.service.GetUserDetailsByEmail(ctx, userInfo["email"].(string))
	id := userInfo["id"].(string)

	dto := dto.UserUpdateDTO{
		UserGoogleId: id,
	}

	if err == nil && existingMailUser.UserID.String() != "" && existingMailUser.UserEmail.String == userInfo["email"].(string) {
		updatedUser, err := s.service.UpdateUser(ctx, existingUser.UserID.String(), dto)
		if err != nil {
			log.Printf("GoogleAuthService: update user error: %v", err)
			return sqlc.GetUserByGoogleIdRow{}, err
		}
		return sqlc.GetUserByGoogleIdRow{
			UserID:        updatedUser.UserID,
			UserEmail:     updatedUser.UserEmail,
			UserFullName:  updatedUser.UserFullName,
			UserCreatedAt: updatedUser.UserCreatedAt,
			UserUpdatedAt: updatedUser.UserUpdatedAt,
		}, nil
	}

	newUser, err := s.service.CreateUser(ctx, userInfo)
	if err != nil {
		log.Printf("GoogleAuthService: create user error: %v", err)
		return sqlc.GetUserByGoogleIdRow{}, err
	}
	return sqlc.GetUserByGoogleIdRow{
		UserID:        newUser.UserID,
		UserEmail:     newUser.UserEmail,
		UserFullName:  newUser.UserFullName,
		UserCreatedAt: newUser.UserCreatedAt,
		UserUpdatedAt: newUser.UserUpdatedAt,
	}, nil
}
