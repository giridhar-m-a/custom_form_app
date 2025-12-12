package serializers

import (
	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
)

// Map sqlc.User to API response
func MapUser(u sqlc.User) dto.User {
	return dto.User{
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
func MapGetUserByEmailRow(u sqlc.GetUserByEmailRow) dto.User {
	return dto.User{
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
func MapGetUserByGoogleIdRow(u sqlc.GetUserByGoogleIdRow) dto.User {
	return dto.User{
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
func MapGetUserByIDRow(u sqlc.GetUserByIDRow) dto.User {
	return dto.User{
		UserID:           u.UserID.String(),
		UserEmail:        u.UserEmail,
		UserFullName:     u.UserFullName,
		UserProfilePic:   utils.NullStringToPtr(u.UserProfilePicName),
		UserProfilePicID: utils.NullUUIDToPtr(u.UserProfilePicID),
		UserCreatedAt:    u.UserCreatedAt.Time,
		UserUpdatedAt:    u.UserUpdatedAt.Time,
	}
}
