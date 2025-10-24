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

// MapUser converts a sqlc.User into the API User DTO.
// Profile picture and profile picture ID are set to nil; created/updated timestamps are taken from the source Time fields.
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

// MapGetUserByEmailRow maps a sqlc.GetUserByEmailRow into a User DTO for API responses.
// The resulting User contains the ID as a string, email, full name, optional profile picture and profile picture ID (may be nil), and created/updated timestamps.
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

// MapGetUserByGoogleIdRow maps a sqlc.GetUserByGoogleIdRow to the User API DTO.
// 
// The resulting User contains stringified UserID, email, full name, optional profile
// picture and profile picture ID (converted to pointers when null), and created/updated timestamps.
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

// MapGetUserByIDRow maps a GetUserByIDRow database result to the API User DTO.
// Profile picture name and ID are converted to pointers when present; created and updated timestamps are taken from the row.
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