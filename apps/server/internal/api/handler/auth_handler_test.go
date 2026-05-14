package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/google/uuid"
	"github.com/resend/resend-go/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Define local mocks for handler testing
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) AuthenticateWithGoogle(ctx context.Context, code string) (sqlc.GetUserByGoogleIdRow, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(sqlc.GetUserByGoogleIdRow), args.Error(1)
}

func (m *MockAuthService) AuthenticateWithEmailPassword(ctx context.Context, payload dto.EmailPasswordAuthRequest) (sqlc.GetUserByEmailRow, error) {
	args := m.Called(ctx, payload)
	return args.Get(0).(sqlc.GetUserByEmailRow), args.Error(1)
}

func (m *MockAuthService) CreateUserWithEmailPassword(ctx context.Context, data dto.EmailPasswordRegisterRequest) (sqlc.User, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(sqlc.User), args.Error(1)
}

func (m *MockAuthService) GenerateTokens(userID string, audience string) (string, string, error) {
	args := m.Called(userID, audience)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockAuthService) VerifyToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) RequestResetPassword(ctx context.Context, email string) (*resend.SendEmailResponse, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*resend.SendEmailResponse), args.Error(1)
}

func (m *MockAuthService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	args := m.Called(ctx, token, newPassword)
	return args.Error(0)
}

func TestEmailPasswordAuthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuth := new(MockAuthService)
	handler := &authHandler{
		authService: mockAuth,
	}

	payload := dto.EmailPasswordAuthRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	expectedUser := sqlc.GetUserByEmailRow{
		UserID: uuid.New(),
	}

	mockAuth.On("AuthenticateWithEmailPassword", mock.Anything, payload).Return(expectedUser, nil)
	mockAuth.On("GenerateTokens", expectedUser.UserID.String(), mock.Anything).Return("access", "refresh", nil)

	handler.EmailPasswordAuthHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Authentication successful", response["message"])
	
	mockAuth.AssertExpectations(t)
}
