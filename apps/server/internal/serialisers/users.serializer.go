package serializers

import (
	"time"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
)

// API response DTO
type User struct {
	UserID           string    `json:"id"`
	UserEmail        string    `json:"email"`
	UserFullName     string    `json:"fullName"`
	UserProfilePic   *string   `json:"profilePic"`
	UserProfilePicID *string   `json:"profilePicId"`
	UserCreatedAt    time.Time `json:"createdAt"`
	UserUpdatedAt    time.Time `json:"updatedAt"`
}

// Map sqlc.User to API response
func MapUser(u sqlc.User) User {
	return User{
		UserID:           u.UserID.String(),
		UserEmail:        u.UserEmail,
		UserFullName:     u.UserFullName,
		UserProfilePic:   nil,
		UserProfilePicID: nil,
		UserCreatedAt:    u.UserCreatedAt.Time,
		UserUpdatedAt:    u.UserUpdatedAt.Time,
	}
}

// Map GetUserByEmailRow to API response
func MapGetUserByEmailRow(u sqlc.GetUserByEmailRow) User {
	return User{
		UserID:           u.UserID.String(),
		UserEmail:        u.UserEmail,
		UserFullName:     u.UserFullName,
		UserProfilePic:   utils.NullStringToPtr(u.UserProfilePicName),
		UserProfilePicID: utils.NullUUIDToPtr(u.UserProfilePicID),
		UserCreatedAt:    u.UserCreatedAt.Time,
		UserUpdatedAt:    u.UserUpdatedAt.Time,
	}
}

// Map GetUserByGoogleIdRow to API response
func MapGetUserByGoogleIdRow(u sqlc.GetUserByGoogleIdRow) User {
	return User{
		UserID:           u.UserID.String(),
		UserEmail:        u.UserEmail,
		UserFullName:     u.UserFullName,
		UserProfilePic:   utils.NullStringToPtr(u.UserProfilePicName),
		UserProfilePicID: utils.NullUUIDToPtr(u.UserProfilePicID),
		UserCreatedAt:    u.UserCreatedAt.Time,
		UserUpdatedAt:    u.UserUpdatedAt.Time,
	}
}

// Map GetUserByIDRow to API response
func MapGetUserByIDRow(u sqlc.GetUserByIDRow) User {
	return User{
		UserID:           u.UserID.String(),
		UserEmail:        u.UserEmail,
		UserFullName:     u.UserFullName,
		UserProfilePic:   utils.NullStringToPtr(u.UserProfilePicName),
		UserProfilePicID: utils.NullUUIDToPtr(u.UserProfilePicID),
		UserCreatedAt:    u.UserCreatedAt.Time,
		UserUpdatedAt:    u.UserUpdatedAt.Time,
	}
}
