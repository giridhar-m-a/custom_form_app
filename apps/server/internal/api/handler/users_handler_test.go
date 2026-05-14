package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserDetailsById(ctx context.Context, userID string) (sqlc.GetUserByIDRow, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(sqlc.GetUserByIDRow), args.Error(1)
}

func (m *MockUserService) GetUserDetailsByEmail(ctx context.Context, email string) (sqlc.GetUserByEmailRow, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(sqlc.GetUserByEmailRow), args.Error(1)
}

func (m *MockUserService) GetUserDetailsByGoogleId(ctx context.Context, googleID string) (sqlc.GetUserByGoogleIdRow, error) {
	args := m.Called(ctx, googleID)
	return args.Get(0).(sqlc.GetUserByGoogleIdRow), args.Error(1)
}

func (m *MockUserService) CreateUser(ctx context.Context, data map[string]any) (sqlc.User, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(sqlc.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, user string, data dto.UserUpdateDTO) (sqlc.User, error) {
	args := m.Called(ctx, user, data)
	return args.Get(0).(sqlc.User), args.Error(1)
}

func (m *MockUserService) UpdateUserProfilePic(ctx context.Context, user string, data dto.FileUploadPayload) (sqlc.UpdateUserProfilePicRow, error) {
	args := m.Called(ctx, user, data)
	return args.Get(0).(sqlc.UpdateUserProfilePicRow), args.Error(1)
}

func (m *MockUserService) CreateUserProfilePic(ctx context.Context, user uuid.UUID, path string, size int64, fileType string) (sqlc.CreateUserProfilePicRow, error) {
	args := m.Called(ctx, user, path, size, fileType)
	return args.Get(0).(sqlc.CreateUserProfilePicRow), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, user string) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) DeleteUserProfilePic(ctx context.Context, user string) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) GetUserPassword(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) CreateTempUser(ctx context.Context, name string) (sqlc.User, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(sqlc.User), args.Error(1)
}

func (m *MockUserService) SoftDeleteUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestGetMeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	handler := &usersHandler{userService: mockService}

	userID := uuid.New().String()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", userID)

	expectedUser := sqlc.GetUserByIDRow{
		UserID:       uuid.MustParse(userID),
		UserFullName: "John Doe",
	}

	mockService.On("GetUserDetailsById", mock.Anything, userID).Return(expectedUser, nil)

	handler.GetMe(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	handler := &usersHandler{userService: mockService}

	userID := uuid.New().String()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", userID)

	mockService.On("SoftDeleteUser", mock.Anything, userID).Return(nil)

	handler.DeleteUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
