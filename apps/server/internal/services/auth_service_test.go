package services

import (
	"context"
	"database/sql"
	"testing"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthenticateWithEmailPassword(t *testing.T) {
	mockUser := new(MockUserService)
	mockHash := new(MockBcryptService)
	service := &authService{
		userService: mockUser,
		hashService: mockHash,
	}

	payload := dto.EmailPasswordAuthRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedUser := sqlc.GetUserByEmailRow{
		UserEmail:    sql.NullString{String: "test@example.com", Valid: true},
		UserPassword: sql.NullString{String: "hashed_password", Valid: true},
	}

	mockUser.On("GetUserDetailsByEmail", mock.Anything, payload.Email).Return(expectedUser, nil)
	mockHash.On("ComparePassword", "hashed_password", payload.Password).Return(true)

	user, err := service.AuthenticateWithEmailPassword(context.Background(), payload)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.UserEmail.String, user.UserEmail.String)
	mockUser.AssertExpectations(t)
	mockHash.AssertExpectations(t)
}

func TestCreateUserWithEmailPassword(t *testing.T) {
	mockUser := new(MockUserService)
	mockHash := new(MockBcryptService)
	service := &authService{
		userService: mockUser,
		hashService: mockHash,
	}

	data := dto.EmailPasswordRegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	hashedPassword := "hashed_password"
	mockHash.On("HashPassword", data.Password).Return(hashedPassword, nil)
	
	expectedUser := sqlc.User{
		UserID:       uuid.New(),
		UserFullName: data.Name,
	}
	mockUser.On("CreateUser", mock.Anything, mock.Anything).Return(expectedUser, nil)

	user, err := service.CreateUserWithEmailPassword(context.Background(), data)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.UserFullName, user.UserFullName)
	mockHash.AssertExpectations(t)
	mockUser.AssertExpectations(t)
}

func TestGenerateTokens(t *testing.T) {
	mockJWT := new(MockJWTService)
	service := &authService{
		jwtService: mockJWT,
	}

	userID := "user-123"
	audience := "app"

	mockJWT.On("GenerateToken", userID, mock.Anything, audience).Return("access_token", nil).Once()
	mockJWT.On("GenerateToken", userID, mock.Anything, audience).Return("refresh_token", nil).Once()

	accessToken, refreshToken, err := service.GenerateTokens(userID, audience)

	assert.NoError(t, err)
	assert.Equal(t, "access_token", accessToken)
	assert.Equal(t, "refresh_token", refreshToken)
	mockJWT.AssertExpectations(t)
}
