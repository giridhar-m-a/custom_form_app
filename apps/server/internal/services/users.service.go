// internal/services/user_service.go
package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
	"github.com/google/uuid"
)

type UserService interface {
	GetUserDetailsById(ctx context.Context, userID string) (sqlc.GetUserByIDRow, error)
	GetUserDetailsByEmail(ctx context.Context, email string) (sqlc.GetUserByEmailRow, error)
	GetUserDetailsByGoogleId(ctx context.Context, googleID string) (sqlc.GetUserByGoogleIdRow, error)
	CreateUser(ctx context.Context, data map[string]any) (sqlc.User, error)
	UpdateUser(ctx context.Context, user string, data dto.UserUpdateDTO) (sqlc.User, error)
	UpdateUserProfilePic(ctx context.Context, user string, data dto.FileUploadPayload) (sqlc.UpdateUserProfilePicRow, error)
	CreateUserProfilePic(ctx context.Context, user uuid.UUID, path string, size int64, fileType string) (sqlc.CreateUserProfilePicRow, error)
	DeleteUser(ctx context.Context, user string) error
	DeleteUserProfilePic(ctx context.Context, user string) error
	GetUserPassword(ctx context.Context, userID string) (string, error)
}

type userService struct {
	repo repositories.UserRepository
}

func UserServiceProvider(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, data map[string]any) (sqlc.User, error) {
	password := ""
	if data["password"] != nil {
		password = data["password"].(string)
	}

	newUser, err := s.repo.Create(ctx, sqlc.CreateUserParams{
		UserFullName: data["name"].(string),
		UserEmail:    data["email"].(string),
		UserGoogleID: utils.ConvertStringToNullString(data["id"].(string)),
		UserPassword: utils.ConvertStringToNullString(password),
	})
	return newUser, err
}

func (s *userService) GetUserDetailsById(ctx context.Context, userID string) (sqlc.GetUserByIDRow, error) {
	user, err := uuid.Parse(userID)
	if err != nil {
		return sqlc.GetUserByIDRow{}, err
	}
	return s.repo.GetByID(ctx, user)
}

func (s *userService) GetUserDetailsByEmail(ctx context.Context, email string) (sqlc.GetUserByEmailRow, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *userService) GetUserDetailsByGoogleId(ctx context.Context, googleID string) (sqlc.GetUserByGoogleIdRow, error) {

	// 2. Fallback to DB
	user, err := s.repo.GetByGoogleID(ctx, googleID)
	if err != nil {
		return sqlc.GetUserByGoogleIdRow{}, err
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, user string, data dto.UserUpdateDTO) (sqlc.User, error) {
	userID, err := utils.ConvertStringToUUID(user)
	if err != nil {
		return sqlc.User{}, err
	}

	return s.repo.UpdateUser(ctx, sqlc.UpdateUserParams{
		UserID:       userID,
		UserFullName: utils.ConvertStringToNullString(data.UserFullName),
		UserPassword: utils.ConvertStringToNullString(data.UserPassword),
		UserGoogleID: utils.ConvertStringToNullString(data.UserGoogleId),
	})
}

func (s *userService) UpdateUserProfilePic(ctx context.Context, user string, data dto.FileUploadPayload) (sqlc.UpdateUserProfilePicRow, error) {
	bucket := utils.GetEnv("MINIO_BUCKET_NAME", "custom-form-app")
	file := data.File
	fileType := data.FileInfo.Header.Get("Content-Type")
	name := data.FileInfo.Filename
	size := data.FileInfo.Size

	isValidImage, _, err := utils.IsValidImage(file)
	if err != nil {
		return sqlc.UpdateUserProfilePicRow{}, err
	}
	if !isValidImage {
		return sqlc.UpdateUserProfilePicRow{}, errors.New("invalid image")
	}

	path := fmt.Sprintf("%s/profile/%s", user, name)
	userUUID, err := utils.ConvertStringToUUID(user)
	if err != nil {
		return sqlc.UpdateUserProfilePicRow{}, err
	}

	profile, err := s.repo.GetProfilePic(ctx, userUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			_, errUpload := MinioUploadFile(bucket, path, file, size, fileType)
			if errUpload != nil {
				return sqlc.UpdateUserProfilePicRow{}, errUpload
			}
			newProfile, errCreate := s.CreateUserProfilePic(ctx, userUUID, path, size, fileType)
			if errCreate != nil {
				return sqlc.UpdateUserProfilePicRow{}, errCreate
			}
			return sqlc.UpdateUserProfilePicRow{
				FileID:   newProfile.FileID,
				FileName: newProfile.FileName,
				FileSize: newProfile.FileSize,
				FileType: newProfile.FileType,
				UserID:   newProfile.UserID,
			}, nil
		} else {
			return sqlc.UpdateUserProfilePicRow{}, err
		}
	}

	err = MinioDeleteFile(bucket, profile.FileName)
	if err != nil {
		fmt.Printf("error deleting profile pic: %v", err)
	}

	_, err = MinioUploadFile(bucket, path, file, size, fileType)
	if err != nil {
		return sqlc.UpdateUserProfilePicRow{}, err
	}

	userRes, err := s.repo.UpdateUserProfile(ctx, sqlc.UpdateUserProfilePicParams{
		FileName:       utils.ConvertStringToNullString(path),
		FileSize:       utils.Int64ToNullInt64(&size),
		UserID:         userUUID,
		FileType:       utils.ConvertStringToNullString(fileType),
		FileUploadedAt: utils.TimeToNullTime(time.Now()),
	})

	if err != nil {
		return sqlc.UpdateUserProfilePicRow{}, err
	}

	return userRes, nil

}

func (s *userService) CreateUserProfilePic(ctx context.Context, user uuid.UUID, path string, size int64, fileType string) (sqlc.CreateUserProfilePicRow, error) {

	userRes, err := s.repo.CreateUserProfilePic(ctx, sqlc.CreateUserProfilePicParams{
		FileName: path,
		FileSize: size,
		UserID:   user,
		FileType: fileType,
	})

	if err != nil {
		return sqlc.CreateUserProfilePicRow{}, err
	}

	return userRes, nil

}

func (s *userService) DeleteUser(ctx context.Context, user string) error {
	userUUID, err := utils.ConvertStringToUUID(user)
	if err != nil {
		return err
	}

	err = s.repo.DeleteUser(ctx, userUUID)

	if err != nil {
		return err
	}

	return nil

}

func (s *userService) DeleteUserProfilePic(ctx context.Context, user string) error {
	userUUID, err := utils.ConvertStringToUUID(user)
	if err != nil {
		return err
	}

	err = s.repo.DeleteUserProfilePic(ctx, userUUID)

	if err != nil {
		return err
	}

	return nil

}

func (s *userService) GetUserPassword(ctx context.Context, userID string) (string, error) {
	userUUID, err := utils.ConvertStringToUUID(userID)
	if err != nil {
		return "", err
	}

	userPassword, err := s.repo.GetUserPassword(ctx, userUUID)

	if err != nil {
		return "", err
	}

	return userPassword.String, nil

}
