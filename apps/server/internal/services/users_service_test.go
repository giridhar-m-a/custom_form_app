package services

import (
	"context"
	"testing"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"database/sql"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByGoogleID(ctx context.Context, googleID string) (sqlc.GetUserByGoogleIdRow, error) {
	args := m.Called(ctx, googleID)
	return args.Get(0).(sqlc.GetUserByGoogleIdRow), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (sqlc.GetUserByEmailRow, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(sqlc.GetUserByEmailRow), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, userID uuid.UUID) (sqlc.GetUserByIDRow, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(sqlc.GetUserByIDRow), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, params sqlc.CreateUserParams) (sqlc.User, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(sqlc.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, data sqlc.UpdateUserParams) (sqlc.User, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(sqlc.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUserProfile(ctx context.Context, data sqlc.UpdateUserProfilePicParams) (sqlc.UpdateUserProfilePicRow, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(sqlc.UpdateUserProfilePicRow), args.Error(1)
}

func (m *MockUserRepository) CreateUserProfilePic(ctx context.Context, data sqlc.CreateUserProfilePicParams) (sqlc.CreateUserProfilePicRow, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(sqlc.CreateUserProfilePicRow), args.Error(1)
}

func (m *MockUserRepository) DeleteUserProfilePic(ctx context.Context, user uuid.UUID) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, user uuid.UUID) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetProfilePic(ctx context.Context, userID uuid.UUID) (sqlc.UserImage, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(sqlc.UserImage), args.Error(1)
}

func (m *MockUserRepository) GetUserPassword(ctx context.Context, userID uuid.UUID) (sql.NullString, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(sql.NullString), args.Error(1)
}

func (m *MockUserRepository) SoftDeleteUser(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepository) CreateTempUser(ctx context.Context, userFullName string) (sqlc.User, error) {
	args := m.Called(ctx, userFullName)
	return args.Get(0).(sqlc.User), args.Error(1)
}

func TestGetUserDetailsById(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := UserServiceProvider(mockRepo)

	userID := uuid.New()
	expectedUser := sqlc.GetUserByIDRow{
		UserID: userID,
		UserFullName: "John Doe",
	}

	mockRepo.On("GetByID", mock.Anything, userID).Return(expectedUser, nil)

	user, err := userService.GetUserDetailsById(context.Background(), userID.String())

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := UserServiceProvider(mockRepo)

	userData := map[string]any{
		"name": "Jane Doe",
		"email": "jane@example.com",
		"id": "google-123",
		"password": "secretpassword",
	}

	expectedUser := sqlc.User{
		UserID: uuid.New(),
		UserFullName: "Jane Doe",
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("sqlc.CreateUserParams")).Return(expectedUser, nil)

	user, err := userService.CreateUser(context.Background(), userData)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestGetUserDetailsByEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := UserServiceProvider(mockRepo)

	email := "test@example.com"
	expectedUser := sqlc.GetUserByEmailRow{
		UserID:       uuid.New(),
		UserFullName: "Test User",
	}

	mockRepo.On("GetByEmail", mock.Anything, email).Return(expectedUser, nil)

	user, err := userService.GetUserDetailsByEmail(context.Background(), email)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := UserServiceProvider(mockRepo)

	userID := uuid.New()
	updateDTO := dto.UserUpdateDTO{
		UserFullName: "Updated Name",
		UserPassword: "newpassword",
	}

	expectedUser := sqlc.User{
		UserID:       userID,
		UserFullName: "Updated Name",
	}

	mockRepo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(p sqlc.UpdateUserParams) bool {
		return p.UserID == userID && p.UserFullName.String == "Updated Name"
	})).Return(expectedUser, nil)

	user, err := userService.UpdateUser(context.Background(), userID.String(), updateDTO)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := UserServiceProvider(mockRepo)

	userID := uuid.New()

	mockRepo.On("DeleteUser", mock.Anything, userID).Return(nil)

	err := userService.DeleteUser(context.Background(), userID.String())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetUserPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := UserServiceProvider(mockRepo)

	userID := uuid.New()
	expectedPassword := "hashedpassword"

	mockRepo.On("GetUserPassword", mock.Anything, userID).Return(sql.NullString{String: expectedPassword, Valid: true}, nil)

	password, err := userService.GetUserPassword(context.Background(), userID.String())

	assert.NoError(t, err)
	assert.Equal(t, expectedPassword, password)
	mockRepo.AssertExpectations(t)
}
